package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/buffalo/generators/newapp"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/plush"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

var rootPath string

// var force bool
// var verbose bool
// var skipPop bool
// var skipWebpack bool
// var withYarn bool
// var dbType = "postgres"
// var ciProvider = "none"

var app = &newapp.App{}

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Creates a new Buffalo application",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !validDbType() {
			return fmt.Errorf("Unknown db-type %s expecting one of postgres, mysql or sqlite3", app.DBType)
		}

		if len(args) == 0 {
			return errors.New("you must enter a name for your new application")
		}

		app.Name = args[0]

		if app.Name == "." {
			app.Name = filepath.Base(app.RootPath)
		} else {
			app.RootPath = filepath.Join(app.RootPath, app.Name)
		}

		err := validateInGoPath()
		if err != nil {
			return err
		}

		s, _ := os.Stat(app.RootPath)
		if s != nil {
			if app.Force {
				os.RemoveAll(app.RootPath)
			} else {
				return fmt.Errorf("%s already exists! Either delete it or use the -f flag to force", app.Name)
			}
		}

		err = genNewFiles()
		if err != nil {
			return err
		}

		fmt.Printf("Congratulations! Your application, %s, has been successfully built!\n\n", app.Name)
		fmt.Println("You can find your new application at:")
		fmt.Println(app.RootPath)
		fmt.Println("\nPlease read the README.md file in your new application for next steps on running your application.")

		return nil
	},
}

func validDbType() bool {
	return app.DBType == "postgres" || app.DBType == "mysql" || app.DBType == "sqlite3"
}

func validateInGoPath() error {
	gp, err := envy.MustGet("GOPATH")
	if err != nil {
		fmt.Println(noGoPath)
		os.Exit(-1)
	}

	var gpMultiple []string

	if runtime.GOOS == "windows" {
		gpMultiple = strings.Split(gp, ";") // Windows uses a different separator
	} else {
		gpMultiple = strings.Split(gp, ":")
	}
	gpMultipleLen := len(gpMultiple)
	foundInPath := false

	for i := 0; i < gpMultipleLen; i++ {
		if strings.HasPrefix(app.RootPath, filepath.Join(gpMultiple[i], "src")) {
			foundInPath = true
			break
		}
	}

	if !foundInPath {
		u, err := user.Current()
		if err != nil {
			return err
		}
		t, err := plush.Render(notInGoWorkspace, plush.NewContextWith(map[string]interface{}{
			"name":     app.Name,
			"gopath":   gp,
			"current":  app.RootPath,
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

func goPath(root string) string {
	var gpMultiple []string
	gp := os.Getenv("GOPATH")

	if runtime.GOOS == "windows" {
		gpMultiple = strings.Split(gp, ";") // Windows uses a different separator
	} else {
		gpMultiple = strings.Split(gp, ":")
	}
	gpMultipleLen := len(gpMultiple)
	path := ""

	for i := 0; i < gpMultipleLen; i++ {
		if strings.HasPrefix(root, filepath.Join(gpMultiple[i], "src")) {
			path = gpMultiple[i]
			break
		}
	}
	return path
}

func packagePath(rootPath string) string {
	gosrcpath := strings.Replace(filepath.Join(goPath(rootPath), "src"), "\\", "/", -1)
	rootPath = strings.Replace(rootPath, "\\", "/", -1)
	return strings.Replace(rootPath, gosrcpath+"/", "", 2)
}

func genNewFiles() error {
	packagePath := packagePath(app.RootPath)

	data := map[string]interface{}{
		"name":        app.Name,
		"titleName":   inflect.Titleize(app.Name),
		"packagePath": packagePath,
		"actionsPath": packagePath + "/actions",
		"modelsPath":  packagePath + "/models",
		"withPop":     !app.SkipPop,
		"withWebpack": !app.SkipWebpack,
		"withYarn":    app.WithYarn,
		"dbType":      app.DBType,
		"version":     Version,
		"ciProvider":  app.CIProvider,
	}

	g, err := app.Generator(data)
	if err != nil {
		return err
	}
	return g.Run(app.RootPath, data)
}

func init() {
	pwd, _ := os.Getwd()

	rootPath = pwd
	app.RootPath = pwd

	RootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVarP(&app.Force, "force", "f", false, "delete and remake if the app already exists")
	newCmd.Flags().BoolVarP(&app.Verbose, "verbose", "v", false, "verbosely print out the go get/install commands")
	newCmd.Flags().BoolVar(&app.SkipPop, "skip-pop", false, "skips adding pop/soda to your app")
	newCmd.Flags().BoolVar(&app.SkipWebpack, "skip-webpack", false, "skips adding Webpack to your app")
	newCmd.Flags().BoolVar(&app.WithYarn, "with-yarn", false, "allows the use of yarn instead of npm as dependency manager")
	newCmd.Flags().StringVar(&app.DBType, "db-type", "postgres", "specify the type of database you want to use [postgres, mysql, sqlite3]")
	newCmd.Flags().StringVar(&app.CIProvider, "ci-provider", "none", "specify the type of ci file you would like buffalo to generate [none, travis]")
}

const notInGoWorkspace = `Oops! It would appear that you are not in your Go Workspace.

Your $GOPATH is set to "<%= gopath %>".

You are currently in "<%= current %>".

The standard location for putting Go projects is something along the lines of "$GOPATH/src/github.com/<%= username %>/<%= name %>" (adjust accordingly).

We recommend you go to "$GOPATH/src/github.com/<%= username %>/" and try "buffalo new <%= name %>" again.`

const noGoPath = `You do not have a $GOPATH set. In order to work with Go, you must set up your $GOPATH and your Go Workspace.

We recommend reading this tutorial on setting everything up: https://www.goinggo.net/2016/05/installing-go-and-your-workspace.html

When you're ready come back and try again. Don't worry, Buffalo will be right here waiting for you. :)`
