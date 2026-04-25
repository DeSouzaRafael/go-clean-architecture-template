package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DeSouzaRafael/go-clean-architecture-template/config"
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/httpserver"
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/logger"
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/postgres"
	postgresModel "github.com/DeSouzaRafael/go-clean-architecture-template/infra/postgres/model"
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/postgres/repository"
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/validator"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/controller/rest"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/usecase"
	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config) {
	l := logger.NewLogger(cfg.Log.Level)
	v := validator.NewValidator()

	pg, err := postgres.NewPostgres(postgres.Options{
		URL:             cfg.PG.URL,
		MaxOpenConns:    cfg.PG.MaxOpenConns,
		MaxIdleConns:    cfg.PG.MaxIdleConns,
		ConnMaxLifetime: cfg.PG.ConnMaxLifetime,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - NewPostgres: %w", err))
	}
	defer pg.Close()

	if cfg.App.Env != "prd" {
		if err := pg.DB.AutoMigrate(&postgresModel.UserModel{}); err != nil {
			l.Fatal(fmt.Errorf("app - Run - AutoMigrate: %w", err))
		}
	}

	userUseCase := usecase.NewUser(repository.NewUserRepo(pg.DB))
	appUseCases := usecase.NewAppUseCases(userUseCase)

	handler := echo.New()
	rest.NewRouter(handler, l, v, appUseCases, cfg.App.Env)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	l.Info("Server is running...")
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
		if err := httpServer.Shutdown(); err != nil {
			l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		if err := httpServer.Shutdown(); err != nil {
			l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	}
}
