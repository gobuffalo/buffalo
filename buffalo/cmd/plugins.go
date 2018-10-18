package cmd

import (
	"github.com/gobuffalo/buffalo-plugins/plugins"
	"github.com/markbates/oncer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var _plugs plugins.List

func plugs() plugins.List {
	oncer.Do("buffalo/cmd/plugins", func() {
		var err error
		_plugs, err = plugins.Available()
		if err != nil {
			_plugs = plugins.List{}
			logrus.Errorf("error loading plugins %s\n", err)
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
