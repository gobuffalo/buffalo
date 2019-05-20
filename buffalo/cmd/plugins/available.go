package plugins

import (
	"github.com/gobuffalo/buffalo-plugins/plugins/plugcmds"
	"github.com/spf13/cobra"
)

var Available = plugcmds.NewAvailable()

var PluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "tools for working with buffalo plugins",
}

func init() {
	PluginsCmd.AddCommand(addCmd)
	PluginsCmd.AddCommand(listCmd)
	PluginsCmd.AddCommand(generateCmd)
	PluginsCmd.AddCommand(removeCmd)
	PluginsCmd.AddCommand(installCmd)
	PluginsCmd.AddCommand(cacheCmd)

	Available.Add("generate", generateCmd)
	Available.ListenFor("buffalo:setup:.+", Listen)
}
