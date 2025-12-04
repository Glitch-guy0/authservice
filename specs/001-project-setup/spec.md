# Feature Specification: Project Setup (Logger + AppContext)

**Feature Branch**: `001-project-setup`  
**Created**: 2025-04-12  
**Status**: Draft  
**Input**: User description: "Setup project structure and tooling. Implement createLogger() util. Implement buildAppContext() with ctx.logger. Central place for all dependencies from day one. Ensures every later milestone uses the same logger + context flow."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Project Structure Setup (Priority: P1)

As a developer, I want to have a well-organized Go project structure following Go best practices so that I can start building features efficiently.

**Why this priority**: A solid foundation is crucial for maintainability and scalability of the codebase.

**Independent Test**: Can be verified by checking if the project structure follows the defined pattern and all Go tooling is properly configured.

**Acceptance Scenarios**:

1. **Given** a new Go project, **When** I initialize the module, **Then** I should see a standard Go project structure with `cmd/`, `internal/`, and `pkg/` directories.
2. **Given** the project structure, **When** I check the root directory, **Then** I should find `go.mod` and `go.sum` files with the project's module name and dependencies.

---

### User Story 2 - Logger Implementation (Priority: P1)

As a developer, I want a centralized logging interface with a `sirupsen/logrus` implementation so that I can consistently log application events and errors.

**Why this priority**: Logging is essential for debugging and monitoring application behavior from day one.

**Independent Test**: Can be tested by importing and using the logger in different parts of the application.

**Acceptance Scenarios**:

1. **Given** the logger package, **When** I call `logger.Create()`, **Then** I should receive a logger instance that implements a standard logging interface.
2. **Given** a logger instance, **When** I log messages at different levels (Info, Warn, Error, Debug), **Then** they should be properly formatted with timestamps and output to the console in a structured format (JSON).

---

### User Story 3 - Application Context (Priority: P1)

As a developer, I want a centralized application context that follows Go's best practices for dependency management so that I can manage and share dependencies across the application without passing the context everywhere.

**Why this priority**: Proper dependency management is crucial for maintainable and testable code.

**Independent Test**: Can be tested by creating and using the application context in different modules.

**Acceptance Scenarios**:

1. **Given** the application, **When** I initialize a new service, **Then** I should be able to inject dependencies through a constructor that takes only what it needs.
2. **Given** a service, **When** it requires a dependency, **Then** it should declare that dependency in its constructor parameters rather than accessing a global context.

---

### User Story 4 - Gin Web Server Setup (Priority: P1)

As a developer, I want to set up a Gin web server with proper routing and middleware so that I can start building API endpoints.

**Why this priority**: The web server is a core component that needs to be set up early in the project.

**Independent Test**: Can be verified by starting the server and making HTTP requests to the endpoints.

**Acceptance Scenarios**:

1. **Given** the Gin server, **When** I start the application, **Then** it should start on the configured port without errors.
2. **Given** a request to a non-existent endpoint, **When** the request is made, **Then** it should return a 404 Not Found response.
3. **Given** a malformed request, **When** the request is validated, **Then** it should return a 400 Bad Request with validation errors.

---

### User Story 5 - Module-based Routing (Priority: P2)

As a developer, I want to organize routes by modules with their own controllers so that the codebase remains maintainable as it grows.

**Why this priority**: Good organization from the start prevents technical debt.

**Independent Test**: Can be verified by checking the module structure and route registration.

**Acceptance Scenarios**:

1. **Given** a new module, **When** I add a controller, **Then** it should be easy to register its routes with the main router.
2. **Given** a module's controller, **When** I make a request to its endpoint, **Then** it should be properly handled by the controller.

---

### User Story 6 - Go Development Workflow (Priority: P2)

As a developer, I want to have a standard Go development workflow with proper tooling so that I can maintain code quality and consistency.

**Why this priority**: A good development workflow improves productivity and reduces errors.

**Independent Test**: Can be verified by running standard Go tools and checking the output.

**Acceptance Scenarios**:

1. **Given** the Go project, **When** I run `go test ./...`, **Then** it should execute all tests in the project.
2. **Given** the project setup, **When** I run `go mod tidy`, **Then** it should clean up and verify the module's dependencies.

