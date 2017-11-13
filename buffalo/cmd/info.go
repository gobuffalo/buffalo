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

		bb.WriteString("\n### Node Version\n")
		if _, err := exec.LookPath("node"); err == nil {
			c = exec.Command("node", "--version")
			c.Stdout = bb
			c.Stderr = bb
			c.Run()
		} else {
			bb.WriteString("Node Not Found\n")
		}

		bb.WriteString("\n### NPM Version\n")
		if _, err := exec.LookPath("npm"); err == nil {
			c = exec.Command("npm", "--version")
			c.Stdout = bb
			c.Stderr = bb
			c.Run()
		} else {
			bb.WriteString("NPM Not Found\n")
		}

		bb.WriteString("\n### Yarn Version\n")
		if _, err := exec.LookPath("yarn"); err == nil {
			c = exec.Command("yarn", "--version")
			c.Stdout = bb
			c.Stderr = bb
			c.Run()
		} else {
			bb.WriteString("Yarn Not Found\n")
		}

		bb.WriteString("\n### Dep Version\n")
		if _, err := exec.LookPath("dep"); err == nil {
			c = exec.Command("dep", "version")
			c.Stdout = bb
			c.Stderr = bb
			c.Run()
		} else {
			bb.WriteString("dep Not Found\n")
		}

		bb.WriteString("\n### Dep Status\n")
		if _, err := exec.LookPath("dep"); err == nil {
			c = exec.Command("dep", "status")
			c.Stdout = bb
			c.Stderr = bb
			c.Run()
		} else {
			bb.WriteString("dep Not Found\n")
		}

		bb.WriteString("\n### PostgreSQL Version\n")
		if _, err := exec.LookPath("pg_ctl"); err == nil {
			c = exec.Command("pg_ctl", "--version")
			c.Stdout = bb
			c.Stderr = bb
			c.Run()
		} else {
			bb.WriteString("PostgreSQL Not Found\n")
		}

		bb.WriteString("\n### MySQL Version\n")
		if _, err := exec.LookPath("mysql"); err == nil {
			c = exec.Command("mysql", "--version")
			c.Stdout = bb
			c.Stderr = bb
			c.Run()
		} else {
			bb.WriteString("MySQL Not Found\n")
		}

		bb.WriteString("\n### SQLite Version\n")
		if _, err := exec.LookPath("sqlite3"); err == nil {
			c = exec.Command("sqlite3", "--version")
			c.Stdout = bb
			c.Stderr = bb
			c.Run()
		} else {
			bb.WriteString("SQLite Not Found\n")
		}

		return nil
	},
}

func init() {
	decorate("info", RootCmd)
	RootCmd.AddCommand(infoCmd)
}
