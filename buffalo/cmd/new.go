package cmd

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

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
	Bootstrap:  3,
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

		logrus.Infof("Congratulations! Your application, %s, has been successfully built!\n\n", app.Name)
		logrus.Infof("You can find your new application at:\n%v", app.Root)
		logrus.Info("\nPlease read the README.md file in your new application for next steps on running your application.")

		return nil
	},
}

func currentUser() (string, error) {
	if _, err := exec.LookPath("git"); err == nil {
		if b, err := exec.Command("git", "config", "github.user").Output(); err != nil {
			return string(b), nil
		}
	}
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	username := u.Username
	if t := strings.Split(username, `\`); len(t) > 0 {
		username = t[len(t)-1]
	}
	return username, nil
}

func notInGoPath(ag newapp.Generator) error {
	username, err := currentUser()
	if err != nil {
		return errors.WithStack(err)
	}
	pwd, _ := os.Getwd()
	t, err := plush.Render(notInGoWorkspace, plush.NewContextWith(map[string]interface{}{
		"name":     ag.Name,
		"gopath":   envy.GoPath(),
		"current":  pwd,
		"username": username,
	}))
	if err != nil {
		return err
	}
	logrus.Error(t)
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
	newCmd.Flags().IntVar(&app.Bootstrap, "bootstrap", app.Bootstrap, "specify version for Bootstrap [3, 4]")
}

const notInGoWorkspace = `Oops! It would appear that you are not in your Go Workspace.

Your $GOPATH is set to "<%= gopath %>".

You are currently in "<%= current %>".

The standard location for putting Go projects is something along the lines of "$GOPATH/src/github.com/<%= username %>/<%= name %>" (adjust accordingly).

We recommend you go to "$GOPATH/src/github.com/<%= username %>/" and try "buffalo new <%= name %>" again.`
