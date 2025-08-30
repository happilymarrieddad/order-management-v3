package repos_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/happilymarrieddad/order-management-v3/api/internal/repos"
	_ "github.com/jackc/pgx/v5/stdlib"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"xorm.io/xorm"
)

var (
	ctx context.Context
	db  *xorm.Engine
	gr  repos.GlobalRepo
)

func TestRepos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repos Suite")
}

var _ = BeforeSuite(func() {
	ctx = context.Background()
	var err error

	// Best practice: use environment variables for database connection details
	// with sensible defaults for local development (e.g., using Docker).
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "postgres") + "-test" // Use a separate test database
	dbSslMode := getEnv("DB_SSL_MODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSslMode)

	log.Printf("Connecting to test database at %s:%s/%s as user %s", dbHost, dbPort, dbName, dbUser)

	// Using pgx as the driver is a modern best practice for performance and features.
	db, err = xorm.NewEngine("pgx", connStr)
	Expect(err).NotTo(HaveOccurred())

	// Ping the database to ensure a connection is established.
	// Skip tests if the database is not available.
	if err = db.Ping(); err != nil {
		Skip(fmt.Sprintf("Skipping repository tests: could not connect to postgres database: %v", err))
	}

	// Run database migrations to ensure the schema is up-to-date.
	// This is a best practice as it ensures the test DB schema matches production.
	// We set the command's working directory to the project root (`api` dir) so the
	// script can find the 'db/migrations' folder correctly.
	/**
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("../../scripts/migrate.cmd", "up")
	} else {
		cmd = exec.Command("../../scripts/migrate.sh", "up")
	}
	// The DB_NAME is overridden to target the test database.
	cmd.Env = append(os.Environ(), "DB_NAME="+dbName)

	output, err := cmd.CombinedOutput()
	Expect(err).NotTo(HaveOccurred(), "Migration script failed to run. Output:\n%s", string(output))
	**/

	gr = repos.NewGlobalRepo(db, nil)
})

var _ = AfterSuite(func() {
	if db != nil {
		Expect(db.Close()).To(Succeed())
	}
})

var _ = BeforeEach(func() {
	// For Postgres, truncating tables is a fast and effective way to ensure test isolation.
	// 'RESTART IDENTITY' resets auto-incrementing primary keys.
	// 'CASCADE' truncates dependent tables via foreign keys.
	_, err := db.Exec("TRUNCATE TABLE users, companies, addresses, commodity_attributes RESTART IDENTITY CASCADE")
	Expect(err).NotTo(HaveOccurred())
})

// getEnv is a helper to read an environment variable or return a default value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
