package actions

import (
	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/examples/html-crud/models"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

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

func UsersList(c buffalo.Context) error {
	users := &models.Users{}
	tx := c.Get("tx").(*pop.Connection)
	err := tx.All(users)
	if err != nil {
		return c.Error(404, errors.WithStack(err))
	}

	c.Set("users", users)
	return c.Render(200, r.HTML("users/index.html"))
}

func UsersShow(c buffalo.Context) error {
	return c.Render(200, r.HTML("users/show.html"))
}

func UsersNew(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(200, r.HTML("users/new.html"))
}

func UsersCreate(c buffalo.Context) error {
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
		c.Set("user", u)
		return c.Render(422, r.HTML("users/new.html"))
	}
	err = tx.Create(u)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(301, "/users/%d", u.ID)
}

func UsersEdit(c buffalo.Context) error {
	return c.Render(200, r.HTML("users/edit.html"))
}

func UsersUpdate(c buffalo.Context) error {
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

func UsersDelete(c buffalo.Context) error {
	tx := c.Get("tx").(*pop.Connection)
	u := c.Get("user").(*models.User)

	err := tx.Destroy(u)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(301, "/users")
}
