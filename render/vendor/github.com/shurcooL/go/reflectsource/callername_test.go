package reflectsource

import (
	"fmt"
	"strings"
)

func ExampleGetExprAsString() {
	var thisIsAFunkyVarName int

	fmt.Println("Name of var:", GetExprAsString(thisIsAFunkyVarName))

	// Output:
	// Name of var: thisIsAFunkyVarName
}

func Example() {
	var thisIsAFunkyVarName int
	var name string = GetExprAsString(thisIsAFunkyVarName)
	fmt.Println("Name of var:", name)
	fmt.Println("Some func name:", GetExprAsString(strings.HasPrefix))
	fmt.Println("Name of second arg:", getMySecondArgExprAsString(5, thisIsAFunkyVarName))

	// Output:
	// Name of var: thisIsAFunkyVarName
	// Some func name: strings.HasPrefix
	// Name of second arg: thisIsAFunkyVarName
}

func Example_trickyCases() {
	var thisIsAFunkyVarName int
	fmt.Println("1 2 3 4:", getMySecondArgExprAsString(1, 2), getMySecondArgExprAsString(3, 4)) // TODO: This should be 2, 4, not 2, 2
	fmt.Println("Name of second arg:",                                                          // TODO: This should work
		getMySecondArgExprAsString(5, thisIsAFunkyVarName))

	// Output:
	// 1 2 3 4: 2 2
	// Name of second arg: <expr not found>
}
