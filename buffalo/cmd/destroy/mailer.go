package destroy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/flect"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// MailerCmd destroys a passed mailer
var MailerCmd = &cobra.Command{
	Use: "mailer [name]",
	// Example: "mailer cars",
	Aliases: []string{"l"},
	Short:   "Destroy mailer files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you need to provide a valid mailer name in order to destroy it")
		}

		name := args[0]

		removeMailer(name)

		return nil
	},
}

func removeMailer(name string) {
	if YesToAll || confirm("Want to remove mailer? (y/N)") {
		mailerFileName := flect.Singularize(flect.Underscore(name))

		os.Remove(filepath.Join("mailers", fmt.Sprintf("%v.go", mailerFileName)))
		os.Remove(filepath.Join("templates/mail", fmt.Sprintf("%v.html", mailerFileName)))

		logrus.Infof("- Deleted %v", fmt.Sprintf("mailers/%v.go", mailerFileName))
		logrus.Infof("- Deleted %v", fmt.Sprintf("templates/mail/%v.html", mailerFileName))
	}
}