### Edge Cases

- What happens when the logger is not properly initialized?
- How does the system handle circular dependencies in the application context?
- What happens when configuration files are missing or malformed?
- How does the system handle concurrent access to the logger instance?
- What happens when an invalid API version is requested?
- How does the system handle malformed JSON in request bodies?
- What happens when a required request parameter is missing?
- How does the system handle invalid configuration values?
- What happens when a required environment variable is not set?
- How does the system handle database connection failures?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a `logger` package that exports a `Logger` interface and a `NewLogger()` function that returns an implementation using `sirupsen/logrus`.
- **FR-009**: System MUST use Gin as the HTTP web framework.
- **FR-010**: System MUST use `go-playground/validator` for request validation.
- **FR-011**: System MUST implement URL path versioning (e.g., `/api/v1/...`).
- **FR-012**: System MUST organize routes by modules with each module having its own controller.
- **FR-002**: System MUST use constructor-based dependency injection for all services and components.
- **FR-003**: The project MUST follow standard Go project layout with `cmd/`, `internal/`, and `pkg/` directories.
- **FR-004**: System MUST include `go.mod` and `go.sum` files with all necessary dependencies.
- **FR-005**: System MUST include a `Makefile` with common development tasks (test, lint, build, etc.).
- **FR-006**: System MUST include a `.gitignore` file with standard Go ignores.
- **FR-007**: System MUST include a basic `README.md` with setup and development instructions.
- **FR-008**: System MUST use Go 1.21 or later (specified in `go.mod`).
- **FR-013**: System MUST use `viper` for configuration management, supporting both environment variables and configuration files.
- **FR-014**: System MUST implement custom error types with appropriate HTTP status codes.
- **FR-015**: System MUST use a standardized JSON response format for all API endpoints.
- **FR-016**: System SHOULD include placeholder comments for logging context (to be implemented later).
- **FR-017**: System MUST be testable with Requestly for API testing and documentation.

### Key Entities

- **Logger**: Interface defining standard logging methods (Info, Warn, Error, Debug) implemented using `sirupsen/logrus`.
- **Gin Engine**: The main Gin router that handles HTTP requests and middleware.
- **Controllers**: Structs that handle HTTP requests and responses, organized by module.
- **Services**: Go structs that implement business logic, with dependencies injected through constructors.
- **Configuration**: Managed by `viper`, loading from environment variables and config files.
- **Error Types**: Custom error types that include HTTP status codes and structured error details.
- **API Response**: Standardized JSON response format for all endpoints.
- **ProjectStructure**: Standard Go project layout with:
  - `cmd/` for main applications
  - `internal/` for private application code
  - `internal/api/` for HTTP handlers and middleware
  - `internal/modules/` for feature modules
  - `pkg/` for reusable packages
  - `configs/` for configuration files and defaults
  - `docs/` for API documentation (Requestly format)

## Implementation Notes

### Logging Context
```go
// TODO: Implement context-based logging with request IDs
// Context should be passed through the request chain
// and include relevant request-scoped information
```

### API Documentation with Requestly
- API documentation will be maintained in Requestly
- Requestly will be used for testing and documenting API endpoints
- Request collections will be versioned alongside the code

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The Go project structure follows standard Go project layout (100% compliance).
- **SC-002**: The logger implementation is complete and passes all tests (100% test coverage for logger package).
- **SC-003**: All services use constructor-based dependency injection (verified by code review).
- **SC-004**: The project builds successfully with `go build ./...` and all tests pass with `go test ./...`.
- **SC-006**: The Gin server starts without errors and handles requests on the configured port.
- **SC-007**: Request validation works as expected with `go-playground/validator`.
- **SC-008**: Module-based routing is implemented with controllers in their respective modules.
- **SC-009**: Configuration is properly loaded using `viper` with environment variable support.
- **SC-010**: All API endpoints return responses in the standardized JSON format.
- **SC-011**: Custom error types are used consistently throughout the application.
- **SC-012**: The application can be tested using Requestly for API testing.
- **SC-005**: The `Makefile` includes all necessary targets for development workflow (test, lint, build, etc.).
