package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/genny/info"
	"github.com/gobuffalo/clara/genny/rx"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/meta"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var infoOptions = struct {
	Clara *rx.Options
	Info  *info.Options
}{
	Clara: &rx.Options{},
	Info:  &info.Options{},
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print diagnostic information (useful for debugging)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		run := genny.WetRunner(ctx)

		copts := infoOptions.Clara
		if err := run.WithNew(rx.New(copts)); err != nil {
			return err
		}

		iopts := infoOptions.Info
		if err := run.WithNew(info.New(iopts)); err != nil {
			return err
		}

		return run.Run()
		// bb := os.Stdout
		//
		// bb.WriteString(fmt.Sprintf("### Buffalo Version\n%s\n", runtime.Version))
		//
		// bb.WriteString("\n### App Information\n")
		// app := meta.New(".")
		// rv := reflect.ValueOf(app)
		// rt := rv.Type()
		//
		// var err error
		// for i := 0; i < rt.NumField(); i++ {
		// 	f := rt.Field(i)
		// 	if !rv.FieldByName(f.Name).CanInterface() {
		// 		continue
		// 	}
		// 	_, err = bb.WriteString(fmt.Sprintf("%s=%v\n", f.Name, rv.FieldByName(f.Name).Interface()))
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		//
		// if err := runInfoCmds(); err != nil {
		// 	return errors.WithStack(err)
		// }
		//
		// if err := configs(app); err != nil {
		// 	return errors.WithStack(err)
		// }
		//
		// return infoGoMod()
	},
}

func configs(app meta.App) error {
	bb := os.Stdout
	root := filepath.Join(app.Root, "config")
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()
		p := strings.TrimPrefix(path, app.Root)
		p = strings.TrimPrefix(p, string(filepath.Separator))
		bb.WriteString(fmt.Sprintf("\n### %s\n", p))
		if _, err := io.Copy(bb, f); err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

type infoCommand struct {
	Name      string
	PathName  string
	Cmd       *exec.Cmd
	InfoLabel string
}

func infoGoMod() error {
	if _, err := os.Stat("go.mod"); err != nil {
		return nil
	}
	f, err := os.Open("go.mod")
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	bb := os.Stdout
	bb.WriteString("\n### go.mod\n")
	_, err = io.Copy(bb, f)

	return err
}

func runInfoCmds() error {

	commands := []infoCommand{
		{"Go", envy.Get("GO_BIN", "go"), exec.Command(envy.Get("GO_BIN", "go"), "version"), "\n### Go Version\n"},
		{"Go", envy.Get("GO_BIN", "go"), exec.Command(envy.Get("GO_BIN", "go"), "env"), "\n### Go Env\n"},
		{"Node", "node", exec.Command("node", "--version"), "\n### Node Version\n"},
		{"NPM", "npm", exec.Command("npm", "--version"), "\n### NPM Version\n"},
		{"Yarn", "yarnpkg", exec.Command("yarn", "--version"), "\n### Yarn Version\n"},
		{"PostgreSQL", "pg_ctl", exec.Command("pg_ctl", "--version"), "\n### PostgreSQL Version\n"},
		{"MySQL", "mysql", exec.Command("mysql", "--version"), "\n### MySQL Version\n"},
		{"SQLite", "sqlite3", exec.Command("sqlite3", "--version"), "\n### SQLite Version\n"},
		{"dep", "dep", exec.Command("dep", "version"), "\n### Dep Version\n"},
		{"dep", "dep", exec.Command("dep", "status"), "\n### Dep Status\n"},
	}

	for _, cmd := range commands {
		err := execIfExists(cmd)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func execIfExists(infoCmd infoCommand) error {
	bb := os.Stdout
	bb.WriteString(infoCmd.InfoLabel)

	if infoCmd.Name == "dep" {
		if _, err := os.Stat("Gopkg.toml"); err != nil {
			bb.WriteString("could not find a Gopkg.toml file\n")
			return nil
		}
	}

	if _, err := exec.LookPath(infoCmd.PathName); err != nil {
		bb.WriteString(fmt.Sprintf("%s Not Found\n", infoCmd.Name))
		return nil
	}

	infoCmd.Cmd.Stdout = bb
	infoCmd.Cmd.Stderr = bb

	err := infoCmd.Cmd.Run()
	return err
}

func init() {
	app := meta.New(".")
	infoOptions.Clara.App = app
	infoOptions.Info.App = app

	decorate("info", RootCmd)
	RootCmd.AddCommand(infoCmd)
}
