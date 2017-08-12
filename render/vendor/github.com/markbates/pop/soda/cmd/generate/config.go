package generate

import (
	"os"
	"path"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/makr"
	"github.com/markbates/going/defaults"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	ConfigCmd.Flags().StringVarP(&dialect, "type", "t", "postgres", "What type of database do you want to use? (postgres, mysql, sqlite3)")
}

var dialect string

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Generates a database.yml file for your project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cflag := cmd.Flag("config")
		cfgFile := defaults.String(cflag.Value.String(), "database.yml")
		dir, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "couldn't get the current directory")
		}
		data := map[string]interface{}{
			"dialect":     dialect,
			"name":        path.Base(dir),
			"packagePath": pkgPath(),
		}
		return GenerateConfig(cfgFile, data)
	},
}

func pkgPath() string {
	pwd, _ := os.Getwd()

	for _, p := range envy.GoPaths() {
		pwd = strings.TrimPrefix(pwd, p)
	}
	return pwd
}

func GenerateConfig(cfgFile string, data map[string]interface{}) error {
	dialect = strings.ToLower(data["dialect"].(string))
	if t, ok := configTemplates[dialect]; ok {
		g := makr.New()
		g.Add(makr.NewFile(cfgFile, t))
		return g.Run(".", data)
	}
	return errors.Errorf("Could not initialize %s!", dialect)
}
