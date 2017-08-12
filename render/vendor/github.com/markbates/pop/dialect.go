package pop

import (
	"encoding/gob"
	"fmt"
	"io"

	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func init() {
	gob.Register(uuid.UUID{})
}

type dialect interface {
	URL() string
	MigrationURL() string
	Details() *ConnectionDetails
	TranslateSQL(string) string
	Create(store, *Model, Columns) error
	Update(store, *Model, Columns) error
	Destroy(store, *Model) error
	SelectOne(store, *Model, Query) error
	SelectMany(store, *Model, Query) error
	CreateDB() error
	DropDB() error
	DumpSchema(io.Writer) error
	LoadSchema(io.Reader) error
	FizzTranslator() fizz.Translator
	Lock(func() error) error
	TruncateAll(*Connection) error
}

func genericCreate(s store, model *Model, cols Columns) error {
	keyType := model.PrimaryKeyType()
	switch keyType {
	case "int", "int64":
		var id int64
		w := cols.Writeable()
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", model.TableName(), w.String(), w.SymbolizedString())
		Log(query)
		res, err := s.NamedExec(query, model.Value)
		if err != nil {
			return errors.Wrapf(err, "create: couldn't execute named statement %s", query)
		}
		id, err = res.LastInsertId()
		if err == nil {
			model.setID(id)
		}
		return errors.Wrap(err, "couldn't get the last inserted id")
	case "UUID":
		model.setID(uuid.NewV4())
		w := cols.Writeable()
		w.Add("id")
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", model.TableName(), w.String(), w.SymbolizedString())
		Log(query)
		stmt, err := s.PrepareNamed(query)
		if err != nil {
			return errors.WithStack(err)
		}
		_, err = stmt.Exec(model.Value)
		return err
	}
	return errors.Errorf("can not use %s as a primary key type!", keyType)
}

func genericUpdate(s store, model *Model, cols Columns) error {
	stmt := fmt.Sprintf("UPDATE %s SET %s where %s", model.TableName(), cols.Writeable().UpdateString(), model.whereID())
	Log(stmt)
	_, err := s.NamedExec(stmt, model.Value)
	return errors.Wrapf(err, "update: couldn't execute named statement %s", stmt)
}

func genericDestroy(s store, model *Model) error {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE %s", model.TableName(), model.whereID())
	return errors.Wrapf(genericExec(s, stmt), "destroy: couldn't execute named statement %s", stmt)
}

func genericExec(s store, stmt string) error {
	Log(stmt)
	_, err := s.Exec(stmt)
	return errors.Wrapf(err, "couldn't execute statement %s", stmt)
}

func genericSelectOne(s store, model *Model, query Query) error {
	sql, args := query.ToSQL(model)
	Log(sql, args...)
	err := s.Get(model.Value, sql, args...)
	return errors.Wrapf(err, "couldn't select one %s %+v", sql, args)
}

func genericSelectMany(s store, models *Model, query Query) error {
	sql, args := query.ToSQL(models)
	Log(sql, args...)
	err := s.Select(models.Value, sql, args...)
	if err != nil {
		return errors.Wrapf(err, "couldn't select many %s %+v", sql, args)
	}
	return nil
}
