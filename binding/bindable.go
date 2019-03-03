package binding

import "net/http"

// Bindable when implemented, on a type
// will override any Binders that have been
// configured when using buffalo#Context.Bind
type Bindable interface {
	Bind(*http.Request) error
}
