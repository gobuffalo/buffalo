package generate

import (
	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	"github.com/markbates/gentronics"
	"github.com/spf13/cobra"
)

var withYarn bool

// WebpackCmd generates a new actions/resource file and a stub test.
var WebpackCmd = &cobra.Command{
	Use:   "webpack [flags]",
	Short: "Generates a webpack asset pipeline.",
	RunE: func(cmd *cobra.Command, args []string) error {
		data := gentronics.Data{
			"withWebpack": true,
			"withYarn":    withYarn,
		}
		wg, err := webpack.New(data)
		if err != nil {
			return err
		}
		return wg.Run(".", data)
	},
}
