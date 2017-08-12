# POP [![GoDoc](https://godoc.org/github.com/markbates/pop?status.svg)](https://godoc.org/github.com/markbates/pop) [![Build Status](https://travis-ci.org/markbates/pop.svg)](https://travis-ci.org/markbates/pop)

## A Tasty Treat For All Your Database Needs

So what does Pop do exactly? Well, it wraps the absolutely amazing [https://github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx) library. It cleans up some of the common patterns and workflows usually associated with dealing with databases in Go.

Pop makes it easy to do CRUD operations, run migrations, and build/execute queries. Is Pop an ORM? I'll leave that up to you, the reader, to decide.

Pop, by default, follows conventions that were defined by the ActiveRecord Ruby gem, http://www.rubyonrails.org. What does this mean?

* Tables must have an "id" column and a corresponding "ID" field on the `struct` being used.
* If there is a timestamp column named "created_at", "CreatedAt" on the `struct`, it will be set with the current time when the record is created.
* If there is a timestamp column named "updated_at", "UpdatedAt" on the `struct`, it will be set with the current time when the record is updated.
* Default database table names are lowercase, plural, and underscored versions of the `struct` name. Examples: User{} is "users", FooBar{} is "foo_bars", etc...

## Supported Databases

* PostgreSQL (>= 9.3)
* MySQL (>= 5.7)
* SQLite (>= 3.x)

## Connecting to Databases

Pop is easily configured using a YAML file. The configuration file should be stored in `config/database.yml` or `database.yml`.

#### Example Configuration File

```yaml
development:
  dialect: "postgres"
  database: "your_db_development"
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "postgres"

test:
  dialect: "mysql"
  database: "your_db_test"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

staging:
  dialect: "sqlite3"
  database: "./staging.sqlite"

production:
  dialect: "postgres"
  url: {{ env "DATABASE_URL" }}
```

Note that the `database.yml` file is also a Go template, so you can use Go template syntax. There are two special functions that are included, `env` and `envOr`.

* `env` - This function will look for the named environment variable and insert it into your file. This is useful for configuring production databases without having to store secret information in your repository. `{{ env "DATABASE_URL" }}`
* `envOr` - This function will look for the named environment variable and use it. If the variable can not be found a default value will be used. `{{ envOr "MYSQL_HOST" "localhost" }}`

You can generate a default configuration file using the `init` command:

```
$ soda g config
```

The default will generate a `database.yml` file in the current directory for a PostgreSQL database. You can override the type of database using the `-t` flag and passing in any of the supported database types: `postgres`, `mysql`, or `sqlite3`.

### In your code

Once you have a configuration file defined you can easily connect to one of these connections in your application.

```go
db, err := pop.Connect("development")
if err != nil {
  log.Panic(err)
}
```

Now that you have your connection to the database you can start executing queries against it.

## CLI Support

Pop features CLI support via the `soda` command for the following operations:

* creating databases
* dropping databases
* migrating databases

### Installing CLI Support

```bash
$ go get github.com/markbates/pop/...
$ go install github.com/markbates/pop/soda
```

### Creating Databases

Assuming you defined a configuration file like that described in the above section you can automatically create those databases using the `soda` command:

#### Create All Databases

```bash
$ soda create -a
```

#### Create a Specific Database

```bash
$ soda create -e development
```

### Dropping Databases

Assuming you defined a configuration file like that described in the above section you can automatically drop those databases using the `soda` command:

#### Drop All Databases

```bash
$ soda drop -a
```

#### Drop a Specific Database

```bash
$ soda drop -e development
```

### Models

The `soda` command supports the generation of models.

A full list of commands available for model generation can be found by asking for help:

```bash
$ soda generate help
```

#### Generate Models

The `soda` command will generate Go models and, optionally, the associated migrations for you.

```bash
$ soda generate model user name:text email:text
```

Running this command with generate the following files:

```text
models/user.go
models/user_test.go
migrations/20170115024143_create_users.up.fizz
migrations/20170115024143_create_users.down.fizz
```

The `models/user.go` file contains a structure named `User` with fields `ID`, `CreatedAt`, `UpdatedAt`, `Name`, and `Email`. The first three correspond to the columns commonly found in ActiveRecord models as mentioned before, and the last two correspond to the additional fields specified on the command line. The known types are:

* `text` (`string` in Go)
* `time` or `timestamp` (`time.Time`)
* `nulls.Text` (`nulls.String`) which corresponds to a nullifyable string, which can be distinguished from an empty string
* `uuid` (`uuid.UUID`)
* Other types are passed thru and are used as [Fizz](./fizz/README.md) types.

The `models/user_test.go` contains tests for the User model and they must be implemented by you.

The other two files correspond to the migrations as explained below.

### Migrations

The `soda` command supports the creation and running of migrations.

A full list of commands available for migration can be found by asking for help:

```bash
$ soda migrate help
```

#### Create Migrations

The `soda` command will generate SQL migrations (both the up and down) files for you.

```bash
$ soda migrate create name_of_migration
```

Running this command with generate the following files:

```text
./migrations/20160815134952_name_of_migration.up.fizz
./migrations/20160815134952_name_of_migration.down.fizz
```

The generated files are `fizz` files. [Fizz](./fizz/README.md) lets you use a common DSL for generating migrations. This means the same `.fizz` file can be run against any of the supported dialects of Pop! Find out more about [Fizz](./fizz/README.md)

If you want to generate old fashion `.sql` files you can use the `-t` flag for that:

```bash
$ soda migrate create name_of_migration -t sql
```

Running this command with generate the following files:

```text
./migrations/20160815134952_name_of_migration.up.sql
./migrations/20160815134952_name_of_migration.down.sql
```

The `soda migrate` command supports both `.fizz` and `.sql` files, so you can mix and match them to suit your needs.

#### Running Migrations

The `soda` command will run the migrations using the following command:

```bash
$ soda migrate up
```

Migrations will be run in sequential order. The previously run migrations will be kept track of in a table named `schema_migrations` in the database.

Migrations can also be run in reverse to rollback the schema.

```bash
$ soda migrate down
```

#### Find
```
user := models.User{}
err := tx.Find(&user, id)
```

#### Query
```
tx := models.DB
query = tx.Where("id = 1").Where("name = 'Mark'")
users := []models.User{}
err := query.All(&users)

err = tx.Where("id in (?)", 1, 2, 3).All(&users)
```

##### Join Query

```
// page: page number
// perpage: limit
roles := []models.UserRole{}
query := models.DB.LeftJoin("roles", "roles.id=user_roles.role_id").
  LeftJoin("users u", "u.id=user_roles.user_id").
  Where(`roles.name like ?`, name).Paginate(page, perpage)

count, _ := query.Count(models.UserRole{})
count, _ := query.CountByField(models.UserRole{}, "*")
sql, args := query.ToSQL(&pop.Model{Value: models.UserRole{}}, "user_roles.*",
  "roles.name as role_name", "u.first_name", "u.last_name")
//log.Printf("sql: %s, args: %v", sql, args)
err := models.DB.RawQuery(sql, args...).All(&roles)
```

