package version

import (
	"fmt"
)

// Version information
var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildTime = "unknown"
)

// GetVersion returns the version string
func GetVersion() string {
	return fmt.Sprintf("lnk version %s (commit: %s, built: %s)", Version, CommitSHA, BuildTime)
}
