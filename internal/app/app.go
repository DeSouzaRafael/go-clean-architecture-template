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
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/postgres/repository"
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/validator"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/controller/rest"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/entity"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/usecase"
	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config) {
	// Logger
	logger := logger.NewLogger(cfg.Log.Level)
	validator := validator.NewValidator()
	// Repository
	pg, err := postgres.NewPostgres(cfg.PG.URL)
	if err != nil {
		logger.Fatal(fmt.Errorf("app - Run - NewPostgresRepository: %w", err))
	}
	defer pg.Close()

	// Auto Migration
	if err := pg.DB.AutoMigrate(
		&entity.UserEntity{},
		// add models
	); err != nil {
		logger.Fatal(fmt.Errorf("app - Run - AutoMigrate: %w", err))
	}

	// UseCases
	userUseCase := usecase.NewUser(
		repository.NewUserRepo(pg.DB),
	)

	appUseCases := usecase.NewAppUseCases(
		userUseCase,
	)

	// HTTP Server
	handler := echo.New()

	rest.NewRouter(handler, logger, validator, appUseCases)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	logger.Info("Server is running...")
	select {
	case s := <-interrupt:
		logger.Error(fmt.Errorf("app - Run - signal: %w" + s.String()))
	case err = <-httpServer.Notify():
		logger.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))

		// Shutdown
		err = httpServer.Shutdown()
		if err != nil {
			logger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	}
}
