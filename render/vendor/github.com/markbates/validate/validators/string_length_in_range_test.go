package validators_test

import (
	"fmt"
	"testing"

	"github.com/markbates/validate"
	. "github.com/markbates/validate/validators"
	"github.com/stretchr/testify/require"
)

func Test_StringLengthInRange(t *testing.T) {
	r := require.New(t)
	var tests = []struct {
		value    string
		min      int
		max      int
		expected bool
	}{
		{"123456", 0, 100, true},
		{"1239999", 0, 0, true},
		{"1239asdfasf99", 100, 200, false},
		{"1239999asdff29", 10, 30, true},
		{"あいうえお", 0, 5, true},
		{"あいうえおか", 0, 5, false},
		{"あいうえお", 0, 0, true},
		{"あいうえ", 5, 10, false},
	}

	for _, test := range tests {
		v := StringLengthInRange{Name: "email", Field: test.value, Min: test.min, Max: test.max}
		errors := validate.NewErrors()
		v.IsValid(errors)
		r.Equal(test.expected, !errors.HasAny(), fmt.Sprintf("Value: %s, Min:%d, Max:%d", test.value, test.min, test.max))
	}
}
