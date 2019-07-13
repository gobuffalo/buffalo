package resource

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/meta"
)

// Options for generating a new resource
type Options struct {
	App           meta.App    `json:"app"`
	Name          string      `json:"name"`
	Model         string      `json:"model"`
	SkipMigration bool        `json:"skip_migration"`
	SkipModel     bool        `json:"skip_model"`
	SkipTemplates bool        `json:"skip_templates"`
	Attrs         attrs.Attrs `json:"props"`
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}

	if len(opts.Name) == 0 {
		return fmt.Errorf("you must provide a name")
	}

	if len(opts.Model) == 0 {
		opts.Model = opts.Name
	}

	if strings.Contains(opts.Model, "/") {
		parts := strings.Split(opts.Model, "/")
		opts.Model = parts[len(parts)-1]
	}

	if opts.App.AsAPI {
		opts.SkipTemplates = true
	}

	return nil
}
