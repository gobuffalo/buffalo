package translators_test

import (
	"fmt"

	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
)

var _ fizz.Translator = (*translators.SQLite)(nil)
var schema = &fauxSchema{schema: map[string]*fizz.Table{}}
var sqt = &translators.SQLite{Schema: schema}

type fauxSchema struct {
	schema map[string]*fizz.Table
}

func (p *fauxSchema) Delete(table string) {
	delete(p.schema, table)
}

func (p *fauxSchema) TableInfo(table string) (*fizz.Table, error) {
	if ti, ok := p.schema[table]; ok {
		return ti, nil
	}
	return nil, fmt.Errorf("Could not find table data for %s!", table)
}

func (p *SQLiteSuite) Test_SQLite_CreateTable() {
	r := p.Require()
	ddl := `CREATE TABLE "users" (
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"first_name" TEXT NOT NULL,
"last_name" TEXT NOT NULL,
"email" TEXT NOT NULL,
"permissions" text,
"age" integer DEFAULT '40',
"id" INTEGER PRIMARY KEY AUTOINCREMENT
);`

	res, _ := fizz.AString(`
	create_table("users", func(t) {
		t.Column("first_name", "string", {})
		t.Column("last_name", "string", {})
		t.Column("email", "string", {"size":20})
		t.Column("permissions", "text", {"null": true})
		t.Column("age", "integer", {"null": true, "default": 40})
	})
	`, sqt)
	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_CreateTable_UUID() {
	r := p.Require()
	ddl := `CREATE TABLE "users" (
"created_at" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL,
"first_name" TEXT NOT NULL,
"last_name" TEXT NOT NULL,
"email" TEXT NOT NULL,
"permissions" text,
"age" integer DEFAULT '40',
"company_id" char(36) NOT NULL DEFAULT lower(hex(randomblob(16))),
"uuid" TEXT PRIMARY KEY
);`

	res, _ := fizz.AString(`
	create_table("users", func(t) {
		t.Column("first_name", "string", {})
		t.Column("last_name", "string", {})
		t.Column("email", "string", {"size":20})
		t.Column("permissions", "text", {"null": true})
		t.Column("age", "integer", {"null": true, "default": 40})
		t.Column("company_id", "uuid", {"default_raw": "lower(hex(randomblob(16)))"})
		t.Column("uuid", "uuid", {"primary": true})
	})
	`, sqt)
	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_DropTable() {
	r := p.Require()

	ddl := `DROP TABLE "users";`

	res, _ := fizz.AString(`drop_table("users")`, sqt)
	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_RenameTable() {
	r := p.Require()

	ddl := `ALTER TABLE "users" RENAME TO "people";`
	schema.schema["users"] = &fizz.Table{}

	res, _ := fizz.AString(`rename_table("users", "people")`, sqt)
	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_RenameTable_NotEnoughValues() {
	r := p.Require()

	_, err := sqt.RenameTable([]fizz.Table{})
	r.Error(err)
}

func (p *SQLiteSuite) Test_SQLite_ChangeColumn() {
	r := p.Require()

	ddl := `ALTER TABLE "users" RENAME TO "_users_tmp";
CREATE TABLE "users" (
"id" INTEGER PRIMARY KEY AUTOINCREMENT,
"created_at" TEXT NOT NULL DEFAULT 'foo',
"updated_at" DATETIME NOT NULL
);
INSERT INTO "users" (id, created_at, updated_at) SELECT id, created_at, updated_at FROM "_users_tmp";
DROP TABLE "_users_tmp";`

	schema.schema["users"] = &fizz.Table{
		Name: "users",
		Columns: []fizz.Column{
			fizz.INT_ID_COL,
			fizz.CREATED_COL,
			fizz.UPDATED_COL,
		},
	}

	res, _ := fizz.AString(`change_column("users", "created_at", "string", {"default": "foo", "size": 50})`, sqt)

	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_AddColumn() {
	r := p.Require()

	ddl := `ALTER TABLE "users" ADD COLUMN "mycolumn" TEXT NOT NULL DEFAULT 'foo';`
	schema.schema["users"] = &fizz.Table{}

	res, _ := fizz.AString(`add_column("users", "mycolumn", "string", {"default": "foo", "size": 50})`, sqt)

	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_DropColumn() {
	r := p.Require()
	ddl := `ALTER TABLE "users" RENAME TO "_users_tmp";
CREATE TABLE "users" (
"id" INTEGER PRIMARY KEY AUTOINCREMENT,
"updated_at" DATETIME NOT NULL
);
INSERT INTO "users" (id, updated_at) SELECT id, updated_at FROM "_users_tmp";
DROP TABLE "_users_tmp";`

	schema.schema["users"] = &fizz.Table{
		Name: "users",
		Columns: []fizz.Column{
			fizz.INT_ID_COL,
			fizz.CREATED_COL,
			fizz.UPDATED_COL,
		},
	}
	res, _ := fizz.AString(`drop_column("users", "created_at")`, sqt)

	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_RenameColumn() {
	r := p.Require()
	ddl := `ALTER TABLE "users" RENAME TO "_users_tmp";
CREATE TABLE "users" (
"id" INTEGER PRIMARY KEY AUTOINCREMENT,
"created_when" DATETIME NOT NULL,
"updated_at" DATETIME NOT NULL
);
INSERT INTO "users" (id, created_when, updated_at) SELECT id, created_at, updated_at FROM "_users_tmp";
DROP TABLE "_users_tmp";`

	schema.schema["users"] = &fizz.Table{
		Name: "users",
		Columns: []fizz.Column{
			fizz.INT_ID_COL,
			fizz.CREATED_COL,
			fizz.UPDATED_COL,
		},
	}
	res, _ := fizz.AString(`rename_column("users", "created_at", "created_when")`, sqt)

	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_AddIndex() {
	r := p.Require()
	ddl := `CREATE INDEX "table_name_column_name_idx" ON "table_name" (column_name);`

	res, _ := fizz.AString(`add_index("table_name", "column_name", {})`, sqt)
	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_AddIndex_Unique() {
	r := p.Require()
	ddl := `CREATE UNIQUE INDEX "table_name_column_name_idx" ON "table_name" (column_name);`

	res, _ := fizz.AString(`add_index("table_name", "column_name", {"unique": true})`, sqt)
	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_AddIndex_MultiColumn() {
	r := p.Require()
	ddl := `CREATE INDEX "table_name_col1_col2_col3_idx" ON "table_name" (col1, col2, col3);`

	res, _ := fizz.AString(`add_index("table_name", ["col1", "col2", "col3"], {})`, sqt)
	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_AddIndex_CustomName() {
	r := p.Require()
	ddl := `CREATE INDEX "custom_name" ON "table_name" (column_name);`

	res, _ := fizz.AString(`add_index("table_name", "column_name", {"name": "custom_name"})`, sqt)
	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_DropIndex() {
	r := p.Require()
	ddl := `DROP INDEX IF EXISTS "my_idx";`

	res, _ := fizz.AString(`drop_index("my_table", "my_idx")`, sqt)
	r.Equal(ddl, res)
}

func (p *SQLiteSuite) Test_SQLite_RenameIndex() {
	r := p.Require()

	ddl := `DROP INDEX IF EXISTS "old_ix";
CREATE UNIQUE INDEX "new_ix" ON "users" (id, created_at);`

	schema.schema["users"] = &fizz.Table{
		Name: "users",
		Columns: []fizz.Column{
			fizz.INT_ID_COL,
			fizz.CREATED_COL,
			fizz.UPDATED_COL,
		},
		Indexes: []fizz.Index{
			{
				Name:    "old_ix",
				Columns: []string{"id", "created_at"},
				Unique:  true,
			},
		},
	}

	res, _ := fizz.AString(`rename_index("users", "old_ix", "new_ix")`, sqt)
	r.Equal(ddl, res)
}
