package generate

import (
	"context"

	"github.com/gobuffalo/buffalo/genny/mail"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/meta"
	"github.com/markbates/inflect"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var mailOptions = struct {
	dryRun bool
	*mail.Options
}{
	Options: &mail.Options{},
}

// MailCmd for generating mailers
var MailCmd = &cobra.Command{
	Use:   "mailer [name]",
	Short: "Generate a new mailer for Buffalo",
	RunE: func(cmd *cobra.Command, args []string) error {
		mailOptions.App = meta.New(".")
		mailOptions.Name = inflect.Name(args[0])
		gg, err := mail.New(mailOptions.Options)
		if err != nil {
			return errors.WithStack(err)
		}

		run := genny.WetRunner(context.Background())
		if mailOptions.dryRun {
			run = genny.DryRunner(context.Background())
		}

		g, err := gotools.GoFmt(mailOptions.App.Root)
		if err != nil {
			return errors.WithStack(err)
		}
		run.With(g)

		gg.With(run)
		return run.Run()

	},
}

func init() {
	MailCmd.Flags().BoolVarP(&mailOptions.dryRun, "dry-run", "d", false, "dry run of the generator")
	MailCmd.Flags().BoolVar(&mailOptions.SkipInit, "skip-init", false, "skip initializing mailers/")
}
