package validators_test

import (
	"testing"
	"time"

	"github.com/markbates/validate"
	. "github.com/markbates/validate/validators"
	"github.com/stretchr/testify/require"
)

func Test_TimeIsPresent(t *testing.T) {
	r := require.New(t)
	v := TimeIsPresent{"Created At", time.Now()}
	es := validate.NewErrors()
	v.IsValid(es)
	r.Equal(0, es.Count())

	v = TimeIsPresent{"Created At", time.Time{}}
	v.IsValid(es)
	r.Equal(1, es.Count())
	r.Equal(es.Get("created_at"), []string{"Created At can not be blank."})
}
