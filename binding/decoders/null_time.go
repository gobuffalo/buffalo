package decoders

import "github.com/gobuffalo/nulls"

// NullTimeDecoderFn is a custom type decoder func for null.Time fields
func NullTimeDecoderFn() func([]string) (interface{}, error) {
	return func(vals []string) (interface{}, error) {
		var ti nulls.Time

		t, err := parseTime(vals)
		if err != nil {
			return ti, err
		}

		ti.Time = t
		ti.Valid = true

		return ti, nil
	}
}
