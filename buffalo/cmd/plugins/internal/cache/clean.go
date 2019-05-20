package cache

import (
	"os"

	"github.com/gobuffalo/buffalo/plugins"
	"github.com/spf13/cobra"
)

// cacheCmd represents the cache command
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "cleans the plugins cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		os.RemoveAll(plugins.CachePath)
		return nil
	},
}
