package db

import (
	"context"
	"database/sql"
	"encoding/json"
)

type UserPool struct {
	ID            string
	TenantID      string
	PoolName      string
	PoolID        string
	Description   string
	SignInMethods []string
	AppType       string
	IsActive      bool
	CreatedAt     string
}

func InsertUserPool(ctx context.Context, tenantID, poolName, poolID, description string, signInMethods []string, appType string) (*UserPool, error) {
	signInJSON, err := json.Marshal(signInMethods)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO user_pools (tenant_id, pool_name, pool_id, description, sign_in_methods, app_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, tenant_id, pool_name, pool_id, description, sign_in_methods, app_type, is_active, created_at
	`

	row := DB.QueryRowContext(ctx, query,
		tenantID,
		poolName,
		poolID,
		nullableString(description),
		signInJSON,
		appType,
	)

	return scanUserPool(row)
}

func GetUserPoolByID(ctx context.Context, id string) (*UserPool, error) {
	query := `
		SELECT id, tenant_id, pool_name, pool_id, description, sign_in_methods, app_type, is_active, created_at
		FROM user_pools
		WHERE id = $1
	`
	row := DB.QueryRowContext(ctx, query, id)
	return scanUserPool(row)
}

func GetUserPoolByPoolID(ctx context.Context, poolID string) (*UserPool, error) {
	query := `
		SELECT id, tenant_id, pool_name, pool_id, description, sign_in_methods, app_type, is_active, created_at
		FROM user_pools
		WHERE pool_id = $1
	`
	row := DB.QueryRowContext(ctx, query, poolID)
	return scanUserPool(row)
}

func GetUserPoolsByTenantID(ctx context.Context, tenantID string) ([]UserPool, error) {
	query := `
		SELECT id, tenant_id, pool_name, pool_id, description, sign_in_methods, app_type, is_active, created_at
		FROM user_pools
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`
	rows, err := DB.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []UserPool
	for rows.Next() {
		p, err := scanUserPoolRow(rows)
		if err != nil {
			return nil, err
		}
		pools = append(pools, *p)
	}
	return pools, rows.Err()
}

func scanUserPool(row *sql.Row) (*UserPool, error) {
	p := &UserPool{}
	var desc sql.NullString
	var signInRaw []byte

	err := row.Scan(
		&p.ID, &p.TenantID, &p.PoolName, &p.PoolID,
		&desc, &signInRaw, &p.AppType, &p.IsActive, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.Description = desc.String
	_ = json.Unmarshal(signInRaw, &p.SignInMethods)
	return p, nil
}

func scanUserPoolRow(rows *sql.Rows) (*UserPool, error) {
	p := &UserPool{}
	var desc sql.NullString
	var signInRaw []byte

	err := rows.Scan(
		&p.ID, &p.TenantID, &p.PoolName, &p.PoolID,
		&desc, &signInRaw, &p.AppType, &p.IsActive, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.Description = desc.String
	_ = json.Unmarshal(signInRaw, &p.SignInMethods)
	return p, nil
}
