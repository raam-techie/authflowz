package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tenantTokenTTL = 20 * time.Minute
const userTokenTTL = 20 * time.Minute

type TenantClaims struct {
	TenantID  string `json:"tenantId"`
	AccountID string `json:"accountId"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	jwt.RegisteredClaims
}

func jwtSecret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}
	return []byte(s), nil
}

func GenerateTenantToken(tenantID, accountID, name, email, status string) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}

	claims := TenantClaims{
		TenantID:  tenantID,
		AccountID: accountID,
		Name:      name,
		Email:     email,
		Status:    status,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tenantTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

type UserClaims struct {
	UserID    string `json:"userId"`
	TenantID  string `json:"tenantId"`
	ProjectID string `json:"projectId"`
	UID       string `json:"uid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateUserToken(userID, tenantID, projectID, uid, name, email, role string) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}

	claims := UserClaims{
		UserID:    userID,
		TenantID:  tenantID,
		ProjectID: projectID,
		UID:       uid,
		Name:      name,
		Email:     email,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(userTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
