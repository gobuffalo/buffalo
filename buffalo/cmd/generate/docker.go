package generate

import (
	"os"

	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo/generators/docker"
	"github.com/gobuffalo/envy"
	"github.com/spf13/cobra"
)

var dockerOptions = struct {
	Style string
}{}

// DockerCmd generates a new Dockerfile
var DockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Generates a Dockerfile",
	RunE: func(cmd *cobra.Command, args []string) error {
		packagePath := envy.CurrentPackage()

		var webpack bool
		if _, err := os.Stat("package.json"); err == nil {
			webpack = true
		}
		data := map[string]interface{}{
			"packagePath": packagePath,
			"version":     Version,
			"docker":      dockerOptions.Style,
			"asWeb":       webpack,
			"withWepack":  webpack,
			"withYarn":    false,
		}

		if _, err := os.Stat("yarn.lock"); err == nil {
			data["withYarn"] = true
		}

		g, err := docker.New()
		if err != nil {
			return errors.WithStack(err)
		}
		return g.Run(".", data)
	},
}

func init() {
	DockerCmd.Flags().StringVar(&dockerOptions.Style, "style", "multi", "what style Dockerfile to generate [multi, standard]")
}
