package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"

	"new-auth-service/db"
	"new-auth-service/payload"

	"golang.org/x/crypto/bcrypt"
)

func generateClientIDHex() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate client id: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func generateClientSecretHex() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate client secret: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// CreateAppClient validates the request, generates credentials, and inserts an app client.
// Returns the app client and the raw (unhashed) client secret — shown only once.
func CreateAppClient(ctx context.Context, req *payload.CreateAppClientRequest) (*db.AppClient, string, error) {
	if !db.IsValidAppType(req.AppType) {
		return nil, "", fmt.Errorf("invalid app type: %s", req.AppType)
	}

	pool, err := db.GetUserPoolByID(ctx, req.UserPoolID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", fmt.Errorf("user pool not found")
		}
		return nil, "", fmt.Errorf("failed to fetch user pool: %w", err)
	}

	clientID, err := generateClientIDHex()
	if err != nil {
		return nil, "", err
	}

	clientSecret, err := generateClientSecretHex()
	if err != nil {
		return nil, "", err
	}

	secretHash, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash client secret: %w", err)
	}

	client, err := db.InsertAppClient(ctx, pool.ID, pool.TenantID, req.ClientName, clientID, string(secretHash), req.AppType)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create app client: %w", err)
	}

	return client, clientSecret, nil
}

func GetAppClientsByPoolID(ctx context.Context, userPoolID string) ([]db.AppClient, error) {
	clients, err := db.GetAppClientsByPoolID(ctx, userPoolID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch app clients: %w", err)
	}
	return clients, nil
}
