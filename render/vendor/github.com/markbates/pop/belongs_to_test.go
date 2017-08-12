package pop_test

import (
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_BelongsTo(t *testing.T) {
	r := require.New(t)

	q := PDB.BelongsTo(&User{ID: 1})

	m := &pop.Model{Value: &Enemy{}}

	sql, _ := q.ToSQL(m)
	r.Equal(ts("SELECT enemies.A FROM enemies AS enemies WHERE user_id = ?"), sql)
}

func Test_BelongsToThrough(t *testing.T) {
	r := require.New(t)

	q := PDB.BelongsToThrough(&User{ID: 1}, &Friend{})
	qs := "SELECT enemies.A FROM enemies AS enemies, good_friends AS good_friends WHERE good_friends.user_id = ? AND enemies.id = good_friends.enemy_id"

	m := &pop.Model{Value: &Enemy{}}
	sql, _ := q.ToSQL(m)
	r.Equal(ts(qs), sql)
}
