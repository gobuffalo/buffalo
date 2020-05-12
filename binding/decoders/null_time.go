package decoders

import "github.com/gobuffalo/nulls"

func NullTimeDecoderFn(formats []string) func([]string) (interface{}, error) {
	return func(vals []string) (interface{}, error) {
		var ti nulls.Time

		t, err := parseTime(vals, formats)
		if err != nil {
			return ti, err
		}

		ti.Time = t
		ti.Valid = true

		return ti, nil
	}
}
