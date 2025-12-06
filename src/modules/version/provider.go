package version

import (
	"runtime"
	"time"

	"github.com/Glitch-guy0/authService/src/modules/core"
	"github.com/Glitch-guy0/authService/src/modules/logger"
)

var (
	version   = "dev"      // Set during build
	commit    = "none"     // Set during build
	buildTime = "unknown"  // Set during build
	buildUser = "unknown"  // Set during build
	buildHost = "unknown"  // Set during build
	dirty     = false      // Set during build
	tags      = []string{} // Set during build
	branch    = "main"     // Set during build
	goMod     = "unknown"  // Set during build
)

type VersionProvider struct {
	appCtx  *core.AppContext
	logger  logger.Logger
	version VersionInfo
}

func NewVersionProvider(appCtx *core.AppContext) *VersionProvider {
	vp := &VersionProvider{
		appCtx:  appCtx,
		logger:  appCtx.GetLogger(),
		version: buildVersionInfo(),
	}

	vp.logger.Info("Version provider initialized",
		"version", vp.version.Version.Version,
		"commit", vp.version.Version.Commit,
		"build_time", vp.version.Version.BuildTime,
	)

	return vp
}

func (vp *VersionProvider) GetVersion() VersionInfo {
	return vp.version
}

func (vp *VersionProvider) GetVersionString() string {
	return vp.version.Version.Version
}

func (vp *VersionProvider) GetCommit() string {
	return vp.version.Version.Commit
}

func (vp *VersionProvider) GetBuildTime() time.Time {
	return vp.version.Version.BuildTime
}

func (vp *VersionProvider) IsDirty() bool {
	return vp.version.Version.Dirty
}

func (vp *VersionProvider) GetTags() []string {
	return vp.version.Version.Tags
}

func (vp *VersionProvider) GetBuildInfo() BuildInfo {
	return vp.version.Build
}

func buildVersionInfo() VersionInfo {
	return VersionInfo{
		Version: Version{
			Version:   version,
			Commit:    commit,
			BuildTime: parseBuildTime(buildTime),
			GoVersion: runtime.Version(),
			BuildUser: buildUser,
			BuildHost: buildHost,
			Dirty:     dirty,
			Tags:      tags,
		},
		Build: BuildInfo{
			BuildTime: parseBuildTime(buildTime),
			BuildUser: buildUser,
			BuildHost: buildHost,
			GoVersion: runtime.Version(),
			GitCommit: commit,
			GitBranch: branch,
			GitTag:    getTagFromTags(tags),
			GoMod:     goMod,
		},
		Environment: getEnvironment(),
		Features:    getFeatures(),
	}
}

func parseBuildTime(buildTimeStr string) time.Time {
	if buildTimeStr == "unknown" {
		return time.Now()
	}

	parsed, err := time.Parse(time.RFC3339, buildTimeStr)
	if err != nil {
		return time.Now()
	}
	return parsed
}

func getTagFromTags(tags []string) string {
	if len(tags) > 0 {
		return tags[0]
	}
	return "none"
}

func getEnvironment() string {
	// This could be read from environment variable in the future
	return "development"
}

func getFeatures() []string {
	// This could be configured based on build flags or environment
	return []string{
		"health-checks",
		"graceful-shutdown",
		"structured-logging",
	}
}
