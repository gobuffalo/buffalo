package pop

import "fmt"

func (q *Query) GroupBy(field string, fields ...string) *Query {
	if q.RawSQL.Fragment != "" {
		fmt.Println("Warning: Query is setup to use raw SQL")
		return q
	}
	q.groupClauses = append(q.groupClauses, GroupClause{field})
	if len(fields) > 0 {
		for i := range fields {
			q.groupClauses = append(q.groupClauses, GroupClause{fields[i]})
		}
	}
	return q
}
