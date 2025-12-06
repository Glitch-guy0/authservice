<!--
Sync Impact Report
- Version change: 0.0.0 → 1.0.0 (Universal Backend Architecture)
- Modified/defined principles:
  - Core Principles → I. core-principles (security-first, explicitness-and-simplicity, modularity-and-plug-and-play, testability, observability-first, consistency)
- Added sections:
  - 0. scope-and-purpose
  - II. architecture-and-project-structure (module-centric)
  - III. dependency-injection-and-modularity
  - IV. design-patterns (mandatory)
  - V. test-driven-development (non-negotiable)
  - VI. idempotency-standard
  - VII. logging-tracing-metrics (opentelemetry)
  - VIII. error-envelope-standard
  - VIII-A. application-exception-model (module-scoped custom exceptions)
  - IX. configuration-and-environment-rules + app-context-pattern
  - X. rate-limiting-and-abuse-protection
  - XI. message-broker-and-event-system
  - XII. outbox-pattern
  - XIII. request-validation-and-sanitization
  - XIV. requestly-centric-workflow
  - XV. controller-and-transport-rules
  - XVI. service-lifecycle
  - XVII. observability-and-slos
  - XVIII. development-practices
  - XIX. naming-and-structure-conventions
  - XX. governance-and-change-control
  - XXI. parser-and-automation-hints
  - XXII. open-or-deferred-decisions
- Removed sections:
  - Generic Core Principles 1-5
  - Generic additional sections 2 and 3
- Templates requiring updates:
  - ✅ .specify/templates/spec-template.md (no backend-structure coupling; remains valid)
  - ⚠ .specify/templates/plan-template.md (update project structure examples to module-centric backend layout)
  - ⚠ .specify/templates/tasks-template.md (update path conventions and examples to module-centric backend layout and testing strategy)
  - ⚠ Commands templates: N/A (directory not present; add/update later if introduced)
- Deferred TODOs:
  - TODO(RATIFICATION_DATE): Original ratification date unknown; set when first adopted
  - TODO(DEFAULT_BROKER_SELECTION): Choose default message broker per environment
  - TODO(ERROR_CODE_TAXONOMY): Define unified error code taxonomy
  - TODO(EVENT_NAMING): Define global event naming scheme
  - TODO(RATE_LIMIT_DEFAULTS): Define global rate-limit defaults
  - TODO(IDEMPOTENCY_TTL): Decide default Idempotency-Key TTL
  - TODO(APP_CONTEXT_EXTENSION_RULES): Document how app-context may be safely extended
  - TODO(CANONICAL_SCHEMA_TOOLING): Standardize canonical schema → persistence tooling
-->

# Universal-Backend-Constitution

## 0. scope-and-purpose

- version: v4
- this constitution defines the universal architecture rules, conventions, and foundational patterns for all backend projects.
- this is NOT a service specification; it is a governing document for:
  - folder structure
  - design patterns
  - dependency injection (di)
  - module system
  - testing (tdd + contract testing)
  - logging, tracing, metrics (opentelemetry)
  - rate limiting
  - idempotency
  - outbox + event brokers
  - error handling and envelope formats
  - code governance and versioning
- service-specific rules must be written in separate specification files.

## I. core-principles

- **security-first**
  - confidentiality, integrity, availability required.
  - no secrets may appear in logs or code.
  - security-sensitive defaults (least privilege, secure-by-default) MUST be applied.

- **explicitness-and-simplicity**
  - avoid hidden automagic behavior.
  - code must be predictable, deterministic, and auditable.
  - configuration and behavior must be explicit and discoverable from code and docs.

- **modularity-and-plug-and-play**
  - all parts of the system must be replaceable via di container.
  - domain modules are self-contained units that can be added/removed without changing core bootstrapping logic.

- **testability**
  - all layers must be mockable (controllers, services, repositories, providers).
  - no hard-coded dependencies that bypass di or app-context.

