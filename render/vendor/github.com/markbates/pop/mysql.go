package pop

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
	"github.com/pkg/errors"
)

var _ dialect = &mysql{}

type mysql struct {
	ConnectionDetails *ConnectionDetails
}

func (m *mysql) Details() *ConnectionDetails {
	return m.ConnectionDetails
}

func (m *mysql) URL() string {
	c := m.ConnectionDetails
	if m.ConnectionDetails.URL != "" {
		return m.ConnectionDetails.URL
	}
	s := "%s:%s@(%s:%s)/%s?parseTime=true&multiStatements=true&readTimeout=1s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database)
}

func (m *mysql) MigrationURL() string {
	return m.URL()
}

func (m *mysql) Create(s store, model *Model, cols Columns) error {
	return errors.Wrap(genericCreate(s, model, cols), "mysql create")
}

func (m *mysql) Update(s store, model *Model, cols Columns) error {
	return errors.Wrap(genericUpdate(s, model, cols), "mysql update")
}

func (m *mysql) Destroy(s store, model *Model) error {
	return errors.Wrap(genericDestroy(s, model), "mysql destroy")
}

func (m *mysql) SelectOne(s store, model *Model, query Query) error {
	return errors.Wrap(genericSelectOne(s, model, query), "mysql select one")
}

func (m *mysql) SelectMany(s store, models *Model, query Query) error {
	return errors.Wrap(genericSelectMany(s, models, query), "mysql select many")
}

func (m *mysql) CreateDB() error {
	c := m.ConnectionDetails
	cmd := exec.Command("mysql", "-u", c.User, "-p"+c.Password, "-h", c.Host, "-P", c.Port, "-e", fmt.Sprintf("create database `%s`", c.Database))
	Log(strings.Join(cmd.Args, " "))
	comboOut, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%s: %s", err.Error(), string(comboOut))
		return errors.Wrapf(err, "error creating MySQL database %s", c.Database)
	}
	fmt.Printf("created database %s\n", c.Database)
	return nil
}

func (m *mysql) DropDB() error {
	c := m.ConnectionDetails
	cmd := exec.Command("mysql", "-u", c.User, "-p"+c.Password, "-h", c.Host, "-P", c.Port, "-e", fmt.Sprintf("drop database `%s`", c.Database))
	Log(strings.Join(cmd.Args, " "))
	comboOut, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%s: %s", err.Error(), string(comboOut))
		return errors.Wrapf(err, "error dropping MySQL database %s", c.Database)
	}
	fmt.Printf("dropped database %s\n", c.Database)
	return nil
}

func (m *mysql) TranslateSQL(sql string) string {
	return sql
}

func (m *mysql) FizzTranslator() fizz.Translator {
	t := translators.NewMySQL(m.URL(), m.Details().Database)
	return t
}

func (m *mysql) Lock(fn func() error) error {
	return fn()
}

func (m *mysql) DumpSchema(w io.Writer) error {
	deets := m.Details()
	cmd := exec.Command("mysqldump", "-d", "-h", deets.Host, "-P", deets.Port, "-u", deets.User, fmt.Sprintf("--password=%s", deets.Password), deets.Database)
	Log(strings.Join(cmd.Args, " "))
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("dumped schema for %s\n", m.Details().Database)
	return nil
}

func (m *mysql) LoadSchema(r io.Reader) error {
	deets := m.Details()
	cmd := exec.Command("mysql", "-u", deets.User, fmt.Sprintf("--password=%s", deets.Password), "-h", deets.Host, "-P", deets.Port, "-D", deets.Database)
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer in.Close()
		io.Copy(in, r)
	}()
	Log(strings.Join(cmd.Args, " "))
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	fmt.Printf("loaded schema for %s\n", m.Details().Database)
	return nil
}

func (m *mysql) TruncateAll(tx *Connection) error {
	stmts := []struct {
		Stmt string `db:"stmt"`
	}{}
	err := tx.RawQuery(mysqlTruncate, m.Details().Database).All(&stmts)
	if err != nil {
		return err
	}
	if len(stmts) == 0 {
		return nil
	}
	qs := []string{}
	for _, x := range stmts {
		qs = append(qs, x.Stmt)
	}
	return tx.RawQuery(strings.Join(qs, " ")).Exec()
}

func newMySQL(deets *ConnectionDetails) dialect {
	cd := &mysql{
		ConnectionDetails: deets,
	}

	return cd
}

const mysqlTruncate = "SELECT concat('TRUNCATE TABLE `', TABLE_NAME, '`;') as stmt FROM INFORMATION_SCHEMA.TABLES where TABLE_SCHEMA = ?"
