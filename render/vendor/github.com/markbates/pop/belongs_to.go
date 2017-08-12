package pop

import "fmt"

// BelongsTo adds a "where" clause based on the "ID" of the
// "model" passed into it.
func (c *Connection) BelongsTo(model interface{}) *Query {
	return Q(c).BelongsTo(model)
}

// BelongsTo adds a "where" clause based on the "ID" of the
// "model" passed into it.
func (q *Query) BelongsTo(model interface{}) *Query {
	m := &Model{Value: model}
	q.Where(fmt.Sprintf("%s = ?", m.associationName()), m.ID())
	return q
}

// BelongsToThrough adds a "where" clause that connects the "bt" model
// through the associated "thru" model.
func (c *Connection) BelongsToThrough(bt, thru interface{}) *Query {
	return Q(c).BelongsToThrough(bt, thru)
}

// BelongsToThrough adds a "where" clause that connects the "bt" model
// through the associated "thru" model.
func (q *Query) BelongsToThrough(bt, thru interface{}) *Query {
	q.belongsToThroughClauses = append(q.belongsToThroughClauses, belongsToThroughClause{
		BelongsTo: &Model{Value: bt},
		Through:   &Model{Value: thru},
	})
	return q
}
