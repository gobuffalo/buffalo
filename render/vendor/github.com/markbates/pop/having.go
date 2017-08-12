package pop

import (
	"fmt"
	"strings"
)

type HavingClause struct {
	Condition string
	Arguments []interface{}
}

type havingClauses []HavingClause

func (c HavingClause) String() string {
	sql := fmt.Sprintf("%s", c.Condition)

	return sql
}

func (c havingClauses) String() string {
	if len(c) == 0 {
		return ""
	}

	cs := []string{}
	for _, cl := range c {
		cs = append(cs, cl.String())
	}
	return strings.Join(cs, " AND ")
}
