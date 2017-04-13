package generate

import (
	"errors"
	"strings"

	"github.com/gobuffalo/buffalo/generators/resource"
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

// ResourceCmd generates a new actions/resource file and a stub test.
var ResourceCmd = &cobra.Command{
	Use:     "resource [name]",
	Example: resourceExamples,
	Aliases: []string{"r"},
	Short:   "Generates a new actions/resource file",
	RunE: func(cmd *cobra.Command, args []string) error {
		var name, modelName, filesPath string
		pkgName := "actions"

		// Allow overwriting modelName with the --use-model flag
		// buffalo generate resource users --use-model people
		if UseResourceModel != "" {
			modelName = inflect.Pluralize(UseResourceModel)
		}

		if len(args) == 0 {
			if UseResourceModel == "" {
				return errors.New("you must specify a resource name")
			}
			// When there is no resource name given and --use-model flag is set
			name = UseResourceModel
			filesPath = name
		} else {
			// When resource name is specified
			parts := strings.Split(args[0], "/")
			name = inflect.Pluralize(parts[len(parts)-1])

			if len(parts) > 1 {
				pkgName = inflect.Underscore(parts[len(parts)-2])
			}

			parts = append(parts[:len(parts)-1], name)
			filesPath = strings.Join(parts, "/")

			// If there is no --use-model flag set use the resource to create the model
			if modelName == "" {
				modelName = name
			}
		}

		modelProps := getModelPropertiesFromArgs(args)
		data := makr.Data{
			"package":          pkgName,
			"path":             inflect.Underscore(filesPath),
			"name":             name,
			"singular":         inflect.Singularize(name),
			"plural":           name,
			"camel":            inflect.Camelize(name),
			"under":            inflect.Underscore(name),
			"underSingular":    inflect.Singularize(inflect.Underscore(name)),
			"downFirstCap":     inflect.CamelizeDownFirst(name),
			"model":            inflect.Singularize(inflect.Camelize(modelName)),
			"modelPlural":      inflect.Camelize(modelName),
			"modelUnder":       inflect.Singularize(inflect.Underscore(modelName)),
			"modelPluralUnder": inflect.Underscore(modelName),
			"varPlural":        inflect.CamelizeDownFirst(modelName),
			"varSingular":      inflect.Singularize(inflect.CamelizeDownFirst(modelName)),
			"actions":          []string{"List", "Show", "New", "Create", "Edit", "Update", "Destroy"},
			"args":             args,
			"modelProps":       modelProps,

			// Flags
			"skipMigration": SkipResourceMigration,
			"skipModel":     SkipResourceModel,
			"useModel":      UseResourceModel,
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

func getModelPropertiesFromArgs(args []string) []modelProp {
	var mProps []modelProp
	if len(args) == 0 {
		return mProps
	}
	for _, a := range args[1:] {
		ax := strings.Split(a, ":")
		p := modelProp{
			Name: inflect.Camelize(ax[0]),
			Type: "string",
		}
		if len(ax) > 1 {
			p.Type = strings.ToLower(strings.TrimPrefix(ax[1], "nulls."))
		}
		mProps = append(mProps, p)
	}
	return mProps
}
