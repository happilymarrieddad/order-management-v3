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
    *   For public endpoints, ensure tests cover non-happy path scenarios such as unsupported HTTP methods, and verify appropriate error responses (e.g., 405 Method Not Allowed).
*   **Individual Test Files**: Each test file (e.g., `create_test.go`, `get_test.go`) should contain specific test cases for a handler, and may include its own helper functions (like `createRequest`, `executeRequest`) if those helpers are specific to that test file's context.

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

*   **Test Execution**: Do not ask to run tests after making changes. Assume tests will be run by the user or as part of a separate verification step.

## Lessons Learned

*   **Windows `make` commands**: When running `make` commands on Windows, `make` might not be directly available. Instead, use `powershell.exe -File scripts/migrate.ps1 <command>` for migration-related tasks, or execute the underlying `goose` commands directly with all necessary arguments (e.g., `goose -dir db/migrations <command>`).
*   **`goose create` command**: Always use the `-dir db/migrations` flag when creating new migration files (e.g., `goose -dir db/migrations create <name> sql`) to ensure they are placed in the correct directory and picked up by the migration tool.
*   **XORM `extends` and JSON tag conflicts**: When using XORM's `extends` tag to embed structs, be aware of potential JSON tag conflicts if both embedded structs have fields with the same JSON tag (e.g., `json:"id"`). Resolve this by adding `json:"-"` to the conflicting field in the embedded struct within the composite struct (e.g., `types.Address `xorm:"extends" json:"-"`). This tells the JSON marshaller to ignore that specific field during serialization.
*   **Testing `unknown` / Zero-Value Enum Constants**: When creating API endpoints that return lists of enum values (e.g., `AllRoles()`, `AllCommodityTypes()`), ensure that `unknown` or zero-value constants are *not* included in the returned list unless they represent a valid, selectable option for the user. Update tests to reflect this expectation, as these values are often internal representations and not meant for external consumption.
