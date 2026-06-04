package main

import (
	"log"
	"os"

	"new-auth-service/db"
	"new-auth-service/server"
)

func main() {

	if _, err := db.Connect("postgresql://postgres.jwubnrczhvebyggmwonu:CoralBay%401234!@aws-1-ap-south-1.pooler.supabase.com:6543/postgres"); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close(db.DB)

	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	router := server.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
