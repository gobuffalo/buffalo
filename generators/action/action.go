package action

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
	"github.com/markbates/inflect"
)

var runningTests bool

// New action generator
func New(name string, actions []string, data makr.Data) (*makr.Generator, error) {
	g := makr.New()

	filePath := filepath.Join("actions", fmt.Sprintf("%v.go", data["filename"]))
	actionsTemplate := buildActionsTemplate(filePath)
	testFilePath := filepath.Join("actions", fmt.Sprintf("%v_test.go", data["filename"]))
	testsTemplate := buildTestsTemplate(testFilePath)
	actionsToAdd := findActionsToAdd(name, filePath, actions)
	testsToAdd := findTestsToAdd(name, testFilePath, actions)

	data["actions"] = actionsToAdd
	data["tests"] = testsToAdd

	g.Add(makr.NewFile(filepath.Join("actions", fmt.Sprintf("%s.go", data["filename"])), actionsTemplate))
	g.Add(makr.NewFile(filepath.Join("actions", fmt.Sprintf("%s_test.go", data["filename"])), testsTemplate))
	g.Add(&makr.Func{
		Should: func(data makr.Data) bool { return true },
		Runner: func(root string, data makr.Data) error {
			routes := []string{}
			for _, a := range actions {
				routes = append(routes, fmt.Sprintf("app.GET(\"/%s/%s\", %s)", name, a, data["namespace"].(string)+inflect.Camelize(a)))
			}
			return generators.AddInsideAppBlock(routes...)
		},
	})
	addTemplateFiles(actionsToAdd, data)

	if !runningTests {
		g.Add(makr.NewCommand(makr.GoFmt()))
	}
	return g, nil
}

func buildActionsTemplate(filePath string) string {
	actionsTemplate := rActionFileT
	fileContents, err := ioutil.ReadFile(filePath)
	if err == nil {
		actionsTemplate = string(fileContents)
	}

	actionsTemplate = actionsTemplate + `
{{ range $action := .actions }}
// {{$.namespace}}{{camelize $action}} default implementation.
func {{$.namespace}}{{camelize $action}}(c buffalo.Context) error {
	return c.Render(200, r.HTML("{{$.filename}}/{{underscore $action}}.html"))
}
{{end}}`
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
{{ range $action := .tests}}
func (as *ActionSuite) Test_{{$.namespace}}_{{camelize $action}}() {
	as.Fail("Not Implemented!")
}

{{end}}`
	return testsTemplate
}

func addTemplateFiles(actionsToAdd []string, data makr.Data) {
	for _, action := range actionsToAdd {
		vg := makr.New()
		viewPath := filepath.Join("templates", fmt.Sprintf("%s", data["filename"]), fmt.Sprintf("%s.html", inflect.Underscore(action)))
		vg.Add(makr.NewFile(viewPath, rViewT))
		vg.Run(".", makr.Data{
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
