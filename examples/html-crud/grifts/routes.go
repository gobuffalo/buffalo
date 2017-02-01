package grifts

import (
	"os"

	"github.com/gobuffalo/buffalo/examples/html-crud/actions"
	"github.com/markbates/grift/grift"
	"github.com/olekukonko/tablewriter"
)

var _ = grift.Add("routes", func(c *grift.Context) error {
	a := actions.App()
	routes := a.Routes()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Method", "Path", "Handler"})
	for _, r := range routes {
		table.Append([]string{r.Method, r.Path, r.HandlerName})
	}
	table.SetCenterSeparator("|")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	return nil
})
