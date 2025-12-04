# Implementation Plan: Project Setup (Logger + AppContext)

**Branch**: `001-project-setup` | **Date**: 2025-04-12 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This implementation plan outlines the setup of a Go-based web service with structured logging, configuration management, and a clean architecture. The project will use Gin as the web framework, logrus for structured logging, and viper for configuration management. The architecture will follow Go best practices with clear separation of concerns and dependency injection.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.21+  
**Primary Dependencies**: 
- Gin v1.9.0+ (Web Framework)
- logrus v1.9.0+ (Structured Logging)
- viper v1.15.0+ (Configuration)
- testify v1.8.2+ (Testing)
- go-playground/validator v10.14.0+ (Request Validation)

**Storage**: N/A (Initial setup only)  
**Testing**: 
- Standard library `testing` package
- Testify for assertions and mocks
- 80%+ test coverage target

**Target Platform**: Linux/amd64 (Docker container)  
**Project Type**: Web API Service  
**Performance Goals**: 
- Response time < 100ms p99
- Handle 1000+ RPS per instance
- Startup time < 2s

**Constraints**: 
- Memory: < 100MB per instance
- No global state
- Zero-downtime deployments

**Scale/Scope**: 
- Initial setup for future features
- Support for 10k+ users
- Module-based architecture for easy scaling

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

[Gates determined based on constitution file]

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
# Standard Go project layout with module support
cmd/
└── app/
    └── main.go           # Application entry point

internal/
├── api/
│   ├── middleware/       # HTTP middleware
│   └── response/         # Response formatters
├── config/              # Configuration loading
├── logger/              # Logging implementation
└── modules/             # Feature modules
    └── example/         # Example module
        ├── controller/   # HTTP handlers
        └── service/      # Business logic

pkg/
└── errors/              # Custom error types

test/
├── integration/         # Integration tests
└── unit/                # Unit tests

configs/                 # Configuration files
├── config.yaml          # Default configuration
└── config.dev.yaml      # Development overrides

docs/
└── api/                 # Requestly API documentation

scripts/                 # Build and deployment scripts

# [REMOVE IF UNUSED] Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# [REMOVE IF UNUSED] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**Structure Decision**: The project follows the standard Go project layout with a clear separation of concerns. Key aspects:
- `cmd/app` contains the application entry point
- `internal` holds all application-specific code
- `pkg` contains reusable packages
- `configs` stores configuration files
- `test` separates unit and integration tests
- `docs` contains API documentation for Requestly

## Phase 0: Research Tasks

1. **Logging Implementation**
   - Research best practices for structured logging with logrus
   - Define log levels and output formats
   - Plan for log rotation and retention

2. **Configuration Management**
   - Research viper configuration patterns
   - Define environment variable naming conventions
   - Plan for configuration validation

3. **API Design**
   - Research Gin middleware patterns
   - Define response format standards including:
     - Standard success response format
     - Error response format with error codes
     - Health check response schema
   - Plan for error handling and validation

## Phase 1: Implementation Tasks

1. **Project Setup**
   - Initialize Go module
   - Set up project structure
   - Configure build and test pipelines

2. **Core Components**
   - Implement logger package
   - Set up configuration management
   - Create base HTTP server with Gin

3. **API Foundation**

#### Health Check Endpoint

**Purpose**: Verify the service is running and responsive

**Implementation Details**:
- Simple endpoint that returns service status
- Includes version information for deployment verification
- Can be extended later with dependency health checks

**Example Request**:
```http
GET /health
Accept: application/json
```

**Success Response**:
```http
HTTP/1.1 200 OK
Content-Type: application/json
{
  "status": "ok",
  "timestamp": "2025-04-12T17:42:00Z",
  "version": "1.0.0"
}
```

**Error Response**:
```http
HTTP/1.1 503 Service Unavailable
Content-Type: application/json
{
  "status": "error",
  "error": "Service Unavailable",
  "message": "Unable to connect to database",
  "timestamp": "2025-04-12T17:42:00Z"
}
```
   - Implement request validation
   - Set up error handling middleware
   - Create health check endpoint with the following specifications:
     - **Endpoint**: `GET /health`
     - **Response**: 
       ```json
       {
         "status": "ok",
         "timestamp": "2025-04-12T17:42:00Z",
         "version": "1.0.0"
       }
       ```
     - **Status Codes**:
       - 200: Service is healthy
       - 503: Service is unhealthy (with error details)
   - Include basic service information in the health response

## Phase 2: Testing & Documentation

1. **Testing**
   - Write unit tests for core components
   - Add integration tests for API endpoints including:
     - Health check endpoint returns 200 when service is healthy
     - Health check response includes required fields
     - Error conditions are properly handled
   - Set up code coverage reporting

2. **Documentation**
   - Create API documentation for Requestly
   - Write developer setup guide
   - Document deployment process

## Success Metrics

- [ ] 100% test coverage for core packages
- [ ] All API endpoints documented in Requestly
- [ ] Build and test pipelines passing
- [ ] Code review completed

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
