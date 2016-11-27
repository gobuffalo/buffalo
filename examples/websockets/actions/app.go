package actions

import (
	"net/http"

	"github.com/markbates/buffalo"
)

func App() http.Handler {
	a := buffalo.Automatic(buffalo.Options{Env: "development"})

	a.ServeFiles("/assets", assetsPath())
	a.GET("/", HomeHandler)
	a.GET("/socket", SocketHandler)

	return a
}
