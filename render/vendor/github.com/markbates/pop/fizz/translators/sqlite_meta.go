package translators

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/pop/fizz"
)

type sqliteIndexListInfo struct {
	Seq     int    `db:"seq"`
	Name    string `db:"name"`
	Unique  bool   `db:"unique"`
	Origin  string `db:"origin"`
	Partial string `db:"partial"`
}

type sqliteIndexInfo struct {
	Seq  int    `db:"seqno"`
	CID  int    `db:"cid"`
	Name string `db:"name"`
}

type sqliteTableInfo struct {
	CID     int         `db:"cid"`
	Name    string      `db:"name"`
	Type    string      `db:"type"`
	NotNull bool        `db:"notnull"`
	Default interface{} `db:"dflt_value"`
	PK      bool        `db:"pk"`
}

func (t sqliteTableInfo) ToColumn() fizz.Column {
	c := fizz.Column{
		Name:    t.Name,
		ColType: t.Type,
		Primary: t.PK,
		Options: fizz.Options{},
	}
	if !t.NotNull {
		c.Options["null"] = true
	}
	if t.Default != nil {
		c.Options["default"] = strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%s", t.Default), "'"), "'")
	}
	return c
}

type sqliteSchema struct {
	URL    string
	db     *sqlx.DB
	schema map[string]*fizz.Table
}

func (p *sqliteSchema) Delete(table string) {
	delete(p.schema, table)
}

func (p *sqliteSchema) TableInfo(table string) (*fizz.Table, error) {
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

func (p *sqliteSchema) buildSchema() error {
	var err error
	p.db, err = sqlx.Open("sqlite3", p.URL)
	if err != nil {
		return err
	}
	defer p.db.Close()

	res, err := p.db.Queryx("SELECT name FROM sqlite_master WHERE type='table';")
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
		if table.Name != "sqlite_sequence" {
			err = p.buildTableData(table)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (p *sqliteSchema) buildTableData(table *fizz.Table) error {
	prag := fmt.Sprintf("PRAGMA table_info(%s)", table.Name)

	res, err := p.db.Queryx(prag)
	if err != nil {
		return nil
	}

	for res.Next() {
		ti := sqliteTableInfo{}
		err = res.StructScan(&ti)
		if err != nil {
			return err
		}
		table.Columns = append(table.Columns, ti.ToColumn())
	}
	err = p.buildTableIndexes(table)
	if err != nil {
		return err
	}
	p.schema[table.Name] = table
	return nil
}

func (p *sqliteSchema) buildTableIndexes(t *fizz.Table) error {
	prag := fmt.Sprintf("PRAGMA index_list(%s)", t.Name)
	res, err := p.db.Queryx(prag)
	if err != nil {
		return err
	}

	for res.Next() {
		li := sqliteIndexListInfo{}
		err = res.StructScan(&li)
		if err != nil {
			return err
		}

		i := fizz.Index{
			Name:    li.Name,
			Unique:  li.Unique,
			Columns: []string{},
		}

		prag = fmt.Sprintf("PRAGMA index_info(%s)", i.Name)
		iires, err := p.db.Queryx(prag)
		if err != nil {
			return err
		}

		for iires.Next() {
			ii := sqliteIndexInfo{}
			err = iires.StructScan(&ii)
			if err != nil {
				return err
			}
			i.Columns = append(i.Columns, ii.Name)
		}

		t.Indexes = append(t.Indexes, i)

	}
	return nil
}
