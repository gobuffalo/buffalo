package goon_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"runtime"

	"github.com/shurcooL/go-goon"
)

func foo(bar int) int { return bar * 2 }

func Example_large() {
	fset := token.NewFileSet()
	if file, err := parser.ParseFile(fset, thisGoSourceFile(), nil, 0); nil == err {
		for _, d := range file.Decls {
			if f, ok := d.(*ast.FuncDecl); ok {
				goon.Dump(f)
				break
			}
		}
	}

	// Output:
	// (*ast.FuncDecl)(&ast.FuncDecl{
	// 	Doc:  (*ast.CommentGroup)(nil),
	// 	Recv: (*ast.FieldList)(nil),
	// 	Name: (*ast.Ident)(&ast.Ident{
	// 		NamePos: (token.Pos)(115),
	// 		Name:    (string)("foo"),
	// 		Obj: (*ast.Object)(&ast.Object{
	// 			Kind: (ast.ObjKind)(5),
	// 			Name: (string)("foo"),
	// 			Decl: (*ast.FuncDecl)(already_shown),
	// 			Data: (interface{})(nil),
	// 			Type: (interface{})(nil),
	// 		}),
	// 	}),
	// 	Type: (*ast.FuncType)(&ast.FuncType{
	// 		Func: (token.Pos)(110),
	// 		Params: (*ast.FieldList)(&ast.FieldList{
	// 			Opening: (token.Pos)(118),
	// 			List: ([]*ast.Field)([]*ast.Field{
	// 				(*ast.Field)(&ast.Field{
	// 					Doc: (*ast.CommentGroup)(nil),
	// 					Names: ([]*ast.Ident)([]*ast.Ident{
	// 						(*ast.Ident)(&ast.Ident{
	// 							NamePos: (token.Pos)(119),
	// 							Name:    (string)("bar"),
	// 							Obj: (*ast.Object)(&ast.Object{
	// 								Kind: (ast.ObjKind)(4),
	// 								Name: (string)("bar"),
	// 								Decl: (*ast.Field)(already_shown),
	// 								Data: (interface{})(nil),
	// 								Type: (interface{})(nil),
	// 							}),
	// 						}),
	// 					}),
	// 					Type: (*ast.Ident)(&ast.Ident{
	// 						NamePos: (token.Pos)(123),
	// 						Name:    (string)("int"),
	// 						Obj:     (*ast.Object)(nil),
	// 					}),
	// 					Tag:     (*ast.BasicLit)(nil),
	// 					Comment: (*ast.CommentGroup)(nil),
	// 				}),
	// 			}),
	// 			Closing: (token.Pos)(126),
	// 		}),
	// 		Results: (*ast.FieldList)(&ast.FieldList{
	// 			Opening: (token.Pos)(0),
	// 			List: ([]*ast.Field)([]*ast.Field{
	// 				(*ast.Field)(&ast.Field{
	// 					Doc:   (*ast.CommentGroup)(nil),
	// 					Names: ([]*ast.Ident)(nil),
	// 					Type: (*ast.Ident)(&ast.Ident{
	// 						NamePos: (token.Pos)(128),
	// 						Name:    (string)("int"),
	// 						Obj:     (*ast.Object)(nil),
	// 					}),
	// 					Tag:     (*ast.BasicLit)(nil),
	// 					Comment: (*ast.CommentGroup)(nil),
	// 				}),
	// 			}),
	// 			Closing: (token.Pos)(0),
	// 		}),
	// 	}),
	// 	Body: (*ast.BlockStmt)(&ast.BlockStmt{
	// 		Lbrace: (token.Pos)(132),
	// 		List: ([]ast.Stmt)([]ast.Stmt{
	// 			(*ast.ReturnStmt)(&ast.ReturnStmt{
	// 				Return: (token.Pos)(134),
	// 				Results: ([]ast.Expr)([]ast.Expr{
	// 					(*ast.BinaryExpr)(&ast.BinaryExpr{
	// 						X: (*ast.Ident)(&ast.Ident{
	// 							NamePos: (token.Pos)(141),
	// 							Name:    (string)("bar"),
	// 							Obj: (*ast.Object)(&ast.Object{
	// 								Kind: (ast.ObjKind)(4),
	// 								Name: (string)("bar"),
	// 								Decl: (*ast.Field)(&ast.Field{
	// 									Doc: (*ast.CommentGroup)(nil),
	// 									Names: ([]*ast.Ident)([]*ast.Ident{
	// 										(*ast.Ident)(&ast.Ident{
	// 											NamePos: (token.Pos)(119),
	// 											Name:    (string)("bar"),
	// 											Obj:     (*ast.Object)(already_shown),
	// 										}),
	// 									}),
	// 									Type: (*ast.Ident)(&ast.Ident{
	// 										NamePos: (token.Pos)(123),
	// 										Name:    (string)("int"),
	// 										Obj:     (*ast.Object)(nil),
	// 									}),
	// 									Tag:     (*ast.BasicLit)(nil),
	// 									Comment: (*ast.CommentGroup)(nil),
	// 								}),
	// 								Data: (interface{})(nil),
	// 								Type: (interface{})(nil),
	// 							}),
	// 						}),
	// 						OpPos: (token.Pos)(145),
	// 						Op:    (token.Token)(14),
	// 						Y: (*ast.BasicLit)(&ast.BasicLit{
	// 							ValuePos: (token.Pos)(147),
	// 							Kind:     (token.Token)(5),
	// 							Value:    (string)("2"),
	// 						}),
	// 					}),
	// 				}),
	// 			}),
	// 		}),
	// 		Rbrace: (token.Pos)(149),
	// 	}),
	// })
	//
}

// thisGoSourceFile returns the full path of the Go source file where this function was called from.
func thisGoSourceFile() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}
