package db

import (
	"context"
	"database/sql"
)

type User struct {
	ID         string
	TenantID   string
	ProjectID  string
	UID        string
	Name       string
	Email      string
	Phone      string
	Department string
	Role       string
	IsActive   bool
	CreatedAt  string
}

func GetUserByEmail(ctx context.Context, projectID, email string) (*User, string, error) {
	query := `
		SELECT id, tenant_id, project_id, uid, name, email, phone, department, role, is_active, password_hash, created_at
		FROM users
		WHERE project_id = $1 AND email = $2
	`
	row := DB.QueryRowContext(ctx, query, projectID, email)

	u := &User{}
	var phone, dept sql.NullString
	var passwordHash string

	err := row.Scan(
		&u.ID, &u.TenantID, &u.ProjectID, &u.UID, &u.Name, &u.Email,
		&phone, &dept, &u.Role, &u.IsActive, &passwordHash, &u.CreatedAt,
	)
	if err != nil {
		return nil, "", err
	}

	u.Phone = phone.String
	u.Department = dept.String

	return u, passwordHash, nil
}

func InsertUser(ctx context.Context, tenantID, projectID, name, email, passwordHash, phone, department, role string) (*User, error) {
	if role == "" {
		role = "USER"
	}

	query := `
		INSERT INTO users (tenant_id, project_id, uid, name, email, password_hash, phone, department, role)
		VALUES ($1, $2, LPAD(nextval('user_uid_seq')::TEXT, 6, '0'), $3, $4, $5, $6, $7, $8)
		RETURNING id, tenant_id, project_id, uid, name, email, phone, department, role, is_active, created_at
	`

	row := DB.QueryRowContext(ctx, query,
		tenantID,
		projectID,
		name,
		email,
		passwordHash,
		nullableString(phone),
		nullableString(department),
		role,
	)

	u := &User{}
	var phone2, dept sql.NullString
	err := row.Scan(&u.ID, &u.TenantID, &u.ProjectID, &u.UID, &u.Name, &u.Email, &phone2, &dept, &u.Role, &u.IsActive, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	u.Phone = phone2.String
	u.Department = dept.String

	return u, nil
}