- **observability-first**
  - logs, traces, and metrics are mandatory for all externally visible operations and infrastructure interactions.
  - all critical paths MUST be observable end-to-end.

- **consistency**
  - all projects MUST share the same structure and conventions defined here.
  - deviations require explicit justification and governance approval.

## II. architecture-and-project-structure (module-centric)

- all backend code MUST live at the project root.

- global folders (kebab-case only):
  - `modules/`
  - `utils/`
  - `helpers/`
  - `middleware/`
  - `persistence/`
  - `repository/`
  - `interface/`
  - `core/`
  - `utils/message-broker/`

- old structure removal:
  - `src/` → removed (Go projects use project root)
  - `service/` → removed
  - `controller/` → removed

- each domain module MUST follow this structure:

```text
modules/<domain>/
  <domain>.interface/      # dtos, validators, contracts
  <domain>.controller/     # endpoints (http/grpc/ws)
  <domain>.service/        # business logic
  <domain>.provider/       # third-party adapters, infra providers
  <domain>.repository/     # db access (repository pattern)
  <domain>.event-bus/      # domain-level event bus adapters (optional)
```

- module responsibilities:
  - modules encapsulate ALL domain logic.
  - modules produce their own events.
  - modules register di bindings and controllers.
  - modules remain plug-and-play: add/remove modules without core changes.

- module constraints:
  - no cross-module business logic; modules communicate via services or events.
  - controllers MUST NOT contain any business logic.
  - repositories MUST ONLY interact with persistence.

- routing inside modules:
  - http/grpc/ws route files live under `<domain>.controller/`.
  - nested routes map to nested folder structure.
  - all file/folder names MUST be kebab-case.

### persistence-and-repository-layers

- purpose:
  - separate **how data is stored** (persistence) from **how the domain accesses data** (repository).
  - enforce dependency flow: `service → repository → persistence → db/cache`.

- persistence (`persistence/`):
  - persistence is the **infrastructure-facing** representation of storage:
    - sql models / orm structs (e.g., gorm models in go).
    - nosql schemas / document models.
    - cache key/value schemas and helpers.
  - persistence MUST NOT contain any business logic.
  - persistence is derived from a **canonical schema/migration source** (tool-agnostic).
    - prisma support example:
      - prisma schema lives under `migrations/prisma/schema.prisma`.
      - a generator converts prisma → gorm structs under `persistence/model/`.
  - recommended folder structure:

```text
persistence/
  model/   # sql / relational orm models
  schema/  # nosql document schemas
  cache/   # cache key/value schemas
```

#### canonical-schema-and-migration-folder

- there MUST be a single canonical schema source for the project:
  - `migrations/prisma/` OR
  - `migrations/sql/`.
- persistence artifacts MUST be generated or derived from this canonical schema.
- manual edits SHOULD be minimized and, when necessary, documented.

#### repository layer (`repository/` or `<domain>.repository/`)

- repositories provide **domain-centric access** to persistence.
- repositories depend on:
  - persistence artifacts (models/schemas).
  - db/cache clients from app-context.
- repositories MUST NOT:
  - contain business rules.
  - use raw db drivers directly bypassing persistence.

- recommended structure:

```text
repository/
  user/
    user-repository.*
  payment/
    payment-repository.*
```

- or module-scoped:

```text
modules/user/user.repository/
modules/payment/payment.repository/
```

- typical usage:
  - `Repository.User.createUser()`
  - `Repository.Payment.makePayment()`
  - `Repository.Payment.Details.getUserDetailsByPaymentID()`

- persistence usage:
  - `Persistence.Model.User.create()`
  - `Persistence.Schema.User.storeMetadata()`
  - `Persistence.Cache.SessionData.storeUserSession()`

- MUST follow dependency rule:
  - `Service → Repository → Persistence → DB/Cache`.

## III. dependency-injection-and-modularity

- a centralized di container MUST be used.

- di container responsibilities:
  - register core providers (db, cache, logger, tracer, metrics, broker).
  - register module providers.
  - route registration.
  - lifecycle hooks.

