package actions_test

import (
	"testing"

	"github.com/markbates/buffalo/examples/json-crud/actions"
	"github.com/markbates/buffalo/examples/json-crud/models"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_UsersList(t *testing.T) {
	r := require.New(t)

	tx(func() {
		w := willie.New(actions.App())

		r.NoError(models.DB.Create(&models.User{
			FirstName: "Mark",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}))

		res := w.JSON("/users").Get()
		r.Empty(200, res.Code)

		users := models.Users{}
		res.Bind(&users)
		r.Len(users, 1)
	})
}

func Test_UsersShow(t *testing.T) {
	r := require.New(t)

	tx(func() {
		w := willie.New(actions.App())

		u := &models.User{
			FirstName: "Mark",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}
		r.NoError(models.DB.Create(u))

		res := w.JSON("/users/%d", u.ID).Get()
		r.Empty(200, res.Code)

		user := &models.User{}
		res.Bind(user)
		r.Equal(u.Email, user.Email)
	})
}

func Test_UsersCreate(t *testing.T) {
	r := require.New(t)

	tx(func() {
		w := willie.New(actions.App())

		ct, err := models.DB.Count("users")
		r.NoError(err)
		r.Equal(0, ct)

		u := &models.User{
			FirstName: "Mark",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}
		res := w.JSON("/users").Put(u)
		r.Empty(201, res.Code)

		ct, err = models.DB.Count("users")
		r.NoError(err)
		r.Equal(1, ct)

	})
}

func Test_UsersUpdate(t *testing.T) {
	r := require.New(t)

	tx(func() {
		w := willie.New(actions.App())

		u := &models.User{
			FirstName: "Mark",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}
		r.NoError(models.DB.Create(u))

		res := w.JSON("/users/%d", u.ID).Put(map[string]string{
			"email": "bates@example.com",
		})
		r.Empty(200, res.Code)

		r.NoError(models.DB.Reload(u))
		r.Equal("bates@example.com", u.Email)
	})
}

func Test_UsersDestroy(t *testing.T) {
	r := require.New(t)

	tx(func() {
		w := willie.New(actions.App())

		u := &models.User{
			FirstName: "Mark",
			LastName:  "Bates",
			Email:     "mark@example.com",
		}
		r.NoError(models.DB.Create(u))

		ct, err := models.DB.Count("users")
		r.NoError(err)
		r.Equal(1, ct)

		res := w.JSON("/users/%d", u.ID).Delete()
		r.Empty(200, res.Code)

		ct, err = models.DB.Count("users")
		r.NoError(err)
		r.Equal(0, ct)
	})
}
