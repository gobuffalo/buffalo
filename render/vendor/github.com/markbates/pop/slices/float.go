package slices

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Float []float64

func (f Float) Interface() interface{} {
	return []float64(f)
}

func (s *Float) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	str := string(b)
	(*s) = strToFloat(str, *s)
	return nil
}

func (s Float) Value() (driver.Value, error) {
	sa := make([]string, len(s))
	for x, i := range s {
		sa[x] = strconv.FormatFloat(i, 'f', -1, 64)
	}
	return fmt.Sprintf("{%s}", strings.Join(sa, ",")), nil
}

func (s *Float) UnmarshalText(text []byte) error {
	ss := []float64{}
	for _, x := range strings.Split(string(text), ",") {
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return errors.WithStack(err)
		}
		ss = append(ss, f)
	}
	(*s) = ss
	return nil
}

func strToFloat(s string, a []float64) []float64 {
	r := strings.Trim(s, "{}")
	a = make([]float64, 0, 10)
	for _, t := range strings.Split(r, ",") {
		i, _ := strconv.ParseFloat(t, 64)
		a = append(a, i)
	}
	return a
}
