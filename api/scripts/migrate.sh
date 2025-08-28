#!/bin/bash
# This script runs database migrations using goose against the PostgreSQL container.
# It runs the command against both the main and test databases.
#
# Usage:
#   ./scripts/migrate.sh up       # Apply all pending migrations on both databases
#   ./scripts/migrate.sh down     # Roll back the last migration on both databases
#   ./scripts/migrate.sh status   # Show migration status for both databases

# Exit immediately if a command exits with a non-zero status.
set -e

# The directory where migration files are stored, relative to the project root.
MIGRATIONS_DIR="db/migrations"

# --- Database Connection Details ---
# Use environment variables if they are set, otherwise use sensible defaults.
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_SSL_MODE="${DB_SSL_MODE:-disable}"

# --- Databases to Migrate ---
# A space-separated list of database names to apply migrations to.
DATABASES_TO_MIGRATE="postgres postgres-test"

# Set Goose driver
export GOOSE_DRIVER="postgres"

# Loop through each database and run the goose command
for DB_NAME in $DATABASES_TO_MIGRATE; do
  echo ""
  echo "--- Targeting database: $DB_NAME ---"
  
  export GOOSE_DBSTRING="host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSL_MODE}"
  
  # Run the goose command, passing all script arguments to it (e.g., up, down, status).
  goose -allow-missing -dir "${MIGRATIONS_DIR}" "$@"
done

echo ""
echo "--- All migrations complete. ---"
