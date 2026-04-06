package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"new-auth-service/db"
	"new-auth-service/payload"
	"new-auth-service/utils"

	"golang.org/x/crypto/bcrypt"
)

var ErrNotImplemented = errors.New("not implemented")

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

func LoginUser(ctx context.Context, req *payload.UserLoginRequest) (string, error) {
	user, passwordHash, err := db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid credentials")
		}
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}

	if !user.IsActive {
		return "", fmt.Errorf("user account is not active")
	}

	project, err := db.GetProjectByID(ctx, user.ProjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("project not found")
		}
		return "", fmt.Errorf("failed to fetch project: %w", err)
	}

	if !project.IsActive {
		return "", fmt.Errorf("project is not active")
	}

	switch project.AuthType {
	case "EMAIL_PASSWORD":
		return loginEmailPassword(user, passwordHash, req)
	case "OTP_EMAIL":
		return "", fmt.Errorf("OTP_EMAIL auth is not yet implemented: %w", ErrNotImplemented)
	case "OTP_PHONE":
		return "", fmt.Errorf("OTP_PHONE auth is not yet implemented: %w", ErrNotImplemented)
	case "OAUTH":
		return "", fmt.Errorf("OAUTH auth is not yet implemented: %w", ErrNotImplemented)
	case "MAGIC_LINK":
		return "", fmt.Errorf("MAGIC_LINK auth is not yet implemented: %w", ErrNotImplemented)
	case "SSO":
		return "", fmt.Errorf("SSO auth is not yet implemented: %w", ErrNotImplemented)
	default:
		return "", fmt.Errorf("unknown auth type: %s", project.AuthType)
	}
}

func loginEmailPassword(user *db.User, passwordHash string, req *payload.UserLoginRequest) (string, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateUserToken(user.ID, user.TenantID, user.ProjectID, user.UID, user.Name, user.Email, user.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
