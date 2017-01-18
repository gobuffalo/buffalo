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
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/velvet"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

var force bool
var verbose bool
var skipPop bool
var skipWebpack bool
var dbType = "postgres"

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Creates a new Buffalo application",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must enter a name for your new application")
		}

		name := args[0]

		err := validateInGoPath(name)
		if err != nil {
			return err
		}

		rootPath, err := rootPath(name)
		if err != nil {
			return err
		}

		s, _ := os.Stat(rootPath)
		if s != nil {
			if force {
				os.RemoveAll(rootPath)
			} else {
				return fmt.Errorf("%s already exists! Either delete it or use the -f flag to force", name)
			}
		}

		return genNewFiles(name, rootPath)
	},
}

func validateInGoPath(name string) error {
	gp, err := envy.MustGet("GOPATH")
	if err != nil {
		fmt.Println(noGoPath)
		os.Exit(-1)
	}

	root, err := rootPath(name)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(root, filepath.Join(gp, "src")) {
		u, err := user.Current()
		if err != nil {
			return err
		}
		t, err := velvet.Render(notInGoWorkspace, velvet.NewContextWith(map[string]interface{}{
			"name":     name,
			"gopath":   gp,
			"current":  root,
			"username": u.Username,
		}))
		if err != nil {
			return err
		}
		fmt.Println(t)
		os.Exit(-1)
	}
	return nil
}

func rootPath(name string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	rootPath := filepath.Join(pwd, name)
	return rootPath, nil
}

func packagePath(rootPath string) string {
	gosrcpath := strings.Replace(filepath.Join(os.Getenv("GOPATH"), "src"), "\\", "/", -1)
	rootPath = strings.Replace(rootPath, "\\", "/", -1)
	return strings.Replace(rootPath, gosrcpath+"/", "", 2)
}

func genNewFiles(name, rootPath string) error {
	packagePath := packagePath(rootPath)

	data := map[string]interface{}{
		"name":        name,
		"titleName":   inflect.Titleize(name),
		"packagePath": packagePath,
		"actionsPath": packagePath + "/actions",
		"modelsPath":  packagePath + "/models",
		"withPop":     !skipPop,
		"withWebpack": !skipWebpack,
		"dbType":      dbType,
		"version":     Version,
	}

	g := newAppGenerator(data)
	return g.Run(rootPath, data)
}

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVarP(&force, "force", "f", false, "delete and remake if the app already exists")
	newCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbosely print out the go get/install commands")
	newCmd.Flags().BoolVar(&skipPop, "skip-pop", false, "skips adding pop/soda to your app")
	newCmd.Flags().BoolVar(&skipWebpack, "skip-webpack", false, "skips adding Webpack to your app")
	newCmd.Flags().StringVar(&dbType, "db-type", "postgres", "specify the type of database you want to use [postgres, mysql, sqlite3]")
}

const notInGoWorkspace = `Oops! It would appear that you are not in your Go Workspace.

Your $GOPATH is set to "{{gopath}}".

You are currently in "{{current}}".

The standard location for putting Go projects is something along the lines of "$GOPATH/src/github.com/{{username}}/{{name}}" (adjust accordingly).

We recommend you go to "$GOPATH/src/github.com/{{username}}/" and try "buffalo new {{name}}" again.`

const noGoPath = `You do not have a $GOPATH set. In order to work with Go, you must set up your $GOPATH and your Go Workspace.

We recommend reading this tutorial on setting everything up: https://www.goinggo.net/2016/05/installing-go-and-your-workspace.html

When you're ready come back and try again. Don't worry, Buffalo will be right here waiting for you. :)`
