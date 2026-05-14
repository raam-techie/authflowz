package db

import (
	"context"
	"time"
)

type RefreshToken struct {
	ID          string
	UserID      string
	AppClientID string
	TenantID    string
	TokenHash   string
	ExpiresAt   time.Time
	CreatedAt   string
}

func InsertUserRefreshToken(ctx context.Context, userID, appClientID, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, app_client_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := DB.ExecContext(ctx, query, userID, appClientID, tokenHash, expiresAt)
	return err
}

func GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*RefreshToken, error) {
	query := `
		SELECT id, COALESCE(user_id::TEXT, ''), COALESCE(app_client_id::TEXT, ''), COALESCE(tenant_id::TEXT, ''), token_hash, expires_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	row := DB.QueryRowContext(ctx, query, tokenHash)

	rt := &RefreshToken{}
	err := row.Scan(&rt.ID, &rt.UserID, &rt.AppClientID, &rt.TenantID, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return rt, nil
}

func DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := DB.ExecContext(ctx, query, tokenHash)
	return err
}
