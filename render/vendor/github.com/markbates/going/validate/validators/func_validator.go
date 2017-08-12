package validators

import (
	"fmt"

	"github.com/markbates/going/validate"
)

type FuncValidator struct {
	Fn      func() bool
	Field   string
	Message string
}

func (f *FuncValidator) IsValid(verrs *validate.Errors) {
	if !f.Fn() {
		verrs.Add(GenerateKey(f.Field), fmt.Sprintf(f.Message, f.Field))
	}
}
