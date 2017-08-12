//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2
import (
	"github.com/mattn/anko/ast"
)

//line parser.go.y:26
type yySymType struct {
	yys          int
	compstmt     []ast.Stmt
	stmt_if      ast.Stmt
	stmt_default ast.Stmt
	stmt_case    ast.Stmt
	stmt_cases   []ast.Stmt
	stmts        []ast.Stmt
	stmt         ast.Stmt
	typ          ast.Type
	expr         ast.Expr
	exprs        []ast.Expr
	expr_many    []ast.Expr
	expr_lets    ast.Expr
	expr_pair    ast.Expr
	expr_pairs   []ast.Expr
	expr_idents  []string
	tok          ast.Token
	term         ast.Token
	terms        ast.Token
	opt_terms    ast.Token
}

const IDENT = 57346
const NUMBER = 57347
const STRING = 57348
const ARRAY = 57349
const VARARG = 57350
const FUNC = 57351
const RETURN = 57352
const VAR = 57353
const THROW = 57354
const IF = 57355
const ELSE = 57356
const FOR = 57357
const IN = 57358
const EQEQ = 57359
const NEQ = 57360
const GE = 57361
const LE = 57362
const OROR = 57363
const ANDAND = 57364
const NEW = 57365
const TRUE = 57366
const FALSE = 57367
const NIL = 57368
const MODULE = 57369
const TRY = 57370
const CATCH = 57371
const FINALLY = 57372
const PLUSEQ = 57373
const MINUSEQ = 57374
const MULEQ = 57375
const DIVEQ = 57376
const ANDEQ = 57377
const OREQ = 57378
const BREAK = 57379
const CONTINUE = 57380
const PLUSPLUS = 57381
const MINUSMINUS = 57382
const POW = 57383
const SHIFTLEFT = 57384
const SHIFTRIGHT = 57385
const SWITCH = 57386
const CASE = 57387
const DEFAULT = 57388
const GO = 57389
const CHAN = 57390
const MAKE = 57391
const OPCHAN = 57392
const ARRAYLIT = 57393
const UNARY = 57394

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENT",
	"NUMBER",
	"STRING",
	"ARRAY",
	"VARARG",
	"FUNC",
	"RETURN",
	"VAR",
	"THROW",
	"IF",
	"ELSE",
	"FOR",
	"IN",
	"EQEQ",
	"NEQ",
	"GE",
	"LE",
	"OROR",
	"ANDAND",
	"NEW",
	"TRUE",
	"FALSE",
	"NIL",
	"MODULE",
	"TRY",
	"CATCH",
	"FINALLY",
	"PLUSEQ",
	"MINUSEQ",
	"MULEQ",
	"DIVEQ",
	"ANDEQ",
	"OREQ",
	"BREAK",
	"CONTINUE",
	"PLUSPLUS",
	"MINUSMINUS",
	"POW",
	"SHIFTLEFT",
	"SHIFTRIGHT",
	"SWITCH",
	"CASE",
	"DEFAULT",
	"GO",
	"CHAN",
	"MAKE",
	"OPCHAN",
	"ARRAYLIT",
	"'='",
	"'?'",
	"':'",
	"','",
	"'>'",
	"'<'",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"UNARY",
	"'{'",
	"'}'",
	"';'",
	"'.'",
	"'!'",
	"'^'",
	"'&'",
	"'('",
	"')'",
	"'['",
	"']'",
	"'|'",
	"'\\n'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:710

//line yacctab:1
var yyExca = [...]int{
	-1, 0,
	1, 3,
	-2, 122,
	-1, 1,
	1, -1,
	-2, 0,
	-1, 2,
	55, 48,
	-2, 1,
	-1, 10,
	55, 49,
	-2, 24,
	-1, 43,
	55, 48,
	-2, 123,
	-1, 85,
	65, 3,
	-2, 122,
	-1, 88,
	55, 49,
	-2, 43,
	-1, 90,
	65, 3,
	-2, 122,
	-1, 97,
	1, 57,
	8, 57,
	45, 57,
	46, 57,
	52, 57,
	54, 57,
	55, 57,
	64, 57,
	65, 57,
	66, 57,
	72, 57,
	74, 57,
	76, 57,
	-2, 52,
	-1, 99,
	1, 59,
	8, 59,
	45, 59,
	46, 59,
	52, 59,
	54, 59,
	55, 59,
	64, 59,
	65, 59,
	66, 59,
	72, 59,
	74, 59,
	76, 59,
	-2, 52,
	-1, 127,
	17, 0,
	18, 0,
	-2, 85,
	-1, 128,
	17, 0,
	18, 0,
	-2, 86,
	-1, 147,
	55, 49,
	-2, 43,
	-1, 149,
	65, 3,
	-2, 122,
	-1, 151,
	65, 3,
	-2, 122,
	-1, 153,
	65, 1,
	-2, 36,
	-1, 156,
	65, 3,
	-2, 122,
	-1, 180,
	65, 3,
	-2, 122,
	-1, 223,
	55, 50,
	-2, 44,
	-1, 224,
	1, 45,
	45, 45,
	46, 45,
	52, 45,
	55, 51,
	65, 45,
	66, 45,
	76, 45,
	-2, 52,
	-1, 231,
	1, 51,
	8, 51,
	45, 51,
	46, 51,
	55, 51,
	65, 51,
	66, 51,
	72, 51,
	74, 51,
	76, 51,
	-2, 52,
	-1, 233,
	65, 3,
	-2, 122,
	-1, 235,
	65, 3,
	-2, 122,
	-1, 248,
	65, 3,
	-2, 122,
	-1, 259,
	1, 106,
	8, 106,
	45, 106,
	46, 106,
	52, 106,
	54, 106,
	55, 106,
	64, 106,
	65, 106,
	66, 106,
	72, 106,
	74, 106,
	76, 106,
	-2, 104,
	-1, 261,
	1, 110,
	8, 110,
	45, 110,
	46, 110,
	52, 110,
	54, 110,
	55, 110,
	64, 110,
	65, 110,
	66, 110,
	72, 110,
	74, 110,
	76, 110,
	-2, 108,
	-1, 271,
	65, 3,
	-2, 122,
	-1, 276,
	65, 3,
	-2, 122,
	-1, 277,
	65, 3,
	-2, 122,
	-1, 282,
	1, 105,
	8, 105,
	45, 105,
	46, 105,
	52, 105,
	54, 105,
	55, 105,
	64, 105,
	65, 105,
	66, 105,
	72, 105,
	74, 105,
	76, 105,
	-2, 103,
	-1, 283,
	1, 109,
	8, 109,
	45, 109,
	46, 109,
	52, 109,
	54, 109,
	55, 109,
	64, 109,
	65, 109,
	66, 109,
	72, 109,
	74, 109,
	76, 109,
	-2, 107,
	-1, 288,
	65, 3,
	-2, 122,
	-1, 289,
	65, 3,
	-2, 122,
	-1, 292,
	45, 3,
	46, 3,
	65, 3,
	-2, 122,
	-1, 296,
	65, 3,
	-2, 122,
	-1, 303,
	45, 3,
	46, 3,
	65, 3,
	-2, 122,
	-1, 316,
	65, 3,
	-2, 122,
	-1, 317,
	65, 3,
	-2, 122,
}

