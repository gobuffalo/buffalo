package generate

import (
	"context"

	"github.com/gobuffalo/buffalo/genny/mail"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gogen"
	"github.com/gobuffalo/meta"
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
		mailOptions.Name = name.New(args[0])
		gg, err := mail.New(mailOptions.Options)
		if err != nil {
			return err
		}

		run := genny.WetRunner(context.Background())
		if mailOptions.dryRun {
			run = genny.DryRunner(context.Background())
		}

		g, err := gogen.Fmt(mailOptions.App.Root)
		if err != nil {
			return err
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
