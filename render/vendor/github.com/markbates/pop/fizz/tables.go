package fizz

type Table struct {
	Name    string `db:"name"`
	Columns []Column
	Indexes []Index
}

func (t *Table) Column(name string, colType string, options map[string]interface{}) {
	var primary bool
	if _, ok := options["primary"]; ok {
		primary = true
	}
	c := Column{
		Name:    name,
		ColType: colType,
		Options: options,
		Primary: primary,
	}
	t.Columns = append(t.Columns, c)
}

func (t *Table) ColumnNames() []string {
	cols := make([]string, len(t.Columns))
	for i, c := range t.Columns {
		cols[i] = c.Name
	}
	return cols
}

func (t *Table) HasColumns(args ...string) bool {
	keys := map[string]struct{}{}
	for _, k := range t.ColumnNames() {
		keys[k] = struct{}{}
	}
	for _, a := range args {
		if _, ok := keys[a]; !ok {
			return false
		}
	}
	return true
}

func (f fizzer) CreateTable() interface{} {
	return func(name string, fn func(t *Table)) {
		t := Table{
			Name:    name,
			Columns: []Column{CREATED_COL, UPDATED_COL},
		}
		fn(&t)
		var foundPrimary bool
		for _, c := range t.Columns {
			if c.Primary {
				foundPrimary = true
				break
			}
		}
		if !foundPrimary {
			t.Columns = append(t.Columns, INT_ID_COL)
		}
		f.add(f.Bubbler.CreateTable(t))
	}
}

func (f fizzer) DropTable() interface{} {
	return func(name string) {
		f.add(f.Bubbler.DropTable(Table{Name: name}))
	}
}

func (f fizzer) RenameTable() interface{} {
	return func(old, new string) {
		f.add(f.Bubbler.RenameTable([]Table{
			{Name: old},
			{Name: new},
		}))
	}
}
