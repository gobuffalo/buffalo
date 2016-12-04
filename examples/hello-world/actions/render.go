package actions

import (
	"net/http"
	"path"
	"runtime"

	"github.com/markbates/buffalo/render"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
		TemplatesPath: fromHere("../templates"),
		HTMLLayout:    "application.html",
	})
}

func assetsPath() http.Dir {
	return http.Dir(fromHere("../assets"))
}

func fromHere(p string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), p)
}
