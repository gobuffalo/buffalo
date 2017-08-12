package pop

import (
	"fmt"

	"github.com/pkg/errors"
)

func CreateDB(c *Connection) error {
	deets := c.Dialect.Details()
	if deets.Database != "" {
		Log(fmt.Sprintf("Create %s (%s)", deets.Database, c.URL()))
		return errors.Wrapf(c.Dialect.CreateDB(), "couldn't create database %s", deets.Database)
	}
	return nil
}

func DropDB(c *Connection) error {
	deets := c.Dialect.Details()
	if deets.Database != "" {
		Log(fmt.Sprintf("Drop %s (%s)", deets.Database, c.URL()))
		return errors.Wrapf(c.Dialect.DropDB(), "couldn't drop database %s", deets.Database)
	}
	return nil
}
