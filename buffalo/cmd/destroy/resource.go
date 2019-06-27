package destroy

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/flect"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// YesToAll means not to ask when destroying but simply confirm all beforehand.
var YesToAll = false

// ResourceCmd destroys a passed resource
var ResourceCmd = &cobra.Command{
	Use: "resource [name]",
	// Example: "resource cars",
	Aliases: []string{"r"},
	Short:   "Destroy resource files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you need to provide a valid resource name in order to destroy it")
		}

		name := args[0]
		fileName := flect.Pluralize(flect.Underscore(name))

		removeTemplates(fileName)
		if err := removeActions(fileName); err != nil {
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
	fmt.Print(msg)
	text, _ := reader.ReadString('\n')

	return (text == "y\n" || text == "Y\n")
}

func removeTemplates(fileName string) {
	if YesToAll || confirm("Want to remove templates? (y/N)") {
		templatesFolder := filepath.Join("templates", fileName)
		logrus.Infof("- Deleted %v folder", templatesFolder)
		os.RemoveAll(templatesFolder)
	}
}

func removeActions(fileName string) error {
	if YesToAll || confirm("Want to remove actions? (y/N)") {
		logrus.Infof("- Deleted %v", fmt.Sprintf("actions/%v.go", fileName))
		os.Remove(filepath.Join("actions", fmt.Sprintf("%v.go", fileName)))

		logrus.Infof("- Deleted %v", fmt.Sprintf("actions/%v_test.go", fileName))
		os.Remove(filepath.Join("actions", fmt.Sprintf("%v_test.go", fileName)))

		content, err := ioutil.ReadFile(filepath.Join("actions", "app.go"))
		if err != nil {
			logrus.Warn("error reading app.go content")
			return err
		}

		resourceExpression := fmt.Sprintf("app.Resource(\"/%v\", %vResource{})", fileName, flect.Pascalize(fileName))
		newContents := strings.Replace(string(content), resourceExpression, "", -1)

		err = ioutil.WriteFile(filepath.Join("actions", "app.go"), []byte(newContents), 0)
		if err != nil {
			logrus.Error("error writing new app.go content")
			return err
		}

		logrus.Infof("- Deleted References for %v in actions/app.go", fileName)
	}

	return nil
}

func removeLocales(fileName string) {
	if YesToAll || confirm("Want to remove locales? (y/N)") {
		removeMatch("locales", fmt.Sprintf("%v.*.yaml", fileName))
	}
}

func removeModel(name string) {
	if YesToAll || confirm("Want to remove model? (y/N)") {
		modelFileName := flect.Singularize(flect.Underscore(name))

		os.Remove(filepath.Join("models", fmt.Sprintf("%v.go", modelFileName)))
		os.Remove(filepath.Join("models", fmt.Sprintf("%v_test.go", modelFileName)))

		logrus.Infof("- Deleted %v", fmt.Sprintf("models/%v.go", modelFileName))
		logrus.Infof("- Deleted %v", fmt.Sprintf("models/%v_test.go", modelFileName))
	}
}

func removeMigrations(fileName string) {
	if YesToAll || confirm("Want to remove migrations? (y/N)") {
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
				logrus.Infof("- Deleted %v", path)
			}
		}
	}
}
