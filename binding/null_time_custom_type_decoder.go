package binding

import "github.com/gobuffalo/nulls"

type NullTimeCustomTypeDecoder struct {
	*TimeCustomTypeDecoder
}

func (td NullTimeCustomTypeDecoder) Decode(values []string) (interface{}, error) {
	var ti nulls.Time

	t, err := td.parseTime(values)
	if err != nil {
		return ti, err
	}
	ti.Time = t
	ti.Valid = true

	return ti, nil
}
