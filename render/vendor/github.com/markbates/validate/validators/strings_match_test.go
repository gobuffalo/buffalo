package validators_test

import (
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_StringsMatch_IsValid(t *testing.T) {
	r := require.New(t)
	var cases = []struct {
		str1     string
		str2     string
		expected bool
	}{
		{"test", "test", true},
		{"test_fail", "test_true", false},
		{"test with space", " test with space ", true},
		{" test with space second", " test with space second       ", true},
	}

	for _, test_case := range cases {
		v := validators.StringsMatch{Name: "strings", Field: test_case.str1, Field2: test_case.str2}
		errors := validate.NewErrors()
		v.IsValid(errors)
		r.Equal(test_case.expected, !errors.HasAny(), "Str1: %s, Str2: %s", test_case.str1, test_case.str2)
	}
}

func BenchmarkStringsMatch_IsValid_Valid(b *testing.B) {
	errors := validate.NewErrors()
	for i := 0; i <= b.N; i++ {
		v := validators.StringsMatch{Name: "strings", Field: " Some string ", Field2: " Some string "}
		v.IsValid(errors)
	}
}

func BenchmarkStringsMatch_IsValid_InValid(b *testing.B) {
	errors := validate.NewErrors()
	for i := 0; i <= b.N; i++ {
		v := validators.StringsMatch{Name: "strings", Field: " Some string ", Field2: " Some string failure"}
		v.IsValid(errors)
	}
}
