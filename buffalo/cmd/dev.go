package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/gobuffalo/buffalo/genny/assets/webpack"
	rg "github.com/gobuffalo/buffalo/genny/refresh"
	"github.com/gobuffalo/events"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/meta"
	"github.com/markbates/refresh/refresh"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	events.NamedListen("buffalo:dev", func(e events.Event) {
		if strings.HasPrefix(e.Kind, "refresh:") {
			e.Kind = strings.Replace(e.Kind, "refresh:", "buffalo:dev:", 1)
			events.Emit(e)
		}
	})
}

var devOptions = struct {
	Debug bool
}{}

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run the Buffalo app in 'development' mode",
	Long: `Run the Buffalo app in 'development' mode.
This includes rebuilding the application when files change.
This behavior can be changed in .buffalo.dev.yml file.`,
	RunE: func(c *cobra.Command, args []string) error {
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
			logrus.Errorf(msg, cause)
		}()
		os.Setenv("GO_ENV", "development")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg, ctx := errgroup.WithContext(ctx)

		wg.Go(func() error {
			return startDevServer(ctx, args)
		})

		wg.Go(func() error {
			app := meta.New(".")
			if !app.WithNodeJs {
				// No need to run dev script
				return nil
			}
			return runDevScript(ctx, app)
		})

		err := wg.Wait()
		if err != context.Canceled {
			return err
		}
		return nil
	},
}

func runDevScript(ctx context.Context, app meta.App) error {
	tool := "yarnpkg"
	if !app.WithYarn {
		tool = "npm"
	}

	if _, err := exec.LookPath(tool); err != nil {
		return fmt.Errorf("could not find %s tool", tool)
	}

	// make sure that the node_modules folder is properly "installed"
	if _, err := os.Stat(filepath.Join(app.Root, "node_modules")); err != nil {
		cmd := exec.CommandContext(ctx, tool, "install")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	cmd := exec.CommandContext(ctx, tool, "run", "dev")
	if _, err := app.NodeScript("dev"); err != nil {
		// Fallback on legacy runner
		cmd = exec.CommandContext(ctx, webpack.BinPath, "--watch")
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func startDevServer(ctx context.Context, args []string) error {
	app := meta.New(".")

	cfgFile := "./.buffalo.dev.yml"
	if _, err := os.Stat(cfgFile); err != nil {
		run := genny.WetRunner(ctx)
		err = run.WithNew(rg.New(&rg.Options{App: app}))
		if err != nil {
			return err
		}

		if err := run.Run(); err != nil {
			return err
		}
	}
	c := &refresh.Configuration{}
	if err := c.Load(cfgFile); err != nil {
		return err
	}
	c.Debug = devOptions.Debug

	bt := app.BuildTags("development")
	if len(bt) > 0 {
		c.BuildFlags = append(c.BuildFlags, "-tags", bt.String())
	}
	r := refresh.NewWithContext(c, ctx)
	r.CommandFlags = args
	return r.Start()
}

func init() {
	devCmd.Flags().BoolVarP(&devOptions.Debug, "debug", "d", false, "use delve to debug the app")
	decorate("dev", devCmd)
	RootCmd.AddCommand(devCmd)
}
