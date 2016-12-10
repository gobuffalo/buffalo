package buffalo

import "errors"

type Resource interface {
	List(Context) error
	Show(Context) error
	New(Context) error
	Create(Context) error
	Edit(Context) error
	Update(Context) error
	Destroy(Context) error
}

// BaseResource fills in the gaps for any Resource interface
// functions you don't want/need to implement.
/*
	type UsersResource struct {
		Resource
	}

	func (ur *UsersResource) List(c Context) error {
		return c.Render(200, render.String("hello")
	}

	// This will fulfill the Resource interface, despite only having
	// one of the functions defined.
	&UsersResource{&BaseResource{})
*/
type BaseResource struct{}

func (v *BaseResource) List(c Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

func (v *BaseResource) Show(c Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

func (v *BaseResource) New(c Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

func (v *BaseResource) Create(c Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

func (v *BaseResource) Edit(c Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

func (v *BaseResource) Update(c Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

func (v *BaseResource) Destroy(c Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}
