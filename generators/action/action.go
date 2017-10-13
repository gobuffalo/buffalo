package action

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/makr"
)

var runningTests bool

// Run action generator
func Run(opts Options, root string, data makr.Data) error {
	g := makr.New()
	defer g.Fmt(root)

	filePath := filepath.Join("actions", fmt.Sprintf("%v.go", opts.Name.File()))
	actionsTemplate := buildActionsTemplate(filePath)
	testFilePath := filepath.Join("actions", fmt.Sprintf("%v_test.go", opts.Name.File()))
	testsTemplate := buildTestsTemplate(testFilePath)
	actionsToAdd := findActionsToAdd(opts.Name, filePath, opts.Actions)
	testsToAdd := findTestsToAdd(opts.Name, testFilePath, opts.Actions)
	handlersToAdd := findHandlersToAdd(opts.Name, filepath.Join("actions", "app.go"), opts.Actions)

	data["opts"] = opts
	data["actions"] = actionsToAdd
	data["tests"] = testsToAdd

	g.Add(makr.NewFile(filePath, actionsTemplate))
	g.Add(makr.NewFile(testFilePath, testsTemplate))
	g.Add(&makr.Func{
		Should: func(data makr.Data) bool { return true },
		Runner: func(root string, data makr.Data) error {
			routes := []string{}
			for _, a := range handlersToAdd {
				routes = append(routes, fmt.Sprintf("app.%s(\"/%s/%s\", %s)", strings.ToUpper(opts.Method), opts.Name, a, opts.Name.Camel()+a.Camel()))
			}
			return generators.AddInsideAppBlock(routes...)
		},
	})
	if !opts.SkipTemplate {
		addTemplateFiles(opts, actionsToAdd, data)
	}
	return g.Run(root, data)
}

func buildActionsTemplate(filePath string) string {
	actionsTemplate := actionsHeaderTmpl
	fileContents, err := ioutil.ReadFile(filePath)
	if err == nil {
		actionsTemplate = string(fileContents)
	}

	actionsTemplate = actionsTemplate + actionsTmpl
	return actionsTemplate
}

func buildTestsTemplate(filePath string) string {
	testsTemplate := testHeaderTmpl

	fileContents, err := ioutil.ReadFile(filePath)
	if err == nil {
		testsTemplate = string(fileContents)
	}

	testsTemplate = testsTemplate + testsTmpl
	return testsTemplate
}

func addTemplateFiles(opts Options, actionsToAdd []meta.Name, data makr.Data) {
	for _, action := range actionsToAdd {
		vg := makr.New()
		viewPath := filepath.Join("templates", fmt.Sprintf("%s", opts.Name.File()), fmt.Sprintf("%s.html", action.File()))
		vg.Add(makr.NewFile(viewPath, viewTmpl))
		vg.Run(".", makr.Data{
			"opts":   opts,
			"action": action.Camel(),
		})
	}
}

func findActionsToAdd(name meta.Name, path string, actions []meta.Name) []meta.Name {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		fileContents = []byte("")
	}

	actionsToAdd := []meta.Name{}

	for _, action := range actions {
		funcSignature := fmt.Sprintf("func %s%s(c buffalo.Context) error", name.Camel(), action.Camel())
		if strings.Contains(string(fileContents), funcSignature) {
			fmt.Printf("--> [warning] skipping %v%v since it already exists\n", name.Camel(), action.Camel())
			continue
		}

		actionsToAdd = append(actionsToAdd, action)
	}

	return actionsToAdd
}

func findHandlersToAdd(name meta.Name, path string, actions []meta.Name) []meta.Name {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		fileContents = []byte("")
	}

	handlersToAdd := []meta.Name{}

	for _, action := range actions {
		funcSignature := fmt.Sprintf("app.GET(\"/%s/%s\", %s%s)", name.URL(), action.URL(), name.Camel(), action.Camel())
		if strings.Contains(string(fileContents), funcSignature) {
			fmt.Printf("--> [warning] skipping %s from app.go since it already exists\n", funcSignature)
			continue
		}

		handlersToAdd = append(handlersToAdd, action)
	}

	return handlersToAdd
}

func findTestsToAdd(name meta.Name, path string, actions []meta.Name) []meta.Name {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		fileContents = []byte("")
	}

	actionsToAdd := []meta.Name{}

	for _, action := range actions {
		funcSignature := fmt.Sprintf("func (as *ActionSuite) Test_%v_%v() {", name.Camel(), action.Camel())
		if strings.Contains(string(fileContents), funcSignature) {
			fmt.Printf("--> [warning] skipping Test_%v_%v since it already exists\n", name.Camel(), action.Camel())
			continue
		}

		actionsToAdd = append(actionsToAdd, action)
	}

	return actionsToAdd
}
