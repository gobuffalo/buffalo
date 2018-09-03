package resource

import (
	"errors"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/genny/movinglater/attrs"
)

// Options for creating a new resource
type Options struct {
	App           meta.App         `json:"app"`
	Attrs         attrs.NamedAttrs `json:"attrs"`
	SkipMigration bool             `json:"skip_migration"`
	SkipModel     bool             `json:"skip_model"`
	SkipTemplates bool             `json:"skip_templates"`
	UseModel      bool             `json:"use_model"`
}

// Validate options
func (opts *Options) Validate() error {
	if (opts.App == meta.App{}) {
		opts.App = meta.New(".")
	}
	if opts.Attrs.Name.String() == "" {
		return errors.New("you must give your resource a name")
	}
	return nil
}
