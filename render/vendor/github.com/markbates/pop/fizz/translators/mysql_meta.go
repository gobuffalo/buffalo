package translators

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/pop/fizz"
)

type mysqlTableInfo struct {
	Field   string      `db:"Field"`
	Type    string      `db:"Type"`
	Null    string      `db:"Null"`
	Key     string      `db:"Key"`
	Default interface{} `db:"Default"`
	Extra   string      `db:"Extra"`
}

func (ti mysqlTableInfo) ToColumn() fizz.Column {
	c := fizz.Column{
		Name:    ti.Field,
		ColType: ti.Type,
		Primary: ti.Key == "PRI",
		Options: map[string]interface{}{},
	}
	if strings.ToLower(ti.Null) == "yes" {
		c.Options["null"] = true
	}
	if ti.Default != nil {
		d := fmt.Sprintf("%s", ti.Default)
		c.Options["default"] = d
	}
	return c
}

type mysqlSchema struct {
	URL    string
	Name   string
	db     *sqlx.DB
	schema map[string]*fizz.Table
}

func (p *mysqlSchema) Delete(table string) {
	delete(p.schema, table)
}

func (p *mysqlSchema) TableInfo(table string) (*fizz.Table, error) {
	if ti, ok := p.schema[table]; ok {
		return ti, nil
	}
	err := p.buildSchema()
	if err != nil {
		return nil, err
	}
	if ti, ok := p.schema[table]; ok {
		return ti, nil
	}
	return nil, fmt.Errorf("Could not find table data for %s!", table)
}

func (p *mysqlSchema) buildSchema() error {
	var err error
	p.db, err = sqlx.Open("mysql", p.URL)
	if err != nil {
		return err
	}
	defer p.db.Close()

	res, err := p.db.Queryx(fmt.Sprintf("select TABLE_NAME as name from information_schema.TABLES where TABLE_SCHEMA = '%s'", p.Name))
	if err != nil {
		return err
	}
	for res.Next() {
		table := &fizz.Table{
			Columns: []fizz.Column{},
			Indexes: []fizz.Index{},
		}
		err = res.StructScan(table)
		if err != nil {
			return err
		}
		err = p.buildTableData(table)
		if err != nil {
			return err
		}

	}
	return nil
}

func (p *mysqlSchema) buildTableData(table *fizz.Table) error {
	prag := fmt.Sprintf("describe %s", table.Name)

	res, err := p.db.Queryx(prag)
	if err != nil {
		return nil
	}

	for res.Next() {
		ti := mysqlTableInfo{}
		err = res.StructScan(&ti)
		if err != nil {
			return err
		}
		table.Columns = append(table.Columns, ti.ToColumn())
	}

	p.schema[table.Name] = table
	return nil
}
