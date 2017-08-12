package pop

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	. "github.com/markbates/pop/columns"
)

type sqlBuilder struct {
	Query      Query
	Model      *Model
	AddColumns []string
	sql        string
	args       []interface{}
}

func newSQLBuilder(q Query, m *Model, addColumns ...string) *sqlBuilder {
	return &sqlBuilder{
		Query:      q,
		Model:      m,
		AddColumns: addColumns,
		args:       []interface{}{},
	}
}

func (sq *sqlBuilder) String() string {
	if sq.sql == "" {
		sq.compile()
	}
	return sq.sql
}

func (sq *sqlBuilder) Args() []interface{} {
	if len(sq.args) == 0 {
		if len(sq.Query.RawSQL.Arguments) > 0 {
			sq.args = sq.Query.RawSQL.Arguments
		} else {
			sq.compile()
		}
	}
	return sq.args
}

func (sq *sqlBuilder) compile() {
	if sq.sql == "" {
		if sq.Query.RawSQL.Fragment != "" {
			sq.sql = sq.Query.RawSQL.Fragment
		} else {
			sq.sql = sq.buildSelectSQL()
		}
		re := regexp.MustCompile(`(?i)in\s*\(\s*\?\s*\)`)
		if re.MatchString(sq.sql) {
			s, _, err := sqlx.In(sq.sql, sq.Args())
			if err == nil {
				sq.sql = s
			}
		}
		sq.sql = sq.Query.Connection.Dialect.TranslateSQL(sq.sql)
	}
}

func (sq *sqlBuilder) buildSelectSQL() string {
	cols := sq.buildColumns()

	fc := sq.buildfromClauses()

	sql := fmt.Sprintf("SELECT %s FROM %s", cols.Readable().SelectString(), fc)

	sql = sq.buildJoinClauses(sql)
	sql = sq.buildWhereClauses(sql)
	sql = sq.buildGroupClauses(sql)
	sql = sq.buildOrderClauses(sql)
	sql = sq.buildPaginationClauses(sql)

	return sql
}

func (sq *sqlBuilder) buildfromClauses() fromClauses {
	models := []*Model{
		sq.Model,
	}
	for _, mc := range sq.Query.belongsToThroughClauses {
		models = append(models, mc.Through)
	}

	fc := sq.Query.fromClauses
	for _, m := range models {
		tableName := m.TableName()
		asName := m.As
		if asName == "" {
			asName = strings.Replace(tableName, ".", "_", -1)
		}
		fc = append(fc, fromClause{
			From: tableName,
			As:   asName,
		})
	}

	return fc
}

func (sq *sqlBuilder) buildWhereClauses(sql string) string {
	mcs := sq.Query.belongsToThroughClauses
	for _, mc := range mcs {
		sq.Query.Where(fmt.Sprintf("%s.%s = ?", mc.Through.TableName(), mc.BelongsTo.associationName()), mc.BelongsTo.ID())
		sq.Query.Where(fmt.Sprintf("%s.id = %s.%s", sq.Model.TableName(), mc.Through.TableName(), sq.Model.associationName()))
	}

	wc := sq.Query.whereClauses
	if len(wc) > 0 {
		sql = fmt.Sprintf("%s WHERE %s", sql, wc.Join(" AND "))
		for _, arg := range wc.Args() {
			sq.args = append(sq.args, arg)
		}
	}
	return sql
}

func (sq *sqlBuilder) buildJoinClauses(sql string) string {
	oc := sq.Query.joinClauses
	if len(oc) > 0 {
		sql += " " + oc.String()
		for i := range oc {
			for _, arg := range oc[i].Arguments {
				sq.args = append(sq.args, arg)
			}
		}
	}

	return sql
}

func (sq *sqlBuilder) buildGroupClauses(sql string) string {
	gc := sq.Query.groupClauses
	if len(gc) > 0 {
		sql = fmt.Sprintf("%s GROUP BY %s", sql, gc.String())

		hc := sq.Query.havingClauses
		if len(hc) > 0 {
			sql = fmt.Sprintf("%s HAVING %s", sql, hc.String())
		}

		for i := range hc {
			for _, arg := range hc[i].Arguments {
				sq.args = append(sq.args, arg)
			}
		}
	}

	return sql
}

func (sq *sqlBuilder) buildOrderClauses(sql string) string {
	oc := sq.Query.orderClauses
	if len(oc) > 0 {
		sql = fmt.Sprintf("%s ORDER BY %s", sql, oc.Join(", "))
		for _, arg := range oc.Args() {
			sq.args = append(sq.args, arg)
		}
	}
	return sql
}

func (sq *sqlBuilder) buildPaginationClauses(sql string) string {
	if sq.Query.limitResults > 0 && sq.Query.Paginator == nil {
		sql = fmt.Sprintf("%s LIMIT %d", sql, sq.Query.limitResults)
	}
	if sq.Query.Paginator != nil {
		sql = fmt.Sprintf("%s LIMIT %d", sql, sq.Query.Paginator.PerPage)
		sql = fmt.Sprintf("%s OFFSET %d", sql, sq.Query.Paginator.Offset)
	}
	return sql
}

var columnCache = map[string]Columns{}

func (sq *sqlBuilder) buildColumns() Columns {
	tableName := sq.Model.TableName()
	acl := len(sq.AddColumns)
	if acl <= 0 {
		cols, ok := columnCache[tableName]
		//if alias is different, remake columns
		if ok && cols.TableAlias == sq.Model.As {
			return cols
		}
		cols = ColumnsForStructWithAlias(sq.Model.Value, tableName, sq.Model.As)
		columnCache[tableName] = cols
		return cols
	} else {
		cols := NewColumns("")
		cols.Add(sq.AddColumns...)
		return cols
	}
}
