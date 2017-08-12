package translators

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/markbates/pop/fizz"
)

type SQLite struct {
	Schema Schema
}

func NewSQLite(url string) *SQLite {
	return &SQLite{
		Schema: &sqliteSchema{
			URL:    url,
			schema: map[string]*fizz.Table{},
		},
	}
}

func (p *SQLite) CreateTable(t fizz.Table) (string, error) {
	sql := []string{}
	cols := []string{}
	var s string
	for _, c := range t.Columns {
		if c.Primary {
			switch strings.ToLower(c.ColType) {
			case "integer":
				s = fmt.Sprintf("\"%s\" INTEGER PRIMARY KEY AUTOINCREMENT", c.Name)
			case "uuid", "string":
				s = fmt.Sprintf("\"%s\" TEXT PRIMARY KEY", c.Name)
			}
		} else {
			s = p.buildColumn(c)
		}
		cols = append(cols, s)
	}
	s = fmt.Sprintf("CREATE TABLE \"%s\" (\n%s\n);", t.Name, strings.Join(cols, ",\n"))
	sql = append(sql, s)

	for _, i := range t.Indexes {
		s, err := p.AddIndex(fizz.Table{
			Name:    t.Name,
			Indexes: []fizz.Index{i},
		})
		if err != nil {
			return "", err
		}
		sql = append(sql, s)
	}

	return strings.Join(sql, "\n"), nil
}

func (p *SQLite) DropTable(t fizz.Table) (string, error) {
	p.Schema.Delete(t.Name)
	s := fmt.Sprintf("DROP TABLE \"%s\";", t.Name)
	return s, nil
}

func (p *SQLite) RenameTable(t []fizz.Table) (string, error) {
	if len(t) < 2 {
		return "", errors.New("Not enough table names supplied!")
	}
	oldName := t[0].Name
	newName := t[1].Name
	tableInfo, err := p.Schema.TableInfo(oldName)
	if err != nil {
		return "", err
	}
	tableInfo.Name = newName
	s := fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\";", oldName, newName)
	return s, nil
}

func (p *SQLite) ChangeColumn(t fizz.Table) (string, error) {
	tableInfo, err := p.Schema.TableInfo(t.Name)

	if err != nil {
		return "", err
	}

	for i := range tableInfo.Columns {
		if tableInfo.Columns[i].Name == t.Columns[0].Name {
			tableInfo.Columns[i] = t.Columns[0]
			break
		}
	}

	sql := []string{}
	s, err := p.withTempTable(t.Name, func(tempTable fizz.Table) (string, error) {
		createTableSQL, err := p.CreateTable(*tableInfo)
		if err != nil {
			return "", err
		}

		ins := fmt.Sprintf("INSERT INTO \"%s\" (%s) SELECT %s FROM \"%s\";", t.Name, strings.Join(tableInfo.ColumnNames(), ", "), strings.Join(tableInfo.ColumnNames(), ", "), tempTable.Name)
		return strings.Join([]string{createTableSQL, ins}, "\n"), nil
	})

	if err != nil {
		return "", err
	}

	sql = append(sql, s)

	return strings.Join(sql, "\n"), nil
}

func (p *SQLite) AddColumn(t fizz.Table) (string, error) {
	if len(t.Columns) == 0 {
		return "", errors.New("Not enough columns supplied!")
	}
	c := t.Columns[0]

	tableInfo, err := p.Schema.TableInfo(t.Name)
	if err != nil {
		return "", err
	}

	tableInfo.Columns = append(tableInfo.Columns, c)

	s := fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN %s;", t.Name, p.buildColumn(c))
	return s, nil
}

func (p *SQLite) DropColumn(t fizz.Table) (string, error) {
	if len(t.Columns) < 1 {
		return "", errors.New("Not enough columns supplied!")
	}

	tableInfo, err := p.Schema.TableInfo(t.Name)
	if err != nil {
		return "", err
	}

	sql := []string{}
	droppedColumn := t.Columns[0]

	newColumns := []fizz.Column{}
	for _, c := range tableInfo.Columns {
		if c.Name != droppedColumn.Name {
			newColumns = append(newColumns, c)
		}
	}
	tableInfo.Columns = newColumns

	newIndexes := []fizz.Index{}
	for _, i := range tableInfo.Indexes {
		s, err := p.DropIndex(fizz.Table{
			Name:    tableInfo.Name,
			Indexes: []fizz.Index{i},
		})
		if err != nil {
			return "", err
		}
		sql = append(sql, s)
		if tableInfo.HasColumns(i.Columns...) {
			newIndexes = append(newIndexes, i)
		}
	}
	tableInfo.Indexes = newIndexes

	s, err := p.withTempTable(t.Name, func(tempTable fizz.Table) (string, error) {
		createTableSQL, err := p.CreateTable(*tableInfo)
		if err != nil {
			return "", err
		}

		s := fmt.Sprintf("INSERT INTO \"%s\" (%s) SELECT %s FROM \"%s\";", tableInfo.Name, strings.Join(tableInfo.ColumnNames(), ", "), strings.Join(tableInfo.ColumnNames(), ", "), tempTable.Name)

		return strings.Join([]string{createTableSQL, s}, "\n"), nil
	})

	if err != nil {
		return "", err
	}
	sql = append(sql, s)

	return strings.Join(sql, "\n"), nil
}

