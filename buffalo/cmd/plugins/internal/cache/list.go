package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gobuffalo/buffalo/plugins"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// ListCmd displays the contents of the plugin cache
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "displays the contents of the plugin cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := ioutil.ReadFile(plugins.CachePath)
		if err != nil {
			return errors.WithStack(err)
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return errors.WithStack(err)
		}
		is, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return errors.WithStack(err)
		}

		fmt.Println(string(is))
		return nil
	},
}
