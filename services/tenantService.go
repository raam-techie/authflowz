package services

import (
	"context"
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
