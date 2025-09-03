# Order Management API

This repository contains a RESTful API for a complex order management system, built with Go. It serves as a demonstration of modern Go development practices, including clean architecture, comprehensive testing, and automated documentation.

## Features

- **RESTful API**: Provides CRUD operations for Users, Companies, Addresses, Locations, Products, and Orders.
- **Tech Stack**: Built with Go, using Gorilla/Mux for routing and PostgreSQL for data storage.
- **Database Layer**: Interacts with the database using the XORM library.
- **Migrations**: Database schema is managed with `goose` migrations.
- **Testing**: Comprehensive integration and repository test suites using Ginkgo & Gomega, including robust role-based authentication testing.
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

## Data Model

The API revolves around a few core concepts that define how products are categorized and described, as well as fundamental entities for managing users, companies, and locations.

*   **`User`**: Represents an individual user of the system, associated with a `Company` and an `Address`. Users have roles that define their permissions.
*   **`Company`**: Represents an organization within the system, associated with an `Address`. Companies own products, locations, and users.
*   **`Address`**: A reusable entity for storing physical addresses, used by `Users`, `Companies`, and `Locations`.
*   **`CompanyAttribute`**: A link between a `Company` and a `CommodityAttribute`, allowing a company to specify which attributes are relevant to its products. It features a `position` field that auto-increments per company, managed by a database trigger.
*   **`Location`**: Represents a specific physical location (e.g., a warehouse, office) belonging to a `Company`, and linked to an `Address`.

*   **`Commodity`**: This is the most general classification. It represents a fundamental good, like "Potatoes" or "Apples". It has a `CommodityType`, such as "Produce".
*   **`CommodityAttribute`**: This defines a *property* that a `Commodity` can have. For example, attributes for the "Produce" type could be "Color", "Size", or "Grade". These attributes are linked to the `CommodityType`, not to a specific `Commodity`.
*   **`Product`**: This is a *specific, sellable item* that belongs to a `Company`. It's an instance of a `Commodity`. For example, a `Product` could be "Organic Russet Potatoes" which is a `Commodity` of "Potatoes", sold by a specific `Company`.
*   **`ProductAttributeValue`**: This is where the concepts connect. It assigns a specific `Value` to a `CommodityAttribute` for a particular `Product`.

### Example Flow

1.  You start with a **`Commodity`**:
    *   `Name`: "Apple"
    *   `CommodityType`: `Produce`

2.  You define **`CommodityAttribute`**s for the `Produce` type:
    *   `Name`: "Variety"
    *   `Name`: "Color"

3.  A `Company` creates a **`Product`**:
    *   It's linked to the "Apple" `Commodity`.
    *   The `Company` gives it a descriptive name, like "Fresh Granny Smith Apples".

4.  Finally, you use **`ProductAttributeValue`** to describe this specific `Product`:
    *   **Value 1**: Links the `Product` to the "Variety" `CommodityAttribute` with the value "Granny Smith".
    *   **Value 2**: Links the `Product` to the "Color" `CommodityAttribute` with the value "Green".

This structure allows for a flexible and detailed product catalog where general commodities can be customized with specific attributes for different products sold by different companies.

## Security & Authorization

The API employs a multi-layered security model to protect endpoints and ensure data isolation between tenants.

*   **JWT Authentication**: All endpoints under `/api` are protected and require a valid JSON Web Token (JWT) to be passed in the `X-App-Token` header. The `/login` endpoint is used to obtain this token.

*   **Role-Based Access Control (RBAC)**: The system defines two primary roles:
    *   **`User`**: Standard users with limited permissions.
    *   **`Admin`**: Superusers who can perform administrative tasks.
    Many endpoints (like creating companies or deleting resources) are restricted to admins only.

*   **Ownership-Based Access (Multi-tenancy)**: This is the core of the security model. A user's actions are scoped to their own `Company`. For example, a standard user can only create new users for their own company and can only update their own user profile. This prevents users from one company from viewing or modifying the data of another. While admins have broader permissions, they are generally not exempt from these ownership checks and can only operate within their own company's data.

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

You will need to edit the `.env` file to provide a `JWT_SECRET` and a valid `GOOGLE_MAPS_API_KEY`. The Google Maps API Key is used for geocoding addresses.

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