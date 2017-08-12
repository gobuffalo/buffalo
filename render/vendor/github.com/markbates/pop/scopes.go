package pop

type ScopeFunc func(q *Query) *Query

// Scope the query by using a `ScopeFunc`
//
//	func ByName(name string) ScopeFunc {
//		return func(q *Query) *Query {
//			return q.Where("name = ?", name)
//		}
//	}
//
//	q.Scope(ByName("mark").Where("id = ?", 1).First(&User{})
func (q *Query) Scope(sf ScopeFunc) *Query {
	return sf(q)
}

// Scope the query by using a `ScopeFunc`
//
//	func ByName(name string) ScopeFunc {
//		return func(q *Query) *Query {
//			return q.Where("name = ?", name)
//		}
//	}
//
//	c.Scope(ByName("mark")).First(&User{})
func (c *Connection) Scope(sf ScopeFunc) *Query {
	return Q(c).Scope(sf)
}
