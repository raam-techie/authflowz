package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"

	"new-auth-service/db"
	"new-auth-service/payload"
)

func generatePoolID() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate pool id: %w", err)
	}
	return "pool_" + hex.EncodeToString(b), nil
}

func CreateUserPool(ctx context.Context, req *payload.CreateUserPoolRequest) (*db.UserPool, error) {
	if !db.IsValidAppType(req.AppType) {
		return nil, fmt.Errorf("invalid app type: %s", req.AppType)
	}
	if len(req.SignInMethods) == 0 {
		return nil, fmt.Errorf("at least one sign-in method is required")
	}
	for _, m := range req.SignInMethods {
		if !db.IsValidSignInMethod(m) {
			return nil, fmt.Errorf("invalid sign-in method: %s", m)
		}
	}

	poolID, err := generatePoolID()
	if err != nil {
		return nil, err
	}

	pool, err := db.InsertUserPool(ctx, req.TenantID, req.PoolName, poolID, req.Description, req.SignInMethods, req.AppType)
	if err != nil {
		return nil, fmt.Errorf("failed to create user pool: %w", err)
	}

	return pool, nil
}

func GetUserPoolByID(ctx context.Context, id string) (*db.UserPool, error) {
	pool, err := db.GetUserPoolByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user pool not found")
		}
		return nil, fmt.Errorf("failed to fetch user pool: %w", err)
	}
	return pool, nil
}

func GetUserPoolsByTenantID(ctx context.Context, tenantID string) ([]db.UserPool, error) {
	pools, err := db.GetUserPoolsByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user pools: %w", err)
	}
	return pools, nil
}
