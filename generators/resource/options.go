package resource

import (
	"errors"
	"os"
	"strings"

	"github.com/gobuffalo/buffalo/meta"
)

// Options for generating a new resource
type Options struct {
	App           meta.App
	Name          meta.Name `json:"name"`
	SkipMigration bool      `json:"skip_migration"`
	SkipModel     bool      `json:"skip_model"`
	MimeType      string    `json:"mime_type"`
	ModelName     meta.Name `json:"model_name"`
	FilesPath     string    `json:"files_path"`
	ActionsPath   string    `json:"actions_path"`
	Props         []Prop    `json:"props"`
	Args          []string  `json:"args"`
}

// NewOptions constructs new options for generating a resource
func NewOptions(modelName string, args ...string) (Options, error) {
	o := Options{
		MimeType: "HTML",
		Args:     args,
	}
	pwd, _ := os.Getwd()
	o.App = meta.New(pwd)

	if len(o.Args) > 0 {
		o.Name = meta.Name(o.Args[0])
		o.ModelName = meta.Name(o.Args[0])
	}
	o.Props = modelPropertiesFromArgs(o.Args)

	o.FilesPath = o.Name.PluralUnder()
	o.ActionsPath = o.FilesPath
	if strings.Contains(string(o.Name), "/") {
		parts := strings.Split(string(o.Name), "/")
		o.ModelName = meta.Name(parts[len(parts)-1])
		o.ActionsPath = strings.Join(parts, "_")
	}
	if modelName != "" {
		o.ModelName = meta.Name(modelName)
	}
	return o, o.Validate()
}

// Validate that the options have what you need to build a new resource
func (o Options) Validate() error {
	mt := o.MimeType
	if mt != "HTML" && mt != "JSON" && mt != "XML" {
		return errors.New("invalid resource type, you need to choose between \"html\", \"xml\" and \"json\"")
	}

	if len(o.Args) == 0 && o.ModelName == "" {
		return errors.New("you must specify a resource name")
	}
	return nil
}
