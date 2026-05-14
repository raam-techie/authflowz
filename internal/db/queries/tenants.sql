-- name: InsertTenant :one
INSERT INTO tenants (account_id, name, email, password, phone, address, website_url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, account_id, name, email, phone, address, website_url, status, created_at;

-- name: GetTenantsByAccountID :one
SELECT id, account_id, name, email, phone, address, website_url, status, password, created_at
FROM tenants
WHERE account_id = $1;

-- name: GetTenantByID :one
SELECT id, account_id, name, email, phone, address, website_url, status, created_at
FROM tenants
WHERE id = $1;

-- name: InsertTenantRefreshToken :exec
INSERT INTO verification_tokens (tenant_id, token_hash, expires_at)
VALUES ($1, $2, $3);
