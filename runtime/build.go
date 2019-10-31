package runtime

import (
	"fmt"
	"sync"
	"time"
)

// BuildInfo holds information about the build
type BuildInfo struct {
	Version string    `json:"version"`
	Time    time.Time `json:"-"`
}

// String implements fmt.String
func (b BuildInfo) String() string {
	return fmt.Sprintf("%s (%s)", b.Version, b.Time)
}

var build = BuildInfo{
	Version: "",
	Time:    time.Time{},
}

// Build returns the information about the current build
// of the application. In development mode this will almost
// always run zero values for BuildInfo.
func Build() BuildInfo {
	return build
}

var so sync.Once

// SetBuild allows the setting of build information only once.
// This is typically managed by the binary built by `buffalo build`.
func SetBuild(b BuildInfo) {
	so.Do(func() {
		build = b
	})
}
