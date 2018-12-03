package resource

import (
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/meta"
)

type Options struct {
	App           meta.App   `json:"app"`
	Name          name.Ident `json:"name"`
	Model         name.Ident `json:"model"`
	SkipMigration bool       `json:"skip_migration"`
	SkipModel     bool       `json:"skip_model"`
	SkipTemplates bool       `json:"skip_templates"`
	UseModel      bool       `json:"use_model"`
	FilesPath     string     `json:"files_path"`
	ActionsPath   string     `json:"actions_path"`
	// Props         []Prop     `json:"props"`
	Args []string `json:"args"`
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}
	return nil
}
