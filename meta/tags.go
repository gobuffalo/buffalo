package meta

import (
	"bytes"
	"io/ioutil"
	"strings"
)

// BuildTags are tags used for building apps
type BuildTags []string

// String returns the tags in the form of:
// "foo bar baz" (with the quotes!)
func (t BuildTags) String() string {
	return strings.Join(t, " ")
}

// BuildTags combines the passed in env, and any additional tags,
// with tags that Buffalo decides the build process requires.
// An example would be adding the "sqlite" build tag if using
// SQLite3.
func (a App) BuildTags(env string, tags ...string) BuildTags {
	tags = append(tags, env)
	if b, err := ioutil.ReadFile("database.yml"); err == nil {
		if bytes.Contains(b, []byte("sqlite")) {
			tags = append(tags, "sqlite")
		}
	}
	return BuildTags(tags)
}
