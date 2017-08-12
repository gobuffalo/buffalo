package pop

import (
	"fmt"
	"strings"
)

type clause struct {
	Fragment  string
	Arguments []interface{}
}

type clauses []clause

func (c clauses) Join(sep string) string {
	out := make([]string, 0, len(c))
	for _, clause := range c {
		out = append(out, clause.Fragment)
	}
	return strings.Join(out, sep)
}

func (c clauses) Args() (args []interface{}) {
	for _, clause := range c {
		for _, arg := range clause.Arguments {
			args = append(args, arg)
		}
	}
	return
}

type fromClause struct {
	From string
	As   string
}

type fromClauses []fromClause

func (c fromClause) String() string {
	return fmt.Sprintf("%s AS %s", c.From, c.As)
}

func (c fromClauses) String() string {
	cs := []string{}
	for _, cl := range c {
		cs = append(cs, cl.String())
	}
	return strings.Join(cs, ", ")
}

type belongsToThroughClause struct {
	BelongsTo *Model
	Through   *Model
}

type belongsToThroughClauses []belongsToThroughClause