- no global mutable singletons outside the di container.

- patterns:
  - constructor-injection only.
  - no service-locators.

- modules must self-register:

```text
modules/<domain>/<domain>.module.*
exports:
  - providers
  - services
  - controllers
  - repositories
  - event-bus adapters
```

- di dependency rules:
  - controller → service → repository → persistence.
  - utils/providers may be injected anywhere.

## IV. design-patterns (mandatory)

- repository-pattern
- provider-pattern
- factory-pattern
- strategy-pattern
- decorator-pattern
- manager-pattern
- domain-driven modularization
- event-bus pattern
- plug-and-play infra patterns

## V. test-driven-development (non-negotiable)

- tests MUST mirror project structure at root.
- unit tests MUST mock all external interactions.
- integration tests MUST use ephemeral postgres/redis/brokers.
- contract tests via requestly/pact (or equivalent) are required for external interfaces.
- no untested business logic is allowed in main branches.

## VI. idempotency-standard

- required for POST/PUT/PATCH.
- header: `Idempotency-Key`.
- store at least:
  - key, hash of request, final response, ttl.
- controllers retrieve and pass key to idempotency middleware / service.
- business logic MUST be retry-safe.

## VII. logging-tracing-metrics (opentelemetry)

- opentelemetry (otel) is mandatory for:
  - logging
  - tracing
  - metrics

- logs:
  - MUST be structured and in json format.
  - MUST NOT contain secrets or pii beyond approved fields.

- tracing:
  - all external calls MUST be traced.
  - trace ids MUST propagate across services where supported.

- metrics:
  - MUST include throughput, latency, error rate, health, and event rate.

## VIII. error-envelope-standard

- responses SHOULD follow this envelope:

```json
{
  "success": true,
  "data": {"...": "..."},
  "error": null
}
```

or, on error:

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "string",
    "message": "string",
    "detail": {"...": "..."},
    "requestId": "string"
  }
}
```

- controllers MUST wrap errors in this envelope.
- http status codes MUST be mapped from error metadata.

## VIII-A. application-exception-model (module-scoped custom exceptions)

- overview:
  - application uses **typed**, **structured**, **module-aware** exceptions.
  - exceptions propagate up the stack until the controller boundary.
  - controllers convert exceptions into the error-envelope (section VIII).

- base application exception (`AppException`):
  - defined under: `core/errors/app-exception.*`.
  - MUST contain at least:
    - `message: string`.
    - `statusCode: number`.
    - `details: any` (structured) with optional `origin` field.

  - recommended structure:

```typescript
type AppException = {
  message: string
  statusCode: number
  details?: {
    origin?: string
    [key: string]: any
  }
}
```

- module-specific exceptions:
  - each module MUST define its own exception family, e.g.:
    - `Exception.User.*`.
    - `Exception.Payment.*`.
    - `Exception.Auth.*`.
  - each module MAY define a module-base exception:
    - `UserException extends AppException`.
    - `PaymentException extends AppException`.

  - folder structure example:

```text
modules/user/user.exception/
  user-exception-base.*
  user-creation-failed-exception.*
  user-not-found-exception.*
