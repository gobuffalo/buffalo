package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/examples/html-resource/models"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

// UsersResource allows CRUD with HTTP against the User model
type UsersResource struct {
	buffalo.BaseResource
}

func findUserMW(n string) buffalo.MiddlewareFunc {
	return func(h buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			id, err := c.ParamInt(n)
			if err == nil {
				u := &models.User{}
				tx := c.Value("tx").(*pop.Connection)
				err = tx.Find(u, id)
				if err != nil {
					return c.Error(404, errors.WithStack(err))
				}
				c.Set("user", u)
			}
			return h(c)
		}
	}
}

// List shows all users in an HTML page
func (ur *UsersResource) List(c buffalo.Context) error {
	users := &models.Users{}
	tx := c.Value("tx").(*pop.Connection)
	err := tx.All(users)
	if err != nil {
		return c.Error(404, errors.WithStack(err))
	}

	c.Set("users", users)
	return c.Render(200, r.HTML("users/index.html"))
}

// Show renders a target user in an HTML page
func (ur *UsersResource) Show(c buffalo.Context) error {
	return c.Render(200, r.HTML("users/show.html"))
}

// New renders a form for adding a new user
func (ur *UsersResource) New(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(200, r.HTML("users/new.html"))
}

// Create is a JSON API endpoint that adds a new user
func (ur *UsersResource) Create(c buffalo.Context) error {
	u := &models.User{}
	err := c.Bind(u)
	if err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.ValidateNew(tx)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("verrs", verrs.Errors)
		c.Set("user", u)
		return c.Render(422, r.HTML("users/new.html"))
	}
	err = tx.Create(u)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(301, "/users/%d", u.ID)
}

// Edit renders an html form for editing a user
func (ur *UsersResource) Edit(c buffalo.Context) error {
	return c.Render(200, r.HTML("users/edit.html"))
}

// Update is a JSON API endpoint that updates a user
func (ur *UsersResource) Update(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	u := c.Value("user").(*models.User)

	err := c.Bind(u)
	if err != nil {
		return errors.WithStack(err)
	}

	verrs, err := u.ValidateUpdate(tx)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("verrs", verrs.Errors)
		c.Set("user", u)
		return c.Render(422, r.HTML("users/edit.html"))
	}
	err = tx.Update(u)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Reload(u)
	if err != nil {
		return errors.WithStack(err)
	}
	return c.Redirect(301, "/users/%d", u.ID)
}

// Destroy is an API endpoint that deletes a user
func (ur *UsersResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	u := c.Value("user").(*models.User)

	err := tx.Destroy(u)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(301, "/users")
}
