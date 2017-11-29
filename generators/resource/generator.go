package resource

import (
	"errors"
	"os"
	"strings"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/markbates/inflect"
)

// Generator for generating a new resource
type Generator struct {
	App           meta.App  `json:"app"`
	Name          meta.Name `json:"name"`
	Model         meta.Name `json:"model"`
	SkipMigration bool      `json:"skip_migration"`
	SkipModel     bool      `json:"skip_model"`
	UseModel      bool      `json:"use_model"`
	MimeType      string    `json:"mime_type"`
	FilesPath     string    `json:"files_path"`
	ActionsPath   string    `json:"actions_path"`
	Props         []Prop    `json:"props"`
	Args          []string  `json:"args"`
}

// New constructs new options for generating a resource
func New(modelName string, args ...string) (Generator, error) {
	o := Generator{
		MimeType: "HTML",
		Args:     args,
	}
	pwd, _ := os.Getwd()
	o.App = meta.New(pwd)

	if len(o.Args) > 0 {
		o.Name = meta.Name(o.Args[0])
		o.Model = meta.Name(o.Args[0])
	}
	o.Props = modelPropertiesFromArgs(o.Args)

	o.FilesPath = o.Name.PluralUnder()
	o.ActionsPath = o.FilesPath
	if strings.Contains(string(o.Name), "/") {
		parts := strings.Split(string(o.Name), "/")
		o.Model = meta.Name(parts[len(parts)-1])
		o.ActionsPath = inflect.Underscore(o.Name.Resource())
	}
	if modelName != "" {
		o.Model = meta.Name(modelName)
	}
	return o, o.Validate()
}

// Validate that the options have what you need to build a new resource
func (o Generator) Validate() error {
	mt := o.MimeType
	if mt != "HTML" && mt != "JSON" && mt != "XML" {
		return errors.New("invalid resource type, you need to choose between \"html\", \"xml\" and \"json\"")
	}

	if len(o.Args) == 0 && o.Model == "" {
		return errors.New("you must specify a resource name")
	}
	return nil
}
