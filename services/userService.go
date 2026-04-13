package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"new-auth-service/db"
	"new-auth-service/payload"
	"new-auth-service/utils"

	"golang.org/x/crypto/bcrypt"
)

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func CreateUser(ctx context.Context, req *payload.CreateUserRequest, project *db.Project) (*db.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := db.InsertUser(ctx, project.TenantID, project.ID, req.Name, req.Email, string(hash), req.Phone, req.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}

func LoginUser(ctx context.Context, req *payload.UserLoginRequest, project *db.Project) (*payload.AuthTokens, error) {
	if !containsAuthType(project.AuthTypes, "EMAIL_PASSWORD") {
		return nil, fmt.Errorf("email/password auth is not enabled for this project")
	}

	user, passwordHash, err := db.GetUserByEmail(ctx, project.ID, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is not active")
	}

	return loginEmailPassword(ctx, user, passwordHash, req)
}

func RefreshAccessToken(ctx context.Context, req *payload.RefreshTokenRequest, project *db.Project) (*payload.AuthTokens, error) {
	tokenHash := sha256Hex(req.RefreshToken)

	rt, err := db.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid refresh token")
		}
		return nil, fmt.Errorf("failed to fetch refresh token: %w", err)
	}

	if rt.ProjectID != project.ID {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if time.Now().After(rt.ExpiresAt) {
		_ = db.DeleteRefreshToken(ctx, tokenHash)
		return nil, fmt.Errorf("refresh token expired")
	}

	if err := db.DeleteRefreshToken(ctx, tokenHash); err != nil {
		return nil, fmt.Errorf("failed to rotate refresh token: %w", err)
	}

	user, err := db.GetUserByID(ctx, rt.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is not active")
	}

	return issueTokenPair(ctx, user, project.ID)
}

func UpdateUserAttributes(ctx context.Context, uid string, req *payload.UpdateAttributesRequest, project *db.Project) error {
	user, err := db.GetUserByUID(ctx, project.ID, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	return db.UpsertUserAttributes(ctx, user.ID, req.Attributes)
}

func SendEmailOTP(ctx context.Context, req *payload.SendOTPRequest, project *db.Project) error {
	// TODO: implement send OTP to email
	return nil
}

func VerifyEmailOTP(ctx context.Context, req *payload.VerifyOTPRequest, project *db.Project) (string, error) {
	// TODO: implement verify OTP and return JWT token
	return "", nil
}

func GetUserByID(ctx context.Context, userID string) (*db.User, error) {
	user, err := db.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	return user, nil
}

func GetUsersByTenantID(ctx context.Context, tenantID string) ([]db.User, error) {
	users, err := db.GetUsersByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}

func containsAuthType(authTypes []string, target string) bool {
	for _, a := range authTypes {
		if a == target {
			return true
		}
	}
	return false
}

func loginEmailPassword(ctx context.Context, user *db.User, passwordHash string, req *payload.UserLoginRequest) (*payload.AuthTokens, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return issueTokenPair(ctx, user, user.ProjectID)
}

func issueTokenPair(ctx context.Context, user *db.User, projectID string) (*payload.AuthTokens, error) {
	attrs, _ := db.GetUserAttributes(ctx, user.ID)

	accessToken, err := utils.GenerateUserToken(user.ID, user.TenantID, user.ProjectID, user.UID, user.Name, user.Email, user.Role, attrs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	rawRefresh, expiresAt, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	tokenHash := sha256Hex(rawRefresh)
	if err := db.InsertUserRefreshToken(ctx, user.ID, projectID, tokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &payload.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}
