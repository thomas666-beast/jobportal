package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"job-portal/internal/service"
)

const ContextKeyUser = "user"

// JWTAuth reads the access token from the httpOnly cookie (not the
// Authorization header anymore). Sets claims in the context.
func JWTAuth(jwtSvc *service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(service.AccessCookieName)
			if err != nil || cookie.Value == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "not authenticated",
				})
			}

			claims, err := jwtSvc.Parse(cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired session",
				})
			}

			c.Set(ContextKeyUser, claims)
			return next(c)
		}
	}
}
