package plugins

import (
	"github.com/gobuffalo/buffalo/buffalo/cmd/plugins/internal/cache"
	"github.com/spf13/cobra"
)

// cacheCmd represents the cache command
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "commands for managing the plugins cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cache.ListCmd.RunE(cmd, args)
	},
}

func init() {
	cacheCmd.AddCommand(cache.CleanCmd)
	cacheCmd.AddCommand(cache.ListCmd)
	cacheCmd.AddCommand(cache.BuildCmd)
	PluginsCmd.AddCommand(cacheCmd)
}
