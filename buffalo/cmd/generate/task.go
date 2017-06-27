package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo/generators/grift"
	"github.com/gobuffalo/makr"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

//TaskCmd is the command called with the generate grift cli.
var TaskCmd = &cobra.Command{
	Use:     "task [name]",
	Aliases: []string{"t", "grift"},
	Short:   "Generates a grift task",
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return errors.New("you need to provide a name for the grift tasks")
		}

		var parts []string
		plain := strings.Contains(args[0], ":") == false
		filename := fmt.Sprintf("%v.go", inflect.Underscore(args[0]))

		if !plain {
			parts = strings.Split(args[0], ":")
			filename = fmt.Sprintf("%v.go", inflect.Underscore(parts[len(parts)-1]))
		}

		data := makr.Data{
			"name":      args[0],
			"taskName":  inflect.Underscore(args[0]),
			"filename":  filename,
			"plainTask": plain,
			"parts":     parts,
			"last":      len(parts) - 1,
		}

		g, err := grift.New(data)
		if err != nil {
			return err
		}

		return g.Run(".", data)
	},
}
