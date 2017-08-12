package parser

import "strings"

type errSlice []string

func (e errSlice) Error() string {
	return strings.Join(e, "\n")
}
