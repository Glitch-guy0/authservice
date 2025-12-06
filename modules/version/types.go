package version

import "time"

// Version represents application version information
type Version struct {
	Version   string    `json:"version"`
	Commit    string    `json:"commit,omitempty"`
	BuildTime time.Time `json:"build_time"`
	GoVersion string    `json:"go_version"`
	BuildUser string    `json:"build_user,omitempty"`
	BuildHost string    `json:"build_host,omitempty"`
	Dirty     bool      `json:"dirty,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
}

// BuildInfo contains build-time information
type BuildInfo struct {
	BuildTime time.Time `json:"build_time"`
	BuildUser string    `json:"build_user,omitempty"`
	BuildHost string    `json:"build_host,omitempty"`
	GoVersion string    `json:"go_version"`
	GitCommit string    `json:"git_commit"`
	GitBranch string    `json:"git_branch"`
	GitTag    string    `json:"git_tag"`
	GoMod     string    `json:"go_mod"`
}

// VersionInfo represents detailed version information
type VersionInfo struct {
	Version     Version   `json:"version"`
	Build       BuildInfo `json:"build"`
	Environment string    `json:"environment,omitempty"`
	Features    []string  `json:"features,omitempty"`
}

// DefaultVersion returns default version information
func DefaultVersion() Version {
	return Version{
		Version:   "dev",
		Commit:    "none",
		BuildTime: time.Now(),
		GoVersion: "1.21+",
	}
}

// DefaultBuildInfo returns default build information
func DefaultBuildInfo() BuildInfo {
	return BuildInfo{
		BuildTime: time.Now(),
		GoVersion: "1.21+",
		GitCommit: "none",
		GitBranch: "main",
		GitTag:    "none",
		GoMod:     "unknown",
	}
}

// DefaultVersionInfo returns default version information
func DefaultVersionInfo() VersionInfo {
	return VersionInfo{
		Version:     DefaultVersion(),
		Build:       DefaultBuildInfo(),
		Environment: "development",
		Features:    []string{},
	}
}
