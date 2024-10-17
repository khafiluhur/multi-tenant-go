package middleware

import (
	"time"

	"jatis_mobile_api/logs"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func PerformanceLogger(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start)

			logs.LogWithFields(logger, logrus.InfoLevel, "Request processed", struct {
				Method   string
				Path     string
				Duration time.Duration
			}{
				Method:   c.Request().Method,
				Path:     c.Request().URL.Path,
				Duration: duration,
			})

			return err
		}
	}
}
