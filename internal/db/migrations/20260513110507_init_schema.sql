-- +goose Up
-- +goose StatementBegin

CREATE TABLE tenants (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id          INT NOT NULL,
    name                TEXT NOT NULL,
    email               TEXT NOT NULL,
    password            TEXT NOT NULL,
    phone               TEXT,
    website_url         TEXT,
    address             JSONB DEFAULT '{}'::JSONB,  -- structured address — queryable by country, city etc.
    status              VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    settings            JSONB DEFAULT '{}'::JSONB, -- {"mfa_enforced": true, "sso_only": false, "max_users": 50}
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE TABLE user_pools (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES tenants (id) ON DELETE RESTRICT,
    pool_name               TEXT NOT NULL,
    description             TEXT,
    sign_in_methods         JSONB NOT NULL DEFAULT '{"email": true}'::JSONB, -- e.g. {"email": true, "google": true, "github": false, "apple": false}
    app_type                VARCHAR(50) NOT NULL DEFAULT 'WEB',
    mfa_required            BOOLEAN NOT NULL DEFAULT FALSE, -- -- MFA settings
    password_policy         JSONB NOT NULL DEFAULT '{"min_length": 8, "require_uppercase": true, "require_numbers": true, "require_symbols": false}'::JSONB,
    token_ttl_seconds       INT NOT NULL DEFAULT 900,           -- 15 minutes
    refresh_ttl_seconds     INT NOT NULL DEFAULT 2592000,       -- 30 days
    max_devices             INT, -- concurrent device limit per user (NULL = unlimited)
    allowed_origins         TEXT[] DEFAULT '{}', -- CORS: which browser origins may make auth requests to this pool
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE app_clients (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    user_pool_id            UUID NOT NULL REFERENCES user_pools (id) ON DELETE RESTRICT,
    tenant_id               UUID NOT NULL REFERENCES tenants (id) ON DELETE RESTRICT,
    client_name             TEXT NOT NULL, -- identity
    client_id               TEXT NOT NULL, -- public-facing OAuth2 client identifier
    -- hashed with Argon2id — raw secret shown once at creation, then discarded
    client_secret_hash      TEXT,
    app_type                 VARCHAR(50) NOT NULL,
    -- per-client token TTL override (NULL = inherit from user_pool)
    token_ttl_override      INT,
    -- arbitrary client metadata: platform, sdk_version, webhook_url etc.
    metadata                JSONB DEFAULT '{}'::JSONB,
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL REFERENCES tenants (id) ON DELETE RESTRICT,
    user_pool_id            UUID NOT NULL REFERENCES user_pools (id) ON DELETE RESTRICT,
    first_name              TEXT,
    last_name               TEXT,
    display_name            VARCHAR(255) NOT NULL,
    email                   TEXT NOT NULL,
    email_verified          BOOLEAN NOT NULL DEFAULT FALSE,
    phone                   TEXT,
    phone_verified          BOOLEAN NOT NULL DEFAULT FALSE,
    avatar_url              TEXT,
    department              TEXT, 
    password_hash           TEXT, -- nullable: social/SSO-only users have no password
    mfa_enabled             BOOLEAN NOT NULL DEFAULT FALSE, -- MFA (TOTP)
    -- base32 seed for TOTP — encrypt at rest with app-level encryption
    mfa_secret              TEXT,
    -- activity tracking
    last_login_at           TIMESTAMPTZ,
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ
);

CREATE TABLE verification_tokens (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    -- SHA-256 hash of the raw one-time token sent to the user
    -- raw token is never stored; compare by hashing the submitted token
    token_hash              TEXT NOT NULL,
    type                    VARCHAR(200) NOT NULL,
    -- expiry windows by type:
    --   email_verification  24 hours
    --   password_reset      15 minutes
    --   magic_link          10 minutes
    --   phone_otp           5  minutes
    --   invite              7  days
    expires_at              TIMESTAMPTZ NOT NULL,
    -- set when consumed — do NOT delete rows; keep for audit trail
    -- guard: WHERE used_at IS NULL AND expires_at > NOW()
    used_at                 TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Audit logs are append-only and grow without bound.
-- We use RANGE partitioning by month so old partitions can be
-- archived to cold storage without touching the live table.
CREATE TABLE audit_logs (
    -- UUIDv7 is time-sortable: sorting by id == sorting by created_at
    -- Use pgcrypto gen_random_uuid() as fallback if UUIDv7 extension unavailable
    id                      UUID NOT NULL DEFAULT gen_random_uuid(),
    tenant_id               UUID NOT NULL,
    -- actor_id nullable: system/scheduled jobs have no human actor
    actor_id                UUID,
    -- actor_type: 'user' | 'system' | 'api_key' | 'service_account'
    actor_type              TEXT NOT NULL,
    -- namespaced action verb: "user.login.success", "pool.created", "client.secret_rotated"
    action                  TEXT NOT NULL,
    -- what entity was acted on
    resource_type           TEXT NOT NULL,
    resource_id             UUID,
    -- full diff for update events
    old_value               JSONB,
    new_value               JSONB,
    -- request context
    ip_address              INET,
    user_agent              TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
 
-- email lookup (login / contact search)
CREATE UNIQUE INDEX idx_tenants_email
    ON tenants (email)
    WHERE deleted_at IS NULL;
 
-- filter active tenants fast
CREATE INDEX idx_tenants_status
    ON tenants (status)
    WHERE deleted_at IS NULL;

-- pool_name unique within a tenant
CREATE UNIQUE INDEX idx_user_pools_name_tenant
    ON user_pools (tenant_id, pool_name)
    WHERE is_active = TRUE;
 
-- fast lookup by tenant (list all pools for a tenant)
CREATE INDEX idx_user_pools_tenant
    ON user_pools (tenant_id);

-- client_id must be globally unique (used in Authorization headers)
CREATE UNIQUE INDEX idx_app_clients_client_id
    ON app_clients (client_id);
 
-- list all clients for a pool
CREATE INDEX idx_app_clients_pool
    ON app_clients (user_pool_id)
    WHERE is_active = TRUE;
 
-- list all clients for a tenant
CREATE INDEX idx_app_clients_tenant
    ON app_clients (tenant_id)
    WHERE is_active = TRUE;

-- email unique per pool (not globally — same email can exist in different pools)
CREATE UNIQUE INDEX idx_users_email_pool
    ON users (email, user_pool_id)
    WHERE deleted_at IS NULL;
 
-- tenant-scoped user listing (admin panels)
CREATE INDEX idx_users_tenant
    ON users (tenant_id)
    WHERE deleted_at IS NULL;
 
-- pool-scoped user listing
CREATE INDEX idx_users_pool
    ON users (user_pool_id)
    WHERE deleted_at IS NULL;

-- primary lookup path: app hashes the submitted token and looks up the hash
CREATE UNIQUE INDEX idx_vtokens_hash
    ON verification_tokens (token_hash)
    WHERE used_at IS NULL;
 
-- detect "resend spam": has this user requested a token of this type recently?
CREATE INDEX idx_vtokens_user_type
    ON verification_tokens (user_id, type)
    WHERE used_at IS NULL;
 
-- background cleanup: purge expired unused tokens
CREATE INDEX idx_vtokens_expires
    ON verification_tokens (expires_at)
    WHERE used_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_vtokens_expires;
DROP INDEX IF EXISTS idx_vtokens_user_type;
DROP INDEX IF EXISTS idx_vtokens_hash;

DROP INDEX IF EXISTS idx_users_pool;
DROP INDEX IF EXISTS idx_users_tenant;
DROP INDEX IF EXISTS idx_users_email_pool;

DROP INDEX IF EXISTS idx_app_clients_tenant;
DROP INDEX IF EXISTS idx_app_clients_pool;
DROP INDEX IF EXISTS idx_app_clients_client_id;

DROP INDEX IF EXISTS idx_user_pools_tenant;
DROP INDEX IF EXISTS idx_user_pools_name_tenant;

DROP INDEX IF EXISTS idx_tenants_status;
DROP INDEX IF EXISTS idx_tenants_email;

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS verification_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS app_clients;
DROP TABLE IF EXISTS user_pools;
DROP TABLE IF EXISTS tenants;

-- +goose StatementEnd