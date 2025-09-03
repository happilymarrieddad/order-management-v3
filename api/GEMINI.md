*   **Explicit Instructions for Schema Changes**: Do not create any new database tables or migration files unless I explicitly ask you to. We will discuss and design data models first, and you will wait for a direct command (e.g., "create the migration for X") before generating any SQL or files.

To ensure our repository tests are consistent, isolated, and easy to understand, please follow these conventions. This is especially important when working with AI assistants.

*   **Preserve Existing Tests**: When asked to update tests, do not remove or rewrite existing tests from scratch. The goal is to *add* to the existing test suite or make targeted updates to existing tests to accommodate new fields or logic.

*   **Dependency Creation**: All required data for a test (e.g., a `Company` or `Address` needed to create a `User`) **must** be created within a `BeforeEach` block. This ensures each test runs with a fresh, known set of data and avoids reliance on a pre-seeded database.

*   **Creation Order**: Pay close attention to dependencies. If one type requires an ID from another (e.g., a `Company` requires an `AddressID`), create the dependency (`Address`) first.
    ```go
    company = &types.Company{Name: "Test Co Inc.", AddressID: address.ID}
    ```

*   **No Test Consolidation**: Do not combine multiple test cases (`It` blocks) into a single one for the sake of brevity. Each test should remain granular and focus on a specific scenario. This is crucial for clarity and debugging.

*   **Suite-Level Setup**: Global test setup (like initializing the database connection `db`, the global repository `gr`, and the context `ctx`) and cleanup (like truncating tables between tests) is handled in the `_suite_test.go` file for the package. Individual test files should not need to repeat this logic.
    *   When adding a new repository, ensure its corresponding table is included in the `TRUNCATE TABLE ... RESTART IDENTITY CASCADE` statement within the `BeforeEach` block of `repos_suite_test.go` to maintain proper test isolation and prevent data conflicts between tests.
    *   Example: `TRUNCATE TABLE users, companies, company_attributes ...`

*   **Default to Visible**: For repositories that have a `visible` column, the standard `Get` and `Find` methods should, by default, only return records where `visible = true`. Provide a separate method (e.g., `GetIncludeInvisible`) for cases where non-visible records need to be accessed.

### Handler (API) Test Conventions

To ensure our API handler tests are robust and consistent, please follow these conventions:

*   **Test Suite (`_suite_test.go`)**: Each handler package (e.g., `v1/users`) must have its own `_suite_test.go` file. This file is responsible for:
    *   Initializing `gomock` and creating mocks for all required repositories (`UsersRepo`, `CompaniesRepo`, etc.).
    *   Setting up the mock chain (e.g., `mockGlobalRepo.EXPECT().Users().Return(mockUsersRepo)`).
    *   Providing a standardized helper function, `newAuthenticatedRequest`, to generate `http.Request` objects for tests.

*   **Mock Injection via Context**: The test request helper function must inject the mocked `GlobalRepo` into the request's context using the correct key (`middleware.RepoKey`). For authenticated requests, it should also inject the user's ID into the context using `middleware.UserIDKey`. This is critical for simulating the application's dependency injection and authentication patterns.
    ```go
    // Example from a test suite helper function, simulating the result of AuthMiddleware
    ctxWithRepo := context.WithValue(req.Context(), middleware.RepoKey, mockGlobalRepo)
    if user != nil {
        ctxWithAuth := context.WithValue(ctxWithRepo, middleware.AuthUserKey, user)
        return req.WithContext(ctxWithAuth)
    }
    return req.WithContext(ctxWithRepo) // For unauthenticated requests
    ```

*   **User Role Setup**: For tests involving authentication and authorization, define `adminUser` and `normalUser` (or similar role-based user types) within the `_suite_test.go` file. This centralizes user setup and ensures consistent role-based testing across all handler tests.

