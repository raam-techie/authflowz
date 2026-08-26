package db

import (
	"context"
	"database/sql"
)

type User struct {
	ID          string
	TenantID    string
	AppClientID string
	UID         string
	Name        string
	Email       string
	Phone       string
	Role        string
	IsActive    bool
	CreatedAt   string
}

func GetUserByEmail(ctx context.Context, appClientID, email string) (*User, string, error) {
	query := `
		SELECT id, tenant_id, app_client_id, uid, name, email, phone, role, is_active, password_hash, created_at
		FROM users
		WHERE app_client_id = $1 AND email = $2
	`
	row := DB.QueryRowContext(ctx, query, appClientID, email)

	u := &User{}
	var phone sql.NullString
	var passwordHash string

	err := row.Scan(
		&u.ID, &u.TenantID, &u.AppClientID, &u.UID, &u.Name, &u.Email,
		&phone, &u.Role, &u.IsActive, &passwordHash, &u.CreatedAt,
	)
	if err != nil {
		return nil, "", err
	}

	u.Phone = phone.String

	return u, passwordHash, nil
}

func GetUserByID(ctx context.Context, userID string) (*User, error) {
	query := `
		SELECT id, tenant_id, app_client_id, uid, name, email, phone, role, is_active, created_at
		FROM users
		WHERE id = $1
	`
	row := DB.QueryRowContext(ctx, query, userID)

	u := &User{}
	var phone sql.NullString
	err := row.Scan(
		&u.ID, &u.TenantID, &u.AppClientID, &u.UID, &u.Name, &u.Email,
		&phone, &u.Role, &u.IsActive, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	u.Phone = phone.String
	return u, nil
}

func GetUserByUID(ctx context.Context, appClientID, uid string) (*User, error) {
	query := `
		SELECT id, tenant_id, app_client_id, uid, name, email, phone, role, is_active, created_at
		FROM users
		WHERE app_client_id = $1 AND uid = $2
	`
	row := DB.QueryRowContext(ctx, query, appClientID, uid)

	u := &User{}
	var phone sql.NullString
	err := row.Scan(
		&u.ID, &u.TenantID, &u.AppClientID, &u.UID, &u.Name, &u.Email,
		&phone, &u.Role, &u.IsActive, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	u.Phone = phone.String
	return u, nil
}

func GetUsersByTenantID(ctx context.Context, tenantID string) ([]User, error) {
	query := `
		SELECT id, tenant_id, app_client_id, uid, name, email, phone, role, is_active, created_at
		FROM users
		WHERE tenant_id = $1
	`
	rows, err := DB.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		u := User{}
		var phone sql.NullString

		err := rows.Scan(
			&u.ID, &u.TenantID, &u.AppClientID, &u.UID, &u.Name,
			&u.Email, &phone, &u.Role, &u.IsActive, &u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		u.Phone = phone.String
		users = append(users, u)
	}
	return users, rows.Err()
}

func UpsertUserAttributes(ctx context.Context, userID string, attrs map[string]string) error {
	query := `
		INSERT INTO user_attributes (user_id, key, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, key) DO UPDATE SET value = EXCLUDED.value
	`
	for k, v := range attrs {
		if _, err := DB.ExecContext(ctx, query, userID, k, v); err != nil {
			return err
		}
	}
	return nil
}

func GetUserAttributes(ctx context.Context, userID string) (map[string]string, error) {
	query := `SELECT key, value FROM user_attributes WHERE user_id = $1`
	rows, err := DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attrs := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		attrs[k] = v
	}
	return attrs, rows.Err()
}

func InsertUser(ctx context.Context, tenantID, appClientID, name, email, passwordHash, phone, role string) (*User, error) {
	if role == "" {
		role = "USER"
	}

	query := `
		INSERT INTO users (tenant_id, app_client_id, uid, name, email, password_hash, phone, role)
		VALUES ($1, $2, LPAD(nextval('user_uid_seq')::TEXT, 6, '0'), $3, $4, $5, $6, $7)
		RETURNING id, tenant_id, app_client_id, uid, name, email, phone, role, is_active, created_at
	`

	row := DB.QueryRowContext(ctx, query,
		tenantID,
		appClientID,
		name,
		email,
		passwordHash,
		nullableString(phone),
		role,
	)

	u := &User{}
	var phone2 sql.NullString
	err := row.Scan(&u.ID, &u.TenantID, &u.AppClientID, &u.UID, &u.Name, &u.Email, &phone2, &u.Role, &u.IsActive, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	u.Phone = phone2.String

	return u, nil
}
