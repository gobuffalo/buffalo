package generate

import (
	"errors"

	"github.com/gobuffalo/buffalo/generators/resource"
	"github.com/gobuffalo/makr"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

const resourceExamples = `

$ buffalo g resource users
Generates:

- actions/users.go
- actions/users_test.go
- models/user.go
- models/user_test.go
- migrations/XXXXX_create_table_users.fizz
- migrations/XXXXX_drop_table_users.fizz

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

`

//SkipResourceMigration allows to generate a resource without the migration.
var SkipResourceMigration = false

//SkipResourceModel allows to generate a resource without the model and Migration.
var SkipResourceModel = false

// ResourceCmd generates a new actions/resource file and a stub test.
var ResourceCmd = &cobra.Command{
	Use:     "resource [name]",
	Example: resourceExamples,
	Aliases: []string{"r"},
	Short:   "Generates a new actions/resource file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must specify a resource name")
		}

		name := args[0]
		data := makr.Data{
			"name":          name,
			"singular":      inflect.Singularize(name),
			"plural":        inflect.Pluralize(name),
			"camel":         inflect.Camelize(name),
			"under":         inflect.Underscore(name),
			"downFirstCap":  inflect.CamelizeDownFirst(name),
			"actions":       []string{"List", "Show", "New", "Create", "Edit", "Update", "Destroy"},
			"args":          args,
			"skipMigration": SkipResourceMigration,
			"skipModel":     SkipResourceModel,
		}

		g, err := resource.New(data)
		if err != nil {
			return err
		}
		return g.Run(".", data)
	},
}
