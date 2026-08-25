package db

import (
	"context"
	"database/sql"
)

type TenantStatus string

const (
	TenantStatusActive  TenantStatus = "ACTIVE"
	TenantStatusHold    TenantStatus = "HOLD"
	TenantStatusBlocked TenantStatus = "BLOCKED"
)

type Tenant struct {
	ID         string
	AccountID  string
	Name       string
	Email      string
	Phone      string
	Address    string
	WebsiteURL string
	Status     TenantStatus
	CreatedAt  string
}

func InsertTenant(ctx context.Context, accountID, name, email, passwordHash, phone, address, websiteURL string) (*Tenant, error) {
	query := `
		INSERT INTO tenants (account_id, name, email, password_hash, phone, address, website_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, account_id, name, email, phone, address, website_url, status, created_at
	`

	row := DB.QueryRowContext(ctx, query,
		accountID,
		name,
		email,
		passwordHash,
		nullableString(phone),
		nullableString(address),
		nullableString(websiteURL),
	)

	t := &Tenant{}
	var phone2, address2, websiteURL2 sql.NullString
	var status string
	err := row.Scan(&t.ID, &t.AccountID, &t.Name, &t.Email, &phone2, &address2, &websiteURL2, &status, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	t.Phone = phone2.String
	t.Address = address2.String
	t.WebsiteURL = websiteURL2.String
	t.Status = TenantStatus(status)

	return t, nil
}

func GetTenantByAccountID(ctx context.Context, accountID string) (*Tenant, string, error) {
	query := `
		SELECT id, account_id, name, email, phone, address, website_url, status, password_hash, created_at
		FROM tenants
		WHERE account_id = $1
	`
	row := DB.QueryRowContext(ctx, query, accountID)

	t := &Tenant{}
	var phone, address, websiteURL sql.NullString
	var passwordHash, status string

	err := row.Scan(
		&t.ID, &t.AccountID, &t.Name, &t.Email,
		&phone, &address, &websiteURL,
		&status, &passwordHash, &t.CreatedAt,
	)
	if err != nil {
		return nil, "", err
	}

	t.Phone = phone.String
	t.Address = address.String
	t.WebsiteURL = websiteURL.String
	t.Status = TenantStatus(status)

	return t, passwordHash, nil
}

func GetTenantByEmail(ctx context.Context, email string) (*Tenant, string, error) {
	query := `
        SELECT id, account_id, name, email, phone, address, website_url, status, password_hash, created_at
        FROM tenants
        WHERE email = $1
    `
	row := DB.QueryRowContext(ctx, query, email)

	t := &Tenant{}
	var phone, address, websiteURL sql.NullString
	var passwordHash string

	err := row.Scan(
		&t.ID, &t.AccountID, &t.Name, &t.Email,
		&phone, &address, &websiteURL,
		&t.Status, &passwordHash, &t.CreatedAt,
	)
	if err != nil {
		return nil, "", err
	}

	t.Phone = phone.String
	t.Address = address.String
	t.WebsiteURL = websiteURL.String

	return t, passwordHash, nil
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
