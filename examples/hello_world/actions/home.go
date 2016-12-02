package actions

import "github.com/markbates/buffalo"

func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}
