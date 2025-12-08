# Tasks: Specification Cleanup

**Input**: Design documents from `/specs/001-spec-cleanup/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

## Phase 1: Discovery & Migration Safeguards

- [X] T001 [US1] Search codebase for `modules/api/health` imports and document all callers
- [X] T002 [US1] Search codebase for `core/health` imports and document all callers
- [X] T003 [US1] Verify no external packages reference `modules/api/health` via grep
- [X] T004 [US1] Run `go test ./modules/api/health/...` to establish baseline coverage

## Phase 2: Health Module Migration to Core

- [X] T005 [US1] Move `modules/api/health/` to `modules/core/health/` with all subdirectories
- [X] T006 [US1] Update all import statements referencing `modules/api/health` to `modules/core/health`
- [X] T007 [US1] Delete empty `modules/api/health/` directory after migration
- [X] T008 [US1] Verify `go build ./...` succeeds after health module migration

## Phase 3: Route Refactoring to Controller Pattern

- [X] T009 [US1] Extract health routes from `modules/server/server.go` to `modules/core/health/controller.go`
- [X] T010 [US1] Create health controller factory function `Routes(appCtx *core.AppContext, router *gin.RouterGroup)` in `modules/core/health/controller.go`
- [X] T011 [US1] Update `modules/server/server.go` to call health controller factory instead of static routes
- [X] T012 [US1] Verify `go build ./...` succeeds after route refactoring

## Phase 4: Module Migration to Core Structure

- [x] T013 [US1] Move `modules/config/` to `modules/core/config/` with all subdirectories
- [x] T014 [US1] Move `modules/logger/` to `modules/core/logger/` with all subdirectories
- [x] T015 [US1] Update all import statements referencing `modules/config` to `modules/core/config`
- [x] T016 [US1] Update all import statements referencing `modules/logger` to `modules/core/logger`
- [x] T017 [US1] Verify `go build ./...` succeeds after module migrations

## Phase 5: Controller Pattern Documentation

- [ ] T018 [US2] Review updated `modules/core/health/controller.go` structure
- [ ] T019 [US2] Document controller export contract in quickstart.md with new health module location
- [ ] T020 [US2] Verify `cmd/app/main.go` follows documented pattern for module registration

## Phase 6: Validation & Testing

- [ ] T021 [US1] Run health endpoints manually: `/health`, `/health/live`, `/health/ready`
- [ ] T022 [US1] Execute full test suite: `go test ./...`
- [ ] T023 [US1] Run linting: `golangci-lint run`

## Phase 7: Documentation Updates

- [ ] T024 [US2] Finalize quickstart.md with controller pattern examples
- [ ] T025 [US2] Update README.md if it references old module structure
- [ ] T026 [US2] Verify all generated docs are consistent with plan.md
