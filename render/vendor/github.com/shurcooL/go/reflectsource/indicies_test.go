package reflectsource

import (
	"fmt"
)

func Example_getLineStartEndIndicies() {
	b := []byte(`this

this is a longer line
and
stuff
last`)

	for lineIndex := 0; ; lineIndex++ {
		s, e := getLineStartEndIndicies(b, lineIndex)
		fmt.Printf("%v: [%v, %v]\n", lineIndex, s, e)
		if s == -1 {
			break
		}
	}

	// Output:
	// 0: [0, 4]
	// 1: [5, 5]
	// 2: [6, 27]
	// 3: [28, 31]
	// 4: [32, 37]
	// 5: [38, 42]
	// 6: [-1, -1]
}
