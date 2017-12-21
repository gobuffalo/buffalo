package action

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/makr"
)

// Run action generator
func (act Generator) Run(root string, data makr.Data) error {
	g := makr.New()
	defer g.Fmt(root)

	filePath := filepath.Join("actions", fmt.Sprintf("%v.go", act.Name.File()))
	actionsTemplate := act.buildActionsTemplate(filePath)
	testFilePath := filepath.Join("actions", fmt.Sprintf("%v_test.go", act.Name.File()))
	testsTemplate := act.buildTestsTemplate(testFilePath)
	actionsToAdd := act.findActionsToAdd(filePath)
	testsToAdd := act.findTestsToAdd(testFilePath)
	handlersToAdd := act.findHandlersToAdd(filepath.Join("actions", "app.go"))

	data["opts"] = act
	data["actions"] = actionsToAdd
	data["tests"] = testsToAdd

	g.Add(makr.NewFile(filePath, actionsTemplate))
	g.Add(makr.NewFile(testFilePath, testsTemplate))
	g.Add(&makr.Func{
		Should: func(data makr.Data) bool { return true },
		Runner: func(root string, data makr.Data) error {
			routes := []string{}
			for _, a := range handlersToAdd {
				routes = append(routes, fmt.Sprintf("app.%s(\"/%s/%s\", %s)", strings.ToUpper(act.Method), act.Name, a, act.Name.Camel()+a.Camel()))
			}
			return generators.AddInsideAppBlock(routes...)
		},
	})
	if !act.SkipTemplate {
		if err := act.addTemplateFiles(actionsToAdd, data); err != nil {
			return errors.WithStack(err)
		}
	}
	return g.Run(root, data)
}

func (act Generator) buildActionsTemplate(filePath string) string {
	actionsTemplate := actionsHeaderTmpl
	fileContents, err := ioutil.ReadFile(filePath)
	if err == nil {
		actionsTemplate = string(fileContents)
	}

	actionsTemplate = actionsTemplate + actionsTmpl
	return actionsTemplate
}

func (act Generator) buildTestsTemplate(filePath string) string {
	testsTemplate := testHeaderTmpl

	fileContents, err := ioutil.ReadFile(filePath)
	if err == nil {
		testsTemplate = string(fileContents)
	}

	testsTemplate = testsTemplate + testsTmpl
	return testsTemplate
}

func (act Generator) addTemplateFiles(actionsToAdd []meta.Name, data makr.Data) error {
	for _, action := range actionsToAdd {
		vg := makr.New()
		viewPath := filepath.Join("templates", fmt.Sprintf("%s", act.Name.File()), fmt.Sprintf("%s.html", action.File()))
		vg.Add(makr.NewFile(viewPath, viewTmpl))
		err := vg.Run(".", makr.Data{
			"opts":   act,
			"action": action.Camel(),
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (act Generator) findActionsToAdd(path string) []meta.Name {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		fileContents = []byte("")
	}

	actionsToAdd := []meta.Name{}

	for _, action := range act.Actions {
		funcSignature := fmt.Sprintf("func %s%s(c buffalo.Context) error", act.Name.Camel(), action.Camel())
		if strings.Contains(string(fileContents), funcSignature) {
			logrus.Warnf("--> skipping %v%v since it already exists\n", act.Name.Camel(), action.Camel())
			continue
		}

		actionsToAdd = append(actionsToAdd, action)
	}

	return actionsToAdd
}

func (act Generator) findHandlersToAdd(path string) []meta.Name {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		fileContents = []byte("")
	}

	handlersToAdd := []meta.Name{}

	for _, action := range act.Actions {
		funcSignature := fmt.Sprintf("app.GET(\"/%s/%s\", %s%s)", act.Name.URL(), action.URL(), act.Name.Camel(), action.Camel())
		if strings.Contains(string(fileContents), funcSignature) {
			logrus.Warnf("--> skipping %s from app.go since it already exists\n", funcSignature)
			continue
		}

		handlersToAdd = append(handlersToAdd, action)
	}

	return handlersToAdd
}

func (act Generator) findTestsToAdd(path string) []meta.Name {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		fileContents = []byte("")
	}

	actionsToAdd := []meta.Name{}

	for _, action := range act.Actions {
		funcSignature := fmt.Sprintf("func (as *ActionSuite) Test_%v_%v() {", act.Name.Camel(), action.Camel())
		if strings.Contains(string(fileContents), funcSignature) {
			logrus.Warnf("--> skipping Test_%v_%v since it already exists\n", act.Name.Camel(), action.Camel())
			continue
		}

		actionsToAdd = append(actionsToAdd, action)
	}

	return actionsToAdd
}
