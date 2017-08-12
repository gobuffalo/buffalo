package annotate

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"text/template"
	"time"
	"unicode/utf8"
)

var saveExp = flag.Bool("exp", false, "overwrite all expected output files with actual output (returning a failure)")
var match = flag.String("m", "", "only run tests whose name contains this string")

func TestAnnotate(t *testing.T) {
	tests := map[string]struct {
		input   string
		anns    Annotations
		want    string
		wantErr error
	}{
		"empty and unannotated": {"", nil, "", nil},
		"unannotated":           {"a⌘b", nil, "a⌘b", nil},

		// The docs say "Annotating an empty byte array always returns an empty
		// byte array.", which is arbitrary but makes implementation easier.
		"empty annotated": {"", Annotations{{0, 0, []byte("["), []byte("]"), 0}}, "", nil},

		"zero-length annotations": {
			"aaaa",
			Annotations{
				{0, 0, []byte("<b>"), []byte("</b>"), 0},
				{0, 0, []byte("<i>"), []byte("</i>"), 0},
				{2, 2, []byte("<i>"), []byte("</i>"), 0},
			},
			"<b></b><i></i>aa<i></i>aa",
			nil,
		},
		"1 annotation": {"a", Annotations{{0, 1, []byte("["), []byte("]"), 0}}, "[a]", nil},
		"nested": {
			"abc",
			Annotations{
				{0, 3, []byte("["), []byte("]"), 0},
				{1, 2, []byte("<"), []byte(">"), 0},
			},
			"[a<b>c]",
			nil,
		},
		"nested 1": {
			"abcd",
			Annotations{
				{0, 4, []byte("<1>"), []byte("</1>"), 0},
				{1, 3, []byte("<2>"), []byte("</2>"), 0},
				{2, 2, []byte("<3>"), []byte("</3>"), 0},
			},
			"<1>a<2>b<3></3>c</2>d</1>",
			nil,
		},
		"same range": {
			"ab",
			Annotations{
				{0, 2, []byte("["), []byte("]"), 0},
				{0, 2, []byte("<"), []byte(">"), 0},
			},
			"[<ab>]",
			nil,
		},
		"same range (with WantInner)": {
			"ab",
			Annotations{
				{0, 2, []byte("["), []byte("]"), 1},
				{0, 2, []byte("<"), []byte(">"), 0},
			},
			"<[ab]>",
			nil,
		},
		"unicode content": {
			"abcdef⌘vwxyz",
			Annotations{
				{6, 9, []byte("<a>"), []byte("</a>"), 0},
				{10, 12, []byte("<b>"), []byte("</b>"), 0},
				{0, 13, []byte("<c>"), []byte("</c>"), 0},
			},
			"<c>abcdef<a>⌘</a>v<b>wx</b>y</c>z",
			nil,
		},
		"remainder": {
			"xyz",
			Annotations{
				{0, 2, []byte("<b>"), []byte("</b>"), 0},
				{0, 1, []byte("<c>"), []byte("</c>"), 0},
			},
			"<b><c>x</c>y</b>z",
			nil,
		},

		// Overlapping
		"overlap simple": {
			"abc",
			Annotations{
				{0, 2, []byte("<X>"), []byte("</X>"), 0},
				{1, 3, []byte("<Y>"), []byte("</Y>"), 0},
			},
			// Without re-opening overlapped annotations, we'd get
			// "<X>a<Y>b</X>c</Y>".
			"<X>a<Y>b</Y></X><Y>c</Y>",
			nil,
		},
		"overlap simple double": {
			"abc",
			Annotations{
				{0, 2, []byte("<X1>"), []byte("</X1>"), 0},
				{0, 2, []byte("<X2>"), []byte("</X2>"), 0},
				{1, 3, []byte("<Y1>"), []byte("</Y1>"), 0},
				{1, 3, []byte("<Y2>"), []byte("</Y2>"), 0},
			},
			"<X1><X2>a<Y1><Y2>b</Y2></Y1></X2></X1><Y1><Y2>c</Y2></Y1>",
			nil,
		},
		"overlap triple complex": {
			"abcd",
			Annotations{
				{0, 2, []byte("<X>"), []byte("</X>"), 0},
				{1, 3, []byte("<Y>"), []byte("</Y>"), 0},
				{2, 4, []byte("<Z>"), []byte("</Z>"), 0},
			},
			"<X>a<Y>b</Y></X><Y><Z>c</Z></Y><Z>d</Z>",
			nil,
		},
		"overlap same start": {
			"abcd",
			Annotations{
				{0, 2, []byte("<X>"), []byte("</X>"), 0},
				{0, 3, []byte("<Y>"), []byte("</Y>"), 0},
				{1, 4, []byte("<Z>"), []byte("</Z>"), 0},
			},
			"<Y><X>a<Z>b</Z></X><Z>c</Z></Y><Z>d</Z>",
			nil,
		},
		"overlap (infinite loop regression #1)": {
			"abcde",
			Annotations{
				{0, 3, []byte("<X>"), []byte("</X>"), 0},
				{1, 5, []byte("<Y>"), []byte("</Y>"), 0},
				{1, 2, []byte("<Z>"), []byte("</Z>"), 0},
			},
			"<X>a<Y><Z>b</Z>c</Y></X><Y>de</Y>",
			nil,
		},

		// Errors
		"start oob": {"a", Annotations{{-1, 1, []byte("<"), []byte(">"), 0}}, "<a>", ErrStartOutOfBounds},
		"start oob (multiple)": {
			"a",
			Annotations{
				{-3, 1, []byte("1"), []byte(""), 0},
				{-3, 1, []byte("2"), []byte(""), 0},
				{-1, 1, []byte("3"), []byte(""), 0},
			},
			"123a",
			ErrStartOutOfBounds,
		},
		"end oob": {"a", Annotations{{0, 3, []byte("<"), []byte(">"), 0}}, "<a>", ErrEndOutOfBounds},
		"end oob (multiple)": {
			"ab",
			Annotations{
				{0, 3, []byte(""), []byte("1"), 0},
				{1, 3, []byte(""), []byte("2"), 0},
				{0, 5, []byte(""), []byte("3"), 0},
			},
			"ab213",
			ErrEndOutOfBounds,
		},
	}
	for label, test := range tests {
		if *match != "" && !strings.Contains(label, *match) {
			continue
		}

		sort.Sort(Annotations(test.anns))

		got, err := Annotate([]byte(test.input), test.anns, nil)
		if err != test.wantErr {
			if test.wantErr == nil {
				t.Errorf("%s: Annotate: %s", label, err)
			} else {
				t.Errorf("%s: Annotate: got error %v, want %v", label, err, test.wantErr)
			}
		}
		if string(got) != test.want {
			t.Errorf("%s: Annotate:\ngot  %q\nwant %q", label, got, test.want)
			continue
		}
	}
}

