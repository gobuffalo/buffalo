package main

import (
	"log"
	"net/http"

	"github.com/markbates/buffalo/examples/hello_world/actions"
)

func main() {
	log.Fatal(http.ListenAndServe(":3000", actions.App()))
}

