package db

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Client struct {
	ID           string
	TenantID     string
	ClientID     string
	Name         string
	Description  string
	AppType      string
	AuthType     string
	RedirectURIs []string
	IsActive     bool
	CreatedAt    string
}

func InsertClient(ctx context.Context, tenantID, clientID, clientSecretHash, name, description, appType, authType string, redirectURIs []string) (*Client, error) {
	redirectURIsJSON, err := json.Marshal(redirectURIs)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO projects (tenant_id, client_id, client_secret_hash, name, description, app_type, auth_type, redirect_uris)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, tenant_id, client_id, name, description, app_type, auth_type, redirect_uris, is_active, created_at
	`

	row := DB.QueryRowContext(ctx, query,
		tenantID,
		clientID,
		clientSecretHash,
		name,
		nullableString(description),
		appType,
		authType,
		redirectURIsJSON,
	)

	c := &Client{}
	var desc sql.NullString
	var redirectURIsRaw []byte
	err = row.Scan(&c.ID, &c.TenantID, &c.ClientID, &c.Name, &desc, &c.AppType, &c.AuthType, &redirectURIsRaw, &c.IsActive, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	c.Description = desc.String
	_ = json.Unmarshal(redirectURIsRaw, &c.RedirectURIs)

	return c, nil
}
