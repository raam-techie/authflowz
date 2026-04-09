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

		`CREATE TABLE IF NOT EXISTS projects (
			id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id          UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			client_id          TEXT NOT NULL UNIQUE,
			client_secret_hash TEXT NOT NULL,
			name               TEXT NOT NULL,
			description        TEXT,
			app_type           TEXT NOT NULL,
			auth_types         JSONB NOT NULL,
			redirect_uris      JSONB NOT NULL DEFAULT '[]',
			is_active          BOOLEAN NOT NULL DEFAULT TRUE,
			created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS users (
			id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			project_id    UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
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
			UNIQUE(project_id, email)
		)`,
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
		`DROP TABLE IF EXISTS users`,
		`DROP TABLE IF EXISTS projects`,
		`DROP TABLE IF EXISTS tenants`,
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
