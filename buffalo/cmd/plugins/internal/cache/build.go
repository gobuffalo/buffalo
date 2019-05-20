package cache

import (
	"os"

	"github.com/gobuffalo/buffalo/plugins"
	"github.com/spf13/cobra"
)

// BuildCmd rebuilds the plugins cache
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "rebuilds the plugins cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		os.RemoveAll(plugins.CachePath)
		_, err := plugins.Available()
		return err
	},
}
