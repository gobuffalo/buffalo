package columns

import "fmt"

type Column struct {
	Name      string
	Writeable bool
	Readable  bool
	SelectSQL string
}

func (c Column) UpdateString() string {
	return fmt.Sprintf("%s = :%s", c.Name, c.Name)
}

func (c *Column) SetSelectSQL(s string) {
	c.SelectSQL = s
	c.Writeable = false
	c.Readable = true
}
