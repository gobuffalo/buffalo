package plush

import (
	"reflect"

	"github.com/pkg/errors"
)

// Iterator type can be implemented and used by the `for` command to build loops in templates
type Iterator interface {
	Next() interface{}
}

type ranger struct {
	pos int
	end int
}

func (r *ranger) Next() interface{} {
	if r.pos < r.end {
		r.pos++
		return r.pos
	}
	return nil
}

func rangeHelper(a, b int) Iterator {
	return &ranger{pos: a - 1, end: b}
}

func betweenHelper(a, b int) Iterator {
	return &ranger{pos: a, end: b - 1}
}

func untilHelper(a int) Iterator {
	return &ranger{pos: -1, end: a - 1}
}

func groupByHelper(size int, underlying interface{}) (*groupBy, error) {
	if size <= 0 {
		return nil, errors.WithStack(errors.New("size must be greater than zero"))
	}
	u := reflect.ValueOf(underlying)
	if u.Kind() == reflect.Ptr {
		u = u.Elem()
	}

	group := []reflect.Value{}
	switch u.Kind() {
	case reflect.Array, reflect.Slice:
		groupSize := u.Len() / size
		if u.Len()%size != 0 {
			groupSize++
		}

		pos := 0
		for pos < u.Len() {
			e := pos + groupSize
			if e > u.Len() {
				e = u.Len()
			}
			group = append(group, u.Slice(pos, e))
			pos += groupSize
		}
	default:
		return nil, errors.WithStack(errors.Errorf("can not use %T in groupBy", underlying))
	}
	g := &groupBy{
		group: group,
	}
	return g, nil
}

type groupBy struct {
	pos   int
	group []reflect.Value
}

func (g *groupBy) Next() interface{} {
	if g.pos >= len(g.group) {
		return nil
	}
	v := g.group[g.pos]
	g.pos++
	return v.Interface()
}
