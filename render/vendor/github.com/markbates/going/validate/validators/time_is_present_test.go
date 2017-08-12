package validators_test

import (
	"testing"
	"time"

	"github.com/markbates/going/validate"
	. "github.com/markbates/going/validate/validators"
	"github.com/stretchr/testify/assert"
)

func Test_TimeIsPresent(t *testing.T) {
	a := assert.New(t)
	v := TimeIsPresent{"Created At", time.Now()}
	es := validate.NewErrors()
	v.IsValid(es)
	a.Equal(0, es.Count())

	v = TimeIsPresent{"Created At", time.Time{}}
	v.IsValid(es)
	a.Equal(1, es.Count())
	a.Equal(es.Get("created_at"), []string{"Created At can not be blank."})
}
