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

$ buffalo g resource users --use-model users
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

// ModelName allows to specify a different model name for the resource.
var ModelName = ""

// ResourceCmd generates a new actions/resource file and a stub test.
var ResourceCmd = &cobra.Command{
	Use:     "resource [name]",
	Example: resourceExamples,
	Aliases: []string{"r"},
	Short:   "Generates a new actions/resource file",
	RunE: func(cmd *cobra.Command, args []string) error {
		var name, modelName, resourceName, filesPath, actionsPath string

		//Check for a valid mime type
		if ResourceMimeType != "html" && ResourceMimeType != "json" && ResourceMimeType != "xml" {
			return errors.New("invalid resource type, you need to choose between \"html\", \"xml\" and \"json\"")
		}

		if len(args) == 0 && UseResourceModel == "" {
			return errors.New("you must specify a resource name")
		}

		name = inflect.Pluralize(args[0])
		modelName = name
		filesPath = name
		actionsPath = name
		resourceName = name

		if strings.Contains(name, "/") {
			parts := strings.Split(name, "/")
			name = parts[len(parts)-1]
			modelName = name

			resourceName = strings.Join(parts, "_")
			actionsPath = resourceName
		}

		// Allow overwriting modelName with the --use-model flag
		// buffalo generate resource users --use-model people
		if UseResourceModel != "" {
			modelName = inflect.Pluralize(UseResourceModel)
			name = UseResourceModel
		}

		if ModelName != "" {
			modelName = inflect.Pluralize(ModelName)
			name = ModelName
		}

		modelProps := modelPropertiesFromArgs(args)

		data := makr.Data{
			"name":     name,
			"singular": inflect.Singularize(name),
			"camel":    inflect.Camelize(name),
			"under":    inflect.Underscore(name),

			"renderFunction": strings.ToUpper(ResourceMimeType),
			"actions":        []string{"List", "Show", "New", "Create", "Edit", "Update", "Destroy"},
			"args":           args,

			"filesPath":   filesPath,
			"actionsPath": actionsPath,

			"model":              inflect.Singularize(inflect.Camelize(name)),
			"modelPlural":        inflect.Pluralize(inflect.Camelize(name)),
			"modelPluralUnder":   inflect.Underscore(modelName),
			"modelFilename":      inflect.Underscore(inflect.Camelize(name)),
			"modelTable":         inflect.Underscore(inflect.Pluralize(name)),
			"modelSingularUnder": inflect.Underscore(inflect.Singularize(name)),
			"modelProps":         modelProps,
			"modelsPath":         packagePath() + "/models",

			"resourceName":          inflect.Camelize(resourceName),
			"resourcePlural":        inflect.Pluralize(inflect.Camelize(resourceName)),
			"resourceURL":           inflect.Pluralize(inflect.Underscore(filesPath)),
			"resourceSingularUnder": inflect.Underscore(inflect.Singularize(resourceName)),

			"routeName":              inflect.Camelize(resourceName),
			"routeNameSingular":      inflect.Camelize(inflect.Singularize(resourceName)),
			"routeFirstDown":         inflect.CamelizeDownFirst(resourceName),
			"routeFirstDownSingular": inflect.CamelizeDownFirst(inflect.Singularize(resourceName)),

			"varPlural":   inflect.CamelizeDownFirst(modelName),
			"varSingular": inflect.Singularize(inflect.CamelizeDownFirst(modelName)),

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
