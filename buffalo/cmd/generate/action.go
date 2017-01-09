// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/bep/inflect"
	"github.com/markbates/gentronics"
	"github.com/spf13/cobra"
)

//ActionCmd is the cmd that generates actions.
var ActionCmd = &cobra.Command{
	Use:     "action [name] [actionName...]",
	Aliases: []string{"a"},
	Short:   "Generates new action(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("you should provide action name and handler name at least")
		}

		if _, err := os.Stat("actions"); err != nil {
			return errors.New("actions directory not found, ensure you're inside your buffalo folder")
		}

		name := args[0]
		actions := args[1:]

		data := gentronics.Data{"under": inflect.Underscore(name)}
		_, err := os.Stat(filepath.Join("actions", fmt.Sprintf("%v.go", data["under"])))
		fileExists := err == nil

		if !fileExists {
			g := gentronics.New()
			g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s.go", data["under"])), rActionFileT))
			g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s_test.go", data["under"])), rActionTest))
			g.Add(Fmt)

			err = g.Run(".", data)

			if err != nil {
				return err
			}
		}

		return generateActionComponents(name, actions)
	},
}

func generateActionComponents(name string, actions []string) error {

	actionData := gentronics.Data{
		"name":  name,
		"under": inflect.Underscore(name),
	}

	path := filepath.Join("actions", fmt.Sprintf("%v.go", actionData["under"]))
	_, err := os.Stat(path)
	fileExists := err == nil
	fileContents, _ := ioutil.ReadFile(path)

	for _, action := range actions {
		actionData["namespace"] = inflect.Camelize(name)
		actionData["action"] = inflect.Camelize(action)
		actionData["action_under"] = inflect.Underscore(action)

		if fileExists {
			funcSignature := fmt.Sprintf("func %s%s(c buffalo.Context) error", actionData["namespace"], actionData["currentAction"])
			if strings.Contains(string(fileContents), funcSignature) {
				fmt.Printf("--> [warning] skipping %v%v since it already exists\n", actionData["namespace"], actionData["currentAction"])
				continue
			}
		}

		if err = appendActionToFile(path, actionData); err != nil {
			return err
		}

		if err = generateTemplate(actionData); err != nil {
			return err
		}
	}

	return nil
}

func appendActionToFile(path string, actionData gentronics.Data) error {
	fileContents, _ := ioutil.ReadFile(path)
	t, _ := raymond.Parse(rActionFuncT)
	fn, _ := t.Exec(actionData)

	fileContents = []byte(string(fileContents) + fn)
	return ioutil.WriteFile(path, fileContents, 0755)
}

func generateTemplate(actionData gentronics.Data) error {
	fg := gentronics.New()
	templatePath := filepath.Join("templates", fmt.Sprintf("%s", actionData["under"]), fmt.Sprintf("%s.html", actionData["action_under"]))
	fg.Add(gentronics.NewFile(templatePath, rViewT))
	return fg.Run(".", actionData)
}

const (
	rActionFileT = `package actions

import "github.com/gobuffalo/buffalo"

`
	rViewT       = `<h1>{{namespace}}#{{action}}</h1>`
	rActionFuncT = `
    
    // {{namespace}}{{action}} default implementation.
    func {{namespace}}{{action}}(c buffalo.Context) error {
	    return c.Render(200, r.String("{{camel}}#{{.}}"))
    }
    `

	rActionTestT = `
func Test_{{namespace}}{{action}}(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}
    `
)
