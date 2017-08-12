package pop

import (
	"github.com/markbates/validate"
	"github.com/pkg/errors"
)

type validateable interface {
	Validate(*Connection) (*validate.Errors, error)
}

func (m *Model) validate(c *Connection) (*validate.Errors, error) {
	if x, ok := m.Value.(validateable); ok {
		return x.Validate(c)
	}
	return validate.NewErrors(), nil
}

type validateCreateable interface {
	ValidateCreate(*Connection) (*validate.Errors, error)
}

func (m *Model) validateCreate(c *Connection) (*validate.Errors, error) {
	verrs, err := m.validate(c)
	if err != nil {
		return verrs, errors.WithStack(err)
	}
	if x, ok := m.Value.(validateCreateable); ok {
		vs, err := x.ValidateCreate(c)
		if vs != nil {
			verrs.Append(vs)
		}
		if err != nil {
			return verrs, errors.WithStack(err)
		}
	}

	return verrs, err
}

type validateSaveable interface {
	ValidateSave(*Connection) (*validate.Errors, error)
}

func (m *Model) validateSave(c *Connection) (*validate.Errors, error) {
	verrs, err := m.validate(c)
	if err != nil {
		return verrs, errors.WithStack(err)
	}
	if x, ok := m.Value.(validateSaveable); ok {
		vs, err := x.ValidateSave(c)
		if vs != nil {
			verrs.Append(vs)
		}
		if err != nil {
			return verrs, errors.WithStack(err)
		}
	}

	return verrs, err
}

type validateUpdateable interface {
	ValidateUpdate(*Connection) (*validate.Errors, error)
}

func (m *Model) validateUpdate(c *Connection) (*validate.Errors, error) {
	verrs, err := m.validate(c)
	if err != nil {
		return verrs, errors.WithStack(err)
	}
	if x, ok := m.Value.(validateUpdateable); ok {
		vs, err := x.ValidateUpdate(c)
		if vs != nil {
			verrs.Append(vs)
		}
		if err != nil {
			return verrs, errors.WithStack(err)
		}
	}

	return verrs, err
}
