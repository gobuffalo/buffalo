package fizz

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Options map[string]interface{}

type fizzer struct {
	Bubbler *Bubbler
}

func (f fizzer) add(s string, err error) error {
	if err != nil {
		panic(err.Error())
	}
	f.Bubbler.data = append(f.Bubbler.data, s)
	return nil
}

func (f fizzer) Exec(out io.Writer) interface{} {
	return func(s string) error {
		args := strings.Split(s, " ")
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}

func AFile(f *os.File, t Translator) (string, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return AString(string(b), t)
}

func AString(s string, t Translator) (string, error) {
	b := NewBubbler(t)
	return b.Bubble(s)
}
