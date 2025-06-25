package version

import (
	"fmt"
)

// Version information
// These variables can be set during build time using ldflags
var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildTime = "unknown"
)

// GetVersion returns the version string
func GetVersion() string {
	return fmt.Sprintf("%s (commit: %s, built: %s)", Version, CommitSHA, BuildTime)
}
