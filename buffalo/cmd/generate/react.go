package generate

import (
	"github.com/gobuffalo/buffalo/generators/assets/react"
	"github.com/gobuffalo/makr"
	"github.com/spf13/cobra"
)

// ReactCmd generates a new actions/resource file and a stub test.
var ReactCmd = &cobra.Command{
	Use:   "react [flags]",
	Short: "Generates a react asset pipeline.",
	RunE: func(cmd *cobra.Command, args []string) error {
		data := makr.Data{
			"withReact": true,
			"withYarn":  withYarn,
		}
		wg, err := react.New(data)
		if err != nil {
			return err
		}

		return wg.Run(".", data)
	},
}

func init() {
	ReactCmd.Flags().BoolVar(&withYarn, "with-yarn", false, "allows the use of yarn instead of npm as dependency manager")
}
