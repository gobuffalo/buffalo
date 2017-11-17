package generate

import (
	"github.com/gobuffalo/buffalo/generators/mail"
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/makr"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var mailer = mail.Generator{}

// MailCmd for generating mailers
var MailCmd = &cobra.Command{
	Use:   "mailer",
	Short: "Generates a new mailer for Buffalo",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must supply a name for your mailer")
		}
		mailer.App = meta.New(".")
		mailer.Name = meta.Name(args[0])
		data := makr.Data{}
		return mailer.Run(".", data)

	},
}

func init() {
	MailCmd.Flags().BoolVar(&mailer.SkipInit, "skip-init", false, "skip initializing mailers/")
}
