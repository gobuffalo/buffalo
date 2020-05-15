package binding

import (
	"net/http"
)

// ContenTypeBinder are those capable of handling a request type like JSON or XML
type ContenTypeBinder interface {
	BinderFunc() Binder
	ContentTypes() []string
}

// Binder takes a request and binds it to an interface.
// If there is a problem it should return an error.
type Binder func(*http.Request, interface{}) error

// CustomTypeDecoder converts a custom type from the request into its exact type.
type CustomTypeDecoder func([]string) (interface{}, error)