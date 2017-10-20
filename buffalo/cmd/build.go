package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo/buffalo/cmd/build"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/markbates/sigtx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var options = build.Options{}
var tags = ""
var debug bool

var xbuildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b", "bill"},
	Short:   "Builds a Buffalo binary, including bundling of assets (packr & webpack)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt)
		defer cancel()

		if options.Debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		b := build.New(ctx, options)
		if tags != "" {
			b.Tags = append(b.Tags, tags)
		}

		go func() {
			<-ctx.Done()
			if ctx.Err() == context.Canceled {
				fmt.Println("~~~ BUILD CANCELLED ~~~")
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

		fmt.Printf("\nYou application was successfully built at %s\n", filepath.Join(b.Root, b.Bin))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(xbuildCmd)

	pwd, _ := os.Getwd()

	options.App = meta.New(pwd)

	xbuildCmd.Flags().StringVarP(&options.Bin, "output", "o", options.Bin, "set the name of the binary")
	xbuildCmd.Flags().StringVarP(&tags, "tags", "t", "", "compile with specific build tags")
	xbuildCmd.Flags().BoolVarP(&options.ExtractAssets, "extract-assets", "e", false, "extract the assets and put them in a distinct archive")
	xbuildCmd.Flags().BoolVarP(&options.SkipAssets, "skip-assets", "k", false, "avoids compiling the assets, usefull if app was generated with --api")
	xbuildCmd.Flags().BoolVarP(&options.Static, "static", "s", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"'")
	xbuildCmd.Flags().StringVar(&options.LDFlags, "ldflags", "", "set any ldflags to be passed to the go build")
	xbuildCmd.Flags().BoolVarP(&options.Debug, "debug", "d", false, "print debugging information")
	xbuildCmd.Flags().BoolVarP(&options.Compress, "compress", "c", true, "compress static files in the binary")

}
