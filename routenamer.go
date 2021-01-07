package buffalo

import (
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/flect/name"
)

// RouteNamer is in charge of naming a route from the
// path assigned, this name typically will be used if no
// name is assined with .Name(...).
type RouteNamer interface {
	// NameRoute receives the path and returns the name
	// for the route.
	NameRoute(string) string
}

// BaseRouteNamer is the default route namer used by apps.
type baseRouteNamer struct{}

func (drn baseRouteNamer) NameRoute(p string) string {
	if p == "/" || p == "" {
		return "root"
	}

	resultParts := []string{}
	parts := strings.Split(p, "/")

	for index, part := range parts {

		originalPart := parts[index]

		var previousPart string
		if index > 0 {
			previousPart = parts[index-1]
		}

		var nextPart string
		if len(parts) > index+1 {
			nextPart = parts[index+1]
		}

		isIdentifierPart := strings.Contains(part, "{") && (strings.Contains(part, flect.Singularize(previousPart)))
		isSimplifiedID := part == `{id}`

		if isIdentifierPart || isSimplifiedID || part == "" {
			continue
		}

		if strings.Contains(nextPart, "{") {
			part = flect.Singularize(part)
		}

		if originalPart == "new" || originalPart == "edit" {
			resultParts = append([]string{part}, resultParts...)
			continue
		}

		if strings.Contains(previousPart, "}") {
			resultParts = append(resultParts, part)
			continue
		}

		resultParts = append(resultParts, part)
	}

	if len(resultParts) == 0 {
		return "unnamed"
	}

	underscore := strings.TrimSpace(strings.Join(resultParts, "_"))
	return name.VarCase(underscore)
}
