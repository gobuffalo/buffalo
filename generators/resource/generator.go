package resource

import (
	"errors"
	"os"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/meta"
)

// Generator for generating a new resource
type Generator struct {
	App           meta.App   `json:"app"`
	Name          name.Ident `json:"name"`
	Model         name.Ident `json:"model"`
	SkipMigration bool       `json:"skip_migration"`
	SkipModel     bool       `json:"skip_model"`
	SkipTemplates bool       `json:"skip_templates"`
	UseModel      bool       `json:"use_model"`
	FilesPath     string     `json:"files_path"`
	ActionsPath   string     `json:"actions_path"`
	Props         []Prop     `json:"props"`
	Args          []string   `json:"args"`
}

// New constructs new options for generating a resource
func New(modelName string, args ...string) (Generator, error) {
	o := Generator{
		Args: args,
	}
	pwd, _ := os.Getwd()
	o.App = meta.New(pwd)

	if len(o.Args) > 0 {
		o.Name = name.New(flect.Singularize(o.Args[0]))
		o.Model = o.Name
	}
	o.Props = modelPropertiesFromArgs(o.Args)

	o.FilesPath = o.Name.File().Pluralize().String()
	o.ActionsPath = o.FilesPath
	if strings.Contains(o.Name.String(), "/") {
		parts := strings.Split(o.Name.String(), "/")
		o.Model = name.New(parts[len(parts)-1])
		o.ActionsPath = o.Name.Resource().Underscore().String()
	}
	if modelName != "" {
		o.Model = name.New(modelName)
	}
	return o, o.Validate()
}

// Validate that the options have what you need to build a new resource
func (o Generator) Validate() error {
	if len(o.Args) == 0 && o.Model.String() == "" {
		return errors.New("you must specify a resource name")
	}
	return nil
}
