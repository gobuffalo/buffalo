package resolvers

// FileResolver interface to try and resolve the location of
// a given file.
type FileResolver interface {
	// Read the location of the given file and return it's
	// contents or an error.
	Read(string) ([]byte, error)
	// Resolve the location of the given file
	Resolve(string) (string, error)
}
