package pop

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/markbates/going/defaults"
	"github.com/pkg/errors"
)

type ConnectionDetails struct {
	// Example: "postgres" or "sqlite3" or "mysql"
	Dialect string
	// The name of your database. Example: "foo_development"
	Database string
	// The host of your database. Example: "127.0.0.1"
	Host string
	// The port of your database. Example: 1234
	// Will default to the "default" port for each dialect.
	Port string
	// The username of the database user. Example: "root"
	User string
	// The password of the database user. Example: "password"
	Password string
	// Instead of specifying each individual piece of the
	// connection you can instead just specify the URL of the
	// database. Example: "postgres://postgres:postgres@localhost:5432/pop_test?sslmode=disable"
	URL string
	// Defaults to 0 "unlimited". See https://golang.org/pkg/database/sql/#DB.SetMaxOpenConns
	Pool    int
	Options map[string]string
}

var dialectX = regexp.MustCompile(`\s+:\/\/`)

// Finalize cleans up the connection details by normalizing names,
// filling in default values, etc...
func (cd *ConnectionDetails) Finalize() error {
	if cd.URL != "" {
		ul := cd.URL
		if cd.Dialect != "" {
			if !dialectX.MatchString(ul) {
				ul = cd.Dialect + "://" + ul
			}
		}
		u, err := url.Parse(ul)
		if err != nil {
			return errors.Wrapf(err, "couldn't parse %s", ul)
		}
		cd.Dialect = u.Scheme
		cd.Database = u.Path

		hp := strings.Split(u.Host, ":")
		cd.Host = hp[0]
		if len(hp) > 1 {
			cd.Port = hp[1]
		}

		if u.User != nil {
			cd.User = u.User.Username()
			cd.Password, _ = u.User.Password()
		}
	}
	switch strings.ToLower(cd.Dialect) {
	case "postgres", "postgresql", "pg":
		cd.Dialect = "postgres"
		cd.Port = defaults.String(cd.Port, "5432")
		cd.Database = strings.TrimPrefix(cd.Database, "/")
	case "mysql":
		cd.Port = defaults.String(cd.Port, "3006")
		cd.Database = strings.TrimPrefix(cd.Database, "/")
	case "sqlite", "sqlite3":
		cd.Dialect = "sqlite3"
	default:
		return errors.Errorf("Unknown dialect %s!", cd.Dialect)
	}
	return nil
}

// Parse is deprecated! Please use `ConnectionDetails.Finalize()` instead!
func (cd *ConnectionDetails) Parse(port string) error {
	fmt.Println("[POP] ConnectionDetails#Parse(port string) has been deprecated!")
	return cd.Finalize()
}

func (cd *ConnectionDetails) RetrySleep() time.Duration {
	d, err := time.ParseDuration(defaults.String(cd.Options["retry_sleep"], "1ms"))
	if err != nil {
		return 1 * time.Millisecond
	}
	return d
}

func (cd *ConnectionDetails) RetryLimit() int {
	i, err := strconv.Atoi(defaults.String(cd.Options["retry_limit"], "1000"))
	if err != nil {
		return 100
	}
	return i
}
