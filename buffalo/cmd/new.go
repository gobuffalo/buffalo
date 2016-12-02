// Copyright © 2016 Mark Bates <mark@markbates.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/markbates/pop/soda/cmd/generate"
	"github.com/spf13/cobra"
)

var force bool
var verbose bool
var skipPop bool
var dbType = "postgres"

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Creates a new Buffalo application",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must enter a name for your new application.")
		}
		name := args[0]
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		rootPath := filepath.Join(pwd, name)

		s, _ := os.Stat(rootPath)
		if s != nil {
			if force {
				os.RemoveAll(rootPath)
			} else {
				return fmt.Errorf("%s already exists! Either delete it or use the -f flag to force.\n", name)
			}
		}

		fmt.Printf("--> ./%s\n", name)
		err = os.MkdirAll(name, 0755)
		if err != nil {
			return err
		}

		err = genNewFiles(name, rootPath)
		if err != nil {
			return err
		}

		err = installDeps(pwd, rootPath)
		if err != nil {
			return err
		}
		return err
	},
}

func installDeps(pwd string, rootPath string) error {
	defer os.Chdir(pwd)
	err := os.Chdir(rootPath)
	if err != nil {
		return err
	}

	cmds := []*exec.Cmd{
		goGet("github.com/markbates/refresh/..."),
		goInstall("github.com/markbates/refresh"),
		goGet("github.com/markbates/grift/..."),
		goInstall("github.com/markbates/grift"),
	}

	if !skipPop {
		cmds = append(cmds,
			goGet("github.com/markbates/pop/..."),
			goInstall("github.com/markbates/pop/soda"),
		)
	}

	cmds = append(cmds, appGoGet())

	err = runCommands(cmds...)

	if !skipPop {
		generate.GenerateConfig(dbType, "./database.yml")
	}

	if err != nil {
		return err
	}

	return err
}

func goInstall(pkg string) *exec.Cmd {
	args := []string{"install"}
	if verbose {
		args = append(args, "-v")
	}
	args = append(args, pkg)
	return exec.Command("go", args...)
}

func goGet(pkg string) *exec.Cmd {
	args := []string{"get", "-u"}
	if verbose {
		args = append(args, "-v")
	}
	args = append(args, pkg)
	return exec.Command("go", args...)
}

func runCommands(cmds ...*exec.Cmd) error {
	for _, cmd := range cmds {
		fmt.Printf("--> %s\n", strings.Join(cmd.Args, " "))
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func genNewFiles(name, rootPath string) error {
	packagePath := strings.Replace(rootPath, filepath.Join(os.Getenv("GOPATH"), "src")+"/", "", 1)

	data := map[string]interface{}{
		"name":        name,
		"packagePath": packagePath,
		"actionsPath": filepath.Join(packagePath, "actions"),
	}

	for fn, tv := range newTemplates {
		dir := filepath.Dir(fn)
		err := os.MkdirAll(filepath.Join(rootPath, dir), 0755)
		if err != nil {
			return err
		}
		t, err := template.New(fn).Parse(tv)
		if err != nil {
			return err
		}
		fmt.Printf("--> ./%s/%s\n", name, fn)
		f, err := os.Create(filepath.Join(rootPath, fn))
		if err != nil {
			return err
		}
		err = t.Execute(f, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func appGoGet() *exec.Cmd {
	appArgs := []string{"get", "-t"}
	if verbose {
		appArgs = append(appArgs, "-v")
	}
	appArgs = append(appArgs, "./...")
	return exec.Command("go", appArgs...)
}

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVarP(&force, "force", "f", false, "delete and remake if the app already exists")
	newCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbosely print out the go get/install commands")
	newCmd.Flags().BoolVar(&skipPop, "skip-pop", false, "skips add pop/soda to your app")
	newCmd.Flags().StringVar(&dbType, "db-type", "postgres", "specify the type of database you want to use [postgres, mysql, sqlite3]")
}