```

- propagation & wrapping rules:
  - modules catching an internal exception MUST either:
    - rethrow it if appropriate; or
    - wrap it into a module-specific exception.
  - wrappers MUST copy:
    - `message`, `statusCode`, `details`.
    - MUST set: `details.origin = "<module-name>"`.

- example:

```typescript
catch (err) {
  throw new UserException.UserCreationFailed({
    message: err.message,
    statusCode: 500,
    details: {
      ...err.details,
      origin: "user-module"
    }
  })
}
```

- matching & handling:
  - exceptions MUST be matched by **type**, not by message strings.

- integration with error-envelope:
  - `AppException.statusCode` → http status.
  - `AppException.message` → `error.message`.
  - `AppException.details` → `error.detail` (including `origin`).
  - unknown errors MUST be wrapped into a safe generic
    `CoreException.UnexpectedError`.

## IX. configuration-and-environment-rules + app-context-pattern

- all env vars MUST be validated on startup.
- no business logic may execute at import time.

### app-context-pattern

- single source of truth: `core/app-context.*`.
- app-context is immutable after initialization.
- app-context contains references to:
  - db
  - cache
  - logger
  - tracer
  - meter
  - broker
  - config
- modules only obtain infrastructure via app-context.

## X. rate-limiting-and-abuse-protection

- platform MUST support ip/user/route limits.
- recommended: redis sliding window implementation.
- abuse events SHOULD be emitted to the event system.

## XI. message-broker-and-event-system

- universal provider lives under `utils/message-broker/`.
- kafka/rabbitmq/sqs/nats (or similar) are supported via provider adapters.
- modules publish only via the domain event bus, not directly to broker
  clients.

## XII. outbox-pattern

- outbox table MUST be transactionally aligned with domain writes.
- a worker or change-data-capture (cdc) mechanism is responsible for delivery.

## XIII. request-validation-and-sanitization

- sanitize html/scripts/control chars from user input where applicable.
- validate email/ids/enums and all critical fields.
- controllers perform pre-validation before invoking services.

## XIV. requestly-centric-workflow

- requestly (or equivalent) is used for:
  - api docs
  - contract tests
  - flows and collections.
- ci/cd pipelines MUST run requestly-based tests where defined.

## XV. controller-and-transport-rules

- controllers do:
  - parsing and validation.
  - dto mapping.
  - service calls.
- controllers MUST NOT contain business logic.

## XVI. service-lifecycle

- startup sequence:
  - load config.
  - initialize app-context.
  - connect to infrastructure (db/cache/broker).
  - run migrations.
  - register modules.
  - start server.

- graceful shutdown:
  - drain requests.
  - flush broker/outbox.
  - close db/cache connections.

## XVII. observability-and-slos

- services MUST expose metrics and health endpoints.
- track at least p95 and p99 latency.
- propagate request-id across logs, traces, and envelopes.

## XVIII. development-practices

- no side-effects on import.
- use dedicated utils for external apis & services.
- create dedicated dto files when a domain has more than 3 dtos.

## XIX. naming-and-structure-conventions

- kebab-case for all files and folders.
- names MUST be descriptive and reflect purpose.

## XX. governance-and-change-control

- semantic versioning applies to this constitution:
  - MAJOR: backward-incompatible governance/principle removals or redefinitions.
  - MINOR: new principle/section added or materially expanded guidance.
  - PATCH: clarifications, wording, typo fixes, non-semantic refinements.

- constitution changes require PR + maintainer approval.
- a constitution-check process MUST validate at least:
  - naming and structure.
  - otel integration.
  - tdd adherence.
  - app-context usage.

## XXI. parser-and-automation-hints

- automation validates:
  - folder structure.
  - naming conventions.
  - module di registration.
  - controllers only inside modules.
  - requestly tests linked where required.
  - persistence–repository dependency order.
  - app-context completeness.

## XXII. open-or-deferred-decisions

- default broker selection.
- unified error code taxonomy.
- event naming scheme.
- rate-limit defaults.
- idempotency-key ttl default.
- app-context extension rules.
- canonical schema → persistence tooling standard.

## Governance

- this constitution supersedes ad-hoc backend practices within this repository.
- all new backend work MUST comply with this constitution.
- amendments require:
  - a documented proposal referencing current version.
  - semantic version bump as per section XX.
  - maintainer approval.

- compliance expectations:
  - all feature specs/plans/tasks SHOULD include a "Constitution Check" gate.
  - reviewers MUST verify structural, testing, and observability rules.
  - violations MUST be explicitly justified and tracked.

**Version**: 1.0.0 | **Ratified**: TODO(RATIFICATION_DATE) | **Last Amended**: 2025-12-04
