package actions_test

import (
	"testing"

	"github.com/markbates/buffalo/examples/html-resource/actions"
	"github.com/markbates/buffalo/examples/html-resource/models"
	"github.com/markbates/pop"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_UsersList(t *testing.T) {
	r := require.New(t)

	tx(func(tx *pop.Connection) {
		w := willie.New(actions.App())
		u := &models.User{
			FirstName: "Mark",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}
		r.NoError(tx.Create(u))

		res := w.Request("/users").Get()
		r.Equal(200, res.Code)

		r.Contains(res.Body.String(), u.Email)
	})
}

func Test_UsersShow(t *testing.T) {
	r := require.New(t)

	tx(func(tx *pop.Connection) {
		w := willie.New(actions.App())

		u := &models.User{
			FirstName: "Mark",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}
		r.NoError(tx.Create(u))

		res := w.Request("/users/%d", u.ID).Get()
		r.Equal(200, res.Code)

		r.Contains(res.Body.String(), u.Email)
	})
}

func Test_UsersCreate(t *testing.T) {
	r := require.New(t)

	tx(func(tx *pop.Connection) {
		w := willie.New(actions.App())

		ct, err := tx.Count("users")
		r.NoError(err)
		r.Equal(0, ct)

		u := &models.User{
			FirstName: "Mark",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}
		res := w.Request("/users").Post(u)
		r.Equal(301, res.Code)

		ct, err = tx.Count("users")
		r.NoError(err)
		r.Equal(1, ct)
	})
}

func Test_UsersCreate_HandlesErrors(t *testing.T) {
	r := require.New(t)

	tx(func(tx *pop.Connection) {
		w := willie.New(actions.App())

		ct, err := tx.Count("users")
		r.NoError(err)
		r.Equal(0, ct)

		u := &models.User{}
		res := w.Request("/users").Post(u)
		r.Equal(422, res.Code)

		r.Contains(res.Body.String(), "First Name can not be blank.")
	})
}

func Test_UsersUpdate(t *testing.T) {
	r := require.New(t)

	tx(func(tx *pop.Connection) {
		w := willie.New(actions.App())

		u := &models.User{
			FirstName: "Mark",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}
		r.NoError(tx.Create(u))

		res := w.Request("/users/%d", u.ID).Put(map[string]string{
			"email": "bates@example.com",
		})
		r.Equal(301, res.Code)

		r.NoError(tx.Reload(u))
		r.Equal("bates@example.com", u.Email)
	})
}

func Test_UsersUpdate_HandlesErrors(t *testing.T) {
	r := require.New(t)

	tx(func(tx *pop.Connection) {
		w := willie.New(actions.App())

		u := &models.User{
			FirstName: "",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}
		r.NoError(tx.Create(u))

		res := w.Request("/users/%d", u.ID).Put(map[string]string{
			"email": "bates@example.com",
		})
		r.Equal(422, res.Code)

		r.Contains(res.Body.String(), "First Name can not be blank.")
	})
}

func Test_UsersDestroy(t *testing.T) {
	r := require.New(t)

	tx(func(tx *pop.Connection) {
		w := willie.New(actions.App())

		u := &models.User{
			FirstName: "Mark",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}
		r.NoError(tx.Create(u))

		ct, err := tx.Count("users")
		r.NoError(err)
		r.Equal(1, ct)

		res := w.Request("/users/%d", u.ID).Delete()
		r.Equal(301, res.Code)

		ct, err = tx.Count("users")
		r.NoError(err)
		r.Equal(0, ct)
	})
}
