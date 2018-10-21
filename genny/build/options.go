package build

import (
	"sync"
	"time"

	"github.com/gobuffalo/meta"
)

type Options struct {
	meta.App
	BuildVersion       string
	BuildTime          time.Time
	ExtractAssets      bool
	WithAssets         bool
	LDFlags            string
	Tags               meta.BuildTags
	BuildFlags         []string
	Static             bool
	Environment        string
	TemplateValidators []TemplateValidator
	rollback           *sync.Map
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
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
	return nil
}
