package rest

import (
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/logger"
	v0 "github.com/DeSouzaRafael/go-clean-architecture-template/internal/controller/rest/routers/v0"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       go-clean-architecture-template API
// @description Template Clean Architecture Golang
// @version     1.0
// @host        localhost:8080
// @BasePath
func NewRouter(handler *echo.Echo, logger logger.Interface, port string, useCases usecase.UseCases) {

	// CORS
	handler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // all
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	handler.Use(middleware.Recover())

	handler.GET("/docs/*", echoSwagger.WrapHandler)

	v0Group := handler.Group("/v0")
	v0.NewUserRoutes(v0Group, logger, useCases.UserUseCase())

	handler.Logger.Fatal(handler.Start(port))
}
