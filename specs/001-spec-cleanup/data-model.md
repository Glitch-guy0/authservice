# Phase 1 Data Model – Specification Cleanup

## Entity Inventory

| Entity | Purpose | Key Fields | Relationships |
|--------|---------|------------|---------------|
| `HealthModule` | Encapsulates all health-related capabilities exposed to the platform. | `name`, `controllers`, `services`, `checkers` | Owns `HealthController`, `HealthService`, and `HealthChecker` registrations.
| `HealthController` | Gin controller responsible for routing (`/health`, `/health/live`, `/health/ready`). Delegates to service layer only. | `routerGroup`, `handler` methods | Injects `HealthService` (DI); exports route factory consumed by `cmd/app/main.go`.
| `HealthService` | Aggregates `HealthChecker`s, calculates uptime/version metadata, returns `HealthResponse`. | `appCtx`, `logger`, `config`, `startTime`, `checkers`, `versionProvider` | Uses `HealthChecker` instances; produces `HealthResponse` consumed by controllers.
| `HealthChecker` | Interface representing a health probe against infrastructure or logical subsystems. | `Name() string`, `Check(ctx context.Context) Check` | Registered with `HealthService`; may depend on adapters (e.g., DB, cache) via AppContext.
| `Check` | Result of a single probe. | `name`, `status` (`healthy`, `degraded`, `unhealthy`), `message`, `timestamp`, `duration`, `metadata` | Aggregated into `HealthResponse.checks`.
| `HealthResponse` | API contract returned to clients. | `status`, `timestamp`, `version`, `checks[]`, `uptime` | Serialized to JSON; consumed by readiness/liveness clients.
| `ControllerContract` | Documentation construct ensuring every module exposes `Routes(appCtx *core.AppContext) gin.HandlerFunc` (or router group). | `domain`, `routeGroup`, `dependencies` | Ensures server bootstrap can dynamically register modules.

## Field-Level Detail

### HealthController
- `routerGroup`: Mounted Gin group (`/health`). Must be created per module to keep DI-lifecycle isolated.
- `HealthHandler` methods:
  - `HealthCheck` → returns `HealthResponse` with aggregated checks.
  - `LivenessProbe` → lightweight OK response.
  - `ReadinessProbe` → returns `ready` or `not ready` depending on aggregate status.
- **Validation rules**: No business logic; must panic/return error if `HealthService` nil.

### HealthService
- `startTime` captured at construction to report `uptime`.
- `checkers` slice guarded by RWMutex; registration happens during bootstrap (default checkers) plus module-specific adapters.
- `config` ensures timeouts + severity thresholds stay configurable; defaults provided via `DefaultHealthCheckConfig()`.
- **Validation rules**: registering duplicate checker names is disallowed (logged + skipped). All checks must complete within configured timeout; degraded status returned if timeout reached.

### HealthChecker Implementations
- Must be stateless; any external clients resolved from AppContext.
- Required metadata per constitution:
  - `name` in kebab-case
  - `status` with allowed enum {`healthy`, `degraded`, `unhealthy`}
  - `message` with actionable detail.

### HealthResponse
- `status`: derived from the worst checker result (healthy > degraded > unhealthy).
- `version`: populated via `modules/version` provider; includes semantic version, commit SHA, build time, Go version.
- `checks`: ordered list mirroring registration order for predictable diffing.
- `uptime`: ISO8601 duration string; ensures monitoring parity after refactor.

## Relationships & Flows

```text
Gin Router (/health) → HealthController → HealthService → []HealthChecker → Check → HealthResponse → HTTP JSON
```

- Controller obtains injected `HealthService` from DI container.
- Service iterates registered `HealthChecker`s concurrently (future enhancement) or sequentially (current behavior) to produce `Check` structs.
- Server bootstrap wires controller route group into global router via exported factory; this is the “controller contract” referenced in the spec.

## Validation & Error Handling
- Controllers must return HTTP 200 for healthy/degraded, 503 for unhealthy, aligning with kubernetes probes.
- When removing `core/health.go`, ensure no code path bypasses controller/service abstractions; compile-time vetting via `go test ./...` + `golangci-lint`.
- Controller pattern documentation mandates envelope compliance if/when error responses introduced.

## State Transitions
- `HealthService` state: `Initialized` → `RegisteringCheckers` → `Ready`.
- Checkers themselves may track internal state (e.g., cached latency), but service treats them as pure functions each request.
- Removing `core/health.go` does not change runtime states; it only eliminates unused code paths.

## Open Questions
- None. Research resolved controller contract + duplication strategy.
