package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// YarnScriptExists determines whether the specified yarn script exists in the package.json file
func YarnScriptExists(script string) (bool, error) {
	f, err := os.Open("./package.json")
	if err != nil {
		return false, errors.WithStack(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return hasYarnScript(script, b)
}

func hasYarnScript(script string, bytes []byte) (bool, error) {
	var pkg map[string]*json.RawMessage
	err := json.Unmarshal(bytes, &pkg)
	if err != nil {
		return false, errors.WithStack(err)
	}

	if scripts, found := pkg["scripts"]; found {
		var data map[string]*json.RawMessage

		// No error here as it would have failed earlier otherwise
		json.Unmarshal(*scripts, &data)

		_, found := data[script]
		return found, nil
	}

	return false, nil
}
