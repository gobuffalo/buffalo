package generate

import (
	"errors"
	"os"
	"path/filepath"

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
var UseResourceModel = false

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
		importPath, err := getImportPath()
		if err != nil {
			return err
		}
		name := args[0]
		data := makr.Data{
			"name":         name,
			"singular":     inflect.Singularize(name),
			"plural":       inflect.Pluralize(name),
			"camel":        inflect.Camelize(name),
			"under":        inflect.Underscore(name),
			"downFirstCap": inflect.CamelizeDownFirst(name),
			"model":        inflect.Singularize(inflect.Camelize(name)),
			"actions":      []string{"List", "Show", "New", "Create", "Edit", "Update", "Destroy"},
			"args":         args,

			// Flags
			"skipMigration": SkipResourceMigration,
			"skipModel":     SkipResourceModel,
			"useModel":      UseResourceModel,

			// System
			"importPath": importPath,
		}

		g, err := resource.New(data)
		if err != nil {
			return err
		}
		return g.Run(".", data)
	},
}

// getImportPath returns the import path of the app created by Buffalo
func getImportPath() (string, error) {
	fp, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	base := filepath.Join(os.Getenv("GOPATH"), "src")
	rel, err := filepath.Rel(base, fp)
	if err != nil {
		return rel, err
	}
	return filepath.Dir(rel), nil
}
