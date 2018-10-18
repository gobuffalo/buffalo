package meta

import (
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
	if a.WithSQLite {
		tags = append(tags, "sqlite")
	}

	m := map[string]string{}
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if len(t) > 0 {
			m[t] = t
		}
	}
	var tt []string
	for k := range m {
		tt = append(tt, k)
	}

	return BuildTags(tt)
}
