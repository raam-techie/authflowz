package services

import (
	"context"
	"fmt"

	"new-auth-service/db"
	"new-auth-service/payload"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx context.Context, req *payload.CreateUserRequest) (*db.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := db.InsertUser(ctx, req.TenantID, req.ProjectID, req.Name, req.Email, string(hash), req.Phone, req.Department, req.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}
