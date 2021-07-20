package plugins

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"

	pluginsin "github.com/gobuffalo/buffalo/plugins"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "a list of installed buffalo plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		list, err := pluginsin.Available()
		if err != nil {
			return err
		}

		var cmds pluginsin.Commands

		for _, l := range list {
			cmds = append(cmds, l...)
		}

		sort.Slice(cmds, func(i, j int) bool {
			c1 := cmds[i]
			c2 := cmds[j]

			return c1.Name+c1.Name < c2.Name+c2.Name
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Bin\tCommand\tDescription")
		fmt.Fprintln(w, "---\t---\t---")

		for _, c := range cmds {
			if c.Name == "" {
				continue
			}
			sb := &bytes.Buffer{}
			sb.WriteString("buffalo ")
			if c.BuffaloCommand != "root" {
				sb.WriteString(c.BuffaloCommand)
				sb.WriteString(" ")
			}
			sb.WriteString(c.Name)
			fmt.Fprintf(w, "%s\t%s\t%s\n", filepath.Base(c.Binary), sb.String(), c.Description)
		}

		return w.Flush()
	},
}
