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

	"github.com/spf13/cobra"
)

var force bool

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

		fmt.Printf("-- ./%s\n", name)
		err = os.MkdirAll(name, 0777)
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

	return runCommands(
		exec.Command("go", "get", "-u", "-v", "github.com/Masterminds/glide"),
		exec.Command("glide", "init", "--non-interactive"),
		exec.Command("glide", "get", "-v", "-u", "--non-interactive", "github.com/markbates/refresh"),
		exec.Command("glide", "get", "-v", "-u", "--non-interactive", "github.com/markbates/pop/"),
		exec.Command("glide", "get", "-v", "-u", "--non-interactive", "github.com/markbates/pop/soda"),
		exec.Command("glide", "get", "-v", "-u", "--non-interactive", "github.com/markbates/grift"),
		exec.Command("glide", "rebuild"),
		exec.Command("refresh", "init"),
	)
}

func runCommands(cmds ...*exec.Cmd) error {
	for _, cmd := range cmds {
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
		err := os.MkdirAll(filepath.Join(rootPath, dir), 0777)
		if err != nil {
			return err
		}
		t, err := template.New(fn).Parse(tv)
		if err != nil {
			return err
		}
		fmt.Printf("-- ./%s/%s\n", name, fn)
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

func init() {
	RootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	newCmd.Flags().BoolVarP(&force, "force", "f", false, "delete and remake if the app already exists")

}
