package pop_test

import (
	"testing"

	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/stretchr/testify/require"
)

func Test_Find(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		a.NoError(err)

		u := User{}
		err = tx.Find(&u, user.ID)
		a.NoError(err)

		a.NotEqual(u.ID, 0)
		a.Equal(u.Name.String, "Mark")
	})
}

func Test_First(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		first := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&first)
		a.NoError(err)

		last := User{Name: nulls.NewString("Mark")}
		err = tx.Create(&last)
		a.NoError(err)

		u := User{}
		err = tx.Where("name = 'Mark'").First(&u)
		a.NoError(err)

		a.Equal(first.ID, u.ID)
	})
}

func Test_Last(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		first := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&first)
		a.NoError(err)

		last := User{Name: nulls.NewString("Mark")}
		err = tx.Create(&last)
		a.NoError(err)

		u := User{}
		err = tx.Where("name = 'Mark'").Last(&u)
		a.NoError(err)

		a.Equal(last.ID, u.ID)
	})
}

func Test_All(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		for _, name := range []string{"Mark", "Joe", "Jane"} {
			user := User{Name: nulls.NewString(name)}
			err := tx.Create(&user)
			a.NoError(err)
		}

		u := Users{}
		err := tx.All(&u)
		a.NoError(err)
		a.Equal(len(u), 3)

		u = Users{}
		err = tx.Where("name = 'Mark'").All(&u)
		a.NoError(err)
		a.Equal(len(u), 1)
	})
}

func Test_Count(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		a.NoError(err)
		c, err := tx.Count(&user)
		a.NoError(err)
		a.Equal(c, 1)

		c, err = tx.Where("1=1").CountByField(&user, "distinct id")
		a.NoError(err)
		a.Equal(c, 1)
		// should ignore order in count

		c, err = tx.Order("id desc").Count(&user)
		a.NoError(err)
		a.Equal(c, 1)
	})
}

func Test_Count_Disregards_Pagination(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		names := []string{
			"Jack",
			"Hurley",
			"Charlie",
			"Desmond",
			"Juliet",
			"Locke",
			"Sawyer",
			"Kate",
			"Benjamin Linus",
		}

		for _, name := range names {
			user := User{Name: nulls.NewString(name)}
			err := tx.Create(&user)
			a.NoError(err)
		}

		first_users := Users{}
		second_users := Users{}

		q := tx.Paginate(1, 3)
		q.All(&first_users)

		a.Equal(3, len(first_users))
		totalFirstPage := q.Paginator.TotalPages

		q = tx.Paginate(2, 3)
		q.All(&second_users)

		a.Equal(3, len(second_users))
		totalSecondPage := q.Paginator.TotalPages

		a.NotEqual(0, totalFirstPage)
		a.NotEqual(0, totalSecondPage)
		a.Equal(totalFirstPage, totalSecondPage)
	})
}

func Test_Count_RawQuery(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		a.NoError(err)

		c, err := tx.RawQuery("select count(*) as row_count from users as users").Count(nil)
		a.NoError(err)
		a.Equal(c, 1)

		c, err = tx.RawQuery("select count(*) as row_count from users as users where id = -1").Count(nil)
		a.NoError(err)
		a.Equal(c, 0)
	})
}

func Test_Exists(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		t, _ := tx.Where("id = ?", 0).Exists("users")
		a.False(t)

		user := User{Name: nulls.NewString("Mark")}
		err := tx.Create(&user)
		a.NoError(err)

		t, _ = tx.Where("id = ?", user.ID).Exists("users")
		a.True(t)
	})
}
