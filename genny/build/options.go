package build

import (
	"os"
	"sync"
	"time"

	"github.com/gobuffalo/meta"
)

// Options for building a Buffalo application
type Options struct {
	meta.App
	// the "timestamp" of the build. defaults to time.Now()
	BuildTime time.Time `json:"build_time,omitempty"`
	// the "version" of the build. defaults to
	// a) git sha of last commit or
	// b) to time.RFC3339 of BuildTime
	BuildVersion string `json:"build_version,omitempty"`
	// CleanAssets will remove the public/assets folder build compiling
	CleanAssets bool `json:"clean_assets"`
	WithAssets  bool `json:"with_assets,omitempty"`
	// places ./public/assets into ./bin/assets.zip.
	// requires WithAssets = true
	ExtractAssets bool `json:"extract_assets,omitempty"`
	// LDFlags to be passed to the final `go build` command
	LDFlags string `json:"ld_flags,omitempty"`
	// Tags to be passed to the final `go build` command
	Tags meta.BuildTags `json:"tags,omitempty"`
	// BuildFlags to be passed to the final `go build` command
	BuildFlags []string `json:"build_flags,omitempty"`
	// Static sets the following flags for the final `go build` command:
	// -linkmode external
	// -extldflags "-static"
	Static bool `json:"static,omitempty"`
	// Environment the binary is meant for. defaults to "development"
	Environment string `json:"environment,omitempty"`
	// TemplateValidators can be used to validate the applications templates.
	// Empty by default
	TemplateValidators []TemplateValidator `json:"-"`
	// Mod is the -mod flag
	Mod string `json:"mod"`
	// GoCommand is the `go X` command to be used. Default is "build".
	GoCommand string `json:"go_command"`
	rollback  *sync.Map
	keep      *sync.Map
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	pwd, _ := os.Getwd()
	if opts.App.IsZero() {
		opts.App = meta.New(pwd)
	}
	if len(opts.Environment) == 0 {
		opts.Environment = "development"
	}
	if opts.BuildTime.IsZero() {
		opts.BuildTime = time.Now()
	}
	if len(opts.BuildVersion) == 0 {
		opts.BuildVersion = opts.BuildTime.Format(time.RFC3339)
	}
	if opts.rollback == nil {
		opts.rollback = &sync.Map{}
	}
	if opts.keep == nil {
		opts.keep = &sync.Map{}
	}
	if len(opts.GoCommand) == 0 {
		opts.GoCommand = "build"
	}
	return nil
}
