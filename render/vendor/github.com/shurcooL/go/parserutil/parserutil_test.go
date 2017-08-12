package parserutil_test

import (
	"fmt"
	"go/ast"
	"os"
	"reflect"
	"testing"

	"github.com/shurcooL/go/parserutil"
)

func Example() {
	stmt, err := parserutil.ParseStmt("var x int")
	if err != nil {
		panic(err)
	}

	ast.Fprint(os.Stdout, nil, stmt, nil)

	// Output:
	//      0  *ast.DeclStmt {
	//      1  .  Decl: *ast.GenDecl {
	//      2  .  .  Doc: nil
	//      3  .  .  TokPos: 31
	//      4  .  .  Tok: var
	//      5  .  .  Lparen: 0
	//      6  .  .  Specs: []ast.Spec (len = 1) {
	//      7  .  .  .  0: *ast.ValueSpec {
	//      8  .  .  .  .  Doc: nil
	//      9  .  .  .  .  Names: []*ast.Ident (len = 1) {
	//     10  .  .  .  .  .  0: *ast.Ident {
	//     11  .  .  .  .  .  .  NamePos: 35
	//     12  .  .  .  .  .  .  Name: "x"
	//     13  .  .  .  .  .  .  Obj: *ast.Object {
	//     14  .  .  .  .  .  .  .  Kind: var
	//     15  .  .  .  .  .  .  .  Name: "x"
	//     16  .  .  .  .  .  .  .  Decl: *(obj @ 7)
	//     17  .  .  .  .  .  .  .  Data: 0
	//     18  .  .  .  .  .  .  .  Type: nil
	//     19  .  .  .  .  .  .  }
	//     20  .  .  .  .  .  }
	//     21  .  .  .  .  }
	//     22  .  .  .  .  Type: *ast.Ident {
	//     23  .  .  .  .  .  NamePos: 37
	//     24  .  .  .  .  .  Name: "int"
	//     25  .  .  .  .  .  Obj: nil
	//     26  .  .  .  .  }
	//     27  .  .  .  .  Values: nil
	//     28  .  .  .  .  Comment: nil
	//     29  .  .  .  }
	//     30  .  .  }
	//     31  .  .  Rparen: 0
	//     32  .  }
	//     33  }
}

func TestParseStmt(t *testing.T) {
	tests := []struct {
		in        string
		want      ast.Stmt
		wantError error
	}{
		{"", &ast.EmptyStmt{Semicolon: 32}, nil},
	}
	for _, tc := range tests {
		stmt, err := parserutil.ParseStmt(tc.in)
		if got, want := err, tc.wantError; !equalError(got, want) {
			t.Errorf("got error: %v, want: %v", got, want)
			continue
		}
		if tc.wantError != nil {
			continue
		}
		if got, want := stmt, tc.want; !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}

func TestParseDecl(t *testing.T) {
	tests := []struct {
		in        string
		want      ast.Decl
		wantError error
	}{
		{"", nil, fmt.Errorf("no declaration")},
	}
	for _, tc := range tests {
		decl, err := parserutil.ParseDecl(tc.in)
		if got, want := err, tc.wantError; !equalError(got, want) {
			t.Errorf("got error: %v, want: %v", got, want)
			continue
		}
		if tc.wantError != nil {
			continue
		}
		if got, want := decl, tc.want; !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}

// equalError reports whether errors a and b are considered equal.
// They're equal if both are nil, or both are not nil and a.Error() == b.Error().
func equalError(a, b error) bool {
	return a == nil && b == nil || a != nil && b != nil && a.Error() == b.Error()
}
