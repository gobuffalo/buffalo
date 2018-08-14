package cmd

import (
	"context"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/buffalo/buffalo/cmd/build"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/markbates/sigtx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var buildOptions = struct {
	build.Options
	SkipAssets bool
	Tags       string
}{
	Options: build.Options{},
}

func bindBuildFlags(fs *pflag.FlagSet) error {
	pwd, _ := os.Getwd()
	buildOptions.App = meta.New(pwd)

	var err error
	// Since the default value for "output" is computed on init,
	// if we change the current directory during runtime (e.g. when
	// running tests), the init value will still be there.
	// Ensure the flag was set manually before using it blindly.
	if fs.Changed("output") {
		buildOptions.Bin, err = fs.GetString("output")
		if err != nil {
			return err
		}
	}
	buildOptions.Tags, err = fs.GetString("tags")
	if err != nil {
		return err
	}
	buildOptions.ExtractAssets, err = fs.GetBool("extract-assets")
	if err != nil {
		return err
	}
	buildOptions.SkipAssets, err = fs.GetBool("skip-assets")
	if err != nil {
		return err
	}
	buildOptions.Static, err = fs.GetBool("static")
	if err != nil {
		return err
	}
	buildOptions.LDFlags, err = fs.GetString("ldflags")
	if err != nil {
		return err
	}
	buildOptions.Debug, err = fs.GetBool("debug")
	if err != nil {
		return err
	}
	buildOptions.Compress, err = fs.GetBool("compress")
	if err != nil {
		return err
	}
	buildOptions.SkipTemplateValidation, err = fs.GetBool("skip-template-validation")
	if err != nil {
		return err
	}
	buildOptions.Environment, err = fs.GetString("environment")
	return err
}

var xbuildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b", "bill"},
	Short:   "Builds a Buffalo binary, including bundling of assets (packr & webpack)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt)
		defer cancel()

		if err := bindBuildFlags(cmd.Flags()); err != nil {
			return err
		}

		buildOptions.Options.WithAssets = !buildOptions.SkipAssets

		if buildOptions.Debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		b := build.New(ctx, buildOptions.Options)
		if buildOptions.Tags != "" {
			b.Tags = append(b.Tags, buildOptions.Tags)
		}

		go func() {
			<-ctx.Done()
			if ctx.Err() == context.Canceled {
				logrus.Info("~~~ BUILD CANCELLED ~~~")
				err := b.Cleanup()
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}()

		err := b.Run()
		if err != nil {
			return errors.WithStack(err)
		}

		bin, _ := filepath.Abs(b.Bin)
		logrus.Infof("\nYour application was successfully built at %s\n", bin)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(xbuildCmd)

	// Ensure the default output is correct
	pwd, _ := os.Getwd()
	bin := filepath.Join("bin", filepath.Base(pwd))
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	xbuildCmd.Flags().StringP("output", "o", bin, "set the name of the binary")
	xbuildCmd.Flags().StringP("tags", "t", "", "compile with specific build tags")
	xbuildCmd.Flags().BoolP("extract-assets", "e", false, "extract the assets and put them in a distinct archive")
	xbuildCmd.Flags().BoolP("skip-assets", "k", false, "skip running webpack and building assets")
	xbuildCmd.Flags().BoolP("static", "s", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"'")
	xbuildCmd.Flags().String("ldflags", "", "set any ldflags to be passed to the go build")
	xbuildCmd.Flags().BoolP("debug", "d", false, "print debugging information")
	xbuildCmd.Flags().BoolP("compress", "c", true, "compress static files in the binary")
	xbuildCmd.Flags().Bool("skip-template-validation", false, "skip validating plush templates")
	xbuildCmd.Flags().StringP("environment", "", "development", "set the environment for the binary")
}
