package db

import (
	"context"
)

type AppClient struct {
	ID         string
	UserPoolID string
	TenantID   string
	ClientName string
	ClientID   string
	AppType    string
	IsActive   bool
	CreatedAt  string
}

func InsertAppClient(ctx context.Context, userPoolID, tenantID, clientName, clientID, secretHash, appType string) (*AppClient, error) {
	query := `
		INSERT INTO app_clients (user_pool_id, tenant_id, client_name, client_id, client_secret_hash, app_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_pool_id, tenant_id, client_name, client_id, app_type, is_active, created_at
	`

	row := DB.QueryRowContext(ctx, query,
		userPoolID,
		tenantID,
		clientName,
		clientID,
		secretHash,
		appType,
	)

	a := &AppClient{}
	err := row.Scan(
		&a.ID, &a.UserPoolID, &a.TenantID, &a.ClientName,
		&a.ClientID, &a.AppType, &a.IsActive, &a.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// GetAppClientByClientID fetches an app client and its secret hash by client_id.
// The secret hash is returned separately and should not be stored on AppClient.
func GetAppClientByClientID(ctx context.Context, clientID string) (*AppClient, string, error) {
	query := `
		SELECT id, user_pool_id, tenant_id, client_name, client_id, client_secret_hash, app_type, is_active, created_at
		FROM app_clients
		WHERE client_id = $1
	`

	row := DB.QueryRowContext(ctx, query, clientID)

	a := &AppClient{}
	var secretHash string
	err := row.Scan(
		&a.ID, &a.UserPoolID, &a.TenantID, &a.ClientName,
		&a.ClientID, &secretHash, &a.AppType, &a.IsActive, &a.CreatedAt,
	)
	if err != nil {
		return nil, "", err
	}
	return a, secretHash, nil
}

func GetAppClientsByPoolID(ctx context.Context, userPoolID string) ([]AppClient, error) {
	query := `
		SELECT id, user_pool_id, tenant_id, client_name, client_id, app_type, is_active, created_at
		FROM app_clients
		WHERE user_pool_id = $1
		ORDER BY created_at DESC
	`

	rows, err := DB.QueryContext(ctx, query, userPoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []AppClient
	for rows.Next() {
		a := AppClient{}
		err := rows.Scan(
			&a.ID, &a.UserPoolID, &a.TenantID, &a.ClientName,
			&a.ClientID, &a.AppType, &a.IsActive, &a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		clients = append(clients, a)
	}
	return clients, rows.Err()
}
