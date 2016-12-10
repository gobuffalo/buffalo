package actions_test

import (
	"testing"

	"github.com/markbates/buffalo/examples/json-resource/actions"
	"github.com/markbates/buffalo/examples/json-resource/models"
	"github.com/markbates/going/validate"
	"github.com/markbates/pop"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_UsersList(t *testing.T) {
	r := require.New(t)

	tx(func(tx *pop.Connection) {
		w := willie.New(actions.App())

		r.NoError(tx.Create(&models.User{
			FirstName: "Mark",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}))

		res := w.JSON("/users").Get()
		r.Equal(200, res.Code)

		users := models.Users{}
		res.Bind(&users)
		r.Len(users, 1)
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

		res := w.JSON("/users/%d", u.ID).Get()
		r.Equal(200, res.Code)

		user := &models.User{}
		res.Bind(user)
		r.Equal(u.Email, user.Email)
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
		res := w.JSON("/users").Post(u)
		r.Equal(201, res.Code)

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
		res := w.JSON("/users").Post(u)
		r.Equal(422, res.Code)

		verrs := validate.NewErrors()
		res.Bind(verrs)
		r.Len(verrs.Errors, 3)
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

		res := w.JSON("/users/%d", u.ID).Put(map[string]string{
			"email": "bates@example.com",
		})
		r.Equal(200, res.Code)

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

		res := w.JSON("/users/%d", u.ID).Put(map[string]string{
			"email": "bates@example.com",
		})
		r.Equal(422, res.Code)

		verrs := validate.NewErrors()
		res.Bind(verrs)
		r.Equal([]string{"First Name can not be blank."}, verrs.Get("first_name"))
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

		res := w.JSON("/users/%d", u.ID).Delete()
		r.Equal(200, res.Code)

		ct, err = tx.Count("users")
		r.NoError(err)
		r.Equal(0, ct)
	})
}
