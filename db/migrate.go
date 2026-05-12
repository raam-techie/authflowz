package db

import "log"

func Migrate() error {
	queries := []string{
		`CREATE SEQUENCE IF NOT EXISTS user_uid_seq START 100000 INCREMENT 1 MAXVALUE 999999 CYCLE`,

		`CREATE TABLE IF NOT EXISTS tenants (
			id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			account_id    TEXT NOT NULL UNIQUE,
			name          TEXT NOT NULL,
			email         TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			phone         TEXT,
			address       TEXT,
			website_url   TEXT,
			status        TEXT NOT NULL DEFAULT 'ACTIVE',
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`DO $$ BEGIN
			CREATE TYPE app_type_enum AS ENUM ('WEB', 'MOBILE', 'SPA', 'NATIVE', 'M2M');
		EXCEPTION WHEN duplicate_object THEN NULL;
		END $$`,

		`CREATE TABLE IF NOT EXISTS user_pools (
			id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id        UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			pool_name        TEXT NOT NULL,
			pool_id          TEXT NOT NULL UNIQUE,
			description      TEXT,
			sign_in_methods  JSONB NOT NULL,
			app_type         app_type_enum NOT NULL,
			is_active        BOOLEAN NOT NULL DEFAULT TRUE,
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`ALTER TABLE user_pools ADD COLUMN IF NOT EXISTS app_type app_type_enum NOT NULL DEFAULT 'WEB'`,

		`CREATE INDEX IF NOT EXISTS idx_user_pools_tenant_id ON user_pools(tenant_id)`,

		`CREATE TABLE IF NOT EXISTS app_clients (
			id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_pool_id       UUID NOT NULL REFERENCES user_pools(id) ON DELETE CASCADE,
			tenant_id          UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			client_name        TEXT NOT NULL,
			client_id          TEXT NOT NULL UNIQUE,
			client_secret_hash TEXT NOT NULL,
			app_type           TEXT NOT NULL,
			is_active          BOOLEAN NOT NULL DEFAULT TRUE,
			created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE INDEX IF NOT EXISTS idx_app_clients_user_pool_id ON app_clients(user_pool_id)`,

		`CREATE TABLE IF NOT EXISTS users (
			id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			app_client_id UUID NOT NULL REFERENCES app_clients(id) ON DELETE CASCADE,
			uid           TEXT NOT NULL,
			name          TEXT NOT NULL,
			email         TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			phone         TEXT,
			department    TEXT,
			role          TEXT NOT NULL DEFAULT 'USER',
			is_active     BOOLEAN NOT NULL DEFAULT TRUE,
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(tenant_id, uid),
			UNIQUE(app_client_id, email)
		)`,

		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id       UUID REFERENCES users(id) ON DELETE CASCADE,
			app_client_id UUID REFERENCES app_clients(id) ON DELETE CASCADE,
			tenant_id     UUID REFERENCES tenants(id) ON DELETE CASCADE,
			token_hash    TEXT NOT NULL UNIQUE,
			expires_at    TIMESTAMPTZ NOT NULL,
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash)`,

		`CREATE TABLE IF NOT EXISTS user_attributes (
			id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			key        TEXT NOT NULL,
			value      TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, key)
		)`,

		`CREATE TABLE IF NOT EXISTS audit_logs (
			id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id      UUID        REFERENCES tenants(id) ON DELETE SET NULL,
			app_client_id  UUID        REFERENCES app_clients(id) ON DELETE SET NULL,
			user_id        UUID        REFERENCES users(id) ON DELETE SET NULL,
			action         TEXT        NOT NULL,
			actor_type     TEXT        NOT NULL,
			ip_address     TEXT,
			user_agent     TEXT,
			status         TEXT        NOT NULL,
			failure_reason TEXT,
			metadata       JSONB,
			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id  ON audit_logs(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id    ON audit_logs(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_action     ON audit_logs(action)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC)`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			return err
		}
	}

	log.Println("Database migration completed successfully")
	return nil
}

func MigrateDown() error {
	queries := []string{
		`DROP TABLE IF EXISTS audit_logs CASCADE`,
		`DROP TABLE IF EXISTS user_attributes CASCADE`,
		`DROP TABLE IF EXISTS refresh_tokens CASCADE`,
		`DROP TABLE IF EXISTS users CASCADE`,
		`DROP TABLE IF EXISTS app_clients CASCADE`,
		`DROP TABLE IF EXISTS user_pools CASCADE`,
		`DROP TABLE IF EXISTS tenants CASCADE`,
		`DROP SEQUENCE IF EXISTS user_uid_seq`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			return err
		}
	}

	log.Println("Database migration rolled back successfully")
	return nil
}