func (p *SQLite) RenameColumn(t fizz.Table) (string, error) {
	if len(t.Columns) < 2 {
		return "", errors.New("Not enough columns supplied!")
	}

	tableInfo, err := p.Schema.TableInfo(t.Name)
	if err != nil {
		return "", err
	}

	oldColumn := t.Columns[0]
	newColumn := t.Columns[1]

	sql := []string{}

	oldColumns := tableInfo.ColumnNames()
	for ic, c := range tableInfo.Columns {
		if c.Name == oldColumn.Name {
			tableInfo.Columns[ic].Name = newColumn.Name
		}
	}

	for _, i := range tableInfo.Indexes {
		s, err := p.DropIndex(fizz.Table{
			Name:    tableInfo.Name,
			Indexes: []fizz.Index{i},
		})
		if err != nil {
			return "", err
		}
		sql = append(sql, s)
		for ic, c := range i.Columns {
			if c == oldColumn.Name {
				i.Columns[ic] = newColumn.Name
			}
		}
	}

	s, err := p.withTempTable(t.Name, func(tempTable fizz.Table) (string, error) {
		createTableSQL, err := p.CreateTable(*tableInfo)
		if err != nil {
			return "", err
		}

		ins := fmt.Sprintf("INSERT INTO \"%s\" (%s) SELECT %s FROM \"%s\";", t.Name, strings.Join(tableInfo.ColumnNames(), ", "), strings.Join(oldColumns, ", "), tempTable.Name)
		return strings.Join([]string{createTableSQL, ins}, "\n"), nil
	})

	if err != nil {
		return "", err
	}

	sql = append(sql, s)

	return strings.Join(sql, "\n"), nil
}

func (p *SQLite) AddIndex(t fizz.Table) (string, error) {
	if len(t.Indexes) == 0 {
		return "", errors.New("Not enough indexes supplied!")
	}
	i := t.Indexes[0]
	s := fmt.Sprintf("CREATE INDEX \"%s\" ON \"%s\" (%s);", i.Name, t.Name, strings.Join(i.Columns, ", "))
	if i.Unique {
		s = strings.Replace(s, "CREATE", "CREATE UNIQUE", 1)
	}
	return s, nil
}

func (p *SQLite) DropIndex(t fizz.Table) (string, error) {
	if len(t.Indexes) == 0 {
		return "", errors.New("Not enough indexes supplied!")
	}
	i := t.Indexes[0]
	s := fmt.Sprintf("DROP INDEX IF EXISTS \"%s\";", i.Name)
	return s, nil
}

func (p *SQLite) RenameIndex(t fizz.Table) (string, error) {
	if len(t.Indexes) < 2 {
		return "", errors.New("Not enough indexes supplied!")
	}

	tableInfo, err := p.Schema.TableInfo(t.Name)
	if err != nil {
		return "", err
	}

	sql := []string{}

	oldIndex := t.Indexes[0]
	newIndex := t.Indexes[1]

	for _, ti := range tableInfo.Indexes {
		if ti.Name == oldIndex.Name {
			ti.Name = newIndex.Name
			newIndex = ti
			break
		}
	}

	s, err := p.DropIndex(fizz.Table{
		Name:    tableInfo.Name,
		Indexes: []fizz.Index{oldIndex},
	})

	if err != nil {
		return "", err
	}

	sql = append(sql, s)

	s, err = p.AddIndex(fizz.Table{
		Name:    t.Name,
		Indexes: []fizz.Index{newIndex},
	})

	if err != nil {
		return "", err
	}

	sql = append(sql, s)

	return strings.Join(sql, "\n"), nil
}

func (p *SQLite) withTempTable(table string, fn func(fizz.Table) (string, error)) (string, error) {
	tempTable := fizz.Table{Name: fmt.Sprintf("_%s_tmp", table)}

	sql := []string{fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\";", table, tempTable.Name)}
	s, err := fn(tempTable)
	if err != nil {
		return "", err
	}
	sql = append(sql, s, fmt.Sprintf("DROP TABLE \"%s\";", tempTable.Name))

	return strings.Join(sql, "\n"), nil
}

func (p *SQLite) buildColumn(c fizz.Column) string {
	s := fmt.Sprintf("\"%s\" %s", c.Name, p.colType(c))
	if c.Options["null"] == nil {
		s = fmt.Sprintf("%s NOT NULL", s)
	}
	if c.Options["default"] != nil {
		s = fmt.Sprintf("%s DEFAULT '%v'", s, c.Options["default"])
	}
	if c.Options["default_raw"] != nil {
		s = fmt.Sprintf("%s DEFAULT %s", s, c.Options["default_raw"])
	}
	return s
}

func (p *SQLite) colType(c fizz.Column) string {
	switch strings.ToLower(c.ColType) {
	case "uuid":
		return "char(36)"
	case "timestamp", "time", "datetime":
		return "DATETIME"
	case "boolean", "date":
		return "NUMERIC"
	case "string":
		return "TEXT"
	default:
		return c.ColType
	}
}
