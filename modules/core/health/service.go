package health

import (
	"context"
	"sync"
	"time"

	"github.com/Glitch-guy0/authService/modules/core"
	"github.com/Glitch-guy0/authService/modules/core/logger"
	"github.com/Glitch-guy0/authService/modules/version"
)

type HealthService struct {
	appCtx          *core.AppContext
	logger          logger.Logger
	config          HealthCheckConfig
	startTime       time.Time
	checkers        []HealthChecker
	versionProvider *version.VersionProvider
	mu              sync.RWMutex
}

func NewHealthService(appCtx *core.AppContext, config HealthCheckConfig) *HealthService {
	return &HealthService{
		appCtx:          appCtx,
		logger:          appCtx.GetLogger(),
		config:          config,
		startTime:       time.Now(),
		checkers:        make([]HealthChecker, 0),
		versionProvider: version.NewVersionProvider(appCtx),
	}
}

func (hs *HealthService) RegisterChecker(checker HealthChecker) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	hs.checkers = append(hs.checkers, checker)
	hs.logger.Info("Health checker registered", "name", checker.Name())
}

func (hs *HealthService) GetHealth(ctx context.Context) HealthResponse {
	overallStatus := StatusHealthy
	hs.mu.RLock()
	checkerCount := len(hs.checkers)
	hs.mu.RUnlock()
	checks := make([]Check, 0, checkerCount)

	hs.mu.RLock()
	for _, checker := range hs.checkers {
		checkStart := time.Now()
		check := checker.Check(ctx)
		check.Duration = time.Since(checkStart).String()
		checks = append(checks, check)

		if check.Status == StatusUnhealthy {
			overallStatus = StatusUnhealthy
		} else if check.Status == StatusDegraded && overallStatus == StatusHealthy {
			overallStatus = StatusDegraded
		}
	}
	hs.mu.RUnlock()

	return HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   hs.getVersionInfo(),
		Checks:    checks,
		Uptime:    time.Since(hs.startTime).String(),
	}
}

func (hs *HealthService) getVersionInfo() VersionInfo {
	versionInfo := hs.versionProvider.GetVersion()
	return VersionInfo{
		Version:   versionInfo.Version.Version,
		Commit:    versionInfo.Version.Commit,
		BuildTime: versionInfo.Version.BuildTime.Format(time.RFC3339),
		GoVersion: versionInfo.Version.GoVersion,
	}
}

type BasicHealthChecker struct {
	name    string
	checker func(ctx context.Context) Check
}

func (bhc *BasicHealthChecker) Name() string {
	return bhc.name
}

func (bhc *BasicHealthChecker) Check(ctx context.Context) Check {
	return bhc.checker(ctx)
}

func NewBasicHealthChecker(name string, checker func(ctx context.Context) Check) HealthChecker {
	return &BasicHealthChecker{
		name:    name,
		checker: checker,
	}
}
