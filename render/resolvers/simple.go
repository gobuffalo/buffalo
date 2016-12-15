package resolvers

import "io/ioutil"

// SimpleResolver is dumb and will just look for the file
// exactly where you ask for it.
type SimpleResolver struct{}

// Read is dumb and will just look for the file
// exactly where you ask for it.
func (s *SimpleResolver) Read(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

// Resolve is dumb and will just look for the file
// exactly where you ask for it.
func (s *SimpleResolver) Resolve(name string) (string, error) {
	return name, nil
}
