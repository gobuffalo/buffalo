package pop

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/markbates/inflect"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var tableMap = map[string]string{}
var tableMapMu = sync.RWMutex{}

// MapTableName allows for the customize table mapping
// between a name and the database. For example the value
// `User{}` will automatically map to "users".
// MapTableName would allow this to change.
//
//	m := &pop.Model{Value: User{}}
//	m.TableName() // "users"
//
//	pop.MapTableName("user", "people")
//	m = &pop.Model{Value: User{}}
//	m.TableName() // "people"
func MapTableName(name string, tableName string) {
	defer tableMapMu.Unlock()
	tableMapMu.Lock()
	tableMap[name] = tableName
}

type Value interface{}

// Model is used throughout Pop to wrap the end user interface
// that is passed in to many functions.
type Model struct {
	Value
	tableName string
	As        string
}

// ID returns the ID of the Model. All models must have an `ID` field this is
// of type `int`,`int64` or of type `uuid.UUID`.
func (m *Model) ID() interface{} {
	fbn, err := m.fieldByName("ID")
	if err != nil {
		return 0
	}
	if m.PrimaryKeyType() == "UUID" {
		return fbn.Interface().(uuid.UUID).String()
	}
	return fbn.Interface()
}

func (m *Model) PrimaryKeyType() string {
	fbn, err := m.fieldByName("ID")
	if err != nil {
		return "int"
	}
	return fbn.Type().Name()
}

// TableName returns the corresponding name of the underlying database table
// for a given `Model`. See also `MapTableName` to change the default name of
// the table.
func (m *Model) TableName() string {
	if m.tableName != "" {
		return m.tableName
	}

	t := reflect.TypeOf(m.Value)
	name := m.typeName(t)

	defer tableMapMu.Unlock()
	tableMapMu.Lock()

	if tableMap[name] == "" {
		m.tableName = inflect.Tableize(name)
		tableMap[name] = m.tableName
	}
	return tableMap[name]
}

func (m *Model) typeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return m.Value.(string)
	case reflect.Slice, reflect.Array:
		el := t.Elem()
		if el.Kind() == reflect.Ptr {
			el = el.Elem()
		}
		return el.Name()
	default:
		return t.Name()
	}
}

func (m *Model) fieldByName(s string) (reflect.Value, error) {
	el := reflect.ValueOf(m.Value).Elem()
	fbn := el.FieldByName(s)
	if !fbn.IsValid() {
		return fbn, errors.Errorf("Model does not have a field named %s", s)
	}
	return fbn, nil
}

func (m *Model) associationName() string {
	tn := inflect.Singularize(m.TableName())
	return fmt.Sprintf("%s_id", tn)
}

func (m *Model) setID(i interface{}) {
	fbn, err := m.fieldByName("ID")
	if err == nil {
		v := reflect.ValueOf(i)
		switch fbn.Kind() {
		case reflect.Int, reflect.Int64:
			fbn.SetInt(v.Int())
		default:
			fbn.Set(reflect.ValueOf(i))
		}
	}
}

func (m *Model) touchCreatedAt() {
	fbn, err := m.fieldByName("CreatedAt")
	if err == nil {
		fbn.Set(reflect.ValueOf(time.Now()))
	}
}

func (m *Model) touchUpdatedAt() {
	fbn, err := m.fieldByName("UpdatedAt")
	if err == nil {
		fbn.Set(reflect.ValueOf(time.Now()))
	}
}

func (m *Model) whereID() string {
	id := m.ID()
	var value string
	switch id.(type) {
	case int, int64:
		value = fmt.Sprintf("%s.id = %d", m.TableName(), id)
	default:
		value = fmt.Sprintf("%s.id ='%s'", m.TableName(), id)
	}
	return value
}
