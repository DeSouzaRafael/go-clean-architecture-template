package rest

import (
	"net/http"

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
func NewRouter(h *echo.Echo, l logger.Interface, v *validator.Validator, uc usecase.UseCases, env string) {
	h.Use(middleware.CORSWithConfig(corsConfig(env)))
	h.Use(middleware.Recover())

	h.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	h.GET("/docs/*", echoSwagger.WrapHandler)

	v0.NewUserRoutes(h, l, v, uc.UserUseCase())
}

func corsConfig(env string) middleware.CORSConfig {
	cc := middleware.CORSConfig{
		AllowHeaders:     []string{echo.HeaderAccept, echo.HeaderAcceptEncoding, echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderXCSRFToken},
		AllowCredentials: true,
		ExposeHeaders:    []string{echo.HeaderAccept, echo.HeaderAcceptEncoding, echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderXCSRFToken},
	}

	if env == "prd" {
		cc.AllowOrigins = []string{"https://*.your.domain.com"}
	} else {
		cc.AllowOrigins = []string{"*"}
	}

	return cc
}
