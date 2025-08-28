# Go Order Management API (Gemini Code Assist Demo)

Welcome! This repository is a demonstration of how **Gemini Code Assist** can be used as a pair-programming partner to build a robust, well-tested, and modern web API in Go. Every component, from the database schema and migrations to the repository layer and integration tests, was developed iteratively with Gemini's help.

## Overview

This project is a simple RESTful API for a conceptual order management system. It includes basic CRUD operations for core entities like Users, Companies, and Addresses, built with a focus on clean architecture and data integrity.

### Tech Stack

- **Language**: Go
- **Web Framework**: [Gorilla/Mux](https://github.com/gorilla/mux)
- **Database**: PostgreSQL
- **ORM/DB Layer**: [XORM](https://xorm.io/)
- **Testing**: [Ginkgo](https://onsi.github.io/ginkgo/) & [Gomega](https://onsi.github.io/gomega/)
- **Mocks**: [GoMock](https://github.com/uber-go/mock)
- **Migrations**: [Goose](https://github.com/pressly/goose)
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)

## Project Philosophy & Structure

This project adheres to several conventions to maintain code quality and developer productivity:

- **Application Entrypoint**: The main application is launched from `cmd/api/main.go`.
- **Mock Generation**: We use `go:generate` comments with `mockgen` to automatically create mocks for our repository interfaces. This ensures that our mocks are always in sync with the interfaces they are mocking.
- **Testing Structure**:
    - Each source file (e.g., `users_repo.go`) has a corresponding test file (`users_repo_test.go`).
    - The Ginkgo test suite bootstrap logic is kept in a separate `_suite_test.go` file within each tested package to keep test files focused on testing logic.

## Getting Started

Follow these steps to get the project running on your local machine.

### Prerequisites

- **Go**: Version 1.21 or later.
- **Docker**: To run the PostgreSQL database.
- **Make**: To use the helper commands in the `Makefile.mk`.

### Installation

1.  **Clone the repository:**
    ```sh
    git clone <your-repo-url>
    cd api
    ```

2.  **Install Go tools:**
    The `install.deps` Make target will install all required Go command-line tools (`ginkgo`, `mockgen`, `goose`).
    ```sh
    make install.deps
    ```

3.  **Install project dependencies:**
    ```sh
    go mod tidy
    ```

## Running the Application

### 1. Start the Database

A `docker-compose.yml` file is provided to easily start both the main and test databases.

```sh
# Start the databases in the background
docker-compose up -d
```

- The main database will be available at `localhost:5432`.
- The test database will be available at `localhost:5433`.

### 2. Run Database Migrations

Apply all database migrations to both the main and test databases using the Make command. This command uses the platform-agnostic scripts in the `/scripts` directory.

```sh
make db.migrate.up
```

### 3. Run the API Server

With the database running and migrations are applied, start the API server.

```sh
# The entrypoint is cmd/api/main.go
go run ./cmd/api
```

## Running Tests

The integration tests require the test database to be running (see step 1 above). The test suite automatically handles running migrations before tests begin.

To run the full test suite, use the `test` Make target, which uses Ginkgo to execute the tests.

```sh
make test
```

This command will run all tests in the `internal/repos` package, providing detailed output on their execution.

## Makefile Commands

The `Makefile.mk` provides several useful commands to streamline development:

- `make install.deps`: Installs required Go tools.
- `make test`: Runs the Ginkgo integration test suite.
- `make db.migrate.up`: Applies all pending migrations to both databases.
- `make db.migrate.down`: Rolls back the last migration on both databases.
- `make db.migrate.reset`: Rolls back all migrations on both databases.
- `make db.migrate.create name=<migration_name>`: Creates a new SQL migration file.

---

*This project was built with Gemini Code Assist.*


## Notes for the A.I.

I want you to notice how cmd/main.go is the entrance point for the application. We are using gorilla mux as the framework for a golang api. Xorm is the driver and postgres is the database. Ginkgo is the testing framework. Please take special note that each file has it's own test file and the ginkgo bootstrap is kept seperate. Mocks are generated using uber mockgen and a mockgen line us put above interfaces using go:generate. Please read the README.md and make updates as necessary.

To ensure our repository tests are consistent, isolated, and easy to understand, please follow these conventions. This is especially important when working with AI assistants.

*   **Preserve Existing Tests**: When asked to update tests, do not remove or rewrite existing tests from scratch. The goal is to *add* to the existing test suite or make targeted updates to existing tests to accommodate new fields or logic.

*   **Dependency Creation**: All required data for a test (e.g., a `Company` or `Address` needed to create a `User`) **must** be created within a `BeforeEach` block. This ensures each test runs with a fresh, known set of data and avoids reliance on a pre-seeded database.

*   **Creation Order**: Pay close attention to dependencies. If one type requires an ID from another (e.g., a `Company` requires an `AddressID`), create the dependency (`Address`) first.
    company = &types.Company{Name: "Test Co Inc.", AddressID: address.ID}
    ```

*   **No Test Consolidation**: Do not combine multiple test cases (`It` blocks) into a single one for the sake of brevity. Each test should remain granular and focus on a specific scenario. This is crucial for clarity and debugging.

*   **Suite-Level Setup**: Global test setup (like initializing the database connection `db`, the global repository `gr`, and the context `ctx`) and cleanup (like truncating tables between tests) is handled in the `_suite_test.go` file for the package. Individual test files should not need to repeat this logic.

### Handler (API) Test Conventions

To ensure our API handler tests are robust and consistent, please follow these conventions:

*   **Test Suite (`_suite_test.go`)**: Each handler package (e.g., `v1/users`) must have its own `_suite_test.go` file. This file is responsible for:
    *   Initializing `gomock` and creating mocks for all required repositories (`UsersRepo`, `CompaniesRepo`, etc.).
    *   Setting up the mock chain (e.g., `mockGlobalRepo.EXPECT().Users().Return(mockUsersRepo)`).
    *   Providing a helper function (e.g., `createRequestWithRepo`) to generate `http.Request` objects for tests.

*   **Mock Injection via Context**: The test request helper function must inject the mocked `GlobalRepo` into the request's context using the correct key (`middleware.RepoKey`). This is critical for simulating the application's dependency injection pattern.
    ```go
    // Example from a test suite helper function
    ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)
    req = req.WithContext(ctxWithRepo)
    ```

*   **Thorough Test Cases**: Each handler should be tested for multiple scenarios using separate `Context` blocks:
    *   The "happy path" (a valid request that succeeds).
    *   Invalid request bodies (malformed JSON, missing required fields, validation errors like mismatched passwords).
    *   Dependency failures (e.g., trying to create a user for a `company_id` that doesn't exist).
    *   Repository errors (simulating database failures like unique constraint violations).

*   **JSON Payload Casing**: Pay close attention to JSON key casing in test payloads. This project uses **`snake_case`**. Mismatched casing will cause validation to fail and result in `400 Bad Request` errors.

## API Documentation (Swagger)

This project uses swaggo/swag to automatically generate interactive API documentation from Go comments.

### Documentation Workflow

1.  **Global Configuration**: The main API definition (e.g., `@title`, `@version`, `@securityDefinitions`) is located in `docs/docs.go`. This file serves as the entry point for the documentation generator.

2.  **Define Payloads**: For `POST` and `PUT` endpoints, create dedicated, exported structs for the request bodies (e.g., `users.CreateUserPayload`). This is crucial for generating clean documentation. Use `example` tags to provide sample data.

3.  **Annotate Handlers**: Add a comment block above each handler function with the following annotations:
    *   `@Summary`: A short summary of what the endpoint does.
    *   `@Description`: A more detailed description.
    *   `@Tags`: A group name for the endpoint (e.g., `users`, `companies`).
    *   `@Accept` & `@Produce`: The content types (usually `json`).
    *   `@Param`: Describes a parameter (path, query, or body).
    *   `@Success`: Describes a successful response (status code and response model).
    *   `@Failure`: Describes a potential error response (status code and error model, e.g., `middleware.ErrorResponse`).
    *   `@Router`: Defines the route path and HTTP method.

4.  **Generate/Update Documentation**: After adding or changing annotations, run the following command from the root of the `api` directory:
    ```sh
    swag init -g docs/docs.go
    ```
*   **Populate All Required Fields**: When creating test data (e.g., `&types.User{...}`), ensure all required fields as defined by the struct and database schema are populated. Missing fields will cause validation or database errors during tests.
