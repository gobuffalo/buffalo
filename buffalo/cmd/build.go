package cmd

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/genny/build"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/meta"
	"github.com/markbates/sigtx"
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
	bin                    string
}{
	Options: &build.Options{
		BuildTime: time.Now(),
	},
}

var xbuildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b", "bill", "install"},
	Short:   "Build the application binary, including bundling of assets (packr & webpack)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt)
		defer cancel()

		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		buildOptions.App = meta.New(pwd)
		if len(buildOptions.bin) > 0 {
			buildOptions.App.Bin = buildOptions.bin
		}

		buildOptions.Options.WithAssets = !buildOptions.SkipAssets

		run := genny.WetRunner(ctx)
		if buildOptions.DryRun {
			run = genny.DryRunner(ctx)
		}

		if buildOptions.Verbose || buildOptions.Debug {
			lg := logger.New(logger.DebugLevel)
			run.Logger = lg
			// plog.Logger = lg
			buildOptions.BuildFlags = append(buildOptions.BuildFlags, "-v")
		}

		opts := buildOptions.Options
		opts.BuildVersion = buildVersion(opts.BuildTime.Format(time.RFC3339))

		if buildOptions.Tags != "" {
			opts.Tags = append(opts.Tags, buildOptions.Tags)
		}

		if !buildOptions.SkipTemplateValidation {
			opts.TemplateValidators = append(opts.TemplateValidators, build.PlushValidator, build.GoTemplateValidator)
		}

		if cmd.CalledAs() == "install" {
			opts.GoCommand = "install"
		}
		clean := build.Cleanup(opts)
		// defer clean(run)
		defer func() {
			if err := clean(run); err != nil {
				log.Fatal("build:clean", err)
			}
		}()
		if err := run.WithNew(build.New(opts)); err != nil {
			return err
		}
		return run.Run()
	},
}

func init() {
	RootCmd.AddCommand(xbuildCmd)

	xbuildCmd.Flags().StringVarP(&buildOptions.bin, "output", "o", buildOptions.Bin, "set the name of the binary")
	xbuildCmd.Flags().StringVarP(&buildOptions.Tags, "tags", "t", "", "compile with specific build tags")
	xbuildCmd.Flags().BoolVarP(&buildOptions.ExtractAssets, "extract-assets", "e", false, "extract the assets and put them in a distinct archive")
	xbuildCmd.Flags().BoolVarP(&buildOptions.SkipAssets, "skip-assets", "k", false, "skip running webpack and building assets")
	xbuildCmd.Flags().BoolVarP(&buildOptions.Static, "static", "s", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"'")
	xbuildCmd.Flags().StringVar(&buildOptions.LDFlags, "ldflags", "", "set any ldflags to be passed to the go build")
	xbuildCmd.Flags().BoolVarP(&buildOptions.Verbose, "verbose", "v", false, "print debugging information")
	xbuildCmd.Flags().BoolVar(&buildOptions.DryRun, "dry-run", false, "runs the build 'dry'")
	xbuildCmd.Flags().BoolVar(&buildOptions.SkipTemplateValidation, "skip-template-validation", false, "skip validating templates")
	xbuildCmd.Flags().BoolVar(&buildOptions.CleanAssets, "clean-assets", false, "will delete public/assets before calling webpack")
	xbuildCmd.Flags().StringVarP(&buildOptions.Environment, "environment", "", "development", "set the environment for the binary")
	xbuildCmd.Flags().StringVar(&buildOptions.Mod, "mod", "", "-mod flag for go build")
}

func buildVersion(version string) string {
	vcs := buildOptions.VCS

	if len(vcs) == 0 {
		return version
	}

	ctx := context.Background()
	run := genny.WetRunner(ctx)
	if buildOptions.DryRun {
		run = genny.DryRunner(ctx)
	}

	_, err := exec.LookPath(vcs)
	if err != nil {
		run.Logger.Warnf("could not find %s; defaulting to version %s", vcs, version)
		return vcs
	}
	var cmd *exec.Cmd
	switch vcs {
	case "git":
		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
	case "bzr":
		cmd = exec.Command("bzr", "revno")
	default:
		run.Logger.Warnf("could not find %s; defaulting to version %s", vcs, version)
		return vcs
	}

	out := &bytes.Buffer{}
	cmd.Stdout = out
	run.WithRun(func(r *genny.Runner) error {
		return r.Exec(cmd)
	})

	if err := run.Run(); err != nil {
		run.Logger.Error(err)
		return version
	}

	if out.String() != "" {
		return strings.TrimSpace(out.String())
	}

	return version
}
