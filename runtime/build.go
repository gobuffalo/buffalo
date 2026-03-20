// Deprecated: This package is deprecated and will be removed in a future version.
// Use runtime/debug.ReadBuildInfo() directly instead.
//
// This package previously provided build information for buffalo applications.
// Starting with Go 1.18, the standard library's runtime/debug package provides
// equivalent functionality through ReadBuildInfo(), which includes VCS information
// (commit hash, time) and module version.
//
// Migration example:
//
//	import "runtime/debug"
//
//	info, ok := debug.ReadBuildInfo()
//	if ok {
//	    // Use info.Main.Version for version
//	    // Use info.Settings for VCS info (vcs.revision, vcs.time)
//	}
//
// For applications built with "buffalo build", build information is now
// automatically embedded by Go's build system and can be accessed via
// runtime/debug.ReadBuildInfo().
package runtime

import (
	"runtime/debug"
	"sync"
	"time"
)

// Version is the current version of the buffalo binary.
// Deprecated: Use runtime/debug.ReadBuildInfo().Main.Version instead.
var Version = "dev"

// BuildInfo holds information about the build.
// Deprecated: Use runtime/debug.BuildInfo instead.
type BuildInfo struct {
	Version string    `json:"version"`
	Time    time.Time `json:"-"`
}

// String implements fmt.Stringer
func (b BuildInfo) String() string {
	if b.Time.IsZero() {
		return b.Version
	}
	return b.Version + " (" + b.Time.Format(time.RFC3339) + ")"
}

var (
	build     BuildInfo
	buildOnce sync.Once
)

// Build returns the information about the current build of the application.
// In development mode this will almost always return zero values for BuildInfo.
//
// Deprecated: Use runtime/debug.ReadBuildInfo() instead. This function now
// returns information derived from the standard library's build info.
//
// For backward compatibility, this function caches the build info after first call.
func Build() BuildInfo {
	buildOnce.Do(func() {
		build = loadBuildInfo()
	})
	return build
}

// loadBuildInfo reads build information from runtime/debug and converts it
// to the legacy BuildInfo format for backward compatibility.
func loadBuildInfo() BuildInfo {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return BuildInfo{
			Version: "dev",
			Time:    time.Time{},
		}
	}

	bi := BuildInfo{
		Version: info.Main.Version,
		Time:    time.Time{},
	}

	// Handle development builds
	if bi.Version == "" || bi.Version == "(devel)" {
		bi.Version = "dev"
	}

	// Try to extract build time from VCS info
	for _, setting := range info.Settings {
		if setting.Key == "vcs.time" {
			if t, err := time.Parse(time.RFC3339, setting.Value); err == nil {
				bi.Time = t
			}
			break
		}
	}

	return bi
}

var so sync.Once

// SetBuild allows the setting of build information only once.
// This is typically managed by the binary built by `buffalo build`.
//
// Deprecated: This function is no longer necessary. Build information is now
// automatically embedded by the Go toolchain and can be accessed via
// runtime/debug.ReadBuildInfo(). This function is kept for backward compatibility
// but has no effect when called after Build() has been called.
func SetBuild(b BuildInfo) {
	so.Do(func() {
		build = b
	})
}
