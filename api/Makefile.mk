install.deps:
	go mod download github.com/onsi/ginkgo/v2
	go install github.com/onsi/ginkgo/v2/ginkgo
	go get github.com/onsi/gomega/...
	go install go.uber.org/mock/mockgen@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

test:
	ginkgo -r -v --fail-fast --randomize-all

db.seed: db.migrate.up
	go run tools/seed.go

# Use 'name' variable for creating new migrations, e.g., `make db.migrate.create name=add_new_field`
db.migrate.create:
	goose -dir db/migrations create $(name) sql

# These targets run migrations against both the main and test databases.
# They use the helper scripts which respect environment variables for connection details.
DB_NAME_MAIN ?= postgres
DB_NAME_TEST ?= postgres-test

# Detect OS and choose the correct script
ifeq ($(OS),Windows_NT)
    MIGRATE_SCRIPT := ./scripts/migrate.cmd
else
    MIGRATE_SCRIPT := ./scripts/migrate.sh
endif

db.migrate.up:
	DB_NAME=$(DB_NAME_MAIN) $(MIGRATE_SCRIPT) up
	DB_NAME=$(DB_NAME_TEST) $(MIGRATE_SCRIPT) up

db.migrate.down:
	DB_NAME=$(DB_NAME_MAIN) $(MIGRATE_SCRIPT) down
	DB_NAME=$(DB_NAME_TEST) $(MIGRATE_SCRIPT) down

db.migrate.reset:
	DB_NAME=$(DB_NAME_MAIN) $(MIGRATE_SCRIPT) reset
	DB_NAME=$(DB_NAME_TEST) $(MIGRATE_SCRIPT) reset

db.migrate.validate:
	goose -dir=./db/migrations -v validate

swagger.generate:
	swag init -g ./cmd/main.go -o ./docs/swagger
