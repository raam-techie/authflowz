package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"

	"new-auth-service/internal"
	"new-auth-service/internal/db"
	"new-auth-service/internal/db/sqlc"
	"new-auth-service/internal/payload"
	"new-auth-service/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func CreateUser(ctx context.Context, req *payload.CreateUserRequest, appClient *db.AppClient) (*db.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := db.InsertUser(ctx, appClient.TenantID, appClient.ID, req.Name, req.Email, string(hash), req.Phone, req.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}

func LoginUser(ctx context.Context, req *payload.UserLoginRequest, appClient *db.AppClient) (*payload.AuthTokens, error) {
	user, passwordHash, err := db.GetUserByEmail(ctx, appClient.ID, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("user account is not active")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return issueTokenPair(ctx, user, appClient.ID)
}

func UpdateUserAttributes(ctx context.Context, uid string, req *payload.UpdateAttributesRequest, appClient *db.AppClient) error {
	user, err := db.GetUserByUID(ctx, appClient.ID, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	return db.UpsertUserAttributes(ctx, user.ID, req.Attributes)
}

func SendEmailOTP(ctx context.Context, req *payload.SendOTPRequest, appClient *db.AppClient) error {
	// TODO: implement send OTP to email
	return nil
}

func VerifyEmailOTP(ctx context.Context, req *payload.VerifyOTPRequest, appClient *db.AppClient) (string, error) {
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

func containsSignInMethod(methods []string, target string) bool {
	for _, m := range methods {
		if m == target {
			return true
		}
	}
	return false
}

func issueTokenPair(ctx context.Context, user *db.User, appClientID string) (*payload.AuthTokens, error) {
	attrs, _ := db.GetUserAttributes(ctx, user.ID)

	accessToken, err := utils.GenerateUserToken(user.ID, user.TenantID, user.AppClientID, user.UID, user.Name, user.Email, user.Role, attrs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	rawRefresh, expiresAt, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	tokenHash := sha256Hex(rawRefresh)
	if err := internal.DB.InsertUserRefreshToken(ctx, sqlc.InsertUserRefreshTokenParams{
		UserID:      utils.StringToPGUUID(user.ID),
		AppClientID: utils.StringToPGUUID(appClientID),
		TokenHash:   tokenHash,
		ExpiresAt:   utils.TimeToPGTimestamptz(expiresAt),
	}); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &payload.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}
