// Package version provides build-time version information for the HawkOp CLI.
package version

import (
	"fmt"
	"runtime"
	"strings"
)

// Build information set via ldflags during compilation
var (
	// Version is the current version of the application
	Version = "dev"
	// BuildTime is when the binary was built
	BuildTime = "unknown"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
)

// Info contains version and build information
type Info struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
	GitCommit string `json:"gitCommit"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
	Arch      string `json:"arch"`
}

// GetInfo returns detailed version information
func GetInfo() Info {
	return Info{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// GetShortVersion returns just the version string
func GetShortVersion() string {
	return Version
}

// GetDetailedVersion returns a formatted version string with build details
func GetDetailedVersion() string {
	var parts []string

	if Version != "" && Version != "dev" {
		parts = append(parts, fmt.Sprintf("HawkOp version %s", Version))
	} else {
		parts = append(parts, "HawkOp version dev")
	}

	if GitCommit != "" && GitCommit != "unknown" {
		if len(GitCommit) > 8 {
			parts = append(parts, fmt.Sprintf("commit %s", GitCommit[:8]))
		} else {
			parts = append(parts, fmt.Sprintf("commit %s", GitCommit))
		}
	}

	if BuildTime != "" && BuildTime != "unknown" {
		parts = append(parts, fmt.Sprintf("built %s", BuildTime))
	}

	parts = append(parts, fmt.Sprintf("go %s %s/%s",
		runtime.Version(), runtime.GOOS, runtime.GOARCH))

	return strings.Join(parts, ", ")
}
