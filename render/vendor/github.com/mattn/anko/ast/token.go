package ast

type Token struct {
	PosImpl // StmtImpl provide Pos() function.
	Tok     int
	Lit     string
}
