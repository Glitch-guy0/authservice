# Quickstart – Controller Pattern & Health Cleanup

## 1. Prerequisites
- Go 1.25.1 toolchain installed
- `gin-gonic` and repo dependencies (`go mod download`)
- Access to `core.AppContext` initialization (see `cmd/app/main.go`)

## 2. Remove duplicate health implementation
1. Delete `core/health.go` once all imports move to `modules/api/health`.
2. Run `grep -R "core/health" -n .` and update any callers to use `modules/api/health`.
3. Re-run `go test ./modules/api/health/...` to confirm handler parity.

## 3. Module controller contract
Each module must expose a controller factory:
```go
package healthcontroller

func Routes(appCtx *core.AppContext, router *gin.RouterGroup) {
    handler := health.NewHealthHandler(appCtx)
    router.GET("/", handler.HealthCheck)
    router.GET("/live", handler.LivenessProbe)
    router.GET("/ready", handler.ReadinessProbe)
}
```
- Controllers stay stateless; all dependencies flow through `appCtx`.
- Use module-local folders (`modules/api/<domain>/<domain>.controller/`).

## 4. Registering controllers with the server
In `cmd/app/main.go`:
1. Initialize `appCtx := core.NewAppContext(...)`.
2. Create Gin router `r := gin.Default()`.
3. Mount health routes:
```go
healthRoutes := r.Group("/health")
healthcontroller.Routes(appCtx, healthRoutes)
```
4. Repeat pattern for new modules to keep DI/responsibility consistent.

## 5. Health checker extensions
- Default checkers cover server, database placeholder, logger.
- To add adapters (e.g., Postgres), create a struct satisfying `HealthChecker` and register via `HealthService.RegisterChecker` inside module bootstrap.

## 6. Verification checklist
1. `go test ./...` passes.
2. `curl localhost:8080/health` → returns status + version info.
3. `curl localhost:8080/health/ready` returns 200 or 503 depending on checker state.
4. `golangci-lint run` has zero issues.

## 7. Next Steps
- After cleanup, run `/speckit.tasks` to generate implementation backlog.
- Coordinate config module relocation under `modules/core` per spec follow-up.
