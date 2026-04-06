package main

import (
	"log"
	"os"

	"new-auth-service/db"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if _, err := db.Connect(databaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close(db.DB)

	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	switch direction {
	case "up":
		if err := db.Migrate(); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
	case "down":
		if err := db.MigrateDown(); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
	default:
		log.Fatalf("Unknown direction: %s. Use 'up' or 'down'", direction)
	}
}
