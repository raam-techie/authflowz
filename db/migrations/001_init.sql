-- PostgreSQL 17: gen_random_uuid() is built-in, no extension required

CREATE TABLE tenants (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id  VARCHAR(50) UNIQUE NOT NULL,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone       VARCHAR(20),
    address     TEXT,
    website_url VARCHAR(500),
    -- subscription_plan VARCHAR(50) DEFAULT 'BASIC',
    status      VARCHAR(20) DEFAULT 'ACTIVE',
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

-- app_type:  WEB | MOBILE | SPA | NATIVE | M2M
-- auth_type: EMAIL_PASSWORD | OTP_EMAIL | OTP_PHONE | OAUTH | MAGIC_LINK | SSO
CREATE TABLE projects (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    client_id         VARCHAR(100) UNIQUE NOT NULL,
    client_secret_hash VARCHAR(255) NOT NULL,
    name              VARCHAR(255) NOT NULL,
    description       TEXT,
    app_type          VARCHAR(20) NOT NULL DEFAULT 'WEB',
    auth_type         VARCHAR(50) NOT NULL DEFAULT 'EMAIL_PASSWORD',
    scopes            JSONB DEFAULT '[]',
    redirect_uris     JSONB DEFAULT '[]',
    is_active         BOOLEAN DEFAULT TRUE,
    last_used_at      TIMESTAMP,
    created_at        TIMESTAMP DEFAULT NOW(),
    updated_at        TIMESTAMP DEFAULT NOW()
);

-- Users are flexible — email/phone/password_hash nullable.
-- The project's auth_type determines which fields are required at signup.
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    uid           VARCHAR(100) NOT NULL,
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255),         -- required when auth_type = EMAIL_PASSWORD | OTP_EMAIL | MAGIC_LINK
    phone         VARCHAR(20),          -- required when auth_type = OTP_PHONE
    department    VARCHAR(100),
    role          VARCHAR(50) DEFAULT 'USER',
    password_hash VARCHAR(255),         -- required when auth_type = EMAIL_PASSWORD
    is_active     BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    UNIQUE(tenant_id, uid),
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_projects_tenant_id ON projects(tenant_id);
CREATE INDEX idx_projects_client_id ON projects(client_id);
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
