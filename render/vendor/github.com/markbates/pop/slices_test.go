package pop_test

import (
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/pop/slices"
)

type Cake struct {
	ID        int           `db:"id"`
	Int       slices.Int    `db:"int_slice"`
	Float     slices.Float  `db:"float_slice"`
	String    slices.String `db:"string_slice"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

func (p *PostgreSQLSuite) Test_String() {
	transaction(func(tx *pop.Connection) {
		r := p.Require()

		c := &Cake{
			String: slices.String{"a", "b", "c"},
		}
		err := tx.Create(c)
		r.NoError(err)

		err = tx.Reload(c)
		r.Equal(slices.String{"a", "b", "c"}, c.String)
	})
}

func (p *PostgreSQLSuite) Test_Int() {
	transaction(func(tx *pop.Connection) {
		r := p.Require()

		c := &Cake{
			Int: slices.Int{1, 2, 3},
		}
		err := tx.Create(c)
		r.NoError(err)

		err = tx.Reload(c)
		r.Equal(slices.Int{1, 2, 3}, c.Int)
	})
}

func (p *PostgreSQLSuite) Test_Float() {
	transaction(func(tx *pop.Connection) {
		r := p.Require()

		c := &Cake{
			Float: slices.Float{1.0, 2.1, 3.2},
		}
		err := tx.Create(c)
		r.NoError(err)

		err = tx.Reload(c)
		r.NoError(err)
		r.Equal(slices.Float{1.0, 2.1, 3.2}, c.Float)
	})
}
