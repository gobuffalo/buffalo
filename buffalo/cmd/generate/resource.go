package generate

import (
	"context"
	"fmt"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/buffalo/genny/resource"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/logger"
	"github.com/spf13/cobra"
)

var resourceOptions = struct {
	*resource.Options
	Verbose bool
	DryRun  bool
}{
	Options: &resource.Options{},
}

// ResourceCmd generates a new actions/resource file and a stub test.
var ResourceCmd = &cobra.Command{
	Use:     "resource [name]",
	Example: resourceExamples,
	Aliases: []string{"r"},
	Short:   "Generate a new actions/resource file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you must supply a name")
		}
		ctx := context.Background()
		run := genny.WetRunner(ctx)
		if resourceOptions.DryRun {
			run = genny.DryRunner(ctx)
		}

		if resourceOptions.Verbose {
			lg := logger.New(logger.DebugLevel)
			run.Logger = lg
		}

		if len(resourceOptions.Name) == 0 {
			resourceOptions.Name = args[0]
		}
		if len(args) > 1 {
			ats, err := attrs.ParseArgs(args[1:]...)
			if err != nil {
				return err
			}
			resourceOptions.Attrs = ats
		}

		if err := run.WithNew(resource.New(resourceOptions.Options)); err != nil {
			return err
		}
		return run.Run()
	},
}

func init() {
	ResourceCmd.Flags().BoolVarP(&resourceOptions.SkipMigration, "skip-migration", "s", false, "tells resource generator not-to add model migration")
	ResourceCmd.Flags().BoolVarP(&resourceOptions.SkipModel, "skip-model", "", false, "tells resource generator not to generate model nor migrations")
	ResourceCmd.Flags().BoolVarP(&resourceOptions.SkipTemplates, "skip-templates", "", false, "tells resource generator not to generate templates for the resource")
	ResourceCmd.Flags().StringVarP(&resourceOptions.Model, "use-model", "", "", "tells resource generator to reference an existing model in generated code")
	ResourceCmd.Flags().StringVarP(&resourceOptions.Name, "name", "n", "", "allows to define a different model name for the resource being generated.")
	ResourceCmd.Flags().BoolVarP(&resourceOptions.DryRun, "dry-run", "d", false, "dry run")
	ResourceCmd.Flags().BoolVarP(&resourceOptions.Verbose, "verbose", "v", false, "verbosely print out the go get commands")
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
