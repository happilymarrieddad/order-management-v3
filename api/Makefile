install.deps:
	go mod download github.com/onsi/ginkgo/v2
	go install github.com/onsi/ginkgo/v2/ginkgo
	go get github.com/onsi/gomega/...
	go install go.uber.org/mock/mockgen@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

db.seed: db.migrate.up
	go run tools/seed.go

db.migrate.create:
	goose -dir db/migrations create $(name) sql

db.migrate.up:
	goose -allow-missing -dir db/migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?connect_timeout=180&sslmode=disable" up
	goose -allow-missing -dir db/migrations postgres "postgres://postgres:postgres@localhost:5432/postgres-test?connect_timeout=180&sslmode=disable" up

db.migrate.down:
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?connect_timeout=180&sslmode=disable" down
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5432/postgres-test?connect_timeout=180&sslmode=disable" down

db.migrate.reset:
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?connect_timeout=180&sslmode=disable" reset
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:5432/postgres-test?connect_timeout=180&sslmode=disable" reset

db.migrate.validate:
	goose -dir=./db/migrations -v validate
	
swagger.generate:
	swag init -g ./cmd/main.go -o ./docs/swagger
