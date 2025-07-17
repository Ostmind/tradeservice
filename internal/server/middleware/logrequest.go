package middleware

import (
	"github.com/labstack/echo/v4"
	"log/slog"
	"time"
)

func LogRequest(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			stop := time.Now()

			logger.Info("Request: ", "Method", c.Request().Method, "URL", c.Request().URL, "Time", stop.Sub(start), "Http Code", c.Response().Status)

			return err
		}
	}
}
