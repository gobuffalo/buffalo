package pop_test

import (
	"fmt"
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_Where(t *testing.T) {
	a := require.New(t)
	m := &pop.Model{Value: &Enemy{}}

	q := PDB.Where("id = ?", 1)
	sql, _ := q.ToSQL(m)
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies WHERE id = ?"), sql)

	q.Where("first_name = ? and last_name = ?", "Mark", "Bates")
	sql, _ = q.ToSQL(m)
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies WHERE id = ? AND first_name = ? and last_name = ?"), sql)

	q = PDB.Where("name = ?", "Mark 'Awesome' Bates")
	sql, _ = q.ToSQL(m)
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies WHERE name = ?"), sql)

	q = PDB.Where("name = ?", "'; truncate users; --")
	sql, _ = q.ToSQL(m)
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies WHERE name = ?"), sql)
}

func Test_Where_In(t *testing.T) {
	r := require.New(t)
	transaction(func(tx *pop.Connection) {
		u1 := &Song{Title: "A"}
		u2 := &Song{Title: "B"}
		u3 := &Song{Title: "C"}
		err := tx.Create(u1)
		r.NoError(err)
		err = tx.Create(u2)
		r.NoError(err)
		err = tx.Create(u3)
		r.NoError(err)

		songs := []Song{}
		err = tx.Where("id in (?)", u1.ID, u3.ID).All(&songs)
		r.Len(songs, 2)
	})
}

func Test_Order(t *testing.T) {
	a := require.New(t)

	m := &pop.Model{Value: &Enemy{}}
	q := PDB.Order("id desc")
	sql, _ := q.ToSQL(m)
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies ORDER BY id desc"), sql)

	q.Order("name desc")
	sql, _ = q.ToSQL(m)
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies ORDER BY id desc, name desc"), sql)
}

func Test_GroupBy(t *testing.T) {
	a := require.New(t)

	m := &pop.Model{Value: &Enemy{}}
	q := PDB.Q()
	q.GroupBy("A")
	sql, _ := q.ToSQL(m)
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies GROUP BY A"), sql)

	q = PDB.Q()
	q.GroupBy("A", "B")
	sql, _ = q.ToSQL(m)
	a.Equal(ts("SELECT enemies.A FROM enemies AS enemies GROUP BY A, B"), sql)

	q = PDB.Q()
	q.GroupBy("A", "B").Having("enemies.A=?", "test")
	sql, _ = q.ToSQL(m)
	if PDB.Dialect.Details().Dialect == "postgres" {
		a.Equal(ts("SELECT enemies.A FROM enemies AS enemies GROUP BY A, B HAVING enemies.A=$1"), sql)
	} else {
		a.Equal(ts("SELECT enemies.A FROM enemies AS enemies GROUP BY A, B HAVING enemies.A=?"), sql)
	}

	q = PDB.Q()
	q.GroupBy("A", "B").Having("enemies.A=?", "test").Having("enemies.B=enemies.A")
	sql, _ = q.ToSQL(m)
	if PDB.Dialect.Details().Dialect == "postgres" {
		a.Equal(ts("SELECT enemies.A FROM enemies AS enemies GROUP BY A, B HAVING enemies.A=$1 AND enemies.B=enemies.A"), sql)
	} else {
		a.Equal(ts("SELECT enemies.A FROM enemies AS enemies GROUP BY A, B HAVING enemies.A=? AND enemies.B=enemies.A"), sql)
	}
}

func Test_ToSQL(t *testing.T) {
	a := require.New(t)
	transaction(func(tx *pop.Connection) {
		user := &pop.Model{Value: &User{}}

		s := "SELECT name as full_name, users.alive, users.bio, users.birth_date, users.created_at, users.email, users.id, users.name, users.price, users.updated_at FROM users AS users"

		query := pop.Q(tx)
		q, _ := query.ToSQL(user)
		a.Equal(s, q)

		query.Order("id desc")
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s ORDER BY id desc", s), q)

		q, _ = query.ToSQL(&pop.Model{Value: &User{}, As: "u"})
		a.Equal("SELECT name as full_name, u.alive, u.bio, u.birth_date, u.created_at, u.email, u.id, u.name, u.price, u.updated_at FROM users AS u ORDER BY id desc", q)

		query = tx.Where("id = 1")
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s WHERE id = 1", s), q)

		query = tx.Where("id = 1").Where("name = 'Mark'")
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s WHERE id = 1 AND name = 'Mark'", s), q)

		query.Order("id desc")
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s WHERE id = 1 AND name = 'Mark' ORDER BY id desc", s), q)

		query.Order("name asc")
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s WHERE id = 1 AND name = 'Mark' ORDER BY id desc, name asc", s), q)

		query = tx.Limit(10)
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s LIMIT 10", s), q)

		query = tx.Paginate(3, 10)
		q, _ = query.ToSQL(user)
		a.Equal(fmt.Sprintf("%s LIMIT 10 OFFSET 20", s), q)

		// join must come first
		query = pop.Q(tx).Where("id = ?", 1).Join("books b", "b.user_id=?", "xx").Order("name asc")
		q, args := query.ToSQL(user)

		if tx.Dialect.Details().Dialect == "postgres" {
			a.Equal(fmt.Sprintf("%s JOIN books b ON b.user_id=$1 WHERE id = $2 ORDER BY name asc", s), q)
		} else {
			a.Equal(fmt.Sprintf("%s JOIN books b ON b.user_id=? WHERE id = ? ORDER BY name asc", s), q)
		}
		// join arguments comes 1st
		a.Equal(args[0], "xx")
		a.Equal(args[1], 1)

		query = pop.Q(tx)
		q, args = query.ToSQL(user, "distinct on (users.name, users.email) users.*", "users.bio")
		a.Equal("SELECT distinct on (users.name, users.email) users.*, users.bio FROM users AS users", q)

		query = pop.Q(tx)
		q, args = query.ToSQL(user, "distinct on (users.id) users.*", "users.bio")
		a.Equal("SELECT distinct on (users.id) users.*, users.bio FROM users AS users", q)

		query = pop.Q(tx)
		q, args = query.ToSQL(user, "id,r", "users.bio,r", "users.email,w")
		a.Equal("SELECT id, users.bio FROM users AS users", q)

		query = pop.Q(tx)
		q, args = query.ToSQL(user, "distinct on (id) id,r", "users.bio,r", "email,w")
		a.Equal("SELECT distinct on (id) id, users.bio FROM users AS users", q)

		query = pop.Q(tx)
		q, args = query.ToSQL(user, "distinct id", "users.bio,r", "email,w")
		a.Equal("SELECT distinct id, users.bio FROM users AS users", q)

		query = pop.Q(tx)
		q, args = query.ToSQL(user, "distinct id", "concat(users.name,'-',users.email)")
		a.Equal("SELECT concat(users.name,'-',users.email), distinct id FROM users AS users", q)

		query = pop.Q(tx)
		q, args = query.ToSQL(user, "id", "concat(users.name,'-',users.email) name_email")
		a.Equal("SELECT concat(users.name,'-',users.email) name_email, id FROM users AS users", q)

		query = pop.Q(tx)
		q, args = query.ToSQL(user, "distinct id", "concat(users.name,'-',users.email),r")
		a.Equal("SELECT concat(users.name,'-',users.email), distinct id FROM users AS users", q)

		query = pop.Q(tx)
		q, args = query.ToSQL(user, "distinct id", "concat(users.name,'-',users.email) AS x")
		a.Equal("SELECT concat(users.name,'-',users.email) AS x, distinct id FROM users AS users", q)

		query = pop.Q(tx)
		q, args = query.ToSQL(user, "distinct id", "users.name as english_name", "email private_email")
		a.Equal("SELECT distinct id, email private_email, users.name as english_name FROM users AS users", q)
	})
}

func Test_ToSQLInjection(t *testing.T) {
	a := require.New(t)
	transaction(func(tx *pop.Connection) {
		user := &pop.Model{Value: &User{}}
		query := tx.Where("name = '?'", "\\\u0027 or 1=1 limit 1;\n-- ")
		q, _ := query.ToSQL(user)
		a.NotEqual("SELECT * FROM users AS users WHERE name = '\\'' or 1=1 limit 1;\n-- '", q)
	})
}

func Test_ToSQL_RawQuery(t *testing.T) {
	a := require.New(t)
	transaction(func(tx *pop.Connection) {
		query := tx.RawQuery("this is some ? raw ?", "random", "query")
		q, args := query.ToSQL(nil)
		a.Equal(q, tx.Dialect.TranslateSQL("this is some ? raw ?"))
		a.Equal(args, []interface{}{"random", "query"})
	})
}
