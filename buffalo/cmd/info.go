package cmd

import (
	"context"
	"os/exec"
	"time"

	"github.com/gobuffalo/buffalo/genny/info"
	"github.com/gobuffalo/clara/genny/rx"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
	"github.com/spf13/cobra"
)

var infoOptions = struct {
	Clara *rx.Options
	Info  *info.Options
}{
	Clara: &rx.Options{},
	Info:  &info.Options{},
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print diagnostic information (useful for debugging)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		run := genny.WetRunner(ctx)

		_, err := run.LookPath("clara")
		if err == nil {
			// use the clara binary if available
			run.WithRun(func(r *genny.Runner) error {
				return r.Exec(exec.Command("clara"))
			})
		} else {
			// no clara binary, so use the one bundled with buffalo
			copts := infoOptions.Clara
			if err := run.WithNew(rx.New(copts)); err != nil {
				return err
			}
		}

		iopts := infoOptions.Info
		if err := run.WithNew(info.New(iopts)); err != nil {
			return err
		}

		return run.Run()
	},
}

func init() {
	app := meta.New(".")
	infoOptions.Clara.App = app
	infoOptions.Info.App = app

	decorate("info", RootCmd)
	RootCmd.AddCommand(infoCmd)
}
