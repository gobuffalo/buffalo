package vm

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mattn/anko/parser"
)

func testInterrupt() {
	env := NewEnv()

	var sleepFunc = func(spec string) {
		if d, err := time.ParseDuration(spec); err != nil {
			panic(err)
		} else {
			time.Sleep(d)
		}
	}

	env.Define("println", fmt.Println)
	env.Define("sleep", sleepFunc)

	script := `
sleep("2s")
# Should interrupt here.
# The next line will not be executed.
println("<this should not be printed>")
`
	stmts, err := parser.ParseSrc(script)
	if err != nil {
		log.Fatal()
	}

	// Interrupts after 1 second.
	go func() {
		time.Sleep(time.Second)
		Interrupt(env)
	}()

	_, err = Run(stmts, env)
	if err != nil {
		log.Fatal()
	}
}

func TestInterruptRaces(t *testing.T) {
	// Run example several times
	for i := 0; i < 100; i++ {
		go testInterrupt()
	}
}
