package slices

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Int []int

func (i Int) Interface() interface{} {
	return []int(i)
}

func (s *Int) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	str := string(b)
	(*s) = strToInt(str)
	return nil
}

func (s Int) Value() (driver.Value, error) {
	sa := make([]string, len(s))
	for x, i := range s {
		sa[x] = strconv.Itoa(i)
	}
	return fmt.Sprintf("{%s}", strings.Join(sa, ",")), nil
}

func (s *Int) UnmarshalText(text []byte) error {
	ss := []int{}
	for _, x := range strings.Split(string(text), ",") {
		f, err := strconv.Atoi(x)
		if err != nil {
			return errors.WithStack(err)
		}
		ss = append(ss, f)
	}
	(*s) = ss
	return nil
}

func strToInt(s string) []int {
	r := strings.Trim(s, "{}")
	a := make([]int, 0, 10)
	for _, t := range strings.Split(r, ",") {
		i, _ := strconv.Atoi(t)
		a = append(a, i)
	}
	return a
}
