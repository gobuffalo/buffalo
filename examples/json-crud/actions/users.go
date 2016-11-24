package actions

import (
	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/examples/json-crud/models"
	"github.com/markbates/buffalo/render"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

func UsersList(c buffalo.Context) error {
	users := &models.Users{}
	tx := c.Get("tx").(*pop.Connection)
	err := tx.All(users)
	if err != nil {
		return c.Error(404, errors.WithStack(err))
	}

	return c.Render(200, render.JSON(users))
}

func UsersShow(c buffalo.Context) error {
	u := &models.User{}
	id, err := c.ParamInt("id")
	if err != nil {
		return errors.WithStack(err)
	}
	tx := c.Get("tx").(*pop.Connection)
	err = tx.Find(u, id)
	if err != nil {
		return c.Error(404, errors.WithStack(err))
	}
	return c.Render(200, render.JSON(u))
}

func UsersCreate(c buffalo.Context) error {
	u := &models.User{}
	err := c.Bind(u)
	if err != nil {
		return errors.WithStack(err)
	}

	tx := c.Get("tx").(*pop.Connection)
	err = tx.Create(u)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(201, render.JSON(u))
}

func UsersUpdate(c buffalo.Context) error {
	u := &models.User{}
	id, err := c.ParamInt("id")
	if err != nil {
		return errors.WithStack(err)
	}

	tx := c.Get("tx").(*pop.Connection)
	err = tx.Find(u, id)
	if err != nil {
		return c.Error(404, errors.WithStack(err))
	}

	err = c.Bind(u)
	if err != nil {
		return errors.WithStack(err)
	}
	// make sure to set the ID back to the one from
	// URL so people can't set it via the JSON to a
	// different value
	u.ID = id

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

func UsersDelete(c buffalo.Context) error {
	u := &models.User{}
	id, err := c.ParamInt("id")
	if err != nil {
		return errors.WithStack(err)
	}

	tx := c.Get("tx").(*pop.Connection)
	err = tx.Find(u, id)
	if err != nil {
		return c.Error(404, errors.WithStack(err))
	}

	err = tx.Destroy(u)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, render.JSON(u))
}
