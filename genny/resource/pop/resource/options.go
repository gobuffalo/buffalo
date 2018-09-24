package resource

import (
	"errors"
	"strings"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/flect/name"
)

// Options for creating a new resource
type Options struct {
	Name          name.Ident
	Model         name.Ident
	App           meta.App `json:"app"`
	Args          []string `json:"args"`
	Attrs         []Prop   `json:"attrs"`
	SkipMigration bool     `json:"skip_migration"`
	SkipModel     bool     `json:"skip_model"`
	SkipTemplates bool     `json:"skip_templates"`
	UseModel      bool     `json:"use_model"`
}

// Validate options
func (opts *Options) Validate() error {
	if (opts.App == meta.App{}) {
		opts.App = meta.New(".")
	}
	if opts.Name.String() == "" {
		if len(opts.Args) > 0 {
			opts.Name = name.New(opts.Args[0])
		}
	}

	if opts.Name.String() == "" {
		return errors.New("you must give your resource a name")
	}

	if len(opts.Model.String()) == 0 {
		opts.Model = opts.Name
	}

	if len(opts.Attrs) == 0 {
		var args []string
		if len(opts.Args) > 1 {
			args = opts.Args[0:]
		}
		opts.Attrs = modelPropertiesFromArgs(args)
	}
	return nil
}

// Prop of a model. Starts as name:type on the command line.
type Prop struct {
	Name name.Ident
	Type string
}

// String representation of Prop
func (m Prop) String() string {
	return m.Name.String()
}

func modelPropertiesFromArgs(args []string) []Prop {
	var props []Prop
	if len(args) == 0 {
		return props
	}
	for _, a := range args[1:] {
		ax := strings.Split(a, ":")
		p := Prop{
			Name: name.New(ax[0]),
			Type: "string",
		}
		if len(ax) > 1 {
			p.Type = strings.ToLower(strings.TrimPrefix(ax[1], "nulls."))
		}
		props = append(props, p)
	}
	return props
}
