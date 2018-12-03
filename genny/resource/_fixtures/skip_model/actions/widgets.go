package actions

import "github.com/gobuffalo/buffalo"

type WidgetsResource struct {
	buffalo.Resource
}

// List default implementation.
func (v WidgetsResource) List(c buffalo.Context) error {
	return c.Render(200, r.String("Widget#List"))
}

// Show default implementation.
func (v WidgetsResource) Show(c buffalo.Context) error {
	return c.Render(200, r.String("Widget#Show"))
}

// New default implementation.
func (v WidgetsResource) New(c buffalo.Context) error {
	return c.Render(200, r.String("Widget#New"))
}

// Create default implementation.
func (v WidgetsResource) Create(c buffalo.Context) error {
	return c.Render(200, r.String("Widget#Create"))
}

// Edit default implementation.
func (v WidgetsResource) Edit(c buffalo.Context) error {
	return c.Render(200, r.String("Widget#Edit"))
}

// Update default implementation.
func (v WidgetsResource) Update(c buffalo.Context) error {
	return c.Render(200, r.String("Widget#Update"))
}

// Destroy default implementation.
func (v WidgetsResource) Destroy(c buffalo.Context) error {
	return c.Render(200, r.String("Widget#Destroy"))
}
