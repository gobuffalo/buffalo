package validators

import (
	"fmt"
	"strings"

	"github.com/serenize/snaker"
)

func init() {
	fmt.Println("This package has been deprecated. Please use github.com/markbates/validate instead.")
}

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
