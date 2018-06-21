package meta

import (
	"bytes"
	"io/ioutil"
	"strings"
)

type Tags []string

func (t Tags) String() string {
	return `"` + strings.Join(t, " ") + `"`
}

func (a App) BuildTags(env string, tags ...string) Tags {
	tags = append(tags, env)
	if b, err := ioutil.ReadFile("database.yml"); err == nil {
		if bytes.Contains(b, []byte("sqlite")) {
			tags = append(tags, "sqlite")
		}
	}
	return Tags(tags)
}
