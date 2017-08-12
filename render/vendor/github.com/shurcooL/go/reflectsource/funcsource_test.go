package reflectsource

import (
	"fmt"
	"reflect"
)

func ExampleGetSourceAsString() {
	var f func()
	f1 := func() {
		panic(123)
	}
	f2 := func() {
		println("Hello from anon func!") // Comments are currently not preserved.
	}
	if 5*5 > 30 {
		f = f1
	} else {
		f = f2
	}

	fmt.Println(GetSourceAsString(f))

	// Output:
	//func() {
	//	println("Hello from anon func!")
	//}
}

func Example_two() {
	f := func(a int, b int) int {
		c := a + b
		return c
	}

	fmt.Println(GetSourceAsString(f))

	// Output:
	//func(a int, b int) int {
	//	c := a + b
	//	return c
	//}
}

func Example_nil() {
	var f func()

	fmt.Println(GetSourceAsString(f))

	// Output:
	//nil
}

func ExampleGetFuncValueSourceAsString() {
	f := func(a int, b int) int {
		c := a + b
		return c
	}

	fv := reflect.ValueOf(f)

	fmt.Println(GetFuncValueSourceAsString(fv))

	// Output:
	//func(a int, b int) int {
	//	c := a + b
	//	return c
	//}
}
