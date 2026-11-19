package service

import (
	"crypto/rand"
	"encoding/base64"
	"crypto/subtle"
)

// GenerateCSRFToken returns a fresh random token for the double-submit pattern.
func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// CSRFValid compares the cookie value to the header value in constant time.
func CSRFValid(cookieValue, headerValue string) bool {
	if cookieValue == "" || headerValue == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(cookieValue), []byte(headerValue)) == 1
}
