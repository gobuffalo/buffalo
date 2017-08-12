package translators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/markbates/pop/fizz"
)

type MySQL struct {
	Schema Schema
}

func NewMySQL(url, name string) *MySQL {
	return &MySQL{
		Schema: &mysqlSchema{URL: url, Name: name, schema: map[string]*fizz.Table{}},
	}
}

func (p *MySQL) CreateTable(t fizz.Table) (string, error) {
	sql := []string{}
	cols := []string{}
	for _, c := range t.Columns {
		cols = append(cols, p.buildColumn(c))
		if c.Primary {
			cols = append(cols, fmt.Sprintf("PRIMARY KEY(%s)", c.Name))
		}
	}
	s := fmt.Sprintf("CREATE TABLE %s (\n%s\n) ENGINE=InnoDB;", t.Name, strings.Join(cols, ",\n"))

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

func (p *MySQL) DropTable(t fizz.Table) (string, error) {
	return fmt.Sprintf("DROP TABLE %s;", t.Name), nil
}

func (p *MySQL) RenameTable(t []fizz.Table) (string, error) {
	if len(t) < 2 {
		return "", errors.New("Not enough table names supplied!")
	}
	return fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", t[0].Name, t[1].Name), nil
}

func (p *MySQL) ChangeColumn(t fizz.Table) (string, error) {
	if len(t.Columns) == 0 {
		return "", errors.New("Not enough columns supplied!")
	}
	c := t.Columns[0]
	s := fmt.Sprintf("ALTER TABLE %s MODIFY %s;", t.Name, p.buildColumn(c))
	return s, nil
}

func (p *MySQL) AddColumn(t fizz.Table) (string, error) {
	if len(t.Columns) == 0 {
		return "", errors.New("Not enough columns supplied!")
	}
	c := t.Columns[0]
	s := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s;", t.Name, p.buildColumn(c))
	return s, nil
}

func (p *MySQL) DropColumn(t fizz.Table) (string, error) {
	if len(t.Columns) == 0 {
		return "", errors.New("Not enough columns supplied!")
	}
	c := t.Columns[0]
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", t.Name, c.Name), nil
}

func (p *MySQL) RenameColumn(t fizz.Table) (string, error) {
	if len(t.Columns) < 2 {
		return "", errors.New("Not enough columns supplied!")
	}
	oc := t.Columns[0]
	nc := t.Columns[1]

	ti, err := p.Schema.TableInfo(t.Name)
	if err != nil {
		return "", err
	}
	var c fizz.Column
	for _, c = range ti.Columns {
		if c.Name == oc.Name {
			break
		}
	}
	col := p.buildColumn(c)
	col = strings.Replace(col, oc.Name, fmt.Sprintf("%s %s", oc.Name, nc.Name), -1)
	s := fmt.Sprintf("ALTER TABLE %s CHANGE %s;", t.Name, col)
	return s, nil
}

func (p *MySQL) AddIndex(t fizz.Table) (string, error) {
	if len(t.Indexes) == 0 {
		return "", errors.New("Not enough indexes supplied!")
	}
	i := t.Indexes[0]
	s := fmt.Sprintf("CREATE INDEX %s ON %s (%s);", i.Name, t.Name, strings.Join(i.Columns, ", "))
	if i.Unique {
		s = strings.Replace(s, "CREATE", "CREATE UNIQUE", 1)
	}
	return s, nil
}

func (p *MySQL) DropIndex(t fizz.Table) (string, error) {
	if len(t.Indexes) == 0 {
		return "", errors.New("Not enough indexes supplied!")
	}
	i := t.Indexes[0]
	return fmt.Sprintf("DROP INDEX %s ON %s;", i.Name, t.Name), nil
}

func (p *MySQL) RenameIndex(t fizz.Table) (string, error) {
	ix := t.Indexes
	if len(ix) < 2 {
		return "", errors.New("Not enough indexes supplied!")
	}
	oi := ix[0]
	ni := ix[1]
	return fmt.Sprintf("ALTER TABLE %s RENAME INDEX %s TO %s;", t.Name, oi.Name, ni.Name), nil
}

func (p *MySQL) buildColumn(c fizz.Column) string {
	s := fmt.Sprintf("%s %s", c.Name, p.colType(c))
	if c.Options["null"] == nil || c.Primary {
		s = fmt.Sprintf("%s NOT NULL", s)
	}
	if c.Options["default"] != nil {
		d := fmt.Sprintf("%#v", c.Options["default"])
		re := regexp.MustCompile("^(\")(.+)(\")$")
		d = re.ReplaceAllString(d, "'$2'")
		s = fmt.Sprintf("%s DEFAULT %s", s, d)
	}
	if c.Options["default_raw"] != nil {
		d := fmt.Sprintf("%s", c.Options["default_raw"])
		s = fmt.Sprintf("%s DEFAULT %s", s, d)
	}

	if c.Primary && c.ColType == "integer" {
		s = fmt.Sprintf("%s AUTO_INCREMENT", s)
	}
	return s
}

func (p *MySQL) colType(c fizz.Column) string {
	switch strings.ToLower(c.ColType) {
	case "string":
		s := "255"
		if c.Options["size"] != nil {
			s = fmt.Sprintf("%d", c.Options["size"])
		}
		return fmt.Sprintf("VARCHAR (%s)", s)
	case "uuid":
		return "char(36)"
	case "timestamp", "time", "datetime":
		return "DATETIME"
	default:
		return c.ColType
	}
}
