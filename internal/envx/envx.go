package envx

import "os"

func Get(name string, alt string) string {
	x := os.Getenv(name)
	if len(x) == 0 {
		return alt
	}
	return x
}
