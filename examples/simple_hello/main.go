package main

import (
	"log"
	"net/http"

	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/render"
)

func main() {
	a := buffalo.Automatic(buffalo.NewOptions())
	a.GET("/", func(c buffalo.Context) error {
		return c.Render(200, render.String("Hello, World!"))
	})
	log.Fatal(http.ListenAndServe(":3000", a))
}
