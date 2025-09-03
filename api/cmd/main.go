package main

import (
	"fmt"
	"log"
	"os"

	"github.com/happilymarrieddad/order-management-v3/api/internal/api"
	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
	"xorm.io/xorm"
)

// appConfig holds all configuration for the application.
type appConfig struct {
	DSN          string
	GoogleAPIKey string
}

// loadConfig reads configuration from environment variables and populates an appConfig struct.
func loadConfig() (*appConfig, error) {
	if os.Getenv("DEV") == "true" {
		// In development, load .env file. In production, env vars are expected to be set directly.
		if err := godotenv.Load(); err != nil {
			// This is a warning, not a fatal error, as env vars could be set by the system.
			log.Printf("Warning: could not load .env file: %v", err)
		}
	}

	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbUser := utils.GetEnv("DB_USER", "postgres")
	dbPassword := utils.GetEnv("DB_PASSWORD", "postgres")
	dbName := utils.GetEnv("DB_NAME", "postgres")
	dbSslMode := utils.GetEnv("DB_SSL_MODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSslMode)

	cfg := &appConfig{
		DSN:          dsn,
		GoogleAPIKey: os.Getenv("GOOGLE_MAPS_API_KEY"), // No fallback, empty string is a valid state we check for later.
	}

	return cfg, nil
}

// In cmd/api/main.go

// @title           Order Management API
// @version         1.0
// @description     This is the API for the Order Management System.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	// Logger setup
	logger := log.New(os.Stdout, "api | ", log.LstdFlags)

	// --- Configuration ---
	cfg, err := loadConfig()
	if err != nil {
		logger.Fatalf("FATAL: failed to load configuration: %v", err)
	}

	// --- Database Connection ---
	db, err := xorm.NewEngine("pgx", cfg.DSN)
	if err != nil {
		logger.Fatalf("FATAL: unable to create db engine: %v", err)
	}

	if err = db.Ping(); err != nil {
		logger.Fatalf("FATAL: unable to connect to db: %v", err)
	}
	logger.Println("database connection established")

	// --- Google Maps Client ---
	var googleClient *maps.Client
	if cfg.GoogleAPIKey != "" {
		var err error
		googleClient, err = maps.NewClient(maps.WithAPIKey(cfg.GoogleAPIKey))
		if err != nil {
			logger.Fatalf("FATAL: unable to create google maps client: %v", err)
		}
	} else {
		logger.Println("warning: GOOGLE_MAPS_API_KEY not set, geocoding will be unavailable")
	}

	// --- Repository Initialization ---
	globalRepo := repos.NewGlobalRepo(db, googleClient)

	// Create a new server instance
	api.Run(globalRepo, logger)
}
