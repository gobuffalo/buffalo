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
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

var force bool
var verbose bool
var skipPop bool
var skipJQuery bool
var skipBootstrap bool
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

		return genNewFiles(name, rootPath)
	},
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
	args := []string{"get"}
	if verbose {
		args = append(args, "-v")
	}
	args = append(args, pkg)
	return exec.Command("go", args...)
}

func genNewFiles(name, rootPath string) error {
	packagePath := strings.Replace(rootPath, filepath.Join(os.Getenv("GOPATH"), "src")+"/", "", 1)

	data := map[string]interface{}{
		"name":          name,
		"titleName":     inflect.Titleize(name),
		"packagePath":   packagePath,
		"actionsPath":   filepath.Join(packagePath, "actions"),
		"modelsPath":    filepath.Join(packagePath, "models"),
		"withPop":       !skipPop,
		"withJQuery":    !skipJQuery,
		"withBootstrap": !skipBootstrap,
		"dbType":        dbType,
	}

	g := newAppGenerator()
	return g.Run(rootPath, data)
}

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVarP(&force, "force", "f", false, "delete and remake if the app already exists")
	newCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbosely print out the go get/install commands")
	newCmd.Flags().BoolVar(&skipPop, "skip-pop", false, "skips adding pop/soda to your app")
	newCmd.Flags().BoolVar(&skipJQuery, "skip-jquery", false, "skips adding jQuery to your app")
	newCmd.Flags().BoolVar(&skipBootstrap, "skip-bootstrap", false, "skips adding Bootstrap to your app")
	newCmd.Flags().StringVar(&dbType, "db-type", "postgres", "specify the type of database you want to use [postgres, mysql, sqlite3]")
}