const yyNprod = 128
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 2279

var yyAct = [...]int{

	81, 169, 240, 10, 45, 241, 253, 40, 172, 213,
	283, 282, 1, 6, 278, 11, 82, 249, 211, 88,
	6, 91, 80, 7, 94, 95, 96, 98, 100, 6,
	7, 92, 154, 93, 89, 174, 105, 93, 108, 7,
	110, 260, 112, 258, 10, 201, 184, 246, 116, 117,
	228, 119, 120, 121, 122, 123, 124, 125, 126, 127,
	128, 129, 130, 131, 132, 133, 134, 135, 136, 137,
	138, 102, 216, 139, 140, 141, 142, 220, 144, 145,
	147, 166, 2, 109, 143, 263, 42, 92, 148, 93,
	148, 115, 148, 148, 161, 146, 216, 216, 152, 106,
	160, 217, 262, 158, 321, 261, 164, 259, 250, 202,
	185, 167, 115, 242, 243, 147, 103, 104, 176, 170,
	264, 64, 65, 66, 67, 68, 69, 155, 320, 207,
	181, 55, 216, 239, 313, 310, 309, 306, 101, 305,
	78, 302, 293, 287, 286, 265, 255, 270, 237, 234,
	232, 190, 198, 192, 10, 194, 195, 49, 147, 150,
	74, 76, 189, 77, 191, 72, 317, 316, 153, 196,
	296, 209, 289, 197, 277, 276, 248, 148, 281, 149,
	223, 90, 218, 219, 227, 111, 221, 222, 229, 230,
	273, 215, 114, 225, 179, 115, 157, 173, 182, 79,
	242, 243, 315, 244, 311, 247, 245, 271, 5, 238,
	151, 8, 84, 44, 251, 208, 256, 170, 257, 226,
	173, 210, 206, 205, 165, 118, 83, 46, 4, 168,
	87, 188, 43, 199, 17, 3, 0, 269, 0, 0,
	0, 177, 200, 272, 178, 0, 267, 0, 268, 0,
	212, 214, 44, 230, 0, 113, 280, 0, 0, 0,
	0, 275, 0, 0, 284, 285, 0, 0, 64, 65,
	66, 67, 68, 69, 0, 0, 0, 0, 55, 0,
	0, 0, 0, 0, 290, 0, 0, 78, 0, 294,
	295, 0, 0, 0, 252, 0, 254, 52, 53, 54,
	308, 300, 301, 0, 49, 304, 0, 74, 76, 307,
	77, 0, 72, 0, 0, 0, 312, 0, 0, 0,
	0, 0, 58, 59, 61, 63, 73, 75, 0, 318,
	319, 0, 0, 0, 0, 0, 64, 65, 66, 67,
	68, 69, 0, 0, 70, 71, 55, 56, 57, 0,
	0, 0, 0, 0, 0, 78, 292, 0, 48, 0,
	299, 60, 62, 50, 51, 52, 53, 54, 0, 0,
	0, 0, 49, 0, 303, 74, 76, 298, 77, 0,
	72, 58, 59, 61, 63, 73, 75, 0, 0, 0,
	0, 0, 0, 0, 0, 64, 65, 66, 67, 68,
	69, 0, 0, 70, 71, 55, 56, 57, 0, 0,
	0, 0, 0, 0, 78, 0, 0, 48, 204, 0,
	60, 62, 50, 51, 52, 53, 54, 0, 0, 0,
	0, 49, 0, 0, 74, 76, 0, 77, 203, 72,
	58, 59, 61, 63, 73, 75, 0, 0, 0, 0,
	0, 0, 0, 0, 64, 65, 66, 67, 68, 69,
	0, 0, 70, 71, 55, 56, 57, 0, 0, 0,
	0, 0, 0, 78, 0, 0, 48, 187, 0, 60,
	62, 50, 51, 52, 53, 54, 0, 0, 0, 0,
	49, 0, 0, 74, 76, 0, 77, 186, 72, 58,
	59, 61, 63, 73, 75, 0, 0, 0, 0, 0,
	0, 0, 0, 64, 65, 66, 67, 68, 69, 0,
	0, 70, 71, 55, 56, 57, 0, 0, 0, 0,
	0, 0, 78, 0, 0, 48, 0, 0, 60, 62,
	50, 51, 52, 53, 54, 0, 0, 0, 0, 49,
	0, 0, 74, 76, 314, 77, 0, 72, 58, 59,
	61, 63, 73, 75, 0, 0, 0, 0, 0, 0,
	0, 0, 64, 65, 66, 67, 68, 69, 0, 0,
	70, 71, 55, 56, 57, 0, 0, 0, 0, 0,
	0, 78, 0, 0, 48, 0, 0, 60, 62, 50,
	51, 52, 53, 54, 0, 0, 0, 0, 49, 0,
	0, 74, 76, 297, 77, 0, 72, 58, 59, 61,
	63, 73, 75, 0, 0, 0, 0, 0, 0, 0,
	0, 64, 65, 66, 67, 68, 69, 0, 0, 70,
	71, 55, 56, 57, 0, 0, 0, 0, 0, 0,
	78, 0, 0, 48, 291, 0, 60, 62, 50, 51,
	52, 53, 54, 0, 0, 0, 0, 49, 0, 0,
	74, 76, 0, 77, 0, 72, 58, 59, 61, 63,
	73, 75, 0, 0, 0, 0, 0, 0, 0, 0,
	64, 65, 66, 67, 68, 69, 0, 0, 70, 71,
	55, 56, 57, 0, 0, 0, 0, 0, 0, 78,
	0, 0, 48, 0, 0, 60, 62, 50, 51, 52,
	53, 54, 0, 288, 0, 0, 49, 0, 0, 74,
	76, 0, 77, 0, 72, 58, 59, 61, 63, 73,
	75, 0, 0, 0, 0, 0, 0, 0, 0, 64,
	65, 66, 67, 68, 69, 0, 0, 70, 71, 55,
	56, 57, 0, 0, 0, 0, 0, 0, 78, 0,
	0, 48, 0, 0, 60, 62, 50, 51, 52, 53,
	54, 0, 0, 0, 0, 49, 0, 0, 74, 76,
	0, 77, 274, 72, 58, 59, 61, 63, 73, 75,
	0, 0, 0, 0, 0, 0, 0, 0, 64, 65,
	66, 67, 68, 69, 0, 0, 70, 71, 55, 56,
	57, 0, 0, 0, 0, 0, 0, 78, 0, 0,
	48, 0, 0, 60, 62, 50, 51, 52, 53, 54,
	0, 0, 0, 0, 49, 0, 0, 74, 76, 0,
	77, 266, 72, 58, 59, 61, 63, 73, 75, 0,
	0, 0, 0, 0, 0, 0, 0, 64, 65, 66,
	67, 68, 69, 0, 0, 70, 71, 55, 56, 57,
	0, 0, 0, 0, 0, 0, 78, 0, 0, 48,
	0, 0, 60, 62, 50, 51, 52, 53, 54, 0,
	0, 0, 236, 49, 0, 0, 74, 76, 0, 77,
	0, 72, 58, 59, 61, 63, 73, 75, 0, 0,
	0, 0, 0, 0, 0, 0, 64, 65, 66, 67,
	68, 69, 0, 0, 70, 71, 55, 56, 57, 0,
	0, 0, 0, 0, 0, 78, 0, 0, 48, 0,
	0, 60, 62, 50, 51, 52, 53, 54, 0, 235,
	0, 0, 49, 0, 0, 74, 76, 0, 77, 0,
	72, 58, 59, 61, 63, 73, 75, 0, 0, 0,
	0, 0, 0, 0, 0, 64, 65, 66, 67, 68,
	69, 0, 0, 70, 71, 55, 56, 57, 0, 0,
	0, 0, 0, 0, 78, 0, 0, 48, 0, 0,
	60, 62, 50, 51, 52, 53, 54, 0, 233, 0,
	0, 49, 0, 0, 74, 76, 0, 77, 0, 72,
	58, 59, 61, 63, 73, 75, 0, 0, 0, 0,
	0, 0, 0, 0, 64, 65, 66, 67, 68, 69,
	0, 0, 70, 71, 55, 56, 57, 0, 0, 0,
	0, 0, 0, 78, 0, 0, 48, 183, 0, 60,
	62, 50, 51, 52, 53, 54, 0, 0, 0, 0,
	49, 0, 0, 74, 76, 0, 77, 0, 72, 58,
	59, 61, 63, 73, 75, 0, 0, 0, 0, 0,
	0, 0, 0, 64, 65, 66, 67, 68, 69, 0,
	0, 70, 71, 55, 56, 57, 0, 0, 0, 0,
	0, 0, 78, 0, 0, 48, 0, 0, 60, 62,
	50, 51, 52, 53, 54, 0, 180, 0, 0, 49,
	0, 0, 74, 76, 0, 77, 0, 72, 58, 59,
	61, 63, 73, 75, 0, 0, 0, 0, 0, 0,
	0, 0, 64, 65, 66, 67, 68, 69, 0, 0,
	70, 71, 55, 56, 57, 0, 0, 0, 0, 0,
	0, 78, 0, 0, 48, 0, 0, 60, 62, 50,
	51, 52, 53, 54, 0, 0, 0, 0, 49, 0,
	0, 74, 76, 171, 77, 0, 72, 58, 59, 61,
	63, 73, 75, 0, 0, 0, 0, 0, 0, 0,
	0, 64, 65, 66, 67, 68, 69, 0, 0, 70,
	71, 55, 56, 57, 0, 0, 0, 0, 0, 0,
	78, 0, 0, 48, 0, 0, 60, 62, 50, 51,
	52, 53, 54, 0, 159, 0, 0, 49, 0, 0,
	74, 76, 0, 77, 0, 72, 58, 59, 61, 63,
	73, 75, 0, 0, 0, 0, 0, 0, 0, 0,
	64, 65, 66, 67, 68, 69, 0, 0, 70, 71,
	55, 56, 57, 0, 0, 0, 0, 0, 0, 78,
	0, 0, 48, 0, 0, 60, 62, 50, 51, 52,
	53, 54, 0, 156, 0, 0, 49, 0, 0, 74,
	76, 0, 77, 0, 72, 21, 22, 28, 0, 0,
	32, 14, 9, 15, 41, 0, 18, 0, 0, 0,
	0, 0, 0, 0, 36, 29, 30, 31, 16, 19,
	0, 0, 0, 0, 0, 0, 0, 0, 12, 13,
	0, 0, 0, 0, 0, 20, 0, 0, 37, 0,
	38, 39, 0, 0, 0, 0, 0, 0, 0, 0,
	23, 27, 0, 0, 0, 34, 0, 6, 0, 24,
	25, 26, 35, 0, 33, 0, 0, 7, 58, 59,
	61, 63, 73, 75, 0, 0, 0, 0, 0, 0,
	0, 0, 64, 65, 66, 67, 68, 69, 0, 0,
	70, 71, 55, 56, 57, 0, 0, 0, 0, 0,
	0, 78, 0, 47, 48, 0, 0, 60, 62, 50,
	51, 52, 53, 54, 0, 0, 0, 0, 49, 0,
	0, 74, 76, 0, 77, 0, 72, 58, 59, 61,
	63, 73, 75, 0, 0, 0, 0, 0, 0, 0,
	0, 64, 65, 66, 67, 68, 69, 0, 0, 70,
	71, 55, 56, 57, 0, 0, 0, 0, 0, 0,
	78, 0, 0, 48, 0, 0, 60, 62, 50, 51,
	52, 53, 54, 0, 0, 0, 0, 49, 0, 0,
	74, 76, 0, 77, 0, 72, 58, 59, 61, 63,
	73, 75, 0, 0, 0, 0, 0, 0, 0, 0,
	64, 65, 66, 67, 68, 69, 0, 0, 70, 71,
	55, 56, 57, 0, 0, 0, 0, 0, 0, 78,
	0, 0, 48, 0, 0, 60, 62, 50, 51, 52,
	53, 54, 0, 0, 0, 0, 49, 0, 0, 74,
	175, 0, 77, 0, 72, 58, 59, 61, 63, 73,
	75, 0, 0, 0, 0, 0, 0, 0, 0, 64,
	65, 66, 67, 68, 69, 0, 0, 70, 71, 55,
	56, 57, 0, 0, 0, 0, 0, 0, 78, 0,
	0, 48, 0, 0, 60, 62, 50, 51, 52, 53,
	54, 0, 0, 0, 0, 163, 0, 0, 74, 76,
	0, 77, 0, 72, 58, 59, 61, 63, 73, 75,
	0, 0, 0, 0, 0, 0, 0, 0, 64, 65,
	66, 67, 68, 69, 0, 0, 70, 71, 55, 56,
	57, 0, 0, 0, 0, 0, 0, 78, 0, 0,
	48, 0, 0, 60, 62, 50, 51, 52, 53, 54,
	58, 59, 61, 63, 162, 75, 0, 74, 76, 0,
	77, 0, 72, 0, 64, 65, 66, 67, 68, 69,
	0, 0, 70, 71, 55, 56, 57, 0, 0, 0,
	0, 0, 0, 78, 0, 0, 0, 0, 0, 60,
	62, 50, 51, 52, 53, 54, 58, 59, 61, 63,
	49, 0, 0, 74, 76, 0, 77, 0, 72, 0,
	64, 65, 66, 67, 68, 69, 0, 0, 70, 71,
	55, 56, 57, 0, 0, 0, 0, 0, 0, 78,
	0, 0, 0, 0, 0, 60, 62, 50, 51, 52,
	53, 54, 0, 0, 0, 0, 49, 0, 0, 74,
	76, 0, 77, 0, 72, 21, 22, 193, 0, 0,
	32, 14, 9, 15, 41, 0, 18, 0, 0, 0,
	0, 0, 0, 0, 36, 29, 30, 31, 16, 19,
	0, 0, 0, 0, 0, 0, 0, 0, 12, 13,
	0, 0, 0, 0, 0, 20, 0, 0, 37, 0,
	38, 39, 0, 0, 0, 0, 0, 0, 0, 0,
	23, 27, 0, 0, 0, 34, 0, 0, 0, 24,
	25, 26, 35, 0, 33, 21, 22, 28, 0, 0,
	32, 14, 9, 15, 41, 0, 18, 0, 0, 0,
	0, 0, 0, 0, 36, 29, 30, 31, 16, 19,
	0, 0, 0, 0, 0, 0, 0, 0, 12, 13,
	0, 0, 0, 0, 0, 20, 0, 0, 37, 0,
	38, 39, 0, 0, 0, 0, 0, 0, 0, 0,
	23, 27, 0, 61, 63, 34, 0, 0, 0, 24,
	25, 26, 35, 0, 33, 64, 65, 66, 67, 68,
	69, 0, 0, 70, 71, 55, 56, 57, 0, 0,
	0, 0, 0, 0, 78, 0, 0, 0, 0, 0,
	60, 62, 50, 51, 52, 53, 54, 0, 0, 0,
	0, 49, 0, 0, 74, 76, 0, 77, 0, 72,
	64, 65, 66, 67, 68, 69, 0, 0, 70, 71,
	55, 0, 0, 231, 22, 28, 0, 0, 32, 78,
	0, 0, 0, 0, 0, 0, 0, 50, 51, 52,
	53, 54, 36, 29, 30, 31, 49, 0, 0, 74,
	76, 0, 77, 0, 72, 21, 22, 28, 0, 0,
	32, 0, 0, 0, 0, 0, 37, 0, 38, 39,
	0, 0, 0, 0, 36, 29, 30, 31, 23, 27,
	0, 0, 0, 34, 0, 0, 0, 24, 25, 26,
	35, 0, 33, 279, 0, 0, 0, 0, 37, 0,
	38, 39, 0, 0, 0, 0, 0, 231, 22, 28,
	23, 27, 32, 0, 0, 34, 0, 0, 0, 24,
	25, 26, 35, 0, 33, 0, 36, 29, 30, 31,
	0, 0, 0, 0, 0, 224, 22, 28, 0, 0,
	32, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	37, 0, 38, 39, 36, 29, 30, 31, 0, 0,
	0, 0, 23, 27, 107, 22, 28, 34, 0, 32,
	0, 24, 25, 26, 35, 0, 33, 0, 37, 0,
	38, 39, 0, 36, 29, 30, 31, 0, 0, 0,
	23, 27, 99, 22, 28, 34, 0, 32, 0, 24,
	25, 26, 35, 0, 33, 0, 0, 37, 0, 38,
	39, 36, 29, 30, 31, 0, 0, 0, 0, 23,
	27, 97, 22, 28, 34, 0, 32, 0, 24, 25,
	26, 35, 0, 33, 0, 37, 0, 38, 39, 0,
	36, 29, 30, 31, 0, 0, 0, 23, 27, 86,
	22, 28, 34, 0, 32, 0, 24, 25, 26, 35,
	0, 33, 0, 0, 37, 0, 38, 39, 36, 29,
	30, 31, 0, 0, 0, 0, 23, 27, 0, 0,
	0, 34, 0, 0, 0, 24, 25, 26, 35, 0,
	33, 0, 37, 0, 38, 39, 0, 0, 0, 0,
	0, 0, 0, 0, 23, 27, 0, 0, 0, 85,
	0, 0, 0, 24, 25, 26, 35, 0, 33,
}
var yyPact = [...]int{

	-53, -1000, 1851, -53, -53, -1000, -1000, -1000, -1000, 223,
	1381, 147, -1000, -1000, 2011, 2011, 222, 198, 2205, 117,
	2011, -40, -1000, 2011, 2011, 2011, 2177, 2148, -1000, -1000,
	-1000, -1000, 67, -53, -53, 2011, 28, 2120, 12, 2011,
	130, 2011, -1000, 1321, -1000, 140, -1000, 2011, 2011, 221,
	2011, 2011, 2011, 2011, 2011, 2011, 2011, 2011, 2011, 2011,
	2011, 2011, 2011, 2011, 2011, 2011, 2011, 2011, 2011, 2011,
	-1000, -1000, 2011, 2011, 2011, 2011, 2011, 2011, 2011, 2011,
	122, 1440, 1440, 115, 146, -53, 16, 61, 1249, 144,
	-53, 1190, 2011, 2011, 90, 90, 90, -40, 1617, -40,
	1558, 220, 10, 2011, 211, 1131, 216, -36, 1499, 193,
	1440, -53, 1072, -1000, 2011, -53, 1440, 1013, -1000, 237,
	237, 90, 90, 90, 1440, 1939, 1939, 1894, 1894, 1939,
	1939, 1939, 1939, 1440, 1440, 1440, 1440, 1440, 1440, 1440,
	1663, 1440, 1709, 38, 423, 1440, -1000, 1440, -53, -53,
	2011, -53, 88, 1781, 2011, 2011, -53, 2011, 87, -53,
	37, 364, 219, 218, 57, 207, 217, -37, -46, -1000,
	137, -1000, 29, -1000, 2011, 2011, 5, 216, 216, 2091,
	-53, -1000, 215, 2011, -22, -1000, -1000, 2011, 2063, 85,
	954, 84, -1000, 137, 895, 836, 83, -1000, 180, 68,
	155, -25, -1000, -1000, 2011, -1000, -1000, 112, -55, 36,
	206, -53, -68, -53, 81, 2011, 214, -1000, 35, 33,
	-1000, 30, 65, 1440, -40, 80, -1000, 1440, -1000, 777,
	1440, -40, -1000, -53, -1000, -53, 2011, -1000, 143, -1000,
	-1000, -1000, 2011, 136, -1000, -1000, -1000, 718, -53, 111,
	110, -58, 1979, -1000, 113, -1000, 1440, -1000, -61, -1000,
	-62, -1000, -1000, 2011, 2011, -1000, -1000, 79, 78, 659,
	108, -53, 600, -53, -1000, 77, -53, -53, 106, -1000,
	-1000, -1000, -1000, -1000, 541, 305, -1000, -1000, -53, -53,
	76, -53, -53, -1000, 74, 72, -53, -1000, -1000, 2011,
	71, 70, 174, -53, -1000, -1000, -1000, 69, 482, -1000,
	172, 103, -1000, -1000, -1000, 102, -53, -53, 63, 39,
	-1000, -1000,
}
var yyPgo = [...]int{

	0, 12, 235, 211, 234, 5, 2, 233, 8, 0,
	7, 15, 230, 1, 229, 4, 82, 228, 208,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 2, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 4, 4, 4, 7, 7,
	7, 7, 7, 6, 5, 13, 14, 14, 14, 15,
	15, 15, 12, 11, 11, 11, 8, 8, 10, 10,
	10, 10, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 9, 9, 9, 9, 9, 9, 9, 9,
	9, 9, 16, 16, 17, 17, 18, 18,
}
var yyR2 = [...]int{

	0, 1, 2, 0, 2, 3, 4, 3, 3, 1,
	1, 2, 2, 5, 1, 4, 7, 9, 5, 13,
	12, 9, 8, 5, 1, 7, 5, 5, 0, 2,
	2, 2, 2, 5, 4, 3, 0, 1, 4, 0,
	1, 4, 3, 1, 4, 4, 1, 3, 0, 1,
	4, 4, 1, 1, 2, 2, 2, 2, 4, 2,
	4, 1, 1, 1, 1, 5, 3, 7, 8, 8,
	9, 5, 6, 5, 6, 3, 4, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 2, 2, 3,
	3, 3, 3, 5, 4, 6, 5, 5, 4, 6,
	5, 4, 4, 6, 6, 4, 5, 7, 7, 9,
	3, 2, 0, 1, 1, 2, 1, 1,
}
var yyChk = [...]int{

	-1000, -1, -16, -2, -17, -18, 66, 76, -3, 11,
	-9, -11, 37, 38, 10, 12, 27, -4, 15, 28,
	44, 4, 5, 59, 68, 69, 70, 60, 6, 24,
	25, 26, 9, 73, 64, 71, 23, 47, 49, 50,
	-10, 13, -16, -17, -18, -15, 4, 52, 53, 67,
	58, 59, 60, 61, 62, 41, 42, 43, 17, 18,
	56, 19, 57, 20, 31, 32, 33, 34, 35, 36,
	39, 40, 75, 21, 70, 22, 71, 73, 50, 52,
	-10, -9, -9, 4, 14, 64, 4, -12, -9, -11,
	64, -9, 71, 73, -9, -9, -9, 4, -9, 4,
	-9, 71, 4, -16, -16, -9, 71, 4, -9, 71,
	-9, 55, -9, -3, 52, 55, -9, -9, 4, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -9, -9, -9, -9, -9, -9, -9,
	-9, -9, -9, -10, -9, -9, -11, -9, 55, 64,
	13, 64, -1, -16, 16, 66, 64, 52, -1, 64,
	-10, -9, 67, 67, -15, 4, 71, -10, -14, -13,
	6, 72, -8, 4, 71, 71, -8, 48, 51, -16,
	64, -11, -16, 54, 8, 72, 74, 54, -16, -1,
	-9, -1, 65, 6, -9, -9, -1, -11, 65, -7,
	-16, 8, 72, 74, 54, 4, 4, 72, 8, -15,
	4, 55, -16, 55, -16, 54, 67, 72, -10, -10,
	72, -8, -8, -9, 4, -1, 4, -9, 72, -9,
	-9, 4, 65, 64, 65, 64, 66, 65, 29, 65,
	-6, -5, 45, 46, -6, -5, 72, -9, 64, 72,
	72, 8, -16, 74, -16, 65, -9, 4, 8, 72,
	8, 72, 72, 55, 55, 65, 74, -1, -1, -9,
	4, 64, -9, 54, 74, -1, 64, 64, 72, 74,
	-13, 65, 72, 72, -9, -9, 65, 65, 64, 64,
	-1, 54, -16, 65, -1, -1, 64, 72, 72, 55,
	-1, -1, 65, -16, -1, 65, 65, -1, -9, 65,
	65, 30, -1, 65, 72, 30, 64, 64, -1, -1,
	65, 65,
}
var yyDef = [...]int{

	-2, -2, -2, 122, 123, 124, 126, 127, 4, 39,
	-2, 0, 9, 10, 48, 0, 0, 14, 48, 0,
	0, 52, 53, 0, 0, 0, 0, 0, 61, 62,
	63, 64, 0, 122, 122, 0, 0, 0, 0, 0,
	0, 0, 2, -2, 125, 0, 40, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	97, 98, 0, 0, 0, 0, 48, 0, 0, 48,
	11, 49, 12, 0, 0, -2, 52, 0, -2, 0,
	-2, 0, 48, 0, 54, 55, 56, -2, 0, -2,
	0, 39, 0, 48, 36, 0, 0, 52, 0, 0,
	121, 122, 0, 5, 48, 122, 7, 0, 66, 77,
	78, 79, 80, 81, 82, 83, 84, -2, -2, 87,
	88, 89, 90, 91, 92, 93, 94, 95, 96, 99,
	100, 101, 102, 0, 0, 120, 8, -2, 122, -2,
	0, -2, 0, -2, 0, 0, -2, 48, 0, 28,
	0, 0, 0, 0, 0, 40, 39, 122, 122, 37,
	0, 75, 0, 46, 48, 48, 0, 0, 0, 0,
	-2, 6, 0, 0, 0, 108, 112, 0, 0, 0,
	0, 0, 15, 61, 0, 0, 0, 42, 0, 0,
	0, 0, 104, 111, 0, 58, 60, 0, 0, 0,
	40, 122, 0, 122, 0, 0, 0, 76, 0, 0,
	115, 0, 0, -2, -2, 0, 41, 65, 107, 0,
	50, -2, 13, -2, 26, -2, 0, 18, 0, 23,
	31, 32, 0, 0, 29, 30, 103, 0, -2, 0,
	0, 0, 0, 71, 0, 73, 35, 47, 0, -2,
	0, -2, 116, 0, 0, 27, 114, 0, 0, 0,
	0, -2, 0, 122, 113, 0, -2, -2, 0, 72,
	38, 74, -2, -2, 0, 0, 25, 16, -2, -2,
	0, 122, -2, 67, 0, 0, -2, 117, 118, 0,
	0, 0, 22, -2, 34, 68, 69, 0, 0, 17,
	21, 0, 33, 70, 119, 0, -2, -2, 0, 0,
	20, 19,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	76, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 68, 3, 3, 3, 62, 70, 3,
	71, 72, 60, 58, 55, 59, 67, 61, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 54, 66,
	57, 52, 56, 53, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 73, 3, 74, 69, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 64, 75, 65,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51,
	63,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:65
		{
			yyVAL.compstmt = nil
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:69
		{
			yyVAL.compstmt = yyDollar[1].stmts
		}
	case 3:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:74
		{
			yyVAL.stmts = nil
			if l, ok := yylex.(*Lexer); ok {
				l.stmts = yyVAL.stmts
			}
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:81
		{
			yyVAL.stmts = []ast.Stmt{yyDollar[2].stmt}
			if l, ok := yylex.(*Lexer); ok {
				l.stmts = yyVAL.stmts
			}
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:88
		{
			if yyDollar[3].stmt != nil {
				yyVAL.stmts = append(yyDollar[1].stmts, yyDollar[3].stmt)
				if l, ok := yylex.(*Lexer); ok {
					l.stmts = yyVAL.stmts
				}
			}
		}
	case 6:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:99
		{
			yyVAL.stmt = &ast.VarStmt{Names: yyDollar[2].expr_idents, Exprs: yyDollar[4].expr_many}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:104
		{
			yyVAL.stmt = &ast.LetsStmt{Lhss: []ast.Expr{yyDollar[1].expr}, Operator: "=", Rhss: []ast.Expr{yyDollar[3].expr}}
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:108
		{
			yyVAL.stmt = &ast.LetsStmt{Lhss: yyDollar[1].expr_many, Operator: "=", Rhss: yyDollar[3].expr_many}
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:112
		{
			yyVAL.stmt = &ast.BreakStmt{}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:117
		{
			yyVAL.stmt = &ast.ContinueStmt{}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:122
		{
			yyVAL.stmt = &ast.ReturnStmt{Exprs: yyDollar[2].exprs}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:127
		{
			yyVAL.stmt = &ast.ThrowStmt{Expr: yyDollar[2].expr}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 13:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:132
		{
			yyVAL.stmt = &ast.ModuleStmt{Name: yyDollar[2].tok.Lit, Stmts: yyDollar[4].compstmt}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:137
		{
			yyVAL.stmt = yyDollar[1].stmt_if
			yyVAL.stmt.SetPosition(yyDollar[1].stmt_if.Position())
		}
	case 15:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:142
		{
			yyVAL.stmt = &ast.LoopStmt{Stmts: yyDollar[3].compstmt}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 16:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:147
		{
			yyVAL.stmt = &ast.ForStmt{Var: yyDollar[2].tok.Lit, Value: yyDollar[4].expr, Stmts: yyDollar[6].compstmt}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 17:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.go.y:152
		{
			yyVAL.stmt = &ast.CForStmt{Expr1: yyDollar[2].expr_lets, Expr2: yyDollar[4].expr, Expr3: yyDollar[6].expr, Stmts: yyDollar[8].compstmt}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 18:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:157
		{
			yyVAL.stmt = &ast.LoopStmt{Expr: yyDollar[2].expr, Stmts: yyDollar[4].compstmt}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 19:
		yyDollar = yyS[yypt-13 : yypt+1]
		//line parser.go.y:162
		{
			yyVAL.stmt = &ast.TryStmt{Try: yyDollar[3].compstmt, Var: yyDollar[6].tok.Lit, Catch: yyDollar[8].compstmt, Finally: yyDollar[12].compstmt}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 20:
		yyDollar = yyS[yypt-12 : yypt+1]
		//line parser.go.y:167
		{
			yyVAL.stmt = &ast.TryStmt{Try: yyDollar[3].compstmt, Catch: yyDollar[7].compstmt, Finally: yyDollar[11].compstmt}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 21:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.go.y:172
		{
			yyVAL.stmt = &ast.TryStmt{Try: yyDollar[3].compstmt, Var: yyDollar[6].tok.Lit, Catch: yyDollar[8].compstmt}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 22:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:177
		{
			yyVAL.stmt = &ast.TryStmt{Try: yyDollar[3].compstmt, Catch: yyDollar[7].compstmt}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 23:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:182
		{
			yyVAL.stmt = &ast.SwitchStmt{Expr: yyDollar[2].expr, Cases: yyDollar[4].stmt_cases}
			yyVAL.stmt.SetPosition(yyDollar[1].tok.Position())
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:187
		{
			yyVAL.stmt = &ast.ExprStmt{Expr: yyDollar[1].expr}
			yyVAL.stmt.SetPosition(yyDollar[1].expr.Position())
		}
	case 25:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:195
		{
			yyDollar[1].stmt_if.(*ast.IfStmt).ElseIf = append(yyDollar[1].stmt_if.(*ast.IfStmt).ElseIf, &ast.IfStmt{If: yyDollar[4].expr, Then: yyDollar[6].compstmt})
			yyVAL.stmt_if.SetPosition(yyDollar[1].stmt_if.Position())
		}
	case 26:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:200
		{
			if yyVAL.stmt_if.(*ast.IfStmt).Else != nil {
				yylex.Error("multiple else statement")
			} else {
				yyVAL.stmt_if.(*ast.IfStmt).Else = append(yyVAL.stmt_if.(*ast.IfStmt).Else, yyDollar[4].compstmt...)
			}
			yyVAL.stmt_if.SetPosition(yyDollar[1].stmt_if.Position())
		}
	case 27:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:209
		{
			yyVAL.stmt_if = &ast.IfStmt{If: yyDollar[2].expr, Then: yyDollar[4].compstmt, Else: nil}
			yyVAL.stmt_if.SetPosition(yyDollar[1].tok.Position())
		}
	case 28:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:215
		{
			yyVAL.stmt_cases = []ast.Stmt{}
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:219
		{
			yyVAL.stmt_cases = []ast.Stmt{yyDollar[2].stmt_case}
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:223
		{
			yyVAL.stmt_cases = []ast.Stmt{yyDollar[2].stmt_default}
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:227
		{
			yyVAL.stmt_cases = append(yyDollar[1].stmt_cases, yyDollar[2].stmt_case)
		}
	case 32:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:231
		{
			for _, stmt := range yyDollar[1].stmt_cases {
				if _, ok := stmt.(*ast.DefaultStmt); ok {
					yylex.Error("multiple default statement")
				}
			}
			yyVAL.stmt_cases = append(yyDollar[1].stmt_cases, yyDollar[2].stmt_default)
		}
	case 33:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:242
		{
			yyVAL.stmt_case = &ast.CaseStmt{Expr: yyDollar[2].expr, Stmts: yyDollar[5].compstmt}
		}
	case 34:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:248
		{
			yyVAL.stmt_default = &ast.DefaultStmt{Stmts: yyDollar[4].compstmt}
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:254
		{
			yyVAL.expr_pair = &ast.PairExpr{Key: yyDollar[1].tok.Lit, Value: yyDollar[3].expr}
		}
	case 36:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:259
		{
			yyVAL.expr_pairs = []ast.Expr{}
		}
	case 37:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:263
		{
			yyVAL.expr_pairs = []ast.Expr{yyDollar[1].expr_pair}
		}
	case 38:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:267
		{
			yyVAL.expr_pairs = append(yyDollar[1].expr_pairs, yyDollar[4].expr_pair)
		}
	case 39:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:272
		{
			yyVAL.expr_idents = []string{}
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:276
		{
			yyVAL.expr_idents = []string{yyDollar[1].tok.Lit}
		}
	case 41:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:280
		{
			yyVAL.expr_idents = append(yyDollar[1].expr_idents, yyDollar[4].tok.Lit)
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:285
		{
			yyVAL.expr_lets = &ast.LetsExpr{Lhss: yyDollar[1].expr_many, Operator: "=", Rhss: yyDollar[3].expr_many}
		}
	case 43:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:291
		{
			yyVAL.expr_many = []ast.Expr{yyDollar[1].expr}
		}
	case 44:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:295
		{
			yyVAL.expr_many = append(yyDollar[1].exprs, yyDollar[4].expr)
		}
	case 45:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:299
		{
			yyVAL.expr_many = append(yyDollar[1].exprs, &ast.IdentExpr{Lit: yyDollar[4].tok.Lit})
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:304
		{
			yyVAL.typ = ast.Type{Name: yyDollar[1].tok.Lit}
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:308
		{
			yyVAL.typ = ast.Type{Name: yyDollar[1].typ.Name + "." + yyDollar[3].tok.Lit}
		}
	case 48:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line parser.go.y:313
		{
			yyVAL.exprs = nil
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:317
		{
			yyVAL.exprs = []ast.Expr{yyDollar[1].expr}
		}
	case 50:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:321
		{
			yyVAL.exprs = append(yyDollar[1].exprs, yyDollar[4].expr)
		}
	case 51:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:325
		{
			yyVAL.exprs = append(yyDollar[1].exprs, &ast.IdentExpr{Lit: yyDollar[4].tok.Lit})
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:331
		{
			yyVAL.expr = &ast.IdentExpr{Lit: yyDollar[1].tok.Lit}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:336
		{
			yyVAL.expr = &ast.NumberExpr{Lit: yyDollar[1].tok.Lit}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 54:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:341
		{
			yyVAL.expr = &ast.UnaryExpr{Operator: "-", Expr: yyDollar[2].expr}
			yyVAL.expr.SetPosition(yyDollar[2].expr.Position())
		}
	case 55:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:346
		{
			yyVAL.expr = &ast.UnaryExpr{Operator: "!", Expr: yyDollar[2].expr}
			yyVAL.expr.SetPosition(yyDollar[2].expr.Position())
		}
	case 56:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:351
		{
			yyVAL.expr = &ast.UnaryExpr{Operator: "^", Expr: yyDollar[2].expr}
			yyVAL.expr.SetPosition(yyDollar[2].expr.Position())
		}
	case 57:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:356
		{
			yyVAL.expr = &ast.AddrExpr{Expr: &ast.IdentExpr{Lit: yyDollar[2].tok.Lit}}
			yyVAL.expr.SetPosition(yyDollar[2].tok.Position())
		}
	case 58:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:361
		{
			yyVAL.expr = &ast.AddrExpr{Expr: &ast.MemberExpr{Expr: yyDollar[2].expr, Name: yyDollar[4].tok.Lit}}
			yyVAL.expr.SetPosition(yyDollar[2].expr.Position())
		}
	case 59:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:366
		{
			yyVAL.expr = &ast.DerefExpr{Expr: &ast.IdentExpr{Lit: yyDollar[2].tok.Lit}}
			yyVAL.expr.SetPosition(yyDollar[2].tok.Position())
		}
	case 60:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:371
		{
			yyVAL.expr = &ast.DerefExpr{Expr: &ast.MemberExpr{Expr: yyDollar[2].expr, Name: yyDollar[4].tok.Lit}}
			yyVAL.expr.SetPosition(yyDollar[2].expr.Position())
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:376
		{
			yyVAL.expr = &ast.StringExpr{Lit: yyDollar[1].tok.Lit}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:381
		{
			yyVAL.expr = &ast.ConstExpr{Value: yyDollar[1].tok.Lit}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 63:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:386
		{
			yyVAL.expr = &ast.ConstExpr{Value: yyDollar[1].tok.Lit}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:391
		{
			yyVAL.expr = &ast.ConstExpr{Value: yyDollar[1].tok.Lit}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 65:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:396
		{
			yyVAL.expr = &ast.TernaryOpExpr{Expr: yyDollar[1].expr, Lhs: yyDollar[3].expr, Rhs: yyDollar[5].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:401
		{
			yyVAL.expr = &ast.MemberExpr{Expr: yyDollar[1].expr, Name: yyDollar[3].tok.Lit}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 67:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:406
		{
			yyVAL.expr = &ast.FuncExpr{Args: yyDollar[3].expr_idents, Stmts: yyDollar[6].compstmt}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 68:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:411
		{
			yyVAL.expr = &ast.FuncExpr{Args: []string{yyDollar[3].tok.Lit}, Stmts: yyDollar[7].compstmt, VarArg: true}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 69:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line parser.go.y:416
		{
			yyVAL.expr = &ast.FuncExpr{Name: yyDollar[2].tok.Lit, Args: yyDollar[4].expr_idents, Stmts: yyDollar[7].compstmt}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 70:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.go.y:421
		{
			yyVAL.expr = &ast.FuncExpr{Name: yyDollar[2].tok.Lit, Args: []string{yyDollar[4].tok.Lit}, Stmts: yyDollar[8].compstmt, VarArg: true}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 71:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:426
		{
			yyVAL.expr = &ast.ArrayExpr{Exprs: yyDollar[3].exprs}
			if l, ok := yylex.(*Lexer); ok {
				yyVAL.expr.SetPosition(l.pos)
			}
		}
	case 72:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:431
		{
			yyVAL.expr = &ast.ArrayExpr{Exprs: yyDollar[3].exprs}
			if l, ok := yylex.(*Lexer); ok {
				yyVAL.expr.SetPosition(l.pos)
			}
		}
	case 73:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:436
		{
			mapExpr := make(map[string]ast.Expr)
			for _, v := range yyDollar[3].expr_pairs {
				mapExpr[v.(*ast.PairExpr).Key] = v.(*ast.PairExpr).Value
			}
			yyVAL.expr = &ast.MapExpr{MapExpr: mapExpr}
			if l, ok := yylex.(*Lexer); ok {
				yyVAL.expr.SetPosition(l.pos)
			}
		}
	case 74:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:445
		{
			mapExpr := make(map[string]ast.Expr)
			for _, v := range yyDollar[3].expr_pairs {
				mapExpr[v.(*ast.PairExpr).Key] = v.(*ast.PairExpr).Value
			}
			yyVAL.expr = &ast.MapExpr{MapExpr: mapExpr}
			if l, ok := yylex.(*Lexer); ok {
				yyVAL.expr.SetPosition(l.pos)
			}
		}
	case 75:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:454
		{
			yyVAL.expr = &ast.ParenExpr{SubExpr: yyDollar[2].expr}
			if l, ok := yylex.(*Lexer); ok {
				yyVAL.expr.SetPosition(l.pos)
			}
		}
	case 76:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:459
		{
			yyVAL.expr = &ast.NewExpr{Type: yyDollar[3].typ.Name}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 77:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:464
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "+", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:469
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "-", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 79:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:474
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "*", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:479
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "/", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:484
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "%", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 82:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:489
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "**", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 83:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:494
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "<<", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 84:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:499
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: ">>", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:504
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "==", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 86:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:509
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "!=", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:514
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: ">", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:519
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: ">=", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 89:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:524
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "<", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 90:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:529
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "<=", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:534
		{
			yyVAL.expr = &ast.AssocExpr{Lhs: yyDollar[1].expr, Operator: "+=", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 92:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:539
		{
			yyVAL.expr = &ast.AssocExpr{Lhs: yyDollar[1].expr, Operator: "-=", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 93:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:544
		{
			yyVAL.expr = &ast.AssocExpr{Lhs: yyDollar[1].expr, Operator: "*=", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:549
		{
			yyVAL.expr = &ast.AssocExpr{Lhs: yyDollar[1].expr, Operator: "/=", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:554
		{
			yyVAL.expr = &ast.AssocExpr{Lhs: yyDollar[1].expr, Operator: "&=", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:559
		{
			yyVAL.expr = &ast.AssocExpr{Lhs: yyDollar[1].expr, Operator: "|=", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 97:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:564
		{
			yyVAL.expr = &ast.AssocExpr{Lhs: yyDollar[1].expr, Operator: "++"}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 98:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:569
		{
			yyVAL.expr = &ast.AssocExpr{Lhs: yyDollar[1].expr, Operator: "--"}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:574
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "|", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 100:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:579
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "||", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:584
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "&", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:589
		{
			yyVAL.expr = &ast.BinOpExpr{Lhs: yyDollar[1].expr, Operator: "&&", Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 103:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:594
		{
			yyVAL.expr = &ast.CallExpr{Name: yyDollar[1].tok.Lit, SubExprs: yyDollar[3].exprs, VarArg: true}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 104:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:599
		{
			yyVAL.expr = &ast.CallExpr{Name: yyDollar[1].tok.Lit, SubExprs: yyDollar[3].exprs}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 105:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:604
		{
			yyVAL.expr = &ast.CallExpr{Name: yyDollar[2].tok.Lit, SubExprs: yyDollar[4].exprs, VarArg: true, Go: true}
			yyVAL.expr.SetPosition(yyDollar[2].tok.Position())
		}
	case 106:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:609
		{
			yyVAL.expr = &ast.CallExpr{Name: yyDollar[2].tok.Lit, SubExprs: yyDollar[4].exprs, Go: true}
			yyVAL.expr.SetPosition(yyDollar[2].tok.Position())
		}
	case 107:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:614
		{
			yyVAL.expr = &ast.AnonCallExpr{Expr: yyDollar[1].expr, SubExprs: yyDollar[3].exprs, VarArg: true}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 108:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:619
		{
			yyVAL.expr = &ast.AnonCallExpr{Expr: yyDollar[1].expr, SubExprs: yyDollar[3].exprs}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 109:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:624
		{
			yyVAL.expr = &ast.AnonCallExpr{Expr: yyDollar[2].expr, SubExprs: yyDollar[4].exprs, VarArg: true, Go: true}
			yyVAL.expr.SetPosition(yyDollar[2].expr.Position())
		}
	case 110:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:629
		{
			yyVAL.expr = &ast.AnonCallExpr{Expr: yyDollar[2].expr, SubExprs: yyDollar[4].exprs, Go: true}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 111:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:634
		{
			yyVAL.expr = &ast.ItemExpr{Value: &ast.IdentExpr{Lit: yyDollar[1].tok.Lit}, Index: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 112:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:639
		{
			yyVAL.expr = &ast.ItemExpr{Value: yyDollar[1].expr, Index: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 113:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:644
		{
			yyVAL.expr = &ast.SliceExpr{Value: &ast.IdentExpr{Lit: yyDollar[1].tok.Lit}, Begin: yyDollar[3].expr, End: yyDollar[5].expr}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 114:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line parser.go.y:649
		{
			yyVAL.expr = &ast.SliceExpr{Value: yyDollar[1].expr, Begin: yyDollar[3].expr, End: yyDollar[5].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 115:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:654
		{
			yyVAL.expr = &ast.MakeExpr{Type: yyDollar[3].typ.Name}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 116:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line parser.go.y:659
		{
			yyVAL.expr = &ast.MakeChanExpr{Type: yyDollar[4].typ.Name, SizeExpr: nil}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 117:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:664
		{
			yyVAL.expr = &ast.MakeChanExpr{Type: yyDollar[4].typ.Name, SizeExpr: yyDollar[6].expr}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 118:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line parser.go.y:669
		{
			yyVAL.expr = &ast.MakeArrayExpr{Type: yyDollar[4].typ.Name, LenExpr: yyDollar[6].expr}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 119:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line parser.go.y:674
		{
			yyVAL.expr = &ast.MakeArrayExpr{Type: yyDollar[4].typ.Name, LenExpr: yyDollar[6].expr, CapExpr: yyDollar[8].expr}
			yyVAL.expr.SetPosition(yyDollar[1].tok.Position())
		}
	case 120:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:679
		{
			yyVAL.expr = &ast.ChanExpr{Lhs: yyDollar[1].expr, Rhs: yyDollar[3].expr}
			yyVAL.expr.SetPosition(yyDollar[1].expr.Position())
		}
	case 121:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:684
		{
			yyVAL.expr = &ast.ChanExpr{Rhs: yyDollar[2].expr}
			yyVAL.expr.SetPosition(yyDollar[2].expr.Position())
		}
	case 124:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:695
		{
		}
	case 125:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:698
		{
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:703
		{
		}
	case 127:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:706
		{
		}
	}
	goto yystack /* stack new state and value */
}
