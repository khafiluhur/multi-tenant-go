package middleware

import (
	"jatis_mobile_api/logs"

	"github.com/labstack/echo/v4"
)

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := logs.SetupLogger()
		c.Set("logger", logger)
		return next(c)
	}
}
