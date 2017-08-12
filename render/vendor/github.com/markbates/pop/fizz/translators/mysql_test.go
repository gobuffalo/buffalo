package translators_test

import (
	"github.com/markbates/pop"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
)

var _ fizz.Translator = (*translators.MySQL)(nil)
var myt = translators.NewMySQL("", "")

func init() {
	myconn, err := pop.Connect("mysql")
	if err != nil {
		panic(err.Error())
	}
	deets := myconn.Dialect.Details()
	myt = translators.NewMySQL(myconn.URL(), deets.Database)
}

func (p *MySQLSuite) Test_MySQL_SchemaMigration() {
	r := p.Require()
	ddl := `CREATE TABLE schema_migrations (
version VARCHAR (255) NOT NULL
) ENGINE=InnoDB;
CREATE UNIQUE INDEX version_idx ON schema_migrations (version);`

	res, err := myt.CreateTable(fizz.Table{
		Name: "schema_migrations",
		Columns: []fizz.Column{
			{Name: "version", ColType: "string"},
		},
		Indexes: []fizz.Index{
			{Name: "version_idx", Columns: []string{"version"}, Unique: true},
		},
	})
	r.NoError(err)
	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_CreateTable() {
	r := p.Require()
	ddl := `CREATE TABLE users (
created_at DATETIME NOT NULL,
updated_at DATETIME NOT NULL,
first_name VARCHAR (255) NOT NULL,
last_name VARCHAR (255) NOT NULL,
email VARCHAR (20) NOT NULL,
permissions text,
age integer DEFAULT 40,
id integer NOT NULL AUTO_INCREMENT,
PRIMARY KEY(id)
) ENGINE=InnoDB;`

	res, _ := fizz.AString(`
	create_table("users", func(t) {
		t.Column("first_name", "string", {})
		t.Column("last_name", "string", {})
		t.Column("email", "string", {"size":20})
		t.Column("permissions", "text", {"null": true})
		t.Column("age", "integer", {"null": true, "default": 40})
	})
	`, myt)
	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_CreateTable_UUID() {
	r := p.Require()
	ddl := `CREATE TABLE users (
created_at DATETIME NOT NULL,
updated_at DATETIME NOT NULL,
first_name VARCHAR (255) NOT NULL,
last_name VARCHAR (255) NOT NULL,
email VARCHAR (20) NOT NULL,
permissions text,
age integer DEFAULT 40,
company_id char(36) NOT NULL DEFAULT 'test',
uuid char(36) NOT NULL,
PRIMARY KEY(uuid)
) ENGINE=InnoDB;`

	res, _ := fizz.AString(`
	create_table("users", func(t) {
		t.Column("first_name", "string", {})
		t.Column("last_name", "string", {})
		t.Column("email", "string", {"size":20})
		t.Column("permissions", "text", {"null": true})
		t.Column("age", "integer", {"null": true, "default": 40})
		t.Column("company_id", "uuid", {"default_raw": "'test'"})
		t.Column("uuid", "uuid", {"primary": true})
	})
	`, myt)
	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_DropTable() {
	r := p.Require()

	ddl := `DROP TABLE users;`

	res, _ := fizz.AString(`drop_table("users")`, myt)
	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_RenameTable() {
	r := p.Require()

	ddl := `ALTER TABLE users RENAME TO people;`

	res, _ := fizz.AString(`rename_table("users", "people")`, myt)
	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_RenameTable_NotEnoughValues() {
	r := p.Require()

	_, err := myt.RenameTable([]fizz.Table{})
	r.Error(err)
}

func (p *MySQLSuite) Test_MySQL_ChangeColumn() {
	r := p.Require()
	ddl := `ALTER TABLE users MODIFY mycolumn VARCHAR (50) NOT NULL DEFAULT 'foo';`

	res, _ := fizz.AString(`change_column("users", "mycolumn", "string", {"default": "foo", "size": 50})`, myt)

	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_AddColumn() {
	r := p.Require()
	ddl := `ALTER TABLE users ADD COLUMN mycolumn VARCHAR (50) NOT NULL DEFAULT 'foo';`

	res, _ := fizz.AString(`add_column("users", "mycolumn", "string", {"default": "foo", "size": 50})`, myt)

	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_DropColumn() {
	r := p.Require()
	ddl := `ALTER TABLE users DROP COLUMN mycolumn;`

	res, _ := fizz.AString(`drop_column("users", "mycolumn")`, myt)

	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_RenameColumn() {
	r := p.Require()
	ddl := `ALTER TABLE users CHANGE email email_address varchar(50) NOT NULL DEFAULT 'foo@example.com';`

	res, _ := fizz.AString(`rename_column("users", "email", "email_address")`, myt)
	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_AddIndex() {
	r := p.Require()
	ddl := `CREATE INDEX users_email_idx ON users (email);`

	res, _ := fizz.AString(`add_index("users", "email", {})`, myt)
	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_AddIndex_Unique() {
	r := p.Require()
	ddl := `CREATE UNIQUE INDEX users_email_idx ON users (email);`

	res, _ := fizz.AString(`add_index("users", "email", {"unique": true})`, myt)
	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_AddIndex_MultiColumn() {
	r := p.Require()
	ddl := `CREATE INDEX users_id_email_idx ON users (id, email);`

	res, _ := fizz.AString(`add_index("users", ["id", "email"], {})`, myt)
	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_AddIndex_CustomName() {
	r := p.Require()
	ddl := `CREATE INDEX email_index ON users (email);`

	res, _ := fizz.AString(`add_index("users", "email", {"name": "email_index"})`, myt)
	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_DropIndex() {
	r := p.Require()
	ddl := `DROP INDEX email_idx ON users;`

	res, _ := fizz.AString(`drop_index("users", "email_idx")`, myt)
	r.Equal(ddl, res)
}

func (p *MySQLSuite) Test_MySQL_RenameIndex() {
	r := p.Require()

	ddl := `ALTER TABLE users RENAME INDEX email_idx TO email_address_ix;`

	res, _ := fizz.AString(`rename_index("users", "email_idx", "email_address_ix")`, myt)
	r.Equal(ddl, res)
}
