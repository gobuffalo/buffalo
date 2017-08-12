package pop_test

import (
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_Scopes(t *testing.T) {
	r := require.New(t)
	oql := "SELECT enemies.A FROM enemies AS enemies"

	m := &pop.Model{Value: &Enemy{}}

	q := PDB.Q()
	s, _ := q.ToSQL(m)
	r.Equal(oql, s)

	q.Scope(func(qy *pop.Query) *pop.Query {
		return qy.Where("id = ?", 1)
	})

	s, _ = q.ToSQL(m)
	r.Equal(ts(oql+" WHERE id = ?"), s)
}
