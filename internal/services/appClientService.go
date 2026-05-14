package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"new-auth-service/internal/db"
	"new-auth-service/internal/payload"

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
func CreateAppClient(ctx context.Context, tenantID string, req *payload.CreateAppClientRequest) (*db.AppClient, string, error) {
	if !db.AppType(req.AppType).IsValid() {
		return nil, "", fmt.Errorf("invalid app type: %s", req.AppType)
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

	client, err := db.InsertAppClient(ctx, tenantID, req.UserPoolID, req.ClientName, clientID, string(secretHash), req.AppType)
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

func GetAppClientsByTenantID(ctx context.Context, tenantID string) ([]db.AppClient, error) {
	clients, err := db.GetAppClientsByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch app clients: %w", err)
	}
	return clients, nil
}
