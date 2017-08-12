package validators

import (
	"strings"

	"github.com/serenize/snaker"
)

var CustomKeys = map[string]string{}

func GenerateKey(s string) string {
	key := CustomKeys[s]
	if key != "" {
		return key
	}
	key = strings.Replace(s, " ", "", -1)
	key = strings.Replace(key, "-", "", -1)
	key = snaker.CamelToSnake(key)
	return key
}
