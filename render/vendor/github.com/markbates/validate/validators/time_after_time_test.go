package validators_test

import (
	"testing"
	"time"

	"github.com/markbates/validate"
	. "github.com/markbates/validate/validators"
	"github.com/stretchr/testify/require"
)

func Test_TimeAfterTime(t *testing.T) {
	r := require.New(t)
	now := time.Now()
	v := TimeAfterTime{
		FirstName: "Opens At", FirstTime: now.Add(100000),
		SecondName: "Now", SecondTime: now,
	}

	es := validate.NewErrors()
	v.IsValid(es)
	r.Equal(0, es.Count())

	v.SecondTime = now.Add(200000)
	v.IsValid(es)

	r.Equal(1, es.Count())
	r.Equal(es.Get("opens_at"), []string{"Opens At must be after Now."})
}
