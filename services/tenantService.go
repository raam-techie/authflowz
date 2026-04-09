package services

import (
	"context"
	"database/sql"
	"fmt"

	"new-auth-service/db"
	"new-auth-service/payload"
	"new-auth-service/utils"

	"golang.org/x/crypto/bcrypt"
)

func CreateTenant(ctx context.Context, req *payload.CreateTenantRequest) (*db.Tenant, error) {

	accountID, err := utils.GenerateAccountID()
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	tenant, err := db.InsertTenant(ctx, accountID, req.Name, req.Email, string(hash), req.Phone, req.Address, req.WebsiteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to insert tenant: %w", err)
	}

	return tenant, nil
}

func LoginTenant(ctx context.Context, req *payload.TenantLoginRequest) (string, error) {
	tenant, passwordHash, err := db.GetTenantByAccountID(ctx, req.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid credentials")
		}
		return "", fmt.Errorf("failed to fetch tenant: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if tenant.Status != db.TenantStatusActive {
		return "", fmt.Errorf("tenant account is not active")
	}

	token, err := utils.GenerateTenantToken(tenant.ID, tenant.AccountID, tenant.Name, tenant.Email, string(tenant.Status))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
