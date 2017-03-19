package actions

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
		HTMLLayout:   "application.html",
		TemplatesBox: packr.NewBox("../templates"),
		Helpers:      render.Helpers{},
	})
}
