package pop

import "fmt"

func (q *Query) Having(condition string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		fmt.Println("Warning: Query is setup to use raw SQL")
		return q
	}
	q.havingClauses = append(q.havingClauses, HavingClause{condition, args})

	return q
}
