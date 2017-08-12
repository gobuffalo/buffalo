package validators

import (
	"fmt"
	"time"

	"github.com/markbates/validate"
)

type TimeIsPresent struct {
	Name  string
	Field time.Time
}

func (v *TimeIsPresent) IsValid(errors *validate.Errors) {
	t := time.Time{}
	if v.Field.UnixNano() == t.UnixNano() {
		errors.Add(GenerateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}
