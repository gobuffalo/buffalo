package columns

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Columns struct {
	Cols       map[string]*Column
	lock       *sync.RWMutex
	TableName  string
	TableAlias string
}

// Add a column to the list.
func (c *Columns) Add(names ...string) []*Column {
	ret := []*Column{}
	c.lock.Lock()

	tableAlias := c.TableAlias
	if tableAlias == "" {
		tableAlias = c.TableName
	}

	for _, name := range names {

		var xs []string
		var col *Column
		ss := ""
		//support for distinct xx, or distinct on (field) table.fields
		if strings.HasSuffix(name, ",r") || strings.HasSuffix(name, ",w") {
			xs = []string{name[0 : len(name)-2], name[len(name)-1 : len(name)]}
		} else {
			xs = []string{name}
		}

		xs[0] = strings.TrimSpace(xs[0])
		//eg: id id2 - select id as id2
		// also distinct columnname
		// and distinct on (column1) column2
		if strings.Contains(strings.ToUpper(xs[0]), " AS ") {
			//eg: select id as id2
			i := strings.LastIndex(strings.ToUpper(xs[0]), " AS ")
			ss = xs[0]
			xs[0] = xs[0][i+4 : len(xs[0])] //get id2
		} else if strings.Contains(xs[0], " ") {
			i := strings.LastIndex(name, " ")
			ss = xs[0]
			xs[0] = xs[0][i+1 : len(xs[0])] //get id2
		}

		col = c.Cols[xs[0]]
		//fmt.Printf("column: %v, col: %v, xs: %v, ss: %v\n", xs[0], col, xs, ss)
		if col == nil {
			if ss == "" {
				ss = xs[0]
				if tableAlias != "" {
					ss = fmt.Sprintf("%s.%s", tableAlias, ss)
				}
			}

			col = &Column{
				Name:      xs[0],
				SelectSQL: ss,
				Readable:  true,
				Writeable: true,
			}

			if len(xs) > 1 {
				if xs[1] == "r" {
					col.Writeable = false
				}
				if xs[1] == "w" {
					col.Readable = false
				}
			} else if col.Name == "id" {
				col.Writeable = false
			}

			c.Cols[col.Name] = col
		}
		ret = append(ret, col)
	}

	c.lock.Unlock()
	return ret
}

// Remove a column from the list.
func (c *Columns) Remove(names ...string) {
	for _, name := range names {
		xs := strings.Split(name, ",")
		name = xs[0]
		delete(c.Cols, name)
	}
}

func (c Columns) Writeable() *WriteableColumns {
	w := &WriteableColumns{NewColumnsWithAlias(c.TableName, c.TableAlias)}
	for _, col := range c.Cols {
		if col.Writeable {
			w.Cols[col.Name] = col
		}
	}
	return w
}

func (c Columns) Readable() *ReadableColumns {
	w := &ReadableColumns{NewColumnsWithAlias(c.TableName, c.TableAlias)}
	for _, col := range c.Cols {
		if col.Readable {
			w.Cols[col.Name] = col
		}
	}
	return w
}

func (c Columns) String() string {
	xs := []string{}
	for _, t := range c.Cols {
		xs = append(xs, t.Name)
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}

func (c Columns) SymbolizedString() string {
	xs := []string{}
	for _, t := range c.Cols {
		xs = append(xs, ":"+t.Name)
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}

func NewColumns(tableName string) Columns {
	return NewColumnsWithAlias(tableName, "")
}

func NewColumnsWithAlias(tableName string, tableAlias string) Columns {
	return Columns{
		lock:       &sync.RWMutex{},
		Cols:       map[string]*Column{},
		TableName:  tableName,
		TableAlias: tableAlias,
	}
}
