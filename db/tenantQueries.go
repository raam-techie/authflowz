package db

import (
	"context"
	"database/sql"
)

type Tenant struct {
	ID         string
	AccountID  string
	Name       string
	Email      string
	Phone      string
	Address    string
	WebsiteURL string
	Status     string
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
	err := row.Scan(&t.ID, &t.AccountID, &t.Name, &t.Email, &phone2, &address2, &websiteURL2, &t.Status, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	t.Phone = phone2.String
	t.Address = address2.String
	t.WebsiteURL = websiteURL2.String

	return t, nil
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
