package resource

import (
	"errors"

	"github.com/gobuffalo/genny/movinglater/attrs"
	"github.com/gobuffalo/meta"
)

type Options struct {
	App           meta.App `json:"app"`
	Name          string   `json:"name"`
	Model         string   `json:"model"`
	SkipMigration bool     `json:"skip_migration"`
	SkipModel     bool     `json:"skip_model"`
	SkipTemplates bool     `json:"skip_templates"`
	UseModel      bool     `json:"use_model"`
	// FilesPath     string     `json:"files_path"`
	// ActionsPath   string     `json:"actions_path"`
	Attrs attrs.Attrs `json:"props"`
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}
	if len(opts.Name) == 0 {
		return errors.New("you must provide a name")
	}
	if len(opts.Model) == 0 {
		opts.Model = opts.Name
	}
	return nil
}
