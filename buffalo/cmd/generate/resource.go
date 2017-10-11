package generate

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo/generators/resource"
	"github.com/gobuffalo/makr"
	"github.com/spf13/cobra"
)

var resourceOptions = struct {
	SkipMigration bool
	SkipModel     bool
	Name          string
	MimeType      string
}{}

// ResourceCmd generates a new actions/resource file and a stub test.
var ResourceCmd = &cobra.Command{
	Use:     "resource [name]",
	Example: resourceExamples,
	Aliases: []string{"r"},
	Short:   "Generates a new actions/resource file",
	RunE: func(cmd *cobra.Command, args []string) error {
		o, err := resource.NewOptions(resourceOptions.Name, args...)
		if err != nil {
			return errors.WithStack(err)
		}
		o.MimeType = strings.ToUpper(resourceOptions.MimeType)
		o.SkipModel = resourceOptions.SkipModel
		o.SkipMigration = resourceOptions.SkipMigration

		if err := o.Validate(); err != nil {
			return err
		}

		data := makr.Data{}
		g, err := resource.New(o, data)
		if err != nil {
			return err
		}
		return g.Run(".", data)
	},
}

var resourceMN string

func init() {
	ResourceCmd.Flags().BoolVarP(&resourceOptions.SkipMigration, "skip-migration", "s", false, "sets resource generator not-to add model migration")
	ResourceCmd.Flags().BoolVarP(&resourceOptions.SkipModel, "skip-model", "", false, "makes resource generator not to generate model nor migrations")
	ResourceCmd.Flags().StringVarP(&resourceOptions.Name, "name", "n", "", "allows to define a different model name for the resource being generated.")
	ResourceCmd.Flags().StringVarP(&resourceOptions.MimeType, "type", "", "html", "sets the resource type (html or json)")
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
