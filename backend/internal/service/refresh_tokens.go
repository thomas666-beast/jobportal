package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"

	"job-portal/internal/db/dbgen"
)

var ErrRefreshTokenInvalid = errors.New("invalid refresh token")

type RefreshTokenService struct {
	Queries *dbgen.Queries
	ttl     time.Duration
}

func NewRefreshTokenService(q *dbgen.Queries) *RefreshTokenService {
	days := 30
	if v := os.Getenv("JWT_REFRESH_EXPIRY_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			days = n
		}
	}
	return &RefreshTokenService{
		Queries: q,
		ttl:     time.Duration(days) * 24 * time.Hour,
	}
}

func (s *RefreshTokenService) TTL() time.Duration {
	return s.ttl
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// Create generates a new refresh token for the user, stores its hash, returns the raw value.
func (s *RefreshTokenService) Create(ctx context.Context, userID uuid.UUID) (string, error) {
	raw, err := generateRefreshToken()
	if err != nil {
		return "", fmt.Errorf("generate: %w", err)
	}

	_, err = s.Queries.CreateRefreshToken(ctx, dbgen.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: hashRefreshToken(raw),
		ExpiresAt: pgTimestamp(time.Now().Add(s.ttl)),
	})
	if err != nil {
		return "", fmt.Errorf("store: %w", err)
	}
	return raw, nil
}

// Validate looks up a raw refresh token and returns the user ID if valid.
// Does NOT consume the token (refresh tokens are long-lived; rotation is on a
// separate mechanism if you want it, but for now single-use-per-login).
func (s *RefreshTokenService) Validate(ctx context.Context, raw string) (uuid.UUID, error) {
	tok, err := s.Queries.GetRefreshTokenByHash(ctx, hashRefreshToken(raw))
	if err != nil {
		return uuid.Nil, ErrRefreshTokenInvalid
	}

	if time.Now().After(tok.ExpiresAt.Time) {
		_ = s.Queries.DeleteRefreshTokenByHash(ctx, tok.TokenHash)
		return uuid.Nil, ErrRefreshTokenInvalid
	}

	return tok.UserID, nil
}

func (s *RefreshTokenService) Revoke(ctx context.Context, raw string) {
	_ = s.Queries.DeleteRefreshTokenByHash(ctx, hashRefreshToken(raw))
}

func (s *RefreshTokenService) RevokeAllForUser(ctx context.Context, userID uuid.UUID) {
	_ = s.Queries.DeleteAllRefreshTokensForUser(ctx, userID)
}
