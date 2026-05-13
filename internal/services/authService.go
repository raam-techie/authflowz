package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"new-auth-service/internal/db"
	"new-auth-service/internal/payload"
)

// RefreshToken is a unified refresh endpoint handler for both users and tenants.
// It detects the token type by checking whether tenant_id or user_id is set on the stored token.
func RefreshToken(ctx context.Context, req *payload.RefreshTokenRequest) (any, error) {
	tokenHash := sha256Hex(req.RefreshToken)

	rt, err := db.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid refresh token")
		}
		return nil, fmt.Errorf("failed to fetch refresh token: %w", err)
	}

	if time.Now().After(rt.ExpiresAt) {
		_ = db.DeleteRefreshToken(ctx, tokenHash)
		return nil, fmt.Errorf("refresh token expired")
	}

	if err := db.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return nil, fmt.Errorf("failed to rotate refresh token: %w", err)
	}

	if rt.TenantID != "" {
		return refreshTenantFromToken(ctx, rt.TenantID)
	}

	return refreshUserFromToken(ctx, rt.UserID, rt.AppClientID)
}

func refreshTenantFromToken(ctx context.Context, tenantID string) (*payload.TenantAuthTokens, error) {
	tenant, err := db.GetTenantByID(ctx, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, fmt.Errorf("failed to fetch tenant: %w", err)
	}

	if tenant.Status != db.TenantStatusActive {
		return nil, fmt.Errorf("tenant account is not active")
	}

	return issueTenantTokenPair(ctx, tenant)
}

func refreshUserFromToken(ctx context.Context, userID, appClientID string) (*payload.AuthTokens, error) {
	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is not active")
	}

	return issueTokenPair(ctx, user, appClientID)
}
