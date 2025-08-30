# Order Management API

This repository contains a RESTful API for a complex order management system, built with Go. It serves as a demonstration of modern Go development practices, including clean architecture, comprehensive testing, and automated documentation.

## Features

- **RESTful API**: Provides CRUD operations for Users, Companies, Addresses, Locations, Products, and Orders.
- **Tech Stack**: Built with Go, using Gorilla/Mux for routing and PostgreSQL for data storage.
- **Database Layer**: Interacts with the database using the XORM library.
- **Migrations**: Database schema is managed with `goose` migrations.
- **Testing**: Comprehensive integration and repository test suites using Ginkgo & Gomega.
- **Mocks**: Automatically generates mocks for interfaces using `go:generate` and GoMock.
- **Authentication**: Secures endpoints using JWT-based authentication.
- **Geocoding**: Integrates with the Google Maps API for address geocoding.
- **API Documentation**: Automatically generates Swagger 2.0 documentation from code comments.

## Project Structure

```
/
├── cmd/main.go         # Application entry point
├── db/migrations/      # SQL database migrations
├── docs/               # Swagger documentation output
├── internal/           # Core application logic
│   ├── api/            # API handlers and routing
│   ├── repos/          # Database repository layer
│   └── types/          # Shared data structures
├── scripts/            # Helper scripts for migrations
├── .env.sample         # Example environment file
├── docker-compose.yml  # Docker configuration for databases
└── Makefile.mk         # Make commands for development
```

## Prerequisites

- **Go**: Version 1.21 or later
- **Docker**: To run the PostgreSQL databases
- **Make**: To use the helper commands in the `Makefile.mk`

## Getting Started

Follow these steps to get the project running on your local machine.

### 1. Clone the Repository

```sh
git clone <your-repo-url>
cd api
```

### 2. Configure Environment

The application is configured using environment variables. Copy the sample file and update it with your specific values.

```sh
cp .env.sample .env
```

You will need to edit the `.env` file to provide a `JWT_SECRET` and a valid `GOOGLE_MAPS_API_KEY`.

### 3. Start Databases

The `docker-compose.yml` file starts the main development database and a separate test database.

```sh
docker-compose up -d
```

- The main database will be available on `localhost:5432`.
- The test database will be available on `localhost:5433`.

### 4. Install Dependencies

Install the required Go tools and project dependencies.

```sh
make install.deps
go mod tidy
```

### 5. Run Database Migrations

Apply all pending database migrations to set up the schema in both databases.

```sh
make db.migrate.up
```

## Running the Application

With the database running and migrations applied, start the API server:

```sh
go run ./cmd/main.go
```

The server will start on the default address `http://localhost:8080`.

## API Documentation

This project uses `swag` to generate interactive API documentation.

1.  **Generate the documentation:**
    ```sh
    make swagger.generate
    ```
2.  **View the documentation:**
    Navigate to [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) in your browser.

## Running Tests

The test suite requires the test database to be running (`docker-compose up -d`). To run all tests:

```sh
make test
```

## Makefile Commands

The `Makefile.mk` provides several useful commands to streamline development:

| Command                                   | Description                                                              |
| ----------------------------------------- | ------------------------------------------------------------------------ |
| `make install.deps`                       | Installs required Go command-line tools (`ginkgo`, `mockgen`, `goose`).    |
| `make test`                               | Runs the complete Ginkgo test suite.                                     |
| `make swagger.generate`                   | Generates Swagger API documentation into the `/docs/swagger` directory.  |
| `make db.migrate.up`                      | Applies all pending migrations to both the main and test databases.      |
| `make db.migrate.down`                    | Rolls back the last migration on both databases.                         |
| `make db.migrate.reset`                   | Rolls back all migrations on both databases.                             |
| `make db.migrate.create name=<file_name>` | Creates a new SQL migration file in `db/migrations`.                     |
| `make db.seed`                            | Seeds the database with initial data.                                    |