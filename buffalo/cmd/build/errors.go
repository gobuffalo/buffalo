package build

import "strings"

type multiError []error

func (m multiError) Error() string {
	s := []string{}
	for _, e := range m {
		s = append(s, e.Error())
	}
	return strings.Join(s, "\n")
}