*   **Thorough Test Cases**: Each handler should be tested for multiple scenarios using separate `Context` blocks:
    *   The "happy path" (a valid request that succeeds).
    *   Invalid request bodies (malformed JSON, missing required fields, validation errors like mismatched passwords).
    *   Dependency failures (e.g., trying to create a user for a `company_id` that doesn't exist).
    *   Repository errors (simulating database failures).
    *   **Authorization**: For endpoints with ownership rules, test that non-admins are forbidden (`403`) from accessing other tenants' data. For admin-only endpoints, test that both non-admins (`403`) and unauthenticated users (`401`) are rejected.

*   **JSON Payload Casing**: Pay close attention to JSON key casing in *all* JSON payloads (request and response). This project uses **`snake_case`**. Mismatched casing will cause validation to fail and result in `400 Bad Request` errors.

### Error Handling Conventions

*   **Standardized Error Responses**: Always use `middleware.WriteError(w, http.Status, message)` to send error responses to ensure consistency in format and status codes across the API.

*   **Invalid Path Parameters**: When testing endpoints with path parameters that expect a specific type (e.g., integer IDs), ensure tests cover scenarios where invalid formats are provided (e.g., non-numeric strings for an integer ID). Such cases typically result in `404 Not Found` responses due to route matching failures, rather than `400 Bad Request` from handler-level validation.
*   **Repository Error Wrapping**: Repository errors should be wrapped in user-friendly messages in the handler. Do not leak raw database errors (e.g., 'duplicate key value...') to the client. Instead, return a generic message like 'unable to update resource' or a specific, user-friendly one like 'resource with this name already exists'.

### Dependency Validation

*   **Validate Related Entities**: Before creating or updating entities that have foreign key relationships (e.g., a `User` requiring a `Company` and `Address`), always validate that the related entities exist. Return appropriate `400 Bad Request` or `404 Not Found` errors if dependencies are not met.

### Data Model & Business Logic

*   **Ownership Checks (Multi-tenancy)**: Handlers that create or modify company-scoped resources (e.g., `User`, `Location`) must verify that the acting user belongs to the correct company. While admins have broad permissions, they are also subject to company ownership rules for most resources. For example, an admin can only find locations within their own company.

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

### Data Model Documentation
When significant changes or additions are made to the core data models (e.g., `types` structs and their relationships), update the 'Data Model' section in `README.md` to reflect these changes. This ensures the high-level documentation stays in sync with the evolving database schema and type definitions.

*   **Populate All Required Fields**: When creating test data (e.g., `&types.User{...}`), ensure all required fields as defined by the struct and database schema are populated. Missing fields will cause validation or database errors during tests.

### A Note on Package Structure

- **`types` package**: The canonical location for all shared data transfer objects (structs like `User`, `Company`, etc.) is `github.com/happilymarrieddad/order-management-v3/api/types`. Please ensure all new code references this path directly to avoid import errors.

- **File Naming Convention**: Files within the `types` package should be named in their plural form (e.g., `users.go`, `products.go`), even if they contain a singular struct definition (e.g., `type User struct {}`).

### Custom Type Mapping with XORM

When a Go type does not map directly to a primitive database type (e.g., mapping a `[]Role` slice to a PostgreSQL `text[]` array), we use XORM's `Conversion` interface.

To teach XORM how to handle a custom type, the type must implement two methods:

1.  **`FromDB(data []byte) error`**: This method is called by XORM when reading from the database. It receives the raw data as a byte slice from the database driver and is responsible for parsing this data and populating the Go type.

    *Example*: Parsing `"{admin,user}"` into a `[]Role` slice.

2.  **`ToDB() ([]byte, error)`**: This method is called by XORM when writing to the database. It is responsible for serializing the Go type into a byte slice that the database can understand.

    *Example*: Converting a `[]Role` slice into the string `"{admin,user}"` and returning it as `[]byte`.

A complete example of this pattern can be found in `types/roles.go`. Please follow this convention for any new types that require custom database serialization.

## Interaction Guidelines

*   **Tone and Style**: Adopt a conversational style akin to J.A.R.V.I.S. from Iron Man.
*   **No colloquialisms**: Avoid informal phrases like 'My Bad'.
*   **Test Execution**: Do not ask to run tests after making changes. Assume tests will be run by the user or as part of a separate verification step.

## Lessons Learned

*   **Windows `make` commands**: When running `make` commands on Windows, `make` might not be directly available. Instead, use `powershell.exe -File scripts/migrate.ps1 <command>` for migration-related tasks, or execute the underlying `goose` commands directly with all necessary arguments (e.g., `goose -dir db/migrations <command>`).
*   **`goose create` command**: Always use the `-dir db/migrations` flag when creating new migration files (e.g., `goose -dir db/migrations create <name> sql`) to ensure they are placed in the correct directory and picked up by the migration tool.
*   **XORM `extends` and JSON tag conflicts**: When using XORM's `extends` tag to embed structs, be aware of potential JSON tag conflicts if both embedded structs have fields with the same JSON tag (e.g., `json:"id"`). Resolve this by adding `json:"-"` to the conflicting field in the embedded struct within the composite struct (e.g., `types.Address `xorm:"extends" json:"-"`). This tells the JSON marshaller to ignore that specific field during serialization.
*   **Testing `unknown` / Zero-Value Enum Constants**: When creating API endpoints that return lists of enum values (e.g., `AllRoles()`, `AllCommodityTypes()`), ensure that `unknown` or zero-value constants are *not* included in the returned list unless they represent a valid, selectable option for the user. Update tests to reflect this expectation, as these values are often internal representations and not meant for external consumption.
*   **Enum Constant Usage**: When using enum constants (like `types.CommodityType`), always verify the available constants by checking the corresponding `types` package file (e.g., `types/commodity_types.go`). Do not assume the existence of a constant without explicit definition.
*   **Variable Shadowing in Tests**: Be cautious with `:=` vs. `=` in `BeforeEach` blocks. Using `:=` can shadow package-level mock variables, causing them to be `nil` in tests and leading to panics.
*   **Mocking Method Signatures**: Ensure `gomock`'s `.Return()` calls match the exact number and type of return values for the mocked method signature to avoid 'wrong number of arguments' panics.
*   **Validation Error Messages**: Be aware that the validation library may generate error messages based on the struct's `json` tag or a lowercased version of the field name. Assertions in tests must match this exact format.
*   **Test Refactoring for Authentication**: Using `adminUser` and `normalUser` variables in test suites (`_suite_test.go`) to clearly separate and manage authenticated user roles in tests. This improves clarity and consistency.
*   **Route-Level Authorization**: Understand that routes can be restricted by middleware (e.g., `adminRouter.Use(middleware.AuthUserAdminRequiredMuxMiddleware())`). Tests should reflect this by expecting `403 Forbidden` for unauthorized access, rather than attempting to test business logic that will not be reached.
*   **`utils.TRef` vs `utils.Ref`**: Ensure correct usage of utility functions for creating pointers to primitive types (e.g., `utils.TRef` if available and intended for this purpose).