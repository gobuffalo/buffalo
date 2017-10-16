package generate

import (
	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo/generators/grift"
	"github.com/gobuffalo/makr"
	"github.com/spf13/cobra"
)

//TaskCmd is the command called with the generate grift cli.
var TaskCmd = &cobra.Command{
	Use:     "task [name]",
	Aliases: []string{"t", "grift"},
	Short:   "Generates a grift task",
	RunE: func(cmd *cobra.Command, args []string) error {
		g, err := grift.New(args...)
		if err != nil {
			return errors.WithStack(err)
		}
		return g.Run(".", makr.Data{})
	},
}
