package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	"github.com/happilymarrieddad/order-management-v3/api/types"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for xorm
	"xorm.io/xorm"
)

func main() {
	logger := log.New(os.Stdout, "seeder | ", log.LstdFlags)
	ctx := context.Background()

	// --- Database Connection ---
	// This uses the same environment variables as your API and migrate script.
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "postgres"
	}
	dbSSLMode := os.Getenv("DB_SSL_MODE")
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := xorm.NewEngine("pgx", dsn)
	if err != nil {
		logger.Fatalf("FATAL: unable to create db engine: %v", err)
	}

	if err = db.Ping(); err != nil {
		logger.Fatalf("FATAL: unable to connect to db: %v", err)
	}
	logger.Println("database connection established")
	db.Exec("TRUNCATE TABLE users, companies, addresses RESTART IDENTITY CASCADE")

	// --- Repository Initialization ---
	// We pass nil for the Google Client as it's not needed for seeding.
	globalRepo := repos.NewGlobalRepo(db, nil)

	// --- Seeding Logic ---
	logger.Println("starting to seed data...")

	// 1. Create an Address
	addr := &types.Address{Line1: "123 Seed St", City: "Seedville", State: "SD", PostalCode: "54321"}
	createdAddr, err := globalRepo.Addresses().Create(ctx, addr)
	if err != nil {
		logger.Fatalf("failed to create address: %v", err)
	}
	logger.Printf("created address with ID: %d\n", createdAddr.ID)

	// 2. Create a Company
	comp := &types.Company{Name: "Seed Company Inc.", AddressID: createdAddr.ID}
	if err := globalRepo.Companies().Create(ctx, comp); err != nil {
		logger.Fatalf("failed to create company: %v", err)
	}
	logger.Printf("created company with ID: %d\n", comp.ID)

	// 3. Create Users
	usersToCreate := []*types.User{
		{FirstName: "Test", LastName: "User", Email: "test@test.com", Password: "password", CompanyID: comp.ID, AddressID: addr.ID},
		{FirstName: "Alice", LastName: "Smith", Email: "alice@example.com", Password: "password123", CompanyID: comp.ID, AddressID: addr.ID},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob@example.com", Password: "password123", CompanyID: comp.ID, AddressID: addr.ID},
	}

	for _, user := range usersToCreate {
		if err := globalRepo.Users().Create(ctx, user); err != nil {
			logger.Fatalf("failed to create user '%s': %v", user.Email, err)
		}
		logger.Printf("created user with email: %s\n", user.Email)
	}

	// --- Output ---
	fmt.Println("\n" + "========================================")
	fmt.Println("🌱 Seeding Complete!")
	fmt.Println("A simple user has been created for API testing.")
	fmt.Printf("Email:    %s\n", usersToCreate[0].Email)
	fmt.Printf("Password: %s\n", "password")
	fmt.Println("========================================")
}
