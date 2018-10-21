package cmd

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo/genny/build"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/meta"
	"github.com/markbates/sigtx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var buildOptions = struct {
	*build.Options
	SkipAssets             bool
	Debug                  bool
	Tags                   string
	SkipTemplateValidation bool
	DryRun                 bool
	Verbose                bool
}{
	Options: &build.Options{},
}

var xbuildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b", "bill"},
	Short:   "Build the application binary, including bundling of assets (packr & webpack)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt)
		defer cancel()

		buildOptions.Options.WithAssets = !buildOptions.SkipAssets

		run := genny.WetRunner(ctx)
		if buildOptions.DryRun {
			run = genny.DryRunner(ctx)
		}

		if buildOptions.Verbose {
			l := logrus.New()
			l.SetLevel(logrus.DebugLevel)
			run.Logger = l
			buildOptions.BuildFlags = append(buildOptions.BuildFlags, "-v")
		}

		opts := buildOptions.Options

		if buildOptions.Tags != "" {
			opts.Tags = append(opts.Tags, buildOptions.Tags)
		}

		if !buildOptions.SkipTemplateValidation {
			opts.TemplateValidators = append(opts.TemplateValidators, build.PlushValidator, build.GoTemplateValidator)
		}

		run.WithNew(build.New(opts))
		return run.Run()
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
	xbuildCmd.Flags().BoolVarP(&buildOptions.Verbose, "verbose", "v", false, "print debugging information")
	xbuildCmd.Flags().BoolVarP(&buildOptions.DryRun, "dry-run", "d", false, "runs the build 'dry'")
	xbuildCmd.Flags().BoolVar(&buildOptions.SkipTemplateValidation, "skip-template-validation", false, "skip validating plush templates")
	xbuildCmd.Flags().StringVarP(&buildOptions.Environment, "environment", "", "development", "set the environment for the binary")
}
