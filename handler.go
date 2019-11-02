package buffalo

// Handler is the basis for all of Buffalo. A Handler
// will be given a Context interface that represents the
// give request/response. It is the responsibility of the
// Handler to handle the request/response correctly. This
// could mean rendering a template, JSON, etc... or it could
// mean returning an error.
/*
	func (c Context) error {
		return c.Render(http.StatusOK, render.String("Hello World!"))
	}

	func (c Context) error {
		return c.Redirect(http.StatusMovedPermanently, "http://github.com/gobuffalo/buffalo")
	}

	func (c Context) error {
		return c.Error(http.StatusUnprocessableEntity, fmt.Errorf("oops!!"))
	}
*/
type Handler func(Context) error
