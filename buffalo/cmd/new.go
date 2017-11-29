package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo/generators/newapp"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/makr"
	"github.com/gobuffalo/plush"
	"github.com/spf13/cobra"
)

var rootPath string

var app = newapp.Generator{
	App:        meta.New("."),
	DBType:     "postgres",
	CIProvider: "none",
	AsWeb:      true,
	Docker:     "multi",
	VCS:        "git",
}

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Creates a new Buffalo application",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			return errors.New("you must enter a name for your new application")
		}

		app.Name = meta.Name(args[0])
		app.Version = Version

		if app.Name == "." {
			app.Name = meta.Name(filepath.Base(app.Root))
		} else {
			app.Root = filepath.Join(app.Root, app.Name.File())
		}
		aa := meta.New(app.Root)
		app.ActionsPkg = aa.ActionsPkg
		app.GriftsPkg = aa.GriftsPkg
		app.ModelsPkg = aa.ModelsPkg
		app.PackagePkg = aa.PackagePkg

		if err := app.Validate(); err != nil {
			if errors.Cause(err) == newapp.ErrNotInGoPath {
				return notInGoPath(app)
			}
			return errors.WithStack(err)
		}

		app.WithPop = !app.SkipPop
		app.WithWebpack = !app.SkipWebpack
		app.WithYarn = !app.SkipYarn
		app.AsWeb = !app.AsAPI

		if err := app.Run(app.Root, makr.Data{}); err != nil {
			return errors.WithStack(err)
		}

		fmt.Printf("Congratulations! Your application, %s, has been successfully built!\n\n", app.Name)
		fmt.Println("You can find your new application at:")
		fmt.Println(app.Root)
		fmt.Println("\nPlease read the README.md file in your new application for next steps on running your application.")

		return nil
	},
}

func notInGoPath(ag newapp.Generator) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	pwd, _ := os.Getwd()
	t, err := plush.Render(notInGoWorkspace, plush.NewContextWith(map[string]interface{}{
		"name":     ag.Name,
		"gopath":   envy.GoPath(),
		"current":  pwd,
		"username": u.Username,
	}))
	if err != nil {
		return err
	}
	fmt.Println(t)
	os.Exit(-1)
	return nil
}

func init() {
	pwd, _ := os.Getwd()

	app.App = meta.New(pwd)

	decorate("new", newCmd)
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVar(&app.AsAPI, "api", false, "skip all front-end code and configure for an API server")
	newCmd.Flags().BoolVarP(&app.Force, "force", "f", false, "delete and remake if the app already exists")
	newCmd.Flags().BoolVarP(&app.Verbose, "verbose", "v", false, "verbosely print out the go get commands")
	newCmd.Flags().BoolVar(&app.SkipPop, "skip-pop", false, "skips adding pop/soda to your app")
	newCmd.Flags().BoolVar(&app.WithDep, "with-dep", false, "adds github.com/golang/dep to your app")
	newCmd.Flags().BoolVar(&app.SkipWebpack, "skip-webpack", false, "skips adding Webpack to your app")
	newCmd.Flags().BoolVar(&app.SkipYarn, "skip-yarn", false, "use npm instead of yarn for frontend dependencies management")
	newCmd.Flags().StringVar(&app.DBType, "db-type", "postgres", "specify the type of database you want to use [postgres, mysql, sqlite3]")
	newCmd.Flags().StringVar(&app.Docker, "docker", "multi", "specify the type of Docker file to generate [none, multi, standard]")
	newCmd.Flags().StringVar(&app.CIProvider, "ci-provider", "none", "specify the type of ci file you would like buffalo to generate [none, travis, gitlab-ci]")
	newCmd.Flags().StringVar(&app.VCS, "vcs", "git", "specify the Version control system you would like to use [none, git, bzr]")
}

const notInGoWorkspace = `Oops! It would appear that you are not in your Go Workspace.

Your $GOPATH is set to "<%= gopath %>".

You are currently in "<%= current %>".

The standard location for putting Go projects is something along the lines of "$GOPATH/src/github.com/<%= username %>/<%= name %>" (adjust accordingly).

We recommend you go to "$GOPATH/src/github.com/<%= username %>/" and try "buffalo new <%= name %>" again.`
