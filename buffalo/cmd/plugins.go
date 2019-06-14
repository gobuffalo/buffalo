package cmd

import (
	pluginscmd "github.com/gobuffalo/buffalo/buffalo/cmd/plugins"
	"github.com/gobuffalo/buffalo/plugins"
	"github.com/markbates/oncer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(pluginscmd.PluginsCmd)
}

var _plugs plugins.List

func plugs() plugins.List {
	oncer.Do("buffalo/cmd/plugins", func() {
		var err error
		_plugs, err = plugins.Available()
		if err != nil {
			_plugs = plugins.List{}
			logrus.Errorf("error loading plugins %s", err)
		}
	})
	return _plugs
}

func decorate(name string, cmd *cobra.Command) {
	pugs := plugs()
	for _, c := range pugs[name] {
		anywhereCommands = append(anywhereCommands, c.Name)
		cc := plugins.Decorate(c)
		cmd.AddCommand(cc)
	}
}
