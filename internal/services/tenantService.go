package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"new-auth-service/internal"
	"new-auth-service/internal/db/sqlc"
	"new-auth-service/internal/enum"
	"new-auth-service/internal/payload"
	"new-auth-service/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func CreateTenant(ctx context.Context, req *payload.CreateTenantRequest) (map[string]any, error) {
	accountID, err := utils.GenerateAccountID()
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	addressJSON, err := json.Marshal(req.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal address: %w", err)
	}

	tenant, err := internal.DB.InsertTenant(ctx, sqlc.InsertTenantParams{
		AccountID:  accountID,
		Name:       req.Name,
		Email:      req.Email,
		Password:   string(hash),
		Phone:      utils.StringToPGText(req.Phone),
		Address:    addressJSON,
		WebsiteUrl: utils.StringToPGText(req.WebsiteURL),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert tenant: %w", err)
	}

	return map[string]any{
		"tenantId":  tenant.ID,
		"accountId": tenant.AccountID,
		"name":      tenant.Name,
		"email":     tenant.Email,
		"status":    tenant.Status,
		"createdAt": tenant.CreatedAt,
	}, nil
}

func LoginTenant(ctx context.Context, req *payload.TenantLoginRequest) (*payload.TenantAuthTokens, error) {
	tenant, err := internal.DB.GetTenantsByAccountID(ctx, req.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to fetch tenant: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(tenant.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if tenant.Status != enum.TenantStatusActive.String() {
		return nil, fmt.Errorf("tenant account is not active")
	}

	return issueTenantTokenPair(ctx, &payload.Tenant{
		ID:        tenant.ID.String(),
		AccountID: tenant.AccountID,
		Name:      tenant.Name,
		Email:     tenant.Email,
		Status:    tenant.Status,
	})
}

func issueTenantTokenPair(ctx context.Context, tenant *payload.Tenant) (*payload.TenantAuthTokens, error) {
	accessToken, err := utils.GenerateTenantToken(tenant.ID, tenant.AccountID, tenant.Name, tenant.Email, string(tenant.Status))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	rawRefresh, expiresAt, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	tokenHash := sha256Hex(rawRefresh)
	if err := internal.DB.InsertTenantRefreshToken(ctx, sqlc.InsertTenantRefreshTokenParams{
		TenantID:  utils.StringToPGUUID(tenant.ID),
		TokenHash: tokenHash,
		ExpiresAt: utils.TimeToPGTimestamptz(expiresAt),
	}); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &payload.TenantAuthTokens{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}
