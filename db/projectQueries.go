package db

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Project struct {
	ID           string
	TenantID     string
	ClientID     string
	Name         string
	Description  string
	AppType      string
	AuthTypes    []string
	RedirectURIs []string
	IsActive     bool
	CreatedAt    string
}

func GetProjectByID(ctx context.Context, projectID string) (*Project, error) {
	query := `
		SELECT id, tenant_id, client_id, name, description, app_type, auth_types, redirect_uris, is_active, created_at
		FROM projects
		WHERE id = $1
	`
	row := DB.QueryRowContext(ctx, query, projectID)

	p := &Project{}
	var desc sql.NullString
	var authTypesRaw, redirectURIsRaw []byte

	err := row.Scan(
		&p.ID, &p.TenantID, &p.ClientID, &p.Name, &desc,
		&p.AppType, &authTypesRaw, &redirectURIsRaw, &p.IsActive, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.Description = desc.String
	_ = json.Unmarshal(authTypesRaw, &p.AuthTypes)
	_ = json.Unmarshal(redirectURIsRaw, &p.RedirectURIs)

	return p, nil
}

func GetProjectByClientID(ctx context.Context, clientID string) (*Project, string, error) {
	query := `
		SELECT id, tenant_id, client_id, client_secret_hash, name, description, app_type, auth_types, redirect_uris, is_active, created_at
		FROM projects
		WHERE client_id = $1
	`
	row := DB.QueryRowContext(ctx, query, clientID)

	p := &Project{}
	var desc sql.NullString
	var secretHash string
	var authTypesRaw, redirectURIsRaw []byte

	err := row.Scan(
		&p.ID, &p.TenantID, &p.ClientID, &secretHash, &p.Name, &desc,
		&p.AppType, &authTypesRaw, &redirectURIsRaw, &p.IsActive, &p.CreatedAt,
	)
	if err != nil {
		return nil, "", err
	}

	p.Description = desc.String
	_ = json.Unmarshal(authTypesRaw, &p.AuthTypes)
	_ = json.Unmarshal(redirectURIsRaw, &p.RedirectURIs)

	return p, secretHash, nil
}

func InsertProject(ctx context.Context, tenantID, clientID, clientSecretHash, name, description, appType string, authTypes []string, redirectURIs []string) (*Project, error) {
	authTypesJSON, err := json.Marshal(authTypes)
	if err != nil {
		return nil, err
	}

	redirectURIsJSON, err := json.Marshal(redirectURIs)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO projects (tenant_id, client_id, client_secret_hash, name, description, app_type, auth_types, redirect_uris)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, tenant_id, client_id, name, description, app_type, auth_types, redirect_uris, is_active, created_at
	`

	row := DB.QueryRowContext(ctx, query,
		tenantID,
		clientID,
		clientSecretHash,
		name,
		nullableString(description),
		appType,
		authTypesJSON,
		redirectURIsJSON,
	)

	p := &Project{}
	var desc sql.NullString
	var authTypesRaw, redirectURIsRaw []byte
	err = row.Scan(&p.ID, &p.TenantID, &p.ClientID, &p.Name, &desc, &p.AppType, &authTypesRaw, &redirectURIsRaw, &p.IsActive, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	p.Description = desc.String
	_ = json.Unmarshal(authTypesRaw, &p.AuthTypes)
	_ = json.Unmarshal(redirectURIsRaw, &p.RedirectURIs)

	return p, nil
}
