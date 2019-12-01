package buffalo

import (
	"fmt"
	"net/http"
)

// Resource interface allows for the easy mapping
// of common RESTful actions to a set of paths. See
// the a.Resource documentation for more details.
// NOTE: When skipping Resource handlers, you need to first declare your
// resource handler as a type of buffalo.Resource for the Skip function to
// properly recognize and match it.
/*
	// Works:
	var cr Resource
	cr = &carsResource{&buffaloBaseResource{}}
	g = a.Resource("/cars", cr)
	g.Use(SomeMiddleware)
	g.Middleware.Skip(SomeMiddleware, cr.Show)

	// Doesn't Work:
	cr := &carsResource{&buffaloBaseResource{}}
	g = a.Resource("/cars", cr)
	g.Use(SomeMiddleware)
	g.Middleware.Skip(SomeMiddleware, cr.Show)
*/
type Resource interface {
	List(Context) error
	Show(Context) error
	Create(Context) error
	Update(Context) error
	Destroy(Context) error
}

// Middler can be implemented to specify additional
// middleware specific to the resource
type Middler interface {
	Use() []MiddlewareFunc
}

// BaseResource fills in the gaps for any Resource interface
// functions you don't want/need to implement.
/*
	type UsersResource struct {
		Resource
	}

	func (ur *UsersResource) List(c Context) error {
		return c.Render(http.StatusOK, render.String("hello")
	}

	// This will fulfill the Resource interface, despite only having
	// one of the functions defined.
	&UsersResource{&BaseResource{})
*/
type BaseResource struct{}

var errResourceNotImplemented = fmt.Errorf("resource not implemented")

// List default implementation. Returns a 404
func (v BaseResource) List(c Context) error {
	return c.Error(http.StatusNotFound, errResourceNotImplemented)
}

// Show default implementation. Returns a 404
func (v BaseResource) Show(c Context) error {
	return c.Error(http.StatusNotFound, errResourceNotImplemented)
}

// Create default implementation. Returns a 404
func (v BaseResource) Create(c Context) error {
	return c.Error(http.StatusNotFound, errResourceNotImplemented)
}

// Update default implementation. Returns a 404
func (v BaseResource) Update(c Context) error {
	return c.Error(http.StatusNotFound, errResourceNotImplemented)
}

// Destroy default implementation. Returns a 404
func (v BaseResource) Destroy(c Context) error {
	return c.Error(http.StatusNotFound, errResourceNotImplemented)
}
