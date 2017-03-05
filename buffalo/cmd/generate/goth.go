package generate

import (
	"errors"

	"github.com/gobuffalo/buffalo/generators/goth"
	"github.com/gobuffalo/makr"
	"github.com/spf13/cobra"
)

// GothCmd generates a actions/goth.go file configured to the specified providers.
var GothCmd = &cobra.Command{
	Use:   "goth [provider provider...]",
	Short: "Generates a actions/goth.go file configured to the specified providers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must specify at least one provider")
		}
		g, err := goth.New()
		if err != nil {
			return err
		}
		return g.Run(".", makr.Data{
			"providers": args,
		})
	},
}
