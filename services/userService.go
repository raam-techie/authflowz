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

	if !containsAuthType(project.AuthTypes, "EMAIL_PASSWORD") {
		return "", fmt.Errorf("email/password auth is not enabled for this project")
	}

	return loginEmailPassword(user, passwordHash, req)
}

func SendEmailOTP(ctx context.Context, req *payload.SendOTPRequest) error {
	// TODO: implement send OTP to email
	return nil
}

func VerifyEmailOTP(ctx context.Context, req *payload.VerifyOTPRequest) (string, error) {
	// TODO: implement verify OTP and return JWT token
	return "", nil
}

func GoogleOAuthLogin(ctx context.Context, req *payload.GoogleOAuthLoginRequest) (string, error) {
	// TODO: implement Google OAuth login
	return "", nil
}

func containsAuthType(authTypes []string, target string) bool {
	for _, a := range authTypes {
		if a == target {
			return true
		}
	}
	return false
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
