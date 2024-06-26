package output

import (
	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Error string `json:"error" example:"message"`
}

func ErrorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(code, ResponseError{Error: msg})
}
