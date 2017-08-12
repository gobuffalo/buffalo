package pop_test

import (
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_Callbacks(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		user := &CallbacksUser{
			BeforeS: "BS",
			BeforeC: "BC",
			BeforeU: "BU",
			BeforeD: "BD",
			AfterS:  "AS",
			AfterC:  "AC",
			AfterU:  "AU",
			AfterD:  "AD",
		}

		err := tx.Save(user)
		a.NoError(err)

		a.Equal("BeforeSave", user.BeforeS)
		a.Equal("BeforeCreate", user.BeforeC)
		a.Equal("AfterSave", user.AfterS)
		a.Equal("AfterCreate", user.AfterC)
		a.Equal("BU", user.BeforeU)
		a.Equal("AU", user.AfterU)

		err = tx.Update(user)
		a.NoError(err)

		a.Equal("BeforeUpdate", user.BeforeU)
		a.Equal("AfterUpdate", user.AfterU)
		a.Equal("BD", user.BeforeD)
		a.Equal("AD", user.AfterD)

		err = tx.Destroy(user)
		a.NoError(err)

		a.Equal("BeforeDestroy", user.BeforeD)
		a.Equal("AfterDestroy", user.AfterD)
	})
}
