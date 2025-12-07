# Implementation Plan: Specification Cleanup

**Branch**: `001-spec-cleanup` | **Date**: December 7, 2025 | **Spec**: [/specs/001-spec-cleanup/spec.md](specs/001-spec-cleanup/spec.md)
**Input**: Feature specification from `/specs/001-spec-cleanup/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

The Specification Cleanup effort consolidates the health module to a single implementation, aligns controller responsibilities with the constitution’s module-centric rules, and documents the controller pattern so future modules follow the same contract. We will remove `core/health.go`, migrate any callers to `modules/api/health`, and ensure controllers expose route capabilities that are cleanly injected into the server bootstrap. Supporting docs (research, data model, contracts, quickstart) will capture guidance for module authors moving forward.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.25.1 (`go.mod`)  
**Primary Dependencies**: Gin (HTTP server), Logrus-based logger wrapper, Viper config, Testify (tests), internal `core.AppContext` DI container  
**Storage**: N/A for health module (pure in-memory checks)  
**Testing**: `go test ./...` with Testify assertions + benchmark helpers in `/test`  
**Target Platform**: Linux containerized microservice (Kubernetes-ready)  
**Project Type**: Backend API service (single project under repo root)  
**Performance Goals**: Health endpoints must respond <50 ms p95 and never block server startup/shutdown  
**Constraints**: Must comply with constitution (module-centric layout, DI-only wiring, OTEL-ready logging), no duplicate implementations, controllers must stay stateless  
**Scale/Scope**: Single module refactor + documentation updates impacting `modules/api/health`, `modules/core`, and server bootstrap

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **G1 – Module-centric layout enforced (Constitution §II & §XV)**  
  Status: PASS. Plan removes `core/health.go`, keeps controllers inside `modules/<domain>.controller`, and documents the controller hand-off to the server.
- **G2 – Dependency Injection boundaries (Constitution §III & §IX)**  
  Status: PASS. All remaining health logic continues to use `core.AppContext`; any future adapters will be injected via DI, no new singletons introduced.
- **G3 – Observability + Testing (Constitution §§V, VII, XVII)**  
  Status: PASS. Health endpoints already emit structured responses; plan retains logging hooks and requires unit/integration coverage to stay green before merge.

## Project Structure

### Documentation (this feature)

```text
specs/001-spec-cleanup/
├── plan.md              # /speckit.plan output (this file)
├── spec.md              # Feature specification
├── research.md          # Phase 0 research deliverable
├── data-model.md        # Phase 1 entity outline
├── quickstart.md        # Phase 1 onboarding guide
├── contracts/           # Phase 1 API contracts (e.g., OpenAPI snippet)
└── tasks.md             # Phase 2 execution plan (/speckit.tasks)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
.
├── cmd/
│   └── app/main.go                # Service entrypoint wiring AppContext + Gin
├── configs/
│   └── config.yaml                # Baseline configuration (Viper loaded)
├── modules/
│   ├── api/
│   │   └── health/                # Controller + service to be consolidated
│   ├── config/                    # Config module (to be moved under modules/core per spec)
│   └── core/                      # Shared primitives (AppContext, logger integrations)
├── pkg/
│   └── errors/                    # Error envelope helpers/tests
├── test/
│   ├── benchmark/                 # Benchmark harnesses (health/config)
│   ├── helpers/                   # Shared test fixtures + assertions
│   └── integration/               # Placeholder for future DI/integration tests
└── specs/                         # Feature specs + planning artifacts
```

**Structure Decision**: Single backend project rooted at repo base. All domain capabilities live under `modules/<domain>/` following constitution rules; controller bindings are exported to `cmd/app/main.go`.

## Complexity Tracking

No constitution violations requiring justification have been identified for this plan.
