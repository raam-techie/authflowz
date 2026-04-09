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

func CreateUser(ctx context.Context, req *payload.CreateUserRequest, project *db.Project) (*db.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := db.InsertUser(ctx, project.TenantID, project.ID, req.Name, req.Email, string(hash), req.Phone, req.Department, req.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}

func LoginUser(ctx context.Context, req *payload.UserLoginRequest, project *db.Project) (string, error) {
	if !containsAuthType(project.AuthTypes, "EMAIL_PASSWORD") {
		return "", fmt.Errorf("email/password auth is not enabled for this project")
	}

	user, passwordHash, err := db.GetUserByEmail(ctx, project.ID, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid credentials")
		}
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}

	if !user.IsActive {
		return "", fmt.Errorf("user account is not active")
	}

	return loginEmailPassword(user, passwordHash, req)
}

func SendEmailOTP(ctx context.Context, req *payload.SendOTPRequest, project *db.Project) error {
	// TODO: implement send OTP to email
	return nil
}

func VerifyEmailOTP(ctx context.Context, req *payload.VerifyOTPRequest, project *db.Project) (string, error) {
	// TODO: implement verify OTP and return JWT token
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
