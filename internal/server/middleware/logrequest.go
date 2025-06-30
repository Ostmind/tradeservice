package middleware

import (
	"fmt"
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
			logger.Info("Request: %s %s %s %d\n", c.Request().Method, c.Request().URL, stop.Sub(start), c.Response().Status)
			if err != nil {
				logger.Error("request is failed", "error", err)
				return fmt.Errorf("error middleware %w", err)
			}
			return nil
		}
	}
}
