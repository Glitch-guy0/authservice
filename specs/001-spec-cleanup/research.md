# Phase 0 Research – Specification Cleanup

## Decision 1: Consolidate health capability under modules/api/health
- **Rationale**: Keeps controller + service logic co-located per constitution §II, eliminates drift between duplicate implementations, enables future adapters (e.g., readiness hooks) through the existing `HealthService` abstraction.
- **Alternatives Considered**:
  1. **Keep both core/health.go and modules/api/health/** – rejected because it violates module-centric rule and increases maintenance burden.
  2. **Move health into modules/core/health immediately** – deferred to follow-up spec once controller contracts are documented; current priority is removing duplication.

## Decision 2: Document and enforce controller export contract
- **Rationale**: Specification demands each module expose route-capable controllers that register through AppContext-powered bootstrap. Documenting the pattern prevents future deviations and clarifies expectation that controllers remain stateless and DI-driven.
- **Alternatives Considered**:
  1. **Implicit knowledge sharing via code reviews** – unreliable for onboarding and does not scale.
  2. **Centralized router file wiring every module manually** – breaks plug-and-play requirement from constitution §§II–III.

## Decision 3: Migration safeguards for dependent imports
- **Rationale**: Removing `core/health.go` requires verifying no other packages import it. Plan mandates: (a) repository-wide search for `core/health` imports, (b) add compile-time guard ensuring only `modules/api/health` is referenced, and (c) expand tests (handler + readiness/liveness) to ensure parity.
- **Alternatives Considered**:
  1. **Soft deprecation warning before removal** – slower feedback; team prefers atomic cleanup with CI coverage.
  2. **Shim forwarding file under core/** – deemed unnecessary once imports are updated and would invite future misuse.

## Decision 4: Controller pattern quickstart contents
- **Rationale**: Quickstart will teach authors to (1) create `<domain>.controller` folders, (2) expose `Routes(appCtx *core.AppContext) *gin.RouterGroup` factories, and (3) register through `cmd/app/main.go`. This addresses the user request for clarity around route capability hand-off.
- **Alternatives Considered**:
  1. **Rely on constitution alone** – too abstract; developers need module-specific example.
  2. **Add ADR instead of quickstart** – ADR would memorialize reasoning but not provide actionable steps for contributors.
