package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints off diagnostic information useful for debugging.",
	RunE: func(cmd *cobra.Command, args []string) error {
		bb := os.Stdout

		bb.WriteString(fmt.Sprintf("### Buffalo Version\n%s\n", Version))

		bb.WriteString("\n### App Information\n")
		app := meta.New(".")
		rv := reflect.ValueOf(app)
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			bb.WriteString(fmt.Sprintf("%s=%v\n", f.Name, rv.FieldByName(f.Name).Interface()))
		}

		err := checkGoInfo()
		if err != nil {
			return errors.WithStack(err)
		}

		checkExternalsTools()

		return nil
	},
}

func checkGoInfo() error {
	bb := os.Stdout

	bb.WriteString("\n### Go Version\n")
	c := exec.Command(envy.Get("GO_BIN", "go"), "version")
	c.Stdout = bb
	err := c.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	bb.WriteString("\n### Go Env\n")
	c = exec.Command(envy.Get("GO_BIN", "go"), "env")
	c.Stdout = bb
	c.Stderr = bb
	c.Run()
}

func checkExternalsTools() {
	bb := os.Stdout

	infoCommands := []struct {
		Name      string
		PathName  string
		Cmd       *exec.Cmd
		InfoLabel string
	}{
		{"Node", "node", exec.Command("node", "--version"), "\n### Node Version\n"},
		{"NPM", "npm", exec.Command("npm", "--version"), "\n### NPM Version\n"},
		{"Yarn", "yarn", exec.Command("yarn", "--version"), "\n### Yarn Version\n"},
		{"dep", "dep", exec.Command("dep", "version"), "\n### Dep Version\n"},
		{"dep", "dep", exec.Command("dep", "status"), "\n### Dep Status\n"},
		{"PostgreSQL", "pg_ctl", exec.Command("pg_ctl", "--version"), "\n### PostgreSQL Version\n"},
		{"MySQL", "mysql", exec.Command("mysql", "--version"), "\n### MySQL Version\n"},
		{"SQLite", "sqlite3", exec.Command("sqlite3", "--version"), "\n### SQLite Version\n"},
	}

	for _, cmd := range infoCommands {
		bb.WriteString(cmd.InfoLabel)
		execIfExists(cmd.Name, cmd.InfoLabel, cmd.Cmd)
	}

}

func execIfExists(name string, pathName string, c *exec.Cmd) {
	bb := os.Stdout

	if _, err := exec.LookPath("mysql"); err != nil {
		bb.WriteString(fmt.Sprintf("%s Not Found\n", name))
		return
	}

	c.Stdout = bb
	c.Stderr = bb
	c.Run()
}

func init() {
	decorate("info", RootCmd)
	RootCmd.AddCommand(infoCmd)
}
