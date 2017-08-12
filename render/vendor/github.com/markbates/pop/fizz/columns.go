package fizz

var INT_ID_COL = Column{
	Name:    "id",
	Primary: true,
	ColType: "integer",
	Options: Options{},
}

var UUID_ID_COL = Column{
	Name:    "id",
	Primary: true,
	ColType: "uuid",
	Options: Options{},
}

var CREATED_COL = Column{Name: "created_at", ColType: "timestamp", Options: Options{}}
var UPDATED_COL = Column{Name: "updated_at", ColType: "timestamp", Options: Options{}}

type Column struct {
	Name    string
	ColType string
	Primary bool
	Options map[string]interface{}
}

func (f fizzer) ChangeColumn() interface{} {
	return func(table, name, ctype string, options Options) {
		t := Table{
			Name: table,
			Columns: []Column{
				{Name: name, ColType: ctype, Options: options},
			},
		}
		f.add(f.Bubbler.ChangeColumn(t))
	}
}

func (f fizzer) AddColumn() interface{} {
	return func(table, name, ctype string, options Options) {
		t := Table{
			Name: table,
			Columns: []Column{
				{Name: name, ColType: ctype, Options: options},
			},
		}
		f.add(f.Bubbler.AddColumn(t))
	}
}

func (f fizzer) DropColumn() interface{} {
	return func(table, name string) {
		t := Table{
			Name: table,
			Columns: []Column{
				{Name: name},
			},
		}
		f.add(f.Bubbler.DropColumn(t))
	}
}

func (f fizzer) RenameColumn() interface{} {
	return func(table, old, new string) error {
		t := Table{
			Name: table,
			Columns: []Column{
				{Name: old},
				{Name: new},
			},
		}
		return f.add(f.Bubbler.RenameColumn(t))
	}
}
