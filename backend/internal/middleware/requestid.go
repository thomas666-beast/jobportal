package middleware

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"job-portal/internal/service"
)

const ContextKeyRequestID = "request_id"

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := c.Request().Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.NewString()
			}
			c.Set(ContextKeyRequestID, reqID)
			c.Response().Header().Set("X-Request-ID", reqID)

			start := time.Now()
			err := next(c)
			duration := time.Since(start)

			attrs := []any{
				"request_id", reqID,
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"duration_ms", duration.Milliseconds(),
				"remote_ip", c.RealIP(),
			}

			if claims, ok := c.Get(ContextKeyUser).(*service.Claims); ok {
				attrs = append(attrs, "user_id", claims.UserID)
			}

			switch {
			case err != nil:
				attrs = append(attrs, "error", err.Error())
				slog.Error("request failed", attrs...)
			case c.Response().Status >= 500:
				slog.Error("request 5xx", attrs...)
			case c.Response().Status >= 400:
				slog.Warn("request 4xx", attrs...)
			default:
				slog.Info("request", attrs...)
			}

			return err
		}
	}
}
