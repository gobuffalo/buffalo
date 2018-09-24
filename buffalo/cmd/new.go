package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/pflag"

	"github.com/markbates/inflect"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo/generators/newapp"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/makr"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configError error

func getAppWithConfig() newapp.Generator {
	pwd, _ := os.Getwd()
	app := newapp.Generator{
		App:         meta.New(pwd),
		AsAPI:       viper.GetBool("api"),
		Force:       viper.GetBool("force"),
		Verbose:     viper.GetBool("verbose"),
		SkipPop:     viper.GetBool("skip-pop"),
		SkipWebpack: viper.GetBool("skip-webpack"),
		SkipYarn:    viper.GetBool("skip-yarn"),
		DBType:      viper.GetString("db-type"),
		CIProvider:  viper.GetString("ci-provider"),
		AsWeb:       true,
		Docker:      viper.GetString("docker"),
		Bootstrap:   viper.GetInt("bootstrap"),
	}
	app.VCS = viper.GetString("vcs")
	app.WithDep = viper.GetBool("with-dep")
	app.WithPop = !app.SkipPop
	app.WithWebpack = !app.SkipWebpack
	app.WithYarn = !app.SkipYarn
	app.AsWeb = !app.AsAPI
	if app.AsAPI {
		app.WithWebpack = false
	}

	return app
}

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Creates a new Buffalo application",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Restore default values after usage (useful for testing)
		defer func() {
			cmd.Flags().Visit(func(f *pflag.Flag) {
				f.Value.Set(f.DefValue)
			})
			viper.BindPFlags(cmd.Flags())
		}()
		if len(args) == 0 {
			return errors.New("you must enter a name for your new application")
		}
		if configError != nil {
			return configError
		}
		app := getAppWithConfig()
		app.Name = inflect.Name(args[0])

		if app.Name == "." {
			app.Name = inflect.Name(filepath.Base(app.Root))
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

		data := makr.Data{
			"version": runtime.Version,
			"db-type": viper.GetString("db-type"),
		}
		if err := app.Run(app.Root, data); err != nil {
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
		if b, err := exec.Command("git", "config", "github.user").Output(); err == nil {
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
	decorate("new", newCmd)
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().Bool("api", false, "skip all front-end code and configure for an API server")
	newCmd.Flags().BoolP("force", "f", false, "delete and remake if the app already exists")
	newCmd.Flags().BoolP("verbose", "v", false, "verbosely print out the go get commands")
	newCmd.Flags().Bool("skip-pop", false, "skips adding pop/soda to your app")
	newCmd.Flags().Bool("with-dep", false, "adds github.com/golang/dep to your app")
	newCmd.Flags().Bool("skip-webpack", false, "skips adding Webpack to your app")
	newCmd.Flags().Bool("skip-yarn", false, "use npm instead of yarn for frontend dependencies management")
	newCmd.Flags().String("db-type", "postgres", fmt.Sprintf("specify the type of database you want to use [%s]", strings.Join(pop.AvailableDialects, ", ")))
	newCmd.Flags().String("docker", "multi", "specify the type of Docker file to generate [none, multi, standard]")
	newCmd.Flags().String("ci-provider", "none", "specify the type of ci file you would like buffalo to generate [none, travis, gitlab-ci]")
	newCmd.Flags().String("vcs", "git", "specify the Version control system you would like to use [none, git, bzr]")
	newCmd.Flags().Int("bootstrap", 4, "specify version for Bootstrap [3, 4]")
	viper.BindPFlags(newCmd.Flags())
	cfgFile := newCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.buffalo.yaml)")
	skipConfig := newCmd.Flags().Bool("skip-config", false, "skips using the config file")
	cobra.OnInitialize(initConfig(skipConfig, cfgFile))
}

func initConfig(skipConfig *bool, cfgFile *string) func() {
	return func() {
		if *skipConfig {
			return
		}

		var err error
		if *cfgFile != "" { // enable ability to specify config file via flag
			viper.SetConfigFile(*cfgFile)
			// Will error only if the --config flag is used
			if err = viper.ReadInConfig(); err != nil {
				configError = err
			}
		} else {
			viper.SetConfigName(".buffalo") // name of config file (without extension)
			viper.AddConfigPath("$HOME")    // adding home directory as first search path
			viper.AutomaticEnv()            // read in environment variables that match
			viper.ReadInConfig()
		}

	}
}

const notInGoWorkspace = `Oops! It would appear that you are not in your Go Workspace.

Your $GOPATH is set to "<%= gopath %>".

You are currently in "<%= current %>".

The standard location for putting Go projects is something along the lines of "$GOPATH/src/github.com/<%= username %>/<%= name %>" (adjust accordingly).

We recommend you go to "$GOPATH/src/github.com/<%= username %>/" and try "buffalo new <%= name %>" again.`
