package resource

import (
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/genny/movinglater/attrs"
)

// Options for creating a new resource
type Options struct {
	App           meta.App    `json:"app"`
	Name          flect.Ident `json:"name"`
	Model         flect.Ident `json:"model"`
	SkipMigration bool        `json:"skip_migration"`
	SkipModel     bool        `json:"skip_model"`
	SkipTemplates bool        `json:"skip_templates"`
	UseModel      bool        `json:"use_model"`
	Attrs         attrs.Attrs `json:"attrs"`
}

// Validate options
func (opts *Options) Validate() error {
	if (opts.App == meta.App{}) {
		opts.App = meta.New(".")
	}
	return nil
}
