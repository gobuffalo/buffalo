package generate

var pgConfig = `development:
  dialect: postgres
  database: {{.name}}_development
  user: postgres
  password: postgres
  host: 127.0.0.1
  pool: 5

test:
  url: {{"{{"}}envOr "TEST_DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/{{.name}}_test?sslmode=disable"}}

production:
  url: {{"{{"}}envOr "DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/{{.name}}_production?sslmode=disable"}}`

var mysqlConfig = `development:
  dialect: "mysql"
  database: "{{.name}}_development"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

test:
  url: {{"{{"}}envOr "TEST_DATABASE_URL" "mysql://root:root@localhost:3306/{{.name}}_test"}}

production:
  url: {{"{{"}}envOr "DATABASE_URL" "mysql://root:root@localhost:3306/{{.name}}_production"}}`

var sqliteConfig = `development:
  dialect: "sqlite3"
  database: {{"{{"}}env "GOPATH" {{"}}"}}/{{.packagePath}}/{{.name}}_development.sqlite

test:
  dialect: "sqlite3"
  database: {{"{{"}}env "GOPATH" {{"}}"}}/{{.packagePath}}/{{.name}}_test.sqlite

production:
  dialect: "sqlite3"
  database: {{"{{"}}env "GOPATH" {{"}}"}}/{{.packagePath}}/{{.name}}_production.sqlite`

var configTemplates = map[string]string{
	"postgres":   pgConfig,
	"postgresql": pgConfig,
	"pg":         pgConfig,
	"mysql":      mysqlConfig,
	"sqlite3":    sqliteConfig,
	"sqlite":     sqliteConfig,
}
