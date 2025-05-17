package decoders

import "github.com/gobuffalo/nulls"

// NullTimeDecoderFn is a custom type decoder func for null.Time fields
func NullTimeDecoderFn() func([]string) (interface{}, error) {
	return func(vals []string) (interface{}, error) {
		var ti nulls.Time

		// If vals is empty, return a nulls.Time with Valid = false (i.e. NULL).
		// The parseTime() function called below does this check as well, but
		// because it doesn't return an error in the case where vals is empty,
		// we have no way to determine from its response that the nulls.Time
		// should actually be NULL.
		if len(vals) == 0 || vals[0] == "" {
			return ti, nil
		}

		t, err := parseTime(vals)
		if err != nil {
			return ti, err
		}

		ti.Time = t
		ti.Valid = true

		return ti, nil
	}
}
