package generate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aymerick/raymond"
	"github.com/markbates/gentronics"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

var runningTests = false

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

		g := gentronics.New()

		if !fileExists {
			g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s.go", data["under"])), rActionFileT))
			g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s_test.go", data["under"])), rActionTest))
			if err = g.Run(".", data); err != nil {
				return err
			}
		}

		for _, action := range actions {
			g.Add(buildActionAppender(name, action))
		}

		if !runningTests {
			g.Add(Fmt)
		}

		return g.Run(".", gentronics.Data{})
	},
}

func buildActionAppender(namespace, action string) gentronics.Runnable {
	aa := actionAppender{namespace, action, inflect.Underscore(namespace)}
	return aa
}

type actionAppender struct {
	Namespace  string
	ActionName string
	FileName   string
}

func (aa actionAppender) Run(rootPath string, data gentronics.Data) error {
	path := filepath.Join("actions", fmt.Sprintf("%v.go", aa.FileName))
	fileContents, _ := ioutil.ReadFile(path)

	funcSignature := fmt.Sprintf("func %s%s(c buffalo.Context) error", inflect.Camelize(aa.Namespace), inflect.Camelize(aa.ActionName))
	if strings.Contains(string(fileContents), funcSignature) {
		fmt.Printf("--> [warning] skipping %v%v since it already exists\n", inflect.Camelize(aa.Namespace), inflect.Camelize(aa.ActionName))
		return nil
	}

	templateData := map[string]string{
		"namespace":       inflect.Camelize(aa.Namespace),
		"action":          inflect.Camelize(aa.ActionName),
		"namespace_under": inflect.Underscore(aa.Namespace),
		"action_under":    inflect.Underscore(aa.ActionName),
	}

	t, _ := raymond.Parse(rActionFuncT)
	fn, _ := t.Exec(templateData)

	fileContents = []byte(string(fileContents) + fn)
	err := ioutil.WriteFile(path, fileContents, 0755)

	if err != nil {
		return err
	}

	t, _ = raymond.Parse(rViewT)
	content, _ := t.Exec(templateData)

	templatePath := filepath.Join("templates", fmt.Sprintf("%s", inflect.Underscore(aa.Namespace)), fmt.Sprintf("%s.html", inflect.Underscore(aa.ActionName)))
	os.MkdirAll(filepath.Join("templates", fmt.Sprintf("%s", inflect.Underscore(aa.Namespace))), 0755)
	return ioutil.WriteFile(templatePath, []byte(content), 0755)
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

	rActionTestT = `
    package actions
    
    func Test_{{namespace}}{{action}}(t *testing.T) {
	    r := require.New(t)
	    r.Fail("Not Implemented!")
    }
    `
)
