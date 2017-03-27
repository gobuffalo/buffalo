package generate

import (
	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	"github.com/gobuffalo/makr"
	"github.com/spf13/cobra"
)

var withYarn bool

// WebpackCmd generates a new actions/resource file and a stub test.
var WebpackCmd = &cobra.Command{
	Use:   "webpack [flags]",
	Short: "Generates a webpack asset pipeline.",
	RunE: func(cmd *cobra.Command, args []string) error {
		data := makr.Data{
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

func init() {
	WebpackCmd.Flags().BoolVar(&withYarn, "with-yarn", false, "allows the use of yarn instead of npm as dependency manager")
}
