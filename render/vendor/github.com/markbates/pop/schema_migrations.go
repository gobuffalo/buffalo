// +build !appengine

package pop

import "github.com/markbates/pop/fizz"

var schemaMigrations = fizz.Table{
	Name: "schema_migration",
	Columns: []fizz.Column{
		{Name: "version", ColType: "string"},
	},
	Indexes: []fizz.Index{
		{Name: "version_idx", Columns: []string{"version"}, Unique: true},
	},
}
