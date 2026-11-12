package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"job-portal/internal/service"
)

const ContextKeyUser = "user"

func JWTAuth(jwtSvc *service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing authorization header",
				})
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header format",
				})
			}

			claims, err := jwtSvc.Parse(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			c.Set(ContextKeyUser, claims)
			return next(c)
		}
	}
}
