package pop

import "fmt"

func (q *Query) Join(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		fmt.Println("Warning: Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"JOIN", table, on, args})
	return q
}

func (q *Query) LeftJoin(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		fmt.Println("Warning: Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"LEFT JOIN", table, on, args})
	return q
}

func (q *Query) RightJoin(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		fmt.Println("Warning: Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"RIGHT JOIN", table, on, args})
	return q
}

func (q *Query) LeftOuterJoin(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		fmt.Println("Warning: Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"LEFT OUTER JOIN", table, on, args})
	return q
}

func (q *Query) RightOuterJoin(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		fmt.Println("Warning: Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"RIGHT OUTER JOIN", table, on, args})
	return q
}

func (q *Query) LeftInnerJoin(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		fmt.Println("Warning: Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"LEFT INNER JOIN", table, on, args})
	return q
}

func (q *Query) RightInnerJoin(table string, on string, args ...interface{}) *Query {
	if q.RawSQL.Fragment != "" {
		fmt.Println("Warning: Query is setup to use raw SQL")
		return q
	}
	q.joinClauses = append(q.joinClauses, joinClause{"RIGHT INNER JOIN", table, on, args})
	return q
}
