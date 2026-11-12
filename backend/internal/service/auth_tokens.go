package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"job-portal/internal/db/dbgen"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenUsed     = errors.New("token already used")
	ErrTokenExpired  = errors.New("token expired")
)

type AuthTokenService struct {
	Queries *dbgen.Queries
}

func NewAuthTokenService(q *dbgen.Queries) *AuthTokenService {
	return &AuthTokenService{Queries: q}
}

// generateRawToken returns a URL-safe random token (32 bytes → 43 chars base64).
func generateRawToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// Create generates a new token, stores its hash, and returns the raw token.
// Any existing unused tokens of the same purpose for the user are NOT deleted
// — they expire naturally or get marked used. This prevents race conditions.
func (s *AuthTokenService) Create(ctx context.Context, userID uuid.UUID, purpose string, ttl time.Duration) (string, error) {
	raw, err := generateRawToken()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	_, err = s.Queries.CreateAuthToken(ctx, dbgen.CreateAuthTokenParams{
		UserID:    userID,
		TokenHash: hashToken(raw),
		Purpose:   purpose,
		ExpiresAt: pgTimestamp(time.Now().Add(ttl)),
	})
	if err != nil {
		return "", fmt.Errorf("store token: %w", err)
	}

	return raw, nil
}

// Consume validates a raw token and marks it used. Returns the user ID.
func (s *AuthTokenService) Consume(ctx context.Context, raw, purpose string) (uuid.UUID, error) {
	tok, err := s.Queries.GetAuthTokenByHash(ctx, dbgen.GetAuthTokenByHashParams{
		TokenHash: hashToken(raw),
		Purpose:   purpose,
	})
	if err != nil {
		return uuid.Nil, ErrTokenNotFound
	}

	if tok.UsedAt.Valid {
		return uuid.Nil, ErrTokenUsed
	}

	if time.Now().After(tok.ExpiresAt.Time) {
		return uuid.Nil, ErrTokenExpired
	}

	if err := s.Queries.MarkAuthTokenUsed(ctx, tok.ID); err != nil {
		return uuid.Nil, fmt.Errorf("mark used: %w", err)
	}

	return tok.UserID, nil
}
