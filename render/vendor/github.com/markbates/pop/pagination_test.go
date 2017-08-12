package pop_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/stretchr/testify/require"
)

func Test_NewPaginator(t *testing.T) {
	a := require.New(t)

	p := pop.NewPaginator(1, 10)
	a.Equal(p.Offset, 0)

	p = pop.NewPaginator(2, 10)
	a.Equal(p.Offset, 10)

	p = pop.NewPaginator(2, 30)
	a.Equal(p.Offset, 30)
}

func Test_NewPaginatorFromParams(t *testing.T) {
	a := require.New(t)

	params := url.Values{}

	p := pop.NewPaginatorFromParams(params)
	a.Equal(p.Page, 1)
	a.Equal(p.PerPage, 20)

	params.Set("page", "2")
	p = pop.NewPaginatorFromParams(params)
	a.Equal(p.Page, 2)
	a.Equal(p.PerPage, 20)

	params.Set("per_page", "30")
	p = pop.NewPaginatorFromParams(params)
	a.Equal(p.Page, 2)
	a.Equal(p.PerPage, 30)
}

func Test_Pagination(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		for _, name := range []string{"Mark", "Joe", "Jane"} {
			user := User{Name: nulls.NewString(name)}
			err := tx.Create(&user)
			a.NoError(err)
		}

		u := Users{}
		q := tx.Paginate(1, 2)
		err := q.All(&u)
		a.NoError(err)
		a.Equal(len(u), 2)

		p := q.Paginator
		a.Equal(p.CurrentEntriesSize, 2)
		a.Equal(p.TotalEntriesSize, 3)
		a.Equal(p.TotalPages, 2)

		u = Users{}
		err = tx.Where("name = 'Mark'").All(&u)
		a.NoError(err)
		a.Equal(reflect.ValueOf(&u).Elem().Len(), 1)
	})
}
