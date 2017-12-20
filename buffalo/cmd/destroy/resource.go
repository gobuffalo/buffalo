package destroy

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/markbates/inflect"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//YesToAll means not to ask when destroying but simply confirm all beforehand.
var YesToAll = false

//ResourceCmd destroys a passed resource
var ResourceCmd = &cobra.Command{
	Use: "resource [name]",
	//Example: "resource cars",
	Aliases: []string{"r"},
	Short:   "Destroys resource files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you need to provide a valid resource name in order to destroy it")
		}

		name := args[0]
		fileName := inflect.Pluralize(inflect.Underscore(name))

		removeTemplates(fileName)
		err := removeActions(fileName)
		if err != nil {
			return err
		}

		removeLocales(fileName)
		removeModel(name)
		removeMigrations(fileName)

		return nil
	},
}

func confirm(msg string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(msg)
	text, _ := reader.ReadString('\n')

	return (text == "y\n" || text == "Y\n")
}

func removeTemplates(fileName string) {
	if YesToAll || confirm("Want to remove templates? (Y/n)") {
		templatesFolder := fmt.Sprintf(filepath.Join("templates", fileName))
		logrus.Infof("- Deleted %v folder\n", templatesFolder)
		os.RemoveAll(templatesFolder)
	}
}

func removeActions(fileName string) error {
	if YesToAll || confirm("Want to remove actions? (Y/n)") {
		logrus.Infof("- Deleted %v\n", fmt.Sprintf("actions/%v.go", fileName))
		os.Remove(filepath.Join("actions", fmt.Sprintf("%v.go", fileName)))

		logrus.Infof("- Deleted %v\n", fmt.Sprintf("actions/%v_test.go", fileName))
		os.Remove(filepath.Join("actions", fmt.Sprintf("%v_test.go", fileName)))

		content, err := ioutil.ReadFile(filepath.Join("actions", "app.go"))
		if err != nil {
			logrus.Warn("error reading app.go content")
			return err
		}

		resourceExpression := fmt.Sprintf("app.Resource(\"/%v\", %vResource{})", fileName, inflect.Camelize(fileName))
		newContents := strings.Replace(string(content), resourceExpression, "", -1)

		err = ioutil.WriteFile(filepath.Join("actions", "app.go"), []byte(newContents), 0)
		if err != nil {
			logrus.Error("error writing new app.go content")
			return err
		}

		logrus.Infof("- Deleted References for %v in actions/app.go\n", fileName)
	}

	return nil
}

func removeLocales(fileName string) {
	if YesToAll || confirm("Want to remove locales? (Y/n)") {
		removeMatch("locales", fmt.Sprintf("%v.*.yaml", fileName))
	}
}

func removeModel(name string) {
	if YesToAll || confirm("Want to remove model? (Y/n)") {
		modelFileName := inflect.Singularize(inflect.Underscore(name))

		os.Remove(filepath.Join("models", fmt.Sprintf("%v.go", modelFileName)))
		os.Remove(filepath.Join("models", fmt.Sprintf("%v_test.go", modelFileName)))

		logrus.Infof("- Deleted %v\n", fmt.Sprintf("models/%v.go", modelFileName))
		logrus.Infof("- Deleted %v\n", fmt.Sprintf("models/%v_test.go", modelFileName))
	}
}

func removeMigrations(fileName string) {
	if YesToAll || confirm("Want to remove migrations? (Y/n)") {
		removeMatch("migrations", fmt.Sprintf("*_create_%v.up.*", fileName))
		removeMatch("migrations", fmt.Sprintf("*_create_%v.down.*", fileName))
	}
}

func removeMatch(folder, pattern string) {
	files, err := ioutil.ReadDir(folder)
	if err == nil {
		for _, f := range files {
			matches, _ := filepath.Match(pattern, f.Name())
			if !f.IsDir() && matches {
				path := filepath.Join(folder, f.Name())
				os.Remove(path)
				logrus.Infof("- Deleted %v\n", path)
			}
		}
	}
}
