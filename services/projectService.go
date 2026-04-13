package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	"new-auth-service/db"
	"new-auth-service/payload"

	"golang.org/x/crypto/bcrypt"
)

func generateClientID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate client id: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func generateClientSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate client secret: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func CreateProject(ctx context.Context, req *payload.CreateProjectRequest) (*db.Project, string, error) {
	clientID, err := generateClientID()
	if err != nil {
		return nil, "", err
	}

	clientSecret, err := generateClientSecret()
	if err != nil {
		return nil, "", err
	}

	secretHash, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash client secret: %w", err)
	}

	redirectURIs := req.RedirectURIs
	if redirectURIs == nil {
		redirectURIs = []string{}
	}

	authTypes := req.AuthTypes
	if authTypes == nil {
		authTypes = []string{}
	}

	project, err := db.InsertProject(ctx, req.TenantID, clientID, string(secretHash), req.Name, req.Description, req.AppType, authTypes, redirectURIs)
	if err != nil {
		if strings.Contains(err.Error(), "projects_tenant_id_fkey") {
			return nil, "", fmt.Errorf("invalid tenant id")
		}
		return nil, "", fmt.Errorf("failed to insert project: %w", err)
	}

	return project, clientSecret, nil
}

func GetProjectByID(ctx context.Context, projectID string) (*db.Project, error) {
	project, err := db.GetProjectByID(ctx, projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to fetch project: %w", err)
	}
	return project, nil
}

func GetProjectsByTenantID(ctx context.Context, tenantID string) ([]db.Project, error) {
	projects, err := db.GetProjectsByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}
	return projects, nil
}
