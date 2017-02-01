package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TravisGeneration(t *testing.T) {
	dir := os.TempDir()
	packagePath := filepath.Join(dir, "src")
	os.MkdirAll(packagePath, 0755)
	os.Chdir(packagePath)

	ciProvider = "travis"
	installDeps = false
	skipPop = true

	genNewFiles("something", filepath.Join(dir, "/something"))

	r := require.New(t)
	content, err := ioutil.ReadFile(filepath.Join(dir, "something", ".travis.yml"))

	r.Equal(err, nil)
	r.Contains(string(content), "language: go")
	r.Contains(string(content), "CREATE DATABASE something_test")
	r.Contains(string(content), "psql -c 'create database something_test;' -U postgres")
	r.Contains(string(content), "go_import_path:")

	ciProvider = "none"
	installDeps = true
}

func Test_NoCIgeneration(t *testing.T) {
	dir := os.TempDir()
	os.RemoveAll(filepath.Join(dir, "src"))
	packagePath := filepath.Join(dir, "src")

	os.MkdirAll(packagePath, 0755)
	os.Chdir(packagePath)

	os.RemoveAll(filepath.Join(dir, "src", "no-travis"))

	ciProvider = "none"
	installDeps = false
	skipPop = true

	genNewFiles("something", filepath.Join(dir, "/no-travis"))

	r := require.New(t)
	_, err := ioutil.ReadFile(filepath.Join(dir, "no-travis", ".travis.yml"))

	r.NotEqual(err, nil)

	ciProvider = "none"
	installDeps = true
}
