package pop

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/satori/go.uuid"
)

// Find the first record of the model in the database with a particular id.
//
//	c.Find(&User{}, 1)
func (c *Connection) Find(model interface{}, id interface{}) error {
	return Q(c).Find(model, id)
}

// Find the first record of the model in the database with a particular id.
//
//	q.Find(&User{}, 1)
func (q *Query) Find(model interface{}, id interface{}) error {
	m := &Model{Value: model}
	idq := fmt.Sprintf("%s.id = ?", m.TableName())
	switch t := id.(type) {
	case uuid.UUID:
		return q.Where(idq, t.String()).First(model)
	case string:
		var err error
		id, err = strconv.Atoi(t)
		if err != nil {
			return q.Where(idq, t).First(model)
		}
	}
	return q.Where(idq, id).First(model)
}

// First record of the model in the database that matches the query.
//
//	c.First(&User{})
func (c *Connection) First(model interface{}) error {
	return Q(c).First(model)
}

// First record of the model in the database that matches the query.
//
//	q.Where("name = ?", "mark").First(&User{})
func (q *Query) First(model interface{}) error {
	return q.Connection.timeFunc("First", func() error {
		q.Limit(1)
		m := &Model{Value: model}
		return q.Connection.Dialect.SelectOne(q.Connection.Store, m, *q)
	})
}

// Last record of the model in the database that matches the query.
//
//	c.Last(&User{})
func (c *Connection) Last(model interface{}) error {
	return Q(c).Last(model)
}

// Last record of the model in the database that matches the query.
//
//	q.Where("name = ?", "mark").Last(&User{})
func (q *Query) Last(model interface{}) error {
	return q.Connection.timeFunc("Last", func() error {
		q.Limit(1)
		q.Order("id desc")
		m := &Model{Value: model}
		return q.Connection.Dialect.SelectOne(q.Connection.Store, m, *q)
	})
}

// All retrieves all of the records in the database that match the query.
//
//	c.All(&[]User{})
func (c *Connection) All(models interface{}) error {
	return Q(c).All(models)
}

// All retrieves all of the records in the database that match the query.
//
//	q.Where("name = ?", "mark").All(&[]User{})
func (q *Query) All(models interface{}) error {
	return q.Connection.timeFunc("All", func() error {
		m := &Model{Value: models}
		err := q.Connection.Dialect.SelectMany(q.Connection.Store, m, *q)
		if err == nil && q.Paginator != nil {
			ct, err := q.Count(models)
			if err == nil {
				q.Paginator.TotalEntriesSize = ct
				st := reflect.ValueOf(models).Elem()
				q.Paginator.CurrentEntriesSize = st.Len()
				q.Paginator.TotalPages = (q.Paginator.TotalEntriesSize / q.Paginator.PerPage)
				if q.Paginator.TotalEntriesSize%q.Paginator.PerPage > 0 {
					q.Paginator.TotalPages = q.Paginator.TotalPages + 1
				}
			}
		}
		return err
	})
}

// Exists returns true/false if a record exists in the database that matches
// the query.
//
// 	q.Where("name = ?", "mark").Exists(&User{})
func (q *Query) Exists(model interface{}) (bool, error) {
	i, err := q.Count(model)
	return i != 0, err
}

// Count the number of records in the database.
//
//	c.Count(&User{})
func (c *Connection) Count(model interface{}) (int, error) {
	return Q(c).Count(model)
}

// Count the number of records in the database.
//
//	q.Where("name = ?", "mark").Count(&User{})
func (q Query) Count(model interface{}) (int, error) {
	return q.CountByField(model, "*")
}

func (q Query) CountByField(model interface{}, field string) (int, error) {
	res := &rowCount{}
	err := q.Connection.timeFunc("Count", func() error {
		q.Paginator = nil
		col := fmt.Sprintf("count(%s) as row_count", field)
		q.orderClauses = clauses{}
		query, args := q.ToSQL(&Model{Value: model}, col)
		Log(query, args...)
		return q.Connection.Store.Get(res, query, args...)
	})
	return res.Count, err
}

type rowCount struct {
	Count int `db:"row_count"`
}
