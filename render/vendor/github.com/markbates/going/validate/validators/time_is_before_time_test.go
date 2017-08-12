package validators_test

import (
	"testing"
	"time"

	"github.com/markbates/going/validate"
	. "github.com/markbates/going/validate/validators"
	"github.com/stretchr/testify/assert"
)

func Test_TimeIsBeforeTime(t *testing.T) {
	a := assert.New(t)
	now := time.Now()
	v := TimeIsBeforeTime{
		FirstName: "Opens At", FirstTime: now,
		SecondName: "Closes At", SecondTime: now.Add(100000),
	}

	es := validate.NewErrors()
	v.IsValid(es)
	a.Equal(0, es.Count())

	v.SecondTime = now.Add(-100000)
	v.IsValid(es)

	a.Equal(1, es.Count())
	a.Equal(es.Get("opens_at"), []string{"Opens At must be before Closes At."})
}
