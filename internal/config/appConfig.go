package config

import (
	"os"

	env "github.com/joho/godotenv"
)

type AppConfig struct {
	Port  string
	DbURL string
}

var GlobalAppConfig *AppConfig

func LoadConfig() string {

	// Load environment variables from a file
	if err := env.Load(".env"); err != nil {
		return err.Error()
	}

	// Retrieve and check the PORT environment variable
	// returns: string containing a PORT number on which application should run
	port := os.Getenv("PORT")
	if port == "" {
		return "PORT number not found"
	}

	// Retrieve and check the DB Connection String variable
	// returns: string containing database connection credentials
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return "Writer URL not found"
	}

	GlobalAppConfig = &AppConfig{
		Port:  port,
		DbURL: dbURL,
	}

	return ""
}
