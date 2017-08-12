package validators

import (
	"fmt"

	"github.com/markbates/going/defaults"
	"github.com/markbates/validate"
)

type FuncValidator struct {
	Fn      func() bool
	Field   string
	Name    string
	Message string
}

func (f *FuncValidator) IsValid(verrs *validate.Errors) {
	// for backwards compatability
	f.Name = defaults.String(f.Name, f.Field)
	if !f.Fn() {
		verrs.Add(GenerateKey(f.Name), fmt.Sprintf(f.Message, f.Field))
	}
}
