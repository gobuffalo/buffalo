package generate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/gentronics"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

var runningTests = false

//ActionCmd is the cmd that generates actions.
var ActionCmd = &cobra.Command{
	Use:     "action [name] [actionName...]",
	Aliases: []string{"a", "actions"},
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

		data := gentronics.Data{
			"filename":  inflect.Underscore(name),
			"namespace": inflect.Camelize(name),
		}

		filePath := filepath.Join("actions", fmt.Sprintf("%v.go", data["filename"]))
		actionsTemplate := buildActionsTemplate(filePath)
		testFilePath := filepath.Join("actions", fmt.Sprintf("%v_test.go", data["filename"]))
		testsTemplate := buildTestsTemplate(testFilePath)
		actionsToAdd := findActionsToAdd(name, filePath, actions)
		testsToAdd := findTestsToAdd(name, testFilePath, actions)

		data["actions"] = actionsToAdd
		data["tests"] = testsToAdd

		g := gentronics.New()
		g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s.go", data["filename"])), actionsTemplate))
		g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s_test.go", data["filename"])), testsTemplate))
		g.Add(&gentronics.Func{
			Should: func(data gentronics.Data) bool { return true },
			Runner: func(root string, data gentronics.Data) error {
				routes := []string{}
				for _, a := range actions {
					routes = append(routes, fmt.Sprintf("app.GET(\"/%s/%s\", %s)", name, a, data["namespace"].(string)+inflect.Camelize(a)))
				}
				return addInsideAppBlock(routes...)
			},
		})
		addTemplateFiles(actionsToAdd, data)

		if !runningTests {
			g.Add(Fmt)
		}

		return g.Run(".", data)
	},
}

func buildActionsTemplate(filePath string) string {
	actionsTemplate := rActionFileT
	fileContents, err := ioutil.ReadFile(filePath)
	if err == nil {
		actionsTemplate = string(fileContents)
	}

	actionsTemplate = actionsTemplate + `
{{#each actions as |action|}}
// {{namespace}}{{camelize action}} default implementation.
func {{namespace}}{{camelize action}}(c buffalo.Context) error {
	return c.Render(200, r.HTML("{{filename}}/{{underscore action}}.html"))
}
{{/each}}`
	return actionsTemplate
}

func buildTestsTemplate(filePath string) string {
	testsTemplate := `package actions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)
	`
	fileContents, err := ioutil.ReadFile(filePath)
	if err == nil {
		testsTemplate = string(fileContents)
	}

	testsTemplate = testsTemplate + `
{{#each tests as |action|}}
func Test_{{namespace}}_{{camelize action}}(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

{{/each}}`
	return testsTemplate
}

func addTemplateFiles(actionsToAdd []string, data gentronics.Data) {
	for _, action := range actionsToAdd {
		vg := gentronics.New()
		viewPath := filepath.Join("templates", fmt.Sprintf("%s", data["filename"]), fmt.Sprintf("%s.html", inflect.Underscore(action)))
		vg.Add(gentronics.NewFile(viewPath, rViewT))
		vg.Run(".", gentronics.Data{
			"namespace": data["namespace"],
			"action":    inflect.Camelize(action),
		})
	}
}

func findActionsToAdd(name, path string, actions []string) []string {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		fileContents = []byte("")
	}

	actionsToAdd := []string{}

	for _, action := range actions {
		funcSignature := fmt.Sprintf("func %s%s(c buffalo.Context) error", inflect.Camelize(name), inflect.Camelize(action))
		if strings.Contains(string(fileContents), funcSignature) {
			fmt.Printf("--> [warning] skipping %v%v since it already exists\n", inflect.Camelize(name), inflect.Camelize(action))
			continue
		}

		actionsToAdd = append(actionsToAdd, action)
	}

	return actionsToAdd
}

func findTestsToAdd(name, path string, actions []string) []string {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		fileContents = []byte("")
	}

	actionsToAdd := []string{}

	for _, action := range actions {
		funcSignature := fmt.Sprintf("func Test_%v_%v(c buffalo.Context) error", inflect.Camelize(name), inflect.Camelize(action))
		if strings.Contains(string(fileContents), funcSignature) {
			fmt.Printf("--> [warning] skipping Test_%v_%v since it already exists\n", inflect.Camelize(name), inflect.Camelize(action))
			continue
		}

		actionsToAdd = append(actionsToAdd, action)
	}

	return actionsToAdd
}

const (
	rActionFileT = `package actions
import "github.com/gobuffalo/buffalo"`

	rViewT       = `<h1>{{namespace}}#{{action}}</h1>`
	rActionFuncT = `
// {{namespace}}{{action}} default implementation.
func {{namespace}}{{action}}(c buffalo.Context) error {
	return c.Render(200, r.HTML("{{namespace_under}}/{{action_under}}.html"))
}
`
)