func TestAnnotate_Files(t *testing.T) {
	annsByFile := map[string]Annotations{
		"hello_world.txt": {
			{0, 5, []byte("<b>"), []byte("</b>"), 0},
			{7, 12, []byte("<i>"), []byte("</i>"), 0},
		},
		"adjacent.txt": {
			{0, 3, []byte("<b>"), []byte("</b>"), 0},
			{3, 6, []byte("<i>"), []byte("</i>"), 0},
		},
		"nested_0.txt": {
			{0, 4, []byte("<1>"), []byte("</1>"), 0},
			{1, 3, []byte("<2>"), []byte("</2>"), 0},
		},
		"nested_2.txt": {
			{0, 2, []byte("<1>"), []byte("</1>"), 0},
			{2, 4, []byte("<2>"), []byte("</2>"), 0},
			{4, 6, []byte("<3>"), []byte("</3>"), 0},
			{7, 8, []byte("<4>"), []byte("</4>"), 0},
		},
		"html.txt": {
			{193, 203, []byte("<1>"), []byte("</1>"), 0},
			{336, 339, []byte("<WOOF>"), []byte("</WOOF>"), 0},
		},
	}

	dir := "testdata"
	tests, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		name := test.Name()
		if !strings.Contains(name, *match) {
			continue
		}
		if strings.HasSuffix(name, ".html") {
			continue
		}
		path := filepath.Join(dir, name)
		input, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatal(err)
			continue
		}

		anns := annsByFile[name]
		sort.Sort(anns)

		got, err := Annotate(input, anns, template.HTMLEscape)
		if err != nil {
			t.Errorf("%s: Annotate: %s", name, err)
			continue
		}

		expPath := path + ".html"
		if *saveExp {
			err = ioutil.WriteFile(expPath, got, 0700)
			if err != nil {
				t.Fatal(err)
			}
			continue
		}

		want, err := ioutil.ReadFile(expPath)
		if err != nil {
			t.Fatal(err)
		}

		want = bytes.TrimSpace(want)
		got = bytes.TrimSpace(got)

		if !bytes.Equal(want, got) {
			t.Errorf("%s: want %q, got %q", name, want, got)
			continue
		}
	}

	if *saveExp {
		t.Fatal("overwrote all expected output files with actual output (run tests again without -exp)")
	}
}

func makeFakeData(size1, size2 int) ([]byte, Annotations) {
	input := []byte(strings.Repeat(strings.Repeat("a", size1)+"⌘", size2))
	inputLength := utf8.RuneCount(input)
	n := len(input)/2 - (size1+1)/2
	anns := make(Annotations, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			anns[i] = &Annotation{Start: 2 * i, End: 2*i + 1}
		} else {
			anns[i] = &Annotation{Start: 2*i - 50, End: 2*i + 50}
			if anns[i].Start < 0 {
				anns[i].Start = 0
				anns[i].End = i
			}
			if anns[i].End >= inputLength {
				anns[i].End = inputLength
			}
		}
		anns[i].Left = []byte("L")  //[]byte(strings.Repeat("L", i%20))
		anns[i].Right = []byte("R") //[]byte(strings.Repeat("R", i%20))
		anns[i].WantInner = i % 5
	}
	sort.Sort(anns)
	return input, anns
}

func TestAnnotate_GeneratedData(t *testing.T) {
	input, anns := makeFakeData(1, 15)

	fail := func(err error) {
		annStrs := make([]string, len(anns))
		for i, a := range anns {
			annStrs[i] = fmt.Sprintf("%v", a)
		}
		t.Fatalf("Annotate: %s\n\nInput was:\n%q\n\nAnnotations:\n%s", err, input, strings.Join(annStrs, "\n"))
	}

	tm := time.NewTimer(time.Millisecond * 500)
	done := make(chan error)

	go func() {
		_, err := Annotate(input, anns, nil)
		done <- err
	}()

	select {
	case <-tm.C:
		fail(errors.New("timed out (is there an infinite loop?)"))
	case err := <-done:
		if err != nil {
			fail(err)
		}
	}
}

func BenchmarkAnnotate(b *testing.B) {
	input, anns := makeFakeData(99, 20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Annotate(input, anns, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}
