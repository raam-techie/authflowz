run:
	go run main.go

build:
	go build -o bin/auth-service ./main.go

# Load the environment variables
ifneq ("$(wildcard .env)","")
    include .env
    export
endif

# Database connection string for Goose
# Note: Goose uses a different driver format than golang-migrate
DB_DSN=user=$(DB_USERNAME) password=$(DB_PASSWORD) dbname=$(DB_DATABASE) host=localhost port=$(DB_PORT) sslmode=disable

# Usage: make goose-create name=init_schema
goose-create:
	goose -dir internal/db/migrations create $(name) sql

# Run all pending migrations
# Run Migration -- `make goose-up`
goose-up:
	GOOSE_DRIVER="$(GOOSE_DRIVER)" GOOSE_DBSTRING="$(DB_DSN)" GOOSE_MIGRATION_DIR="$(GOOSE_MIGRATIONS_DIR)" goose up

# Rollback the last migration
goose-down:
	GOOSE_DRIVER="$(GOOSE_DRIVER)" GOOSE_DBSTRING="$(DB_DSN)" GOOSE_MIGRATION_DIR="$(GOOSE_MIGRATIONS_DIR)" goose down

# Check migration status
goose-status:
	GOOSE_DRIVER="$(GOOSE_DRIVER)" GOOSE_DBSTRING="$(DB_DSN)" GOOSE_MIGRATION_DIR="$(GOOSE_MIGRATIONS_DIR)" goose status
