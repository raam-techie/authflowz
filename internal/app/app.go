package app

import (
	"log"
	"new-auth-service/internal"
	"new-auth-service/internal/config"
	"new-auth-service/internal/db/sqlc"
	"new-auth-service/internal/server"
)

// App initializes application dependencies such as configuration,
// database connection, sqlc store, router setup, and starts the HTTP server.
//
// Flow:
//  1. Load application configuration
//  2. Initialize PostgreSQL connection pool
//  3. Initialize sqlc query store
//  4. Register database dependencies globally
//  5. Setup HTTP router
//  6. Start HTTP server
func App() {

	// Load application configuration
	configErr := config.LoadConfig()
	if configErr != "" {
		log.Fatalf("failed to load application config: %s", configErr)
	}

	// Initialize PostgreSQL connection pool
	dbPool, err := config.NewStore(config.GlobalAppConfig.DbURL)
	if err != nil {
		log.Fatalf("failed to initialize database connection: %v", err)
	}

	// Initialize sqlc generated store
	sqlStore := sqlc.NewStore(dbPool)

	// Register database dependencies globally
	err = internal.SetDB(dbPool, sqlStore)
	if err != nil {
		log.Fatalf("failed to register database dependencies: %v", err)
	}

	// Setup application routes
	router := server.SetupRouter()

	if err := router.Run(config.GlobalAppConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
