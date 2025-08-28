<#
.SYNOPSIS
    Runs database migrations using goose against the PostgreSQL containers.
.DESCRIPTION
    This script provides a consistent, single command for all developers to manage the database schema
    for both the main and test databases.
.PARAMETER Args
    All arguments passed to the script will be forwarded to the goose command (e.g., up, down, status).
.EXAMPLE
    ./scripts/migrate.ps1 up
    Applies all pending migrations to both databases.
.EXAMPLE
    ./scripts/migrate.ps1 status
    Shows the migration status for both databases.
#>
# Exit on any error.
$ErrorActionPreference = "Stop"

# The directory where migration files are stored, relative to the project root.
$MigrationsDir = "db/migrations"

# --- Database Connection Details ---
# Use environment variables if they are set, otherwise use sensible defaults.
$DbHost = if ($env:DB_HOST) { $env:DB_HOST } else { "localhost" }
$DbPort = if ($env:DB_PORT) { $env:DB_PORT } else { "5432" }
$DbUser = if ($env:DB_USER) { $env:DB_USER } else { "postgres" }
$DbPassword = if ($env:DB_PASSWORD) { $env:DB_PASSWORD } else { "postgres" }
$DbSslMode = if ($env:DB_SSL_MODE) { $env:DB_SSL_MODE } else { "disable" }

# --- Databases to Migrate ---
# A list of database names to apply migrations to.
$DatabasesToMigrate = @("postgres", "postgres-test")

# Set Goose driver and loop through each database to run the goose command.
$env:GOOSE_DRIVER = "postgres"
foreach ($DbName in $DatabasesToMigrate) {
    Write-Host ""
    Write-Host "--- Targeting database: $DbName ---" -ForegroundColor Green
    $env:GOOSE_DBSTRING = "host=$DbHost port=$DbPort user=$DbUser password=$DbPassword dbname=$DbName sslmode=$DbSslMode"
    # Pass all script arguments ($args) to the goose command
    goose -allow-missing -dir $MigrationsDir $args
}

Write-Host ""
Write-Host "--- All migrations complete. ---" -ForegroundColor Green