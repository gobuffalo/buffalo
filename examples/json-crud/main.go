package main

import (
	"log"
	"net/http"

	"github.com/markbates/buffalo/examples/json-crud/actions"
)

func main() {
	log.Fatal(http.ListenAndServe(":3000", actions.App()))
}
