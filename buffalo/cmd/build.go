package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo/buffalo/cmd/build"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/markbates/sigtx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var buildOptions = struct {
	build.Options
	SkipAssets bool
	Tags       string
}{
	Options: build.Options{},
}

var xbuildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b", "bill"},
	Short:   "Builds a Buffalo binary, including bundling of assets (packr & webpack)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt)
		defer cancel()

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

		logrus.Infof("\nYou application was successfully built at %s\n", filepath.Join(b.Root, b.Bin))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(xbuildCmd)

	pwd, _ := os.Getwd()

	buildOptions.App = meta.New(pwd)

	xbuildCmd.Flags().StringVarP(&buildOptions.Bin, "output", "o", buildOptions.Bin, "set the name of the binary")
	xbuildCmd.Flags().StringVarP(&buildOptions.Tags, "tags", "t", "", "compile with specific build tags")
	xbuildCmd.Flags().BoolVarP(&buildOptions.ExtractAssets, "extract-assets", "e", false, "extract the assets and put them in a distinct archive")
	xbuildCmd.Flags().BoolVarP(&buildOptions.SkipAssets, "skip-assets", "k", false, "skip running webpack and building assets")
	xbuildCmd.Flags().BoolVarP(&buildOptions.Static, "static", "s", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"'")
	xbuildCmd.Flags().StringVar(&buildOptions.LDFlags, "ldflags", "", "set any ldflags to be passed to the go build")
	xbuildCmd.Flags().BoolVarP(&buildOptions.Debug, "debug", "d", false, "print debugging information")
	xbuildCmd.Flags().BoolVarP(&buildOptions.Compress, "compress", "c", true, "compress static files in the binary")

}
