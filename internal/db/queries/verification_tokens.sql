-- name: InsertTenantRefreshToken :exec
INSERT INTO verification_tokens (tenant_id, token_hash, expires_at)
VALUES ($1, $2, $3);

-- name: InsertUserRefreshToken :exec
INSERT INTO verification_tokens (user_id, app_client_id, token_hash, expires_at)
VALUES ($1, $2, $3, $4);

-- name: GetRefreshTokenByHash :one
SELECT id, COALESCE(user_id::TEXT, '') as user_id, COALESCE(app_client_id::TEXT, '') as app_client_id, COALESCE(tenant_id::TEXT, '') as tenant_id, token_hash, expires_at, created_at
FROM verification_tokens
WHERE token_hash = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM verification_tokens WHERE token_hash = $1;