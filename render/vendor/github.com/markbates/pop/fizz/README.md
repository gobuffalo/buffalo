# Fizz

## A Common DSL for Migrating Databases

## Create a Table

``` javascript
create_table("users", func(t) {
  t.Column("email", "string", {})
  t.Column("twitter_handle", "string", {"size": 50})
  t.Column("age", "integer", {"default": 0})
  t.Column("admin", "boolean", {"default": false})
  t.Column("company_id", "uuid", {"default_raw": "uuid_generate_v1()"})
  t.Column("bio", "text", {"null": true})
  t.Column("joined_at", "timestamp", {})
})
```

The `create_table` function will also generate an `id` column of type `integer` that will auto-increment. It will also generate two `timestamp` columns; `created_at` and `updated_at`.

Columns all have the same syntax. First is the name of the column. Second is the type of the field. Third is any options you want to set on that column.

#### <a name="column-info"></a> "Common" Types:

* `string`
* `text`
* `timestamp`, `time`, `datetime`
* `integer`
* `boolean`

Any other type passed it will be be passed straight through to the underlying database. For example for PostgreSQL you could pass `jsonb`and it will be supported, however, SQLite will yell very loudly at you if you do the same thing!

#### Supported Options:

* `size` - The size of the column. For example if you wanted a `varchar(50)` in Postgres you would do: `t.Column("column_name", "string", {"size": 50})`
* `null` - By default columns are not allowed to be `null`.
* `default` - The default value you want for this column. By default this is `null`.
* `default_raw` - The default value defined as a database function.

## Drop a Table

``` javascript
drop_table("table_name")
```

## Rename a Table

``` javascript
rename_table("old_table_name", "new_table_name")
```

## Add a Column

``` javascript
add_column("table_name", "column_name", "string", {})
```

See [above](#column-info) for more details on column types and options.

## Alter a column

``` javascript
change_column("table_name", "column_name", "string", {})
```

## Rename a Column

``` javascript
rename_column("table_name", "old_column_name", "new_column_name")
```

## Drop a Column

``` javascript
drop_column("table_name", "column_name")
```

## Add an Index

#### Supported Options:

* `name` - This defaults to `table_name_column_name_idx`
* `unique`

### Simple Index:

``` javascript
add_index("table_name", "column_name", {})
```

### Multi-Column Index:

``` javascript
add_index("table_name", ["column_1", "column_2"], {})
```

### Unique Index:

``` javascript
add_index("table_name", "column_name", {"unique": true})
```

### Index Names:

``` javascript
add_index("table_name", "column_name", {}) # name => table_name_column_name_idx
add_index("table_name", "column_name", {"name": "custom_index_name"})
```

## Rename an Index

``` javascript
rename_index("table_name", "old_index_name", "new_index_name")
```

## Drop an Index

``` javascript
drop_index("table_name", "index_name")
```

## Raw SQL

``` javascript
raw("select * from users;")
```

*All calls to `raw` must end with a `;`!*

## Execute an External Command

Sometimes during a migration you need to shell out to an external command.

```javascript
exec("echo hello")
```
