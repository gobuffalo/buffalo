package actions

import (
	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/examples/json-resource/models"
	"github.com/markbates/buffalo/render"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

type UsersResource struct {
	buffalo.BaseResource
}

func findUserMW(h buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		id, err := c.ParamInt("user_id")
		if err == nil {
			u := &models.User{}
			tx := c.Get("tx").(*pop.Connection)
			err = tx.Find(u, id)
			if err != nil {
				return c.Error(404, errors.WithStack(err))
			}
			c.Set("user", u)
		}
		return h(c)
	}
}

func (ur *UsersResource) List(c buffalo.Context) error {
	users := &models.Users{}
	tx := c.Get("tx").(*pop.Connection)
	err := tx.All(users)
	if err != nil {
		return c.Error(404, errors.WithStack(err))
	}

	return c.Render(200, render.JSON(users))
}

func (ur *UsersResource) Show(c buffalo.Context) error {
	return c.Render(200, render.JSON(c.Get("user")))
}

func (ur *UsersResource) Create(c buffalo.Context) error {
	u := &models.User{}
	err := c.Bind(u)
	if err != nil {
		return errors.WithStack(err)
	}

	tx := c.Get("tx").(*pop.Connection)
	verrs, err := u.ValidateNew(tx)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("verrs", verrs.Errors)
		return c.Render(422, render.JSON(verrs))
	}
	err = tx.Create(u)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(201, render.JSON(u))
}

func (ur *UsersResource) Update(c buffalo.Context) error {
	tx := c.Get("tx").(*pop.Connection)
	u := c.Get("user").(*models.User)

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
		return c.Render(422, render.JSON(verrs))
	}
	err = tx.Update(u)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Reload(u)
	if err != nil {
		return errors.WithStack(err)
	}
	return c.Render(200, render.JSON(u))
}

func (ur *UsersResource) Destroy(c buffalo.Context) error {
	tx := c.Get("tx").(*pop.Connection)
	u := c.Get("user").(*models.User)

	err := tx.Destroy(u)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, render.JSON(u))
}
