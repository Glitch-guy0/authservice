# Feature Specification: Specification Cleanup

**Feature Branch**: `001-spec-cleanup`  
**Created**: December 7, 2025  
**Status**: Draft  
**Input**: User description: "Specification cleanup - there are two files of same name /api/health and in core health.go which one is it. if you ask me then remove the core/health.go and for routes there is no controller file all route capability of a module comes from controller file in the module and this is passed to the server"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Remove Duplicate Health Implementation (Priority: P1)

As a developer, I want to eliminate duplicate health check implementations so that the codebase is cleaner and easier to maintain.

**Why this priority**: Having duplicate implementations creates confusion, maintenance overhead, and potential for inconsistencies between the two versions.

**Independent Test**: Can be verified by confirming that only one health implementation exists and that all health check functionality continues to work properly.

**Acceptance Scenarios**:

1. **Given** the system has duplicate health implementations in core/health.go and api/health/, **When** the core/health.go file is removed, **Then** the system should still function correctly with only the api/health/ implementation
2. **Given** health checks are called throughout the application, **When** the duplicate is removed, **Then** all health check calls should resolve to the remaining implementation without errors

---

### User Story 2 - Clarify Controller Structure (Priority: P2)

As a developer, I want to understand the controller pattern so that I can correctly implement new modules following the established architecture.

**Why this priority**: Clear documentation of the controller pattern ensures consistency across modules and reduces onboarding time for new developers.

**Independent Test**: Can be verified by examining the existing module structure and confirming that the controller pattern is well-defined and consistently applied.

**Acceptance Scenarios**:

1. **Given** the current module structure, **When** examining any module, **Then** there should be a clear controller file that handles all route capabilities for that module
2. **Given** the server setup, **When** routes are registered, **Then** each module's controller should be passed to the server for route registration

---

### Edge Cases

- What happens when other parts of the codebase import from the removed core/health.go?
- How does the system handle health check failures during the transition?

## Clarifications

### Session 2025-12-07

- Q: How should modules define their core business logic interfaces to ensure vendor-agnostic plug-and-play capability? → A: Define interfaces in each module's core package with adapter implementations for external dependencies
- Q: How should modules handle external dependencies like databases, message queues, or external APIs while maintaining vendor-agnostic design? → A: Use adapter pattern with interfaces defined in modules and concrete implementations in separate infrastructure packages
- Q: What should be the standard directory structure for each module to support hexagonal architecture? → A: Three-layer structure: core/ (business logic), adapters/ (external interfaces), infrastructure/ (vendor implementations)
- Q: How should the server discover and register modules while maintaining vendor-agnostic plug-and-play capability? → A: Use dependency injection container with module registration interfaces
- Q: How does dependency injection work with module dependencies? → A: Each module defines dependencies as interfaces, injected implementations must satisfy those interfaces
- Q: Which dependencies are fixed and cannot be changed? → A: Gin (server), Logger (base implementation), Config (Viper), AppContext (concept remains same)
- Q: How do modules interact with AppContext during initialization? → A: New modules take AppContext and pass it to submodules during initialization
- Q: How should the health module be refactored to follow hexagonal architecture? → A: Extract health check logic to core package, move infrastructure-specific implementations to infrastructure package
- Q: Should the health and config modules be moved to modules/core/health and modules/core/config respectively, with all import statements updated to reflect this new core module structure? → A: Move to modules/core/health and modules/core/config, update all imports

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST eliminate duplicate health check implementations
- **FR-002**: System MUST preserve all existing health check functionality after cleanup
- **FR-003**: System MUST maintain hexagonal architecture with vendor-agnostic modules
- **FR-004**: System MUST ensure each module follows three-layer structure: core/, adapters/, infrastructure/
- **FR-005**: System MUST define interfaces in each module's core package with adapter implementations
- **FR-006**: System MUST ensure module controllers are properly passed to the server for route registration
- **FR-007**: System MUST not break any existing functionality during the cleanup process
- **FR-008**: Health module MUST be refactored to separate core logic from infrastructure concerns
- **FR-009**: Health module MUST be moved to modules/core/health with updated import statements
- **FR-010**: Config module MUST be moved to modules/core/config with updated import statements

### Architectural Requirements

- **AR-001**: All modules MUST be vendor-agnostic except for fixed core dependencies (Gin, Logger, Viper, AppContext)
- **AR-002**: Modules MUST use adapter pattern for external dependencies with interface-based DI
- **AR-003**: Business logic MUST be isolated in core/ packages without external dependencies
- **AR-004**: Infrastructure implementations MUST be in separate packages from core business logic
- **AR-005**: Modules MUST provide plug-and-play capability for microservice foundation
- **AR-006**: Server MUST use dependency injection container for module discovery and registration
- **AR-007**: Controllers MUST use Gin as the web server framework
- **AR-008**: Logger base implementation will be extended to OTEL (future work)
- **AR-009**: Config MUST use Viper for configuration management
- **AR-010**: AppContext concept remains unchanged; modules receive and pass to submodules

### Key Entities *(include if feature involves data)*

- **Health Module**: The module responsible for health check endpoints and monitoring
- **Hexagonal Architecture**: The architectural pattern ensuring vendor-agnostic, plug-and-play modules
- **Module Structure**: Three-layer directory organization (core/, adapters/, infrastructure/)
- **Adapter Pattern**: Design pattern for handling external dependencies while maintaining vendor agnosticity
- **Core Dependencies**: Fixed vendor dependencies (Gin for server, Logger base implementation, Viper for config, AppContext pattern)
- **Dependency Injection Container**: System for module discovery and registration with plug-and-play capability

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Codebase has zero duplicate health implementations after cleanup
- **SC-002**: All health endpoints continue to function correctly without any regressions
- **SC-003**: Developer documentation clearly explains hexagonal architecture and module structure
- **SC-004**: Build and test suite passes without any failures after the cleanup
- **SC-005**: Code review shows consistent three-layer module structure (core/, adapters/, infrastructure/)
- **SC-006**: All modules demonstrate vendor-agnostic design with proper adapter pattern implementation
- **SC-007**: Health module refactored to follow hexagonal architecture principles
- **SC-008**: Server implements dependency injection container for module registration
- **SC-009**: Health module demonstrates proper separation of core logic and infrastructure
- **SC-010**: Health module successfully moved to modules/core/health with all imports updated
- **SC-011**: Config module successfully moved to modules/core/config with all imports updated
- **SC-012**: All tests pass after module restructuring and import updates
