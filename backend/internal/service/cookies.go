package service

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	AccessCookieName  = "access_token"
	RefreshCookieName = "refresh_token"
	CSRFCookieName    = "csrf_token"
	CSRFHeaderName    = "X-CSRF-Token"
	RefreshPath       = "/api/auth/refresh"
)

func cookieSecure() bool {
	return os.Getenv("COOKIE_SECURE") == "true"
}

func cookieDomain() string {
	return os.Getenv("COOKIE_DOMAIN")
}

// SetAccessCookie writes the JWT access token cookie.
func SetAccessCookie(c echo.Context, token string, ttl time.Duration) {
	c.SetCookie(&http.Cookie{
		Name:     AccessCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(ttl.Seconds()),
		Secure:   cookieSecure(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Domain:   cookieDomain(),
	})
}

// SetRefreshCookie writes the opaque refresh token cookie (scoped to refresh path).
func SetRefreshCookie(c echo.Context, token string, ttl time.Duration) {
	c.SetCookie(&http.Cookie{
		Name:     RefreshCookieName,
		Value:    token,
		Path:     RefreshPath,
		MaxAge:   int(ttl.Seconds()),
		Secure:   cookieSecure(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   cookieDomain(),
	})
}

// SetCSRFCookie writes the CSRF token cookie — NOT HttpOnly, because the
// browser needs to read it and echo it back in a header.
func SetCSRFCookie(c echo.Context, token string, ttl time.Duration) {
	c.SetCookie(&http.Cookie{
		Name:     CSRFCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(ttl.Seconds()),
		Secure:   cookieSecure(),
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Domain:   cookieDomain(),
	})
}

// ClearAuthCookies removes all three cookies (used on logout).
func ClearAuthCookies(c echo.Context) {
	for _, name := range []string{AccessCookieName, RefreshCookieName, CSRFCookieName} {
		c.SetCookie(&http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   cookieSecure(),
			HttpOnly: name != CSRFCookieName,
			SameSite: http.SameSiteLaxMode,
			Domain:   cookieDomain(),
		})
	}
	// Also clear the refresh cookie at its scoped path.
	c.SetCookie(&http.Cookie{
		Name:     RefreshCookieName,
		Value:    "",
		Path:     RefreshPath,
		MaxAge:   -1,
		Secure:   cookieSecure(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   cookieDomain(),
	})
}
