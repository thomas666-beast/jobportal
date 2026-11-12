package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"job-portal/internal/service"
)

// RequireRole returns middleware that allows the request only if the
// JWT claims carry one of the allowed roles. Must run AFTER JWTAuth.
func RequireRole(allowed ...string) echo.MiddlewareFunc {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, r := range allowed {
		allowedSet[r] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(ContextKeyUser).(*service.Claims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "no user in context",
				})
			}

			if _, ok := allowedSet[claims.Role]; !ok {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "insufficient permissions",
				})
			}

			if !claims.OnboardingCompleted {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "onboarding not completed",
				})
			}

			return next(c)
		}
	}
}
