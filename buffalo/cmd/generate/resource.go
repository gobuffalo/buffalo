package generate

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo/generators/resource"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/makr"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

const resourceExamples = `$ buffalo g resource users
Generates:

- actions/users.go
- actions/users_test.go
- models/user.go
- models/user_test.go
- migrations/2016020216301234_create_users.up.fizz
- migrations/2016020216301234_create_users.down.fizz

$ buffalo g resource users --skip-migration
Generates:

- actions/users.go
- actions/users_test.go
- models/user.go
- models/user_test.go

$ buffalo g resource users --skip-model
Generates:

- actions/users.go
- actions/users_test.go

$ buffalo g resource users --use-model
Generates:

- actions/users.go
- actions/users_test.go`

// SkipResourceMigration allows to generate a resource without the migration.
var SkipResourceMigration = false

// SkipResourceModel allows to generate a resource without the model and Migration.
var SkipResourceModel = false

// UseResourceModel allows to generate a resource with a working model.
var UseResourceModel = ""

// ResourceMimeType allows to generate a typed resource (HTML by default, JSON...).
var ResourceMimeType = "html"

// ResourceCmd generates a new actions/resource file and a stub test.
var ResourceCmd = &cobra.Command{
	Use:     "resource [name]",
	Example: resourceExamples,
	Aliases: []string{"r"},
	Short:   "Generates a new actions/resource file",
	RunE: func(cmd *cobra.Command, args []string) error {
		var name, modelName, folderPath string

		if len(args) == 0 && UseResourceModel == "" {
			return errors.New("you must specify a resource name")
		}

		name = inflect.Pluralize(args[0])

		// Allow overwriting modelName with the --use-model flag
		// buffalo generate resource users --use-model people
		if UseResourceModel != "" {
			modelName = inflect.Pluralize(UseResourceModel)
			name = UseResourceModel
		}

		if modelName == "" {
			modelName = name
		}

		if ResourceMimeType != "html" && ResourceMimeType != "json" && ResourceMimeType != "xml" {
			return errors.New("invalid resource type, you need to choose between \"html\", \"xml\" and \"json\"")
		}

		folderPath = name

		if strings.Contains(name, "/") {
			parts := strings.Split(name, "/")
			name = parts[len(parts)-1]
			modelName = strings.Join(parts, "_")
		}

		modelProps := modelPropertiesFromArgs(args)

		data := makr.Data{
			"name":          name,
			"singular":      inflect.Singularize(name),
			"camel":         inflect.Camelize(name),
			"under":         inflect.Underscore(name),
			"underSingular": inflect.Singularize(inflect.Underscore(name)),
			"model":         inflect.Singularize(inflect.Camelize(modelName)),
			"modelPlural":   inflect.Camelize(modelName),

			"varPlural":   inflect.CamelizeDownFirst(modelName),
			"varSingular": inflect.Singularize(inflect.CamelizeDownFirst(modelName)),

			"renderFunction": strings.ToUpper(ResourceMimeType),
			"actions":        []string{"List", "Show", "New", "Create", "Edit", "Update", "Destroy"},
			"args":           args,
			"modelProps":     modelProps,
			"modelsPath":     packagePath() + "/models",

			"modelPluralUnder": inflect.Underscore(modelName),
			"folderPath":       folderPath,

			// Flags
			"skipMigration": SkipResourceMigration,
			"skipModel":     SkipResourceModel,
			"useModel":      UseResourceModel,
			"mimeType":      ResourceMimeType,
		}
		g, err := resource.New(data)
		if err != nil {
			return err
		}
		return g.Run(".", data)
	},
}

type modelProp struct {
	Name string
	Type string
}

func (m modelProp) String() string {
	return m.Name
}

func modelPropertiesFromArgs(args []string) []modelProp {
	var mProps []modelProp
	if len(args) == 0 {
		return mProps
	}
	for _, a := range args[1:] {
		ax := strings.Split(a, ":")
		p := modelProp{
			Name: inflect.ForeignKeyToAttribute(ax[0]),
			Type: "string",
		}
		if len(ax) > 1 {
			p.Type = strings.ToLower(strings.TrimPrefix(ax[1], "nulls."))
		}
		mProps = append(mProps, p)
	}
	return mProps
}

func goPath(root string) string {
	gpMultiple := envy.GoPaths()
	path := ""

	for i := 0; i < len(gpMultiple); i++ {
		if strings.HasPrefix(root, filepath.Join(gpMultiple[i], "src")) {
			path = gpMultiple[i]
			break
		}
	}
	return path
}

func packagePath() string {
	rootPath, _ := os.Getwd()
	gosrcpath := strings.Replace(filepath.Join(goPath(rootPath), "src"), "\\", "/", -1)
	rootPath = strings.Replace(rootPath, "\\", "/", -1)
	return strings.Replace(rootPath, gosrcpath+"/", "", 2)
}
