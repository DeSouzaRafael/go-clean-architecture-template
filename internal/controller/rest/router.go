package rest

import (
	"os"

	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/logger"
	"github.com/DeSouzaRafael/go-clean-architecture-template/infra/validator"
	v0 "github.com/DeSouzaRafael/go-clean-architecture-template/internal/controller/rest/routers/v0"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Architecture Template API
// @description Template Golang
// @version     1.0
// @host        localhost:8080
// @BasePath
func NewRouter(h *echo.Echo, l logger.Interface, v *validator.Validator, uc usecase.UseCases) {

	h.Use(middleware.CORSWithConfig(corsConfig()))
	h.Use(middleware.Recover())

	// Swagger docs
	h.GET("/docs/*", echoSwagger.WrapHandler)

	// REST versioning
	v0.NewUserRoutes(h, l, v, uc.UserUseCase())
}

func corsConfig() middleware.CORSConfig {
	cc := middleware.CORSConfig{
		AllowHeaders:     []string{echo.HeaderAccept, echo.HeaderAcceptEncoding, echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderXCSRFToken},
		AllowCredentials: true,
		ExposeHeaders:    []string{echo.HeaderAccept, echo.HeaderAcceptEncoding, echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderXCSRFToken},
	}

	if os.Getenv("ENV") == "prd" {
		cc.AllowOrigins = []string{"https://*.your.domain.com"}
	} else {
		cc.AllowOrigins = []string{"*"}
	}

	return cc
}
