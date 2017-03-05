package generate

import (
	"errors"

	"github.com/gobuffalo/buffalo/generators/resource"
	"github.com/gobuffalo/makr"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

// ResourceCmd generates a new actions/resource file and a stub test.
var ResourceCmd = &cobra.Command{
	Use:     "resource [name]",
	Aliases: []string{"r"},
	Short:   "Generates a new actions/resource file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must specify a resource name")
		}
		name := args[0]
		data := makr.Data{
			"name":         name,
			"singular":     inflect.Singularize(name),
			"plural":       inflect.Pluralize(name),
			"camel":        inflect.Camelize(name),
			"under":        inflect.Underscore(name),
			"downFirstCap": inflect.CamelizeDownFirst(name),
			"actions":      []string{"List", "Show", "New", "Create", "Edit", "Update", "Destroy"},
		}
		g, err := resource.New(data)
		if err != nil {
			return err
		}
		return g.Run(".", data)
	},
}
