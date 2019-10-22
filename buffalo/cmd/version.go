package cmd

import (
	"encoding/json"
	"os"

	"github.com/gobuffalo/buffalo/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var jsonOutput bool

func init() {
	decorate("version", versionCmd)
	versionCmd.Flags().BoolVar(&jsonOutput, "json", false, "Print information in json format")

	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(c *cobra.Command, args []string) {
		if jsonOutput {
			build := runtime.BuildInfo{}
			build.Version = runtime.Version
			enc := json.NewEncoder(os.Stderr)
			enc.SetIndent("", "    ")
			enc.Encode(build)
			return
		}

		logrus.Infof("Buffalo version is: %s", runtime.Version)
	},
	// needed to override the root level pre-run func
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		return nil
	},
}
