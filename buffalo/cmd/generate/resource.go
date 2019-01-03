package generate

import (
	"context"

	"github.com/gobuffalo/buffalo/genny/resource"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/attrs"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var resourceOptions = &resource.Options{}

// ResourceCmd generates a new actions/resource file and a stub test.
var ResourceCmd = &cobra.Command{
	Use:     "resource [name]",
	Example: resourceExamples,
	Aliases: []string{"r"},
	Short:   "Generate a new actions/resource file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must supply a name")
		}

		run := genny.WetRunner(context.Background())

		resourceOptions.Name = args[0]
		ats, err := attrs.ParseArgs(args[0:]...)
		if err != nil {
			return errors.WithStack(err)
		}
		resourceOptions.Attrs = ats

		if err := run.WithNew(resource.New(resourceOptions)); err != nil {
			return err
		}
		return run.Run()
		// o, err := resource.New(resourceOptions.Name, args...)
		// if err != nil {
		// 	return errors.WithStack(err)
		// }
		// if o.App.AsAPI {
		// 	resourceOptions.SkipTemplates = true
		// }
		// o.SkipModel = resourceOptions.SkipModel
		// o.SkipMigration = resourceOptions.SkipMigration
		// o.SkipTemplates = resourceOptions.SkipTemplates
		// if resourceOptions.ModelName != "" {
		// 	o.UseModel = true
		// 	o.Model = name.New(resourceOptions.ModelName)
		// }
		//
		// if err := o.Validate(); err != nil {
		// 	return err
		// }
		//
		// return o.Run(".", makr.Data{})
	},
}

func init() {
	ResourceCmd.Flags().BoolVarP(&resourceOptions.SkipMigration, "skip-migration", "s", false, "tells resource generator not-to add model migration")
	ResourceCmd.Flags().BoolVarP(&resourceOptions.SkipModel, "skip-model", "", false, "tells resource generator not to generate model nor migrations")
	ResourceCmd.Flags().BoolVarP(&resourceOptions.SkipTemplates, "skip-templates", "", false, "tells resource generator not to generate templates for the resource")
	ResourceCmd.Flags().StringVarP(&resourceOptions.Model, "use-model", "", "", "tells resource generator to reference an existing model in generated code")
	ResourceCmd.Flags().StringVarP(&resourceOptions.Name, "name", "n", "", "allows to define a different model name for the resource being generated.")
}

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
