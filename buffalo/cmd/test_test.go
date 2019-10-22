package cmd

import (
	"reflect"
	"testing"
)

func Test_CutArg(t *testing.T) {
	var tests = []struct {
		arg      string
		args     []string
		expected []string
	}{
		{"b", []string{"a", "b", "c"}, []string{"a", "c"}},
		{"--is-not-in-args", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"--foo", []string{"--foo", "--bar", "--baz"}, []string{"--bar", "--baz"}},
		{"--force-migrations", []string{"./actions/", "--force-migrations"}, []string{"./actions/"}},
		{"--force-migrations", []string{"./actions/", "--force-migrations", "-m", "Test_HomeHandler"}, []string{"./actions/", "-m", "Test_HomeHandler"}},
	}

	for _, tt := range tests {
		result := cutArg(tt.arg, tt.args)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("got %s, want %s when cutting %s from %s", result, tt.expected, tt.arg, tt.args)
		}
	}
}
