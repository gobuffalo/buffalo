// +build !appengine

package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	rg "github.com/gobuffalo/buffalo/generators/refresh"
	"github.com/markbates/refresh/refresh"
	"github.com/spf13/cobra"
)

var devOptions = struct {
	Debug bool
}{}

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs your Buffalo app in 'development' mode",
	Long: `Runs your Buffalo app in 'development' mode.
This includes rebuilding your application when files change.
This behavior can be changed in your .buffalo.dev.yml file.`,
	Run: func(c *cobra.Command, args []string) {
		if runtime.GOOS == "windows" {
			color.NoColor = true
		}
		defer func() {
			msg := "There was a problem starting the dev server, Please review the troubleshooting docs: %s\n"
			cause := "Unknown"
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					cause = err.Error()
				}
			}
			fmt.Printf(msg, cause)
		}()
		os.Setenv("GO_ENV", "development")
		ctx := context.Background()
		ctx, cancelFunc := context.WithCancel(ctx)
		go func() {
			err := startDevServer(ctx)
			if err != nil {
				cancelFunc()
				log.Fatal(err)
			}
		}()
		go func() {
			err := startWebpack(ctx)
			if err != nil {
				cancelFunc()
				log.Fatal(err)
			}
		}()
		// wait for the ctx to finish
		<-ctx.Done()
	},
}

func startWebpack(ctx context.Context) error {
	cfgFile := "./webpack.config.js"
	_, err := os.Stat(cfgFile)
	if err != nil {
		// there's no webpack, so don't do anything
		return nil
	}
	cmd := exec.Command(webpack.BinPath, "--watch")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func startDevServer(ctx context.Context) error {
	cfgFile := "./.buffalo.dev.yml"
	_, err := os.Stat(cfgFile)
	if err != nil {
		g, err := rg.New()
		if err != nil {
			return err
		}
		err = g.Run("./", map[string]interface{}{
			"name": "buffalo",
		})
	}
	c := &refresh.Configuration{}
	err = c.Load(cfgFile)
	if err != nil {
		return err
	}
	c.Debug = devOptions.Debug
	r := refresh.NewWithContext(c, ctx)
	return r.Start()
}

func init() {
	devCmd.Flags().BoolVarP(&devOptions.Debug, "debug", "d", false, "use delve to debug the app")
	decorate("dev", devCmd)
	RootCmd.AddCommand(devCmd)
}
