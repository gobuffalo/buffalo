package actions

import "github.com/gobuffalo/buffalo"

// HomeHandler renders "index.html"
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}
