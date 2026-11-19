package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"job-portal/internal/service"
)

// CSRFProtect enforces the double-submit cookie pattern on state-changing
// methods. Safe methods (GET, HEAD, OPTIONS) pass through.
//
// The browser sends the CSRF cookie automatically. A malicious site cannot
// read the cookie to set the header. So if cookie == header, the request
// came from our own JS.
func CSRFProtect() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method
			if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
				return next(c)
			}

			cookie, err := c.Cookie(service.CSRFCookieName)
			if err != nil || cookie.Value == "" {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "missing csrf cookie",
				})
			}

			header := c.Request().Header.Get(service.CSRFHeaderName)
			if !service.CSRFValid(cookie.Value, header) {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "invalid csrf token",
				})
			}

			return next(c)
		}
	}
}
