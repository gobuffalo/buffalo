package generate

import (
	"github.com/gobuffalo/buffalo/generators/docker"
	"github.com/spf13/cobra"
)

var dockerOptions = docker.NewOptions()

// DockerCmd generates a new Dockerfile
var DockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Generates a Dockerfile",
	RunE: func(cmd *cobra.Command, args []string) error {
		dockerOptions.Version = Version
		return docker.Run(".", dockerOptions)
	},
}

func init() {
	DockerCmd.Flags().StringVar(&dockerOptions.Style, "style", "multi", "what style Dockerfile to generate [multi, standard]")
}
