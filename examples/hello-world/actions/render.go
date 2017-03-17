package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/buffalo/render/resolvers"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
		HTMLLayout:     "application.html",
		TemplateEngine: plush.BuffaloRenderer,
		FileResolverFunc: func() resolvers.FileResolver {
			return &resolvers.PackrBox{
				Box: packr.NewBox("../templates"),
			}
		},
	})
}

func assetsPath() http.FileSystem {
	box := packr.NewBox("../assets")
	return box.HTTPBox()
}
