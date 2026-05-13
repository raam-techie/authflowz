-- name: TestTenants :exec
INSERT INTO "tenants"(
    name,
    email,
    password,
    phone,
    website_url,
    address,
    status,
    settings
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);