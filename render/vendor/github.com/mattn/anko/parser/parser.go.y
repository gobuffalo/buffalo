%{
package parser

import (
	"github.com/mattn/anko/ast"
)

%}

%type<compstmt> compstmt
%type<stmts> stmts
%type<stmt> stmt
%type<stmt_if> stmt_if
%type<stmt_default> stmt_default
%type<stmt_case> stmt_case
%type<stmt_cases> stmt_cases
%type<typ> typ
%type<expr> expr
%type<exprs> exprs
%type<expr_many> expr_many
%type<expr_lets> expr_lets
%type<expr_pair> expr_pair
%type<expr_pairs> expr_pairs
%type<expr_idents> expr_idents

%union{
	compstmt               []ast.Stmt
	stmt_if                ast.Stmt
	stmt_default           ast.Stmt
	stmt_case              ast.Stmt
	stmt_cases             []ast.Stmt
	stmts                  []ast.Stmt
	stmt                   ast.Stmt
	typ                    ast.Type
	expr                   ast.Expr
	exprs                  []ast.Expr
	expr_many              []ast.Expr
	expr_lets              ast.Expr
	expr_pair              ast.Expr
	expr_pairs             []ast.Expr
	expr_idents            []string
	tok                    ast.Token
	term                   ast.Token
	terms                  ast.Token
	opt_terms              ast.Token
}

%token<tok> IDENT NUMBER STRING ARRAY VARARG FUNC RETURN VAR THROW IF ELSE FOR IN EQEQ NEQ GE LE OROR ANDAND NEW TRUE FALSE NIL MODULE TRY CATCH FINALLY PLUSEQ MINUSEQ MULEQ DIVEQ ANDEQ OREQ BREAK CONTINUE PLUSPLUS MINUSMINUS POW SHIFTLEFT SHIFTRIGHT SWITCH CASE DEFAULT GO CHAN MAKE OPCHAN ARRAYLIT

%right '='
%right '?' ':'
%left OROR
%left ANDAND
%left IDENT
%nonassoc EQEQ NEQ ','
%left '>' GE '<' LE SHIFTLEFT SHIFTRIGHT

%left '+' '-' PLUSPLUS MINUSMINUS
%left '*' '/' '%'
%right UNARY

%%

compstmt : opt_terms
	{
		$$ = nil
	}
	| stmts opt_terms
	{
		$$ = $1
	}

stmts :
	{
		$$ = nil
		if l, ok := yylex.(*Lexer); ok {
			l.stmts = $$
		}
	}
	| opt_terms stmt
	{
		$$ = []ast.Stmt{$2}
		if l, ok := yylex.(*Lexer); ok {
			l.stmts = $$
		}
	}
	| stmts terms stmt
	{
		if $3 != nil {
			$$ = append($1, $3)
			if l, ok := yylex.(*Lexer); ok {
				l.stmts = $$
			}
		}
	}

stmt :
	VAR expr_idents '=' expr_many
	{
		$$ = &ast.VarStmt{Names: $2, Exprs: $4}
		$$.SetPosition($1.Position())
	}
	| expr '=' expr
	{
		$$ = &ast.LetsStmt{Lhss: []ast.Expr{$1}, Operator: "=", Rhss: []ast.Expr{$3}}
	}
	| expr_many '=' expr_many
	{
		$$ = &ast.LetsStmt{Lhss: $1, Operator: "=", Rhss: $3}
	}
	| BREAK
	{
		$$ = &ast.BreakStmt{}
		$$.SetPosition($1.Position())
	}
	| CONTINUE
	{
		$$ = &ast.ContinueStmt{}
		$$.SetPosition($1.Position())
	}
	| RETURN exprs
	{
		$$ = &ast.ReturnStmt{Exprs: $2}
		$$.SetPosition($1.Position())
	}
	| THROW expr
	{
		$$ = &ast.ThrowStmt{Expr: $2}
		$$.SetPosition($1.Position())
	}
	| MODULE IDENT '{' compstmt '}'
	{
		$$ = &ast.ModuleStmt{Name: $2.Lit, Stmts: $4}
		$$.SetPosition($1.Position())
	}
	| stmt_if
	{
		$$ = $1
		$$.SetPosition($1.Position())
	}
	| FOR '{' compstmt '}'
	{
		$$ = &ast.LoopStmt{Stmts: $3}
		$$.SetPosition($1.Position())
	}
	| FOR IDENT IN expr '{' compstmt '}'
	{
		$$ = &ast.ForStmt{Var: $2.Lit, Value: $4, Stmts: $6}
		$$.SetPosition($1.Position())
	}
	| FOR expr_lets ';' expr ';' expr '{' compstmt '}'
	{
		$$ = &ast.CForStmt{Expr1: $2, Expr2: $4, Expr3: $6, Stmts: $8}
		$$.SetPosition($1.Position())
	}
	| FOR expr '{' compstmt '}'
	{
		$$ = &ast.LoopStmt{Expr: $2, Stmts: $4}
		$$.SetPosition($1.Position())
	}
	| TRY '{' compstmt '}' CATCH IDENT '{' compstmt '}' FINALLY '{' compstmt '}'
	{
		$$ = &ast.TryStmt{Try: $3, Var: $6.Lit, Catch: $8, Finally: $12}
		$$.SetPosition($1.Position())
	}
	| TRY '{' compstmt '}' CATCH '{' compstmt '}' FINALLY '{' compstmt '}'
	{
		$$ = &ast.TryStmt{Try: $3, Catch: $7, Finally: $11}
		$$.SetPosition($1.Position())
	}
	| TRY '{' compstmt '}' CATCH IDENT '{' compstmt '}'
	{
		$$ = &ast.TryStmt{Try: $3, Var: $6.Lit, Catch: $8}
		$$.SetPosition($1.Position())
	}
	| TRY '{' compstmt '}' CATCH '{' compstmt '}'
	{
		$$ = &ast.TryStmt{Try: $3, Catch: $7}
		$$.SetPosition($1.Position())
	}
	| SWITCH expr '{' stmt_cases '}'
	{
		$$ = &ast.SwitchStmt{Expr: $2, Cases: $4}
		$$.SetPosition($1.Position())
	}
	| expr
	{
		$$ = &ast.ExprStmt{Expr: $1}
		$$.SetPosition($1.Position())
	}


stmt_if :
	stmt_if ELSE IF expr '{' compstmt '}'
	{
		$1.(*ast.IfStmt).ElseIf = append($1.(*ast.IfStmt).ElseIf, &ast.IfStmt{If: $4, Then: $6})
		$$.SetPosition($1.Position())
	}
	| stmt_if ELSE '{' compstmt '}'
	{
		if $$.(*ast.IfStmt).Else != nil {
			yylex.Error("multiple else statement")
		} else {
			$$.(*ast.IfStmt).Else = append($$.(*ast.IfStmt).Else, $4...)
		}
		$$.SetPosition($1.Position())
	}
	| IF expr '{' compstmt '}'
	{
		$$ = &ast.IfStmt{If: $2, Then: $4, Else: nil}
		$$.SetPosition($1.Position())
	}

stmt_cases :
	{
		$$ = []ast.Stmt{}
	}
	| opt_terms stmt_case
	{
		$$ = []ast.Stmt{$2}
	}
	| opt_terms stmt_default
	{
		$$ = []ast.Stmt{$2}
	}
	| stmt_cases stmt_case
	{
		$$ = append($1, $2)
	}
	| stmt_cases stmt_default
	{
		for _, stmt := range $1 {
			if _, ok := stmt.(*ast.DefaultStmt); ok {
				yylex.Error("multiple default statement")
			}
		}
		$$ = append($1, $2)
	}

stmt_case :
	CASE expr ':' opt_terms compstmt
	{
		$$ = &ast.CaseStmt{Expr: $2, Stmts: $5}
	}

stmt_default :
	DEFAULT ':' opt_terms compstmt
	{
		$$ = &ast.DefaultStmt{Stmts: $4}
	}

expr_pair :
	STRING ':' expr
	{
		$$ = &ast.PairExpr{Key: $1.Lit, Value: $3}
	}

expr_pairs :
	{
		$$ = []ast.Expr{}
	}
	| expr_pair
	{
		$$ = []ast.Expr{$1}
	}
	| expr_pairs ',' opt_terms expr_pair
	{
		$$ = append($1, $4)
	}

expr_idents :
	{
		$$ = []string{}
	}
	| IDENT
	{
		$$ = []string{$1.Lit}
	}
	| expr_idents ',' opt_terms IDENT
	{
		$$ = append($1, $4.Lit)
	}

expr_lets : expr_many '=' expr_many
	{
		$$ = &ast.LetsExpr{Lhss: $1, Operator: "=", Rhss: $3}
	}

expr_many :
	expr
	{
		$$ = []ast.Expr{$1}
	}
	| exprs ',' opt_terms expr
	{
		$$ = append($1, $4)
	}
	| exprs ',' opt_terms IDENT
	{
		$$ = append($1, &ast.IdentExpr{Lit: $4.Lit})
	}

typ : IDENT
	{
		$$ = ast.Type{Name: $1.Lit}
	}
	| typ '.' IDENT
	{
		$$ = ast.Type{Name: $1.Name + "." + $3.Lit}
	}

exprs :
	{
		$$ = nil
	}
	| expr 
	{
		$$ = []ast.Expr{$1}
	}
	| exprs ',' opt_terms expr
	{
		$$ = append($1, $4)
	}
	| exprs ',' opt_terms IDENT
	{
		$$ = append($1, &ast.IdentExpr{Lit: $4.Lit})
	}

expr :
	IDENT
	{
		$$ = &ast.IdentExpr{Lit: $1.Lit}
		$$.SetPosition($1.Position())
	}
	| NUMBER
	{
		$$ = &ast.NumberExpr{Lit: $1.Lit}
		$$.SetPosition($1.Position())
	}
	| '-' expr %prec UNARY
	{
		$$ = &ast.UnaryExpr{Operator: "-", Expr: $2}
		$$.SetPosition($2.Position())
	}
	| '!' expr %prec UNARY
	{
		$$ = &ast.UnaryExpr{Operator: "!", Expr: $2}
		$$.SetPosition($2.Position())
	}
	| '^' expr %prec UNARY
	{
		$$ = &ast.UnaryExpr{Operator: "^", Expr: $2}
		$$.SetPosition($2.Position())
	}
	| '&' IDENT %prec UNARY
	{
		$$ = &ast.AddrExpr{Expr: &ast.IdentExpr{Lit: $2.Lit}}
		$$.SetPosition($2.Position())
	}
	| '&' expr '.' IDENT %prec UNARY
	{
		$$ = &ast.AddrExpr{Expr: &ast.MemberExpr{Expr: $2, Name: $4.Lit}}
		$$.SetPosition($2.Position())
	}
	| '*' IDENT %prec UNARY
	{
		$$ = &ast.DerefExpr{Expr: &ast.IdentExpr{Lit: $2.Lit}}
		$$.SetPosition($2.Position())
	}
	| '*' expr '.' IDENT %prec UNARY
	{
		$$ = &ast.DerefExpr{Expr: &ast.MemberExpr{Expr: $2, Name: $4.Lit}}
		$$.SetPosition($2.Position())
	}
	| STRING
	{
		$$ = &ast.StringExpr{Lit: $1.Lit}
		$$.SetPosition($1.Position())
	}
	| TRUE
	{
		$$ = &ast.ConstExpr{Value: $1.Lit}
		$$.SetPosition($1.Position())
	}
	| FALSE
	{
		$$ = &ast.ConstExpr{Value: $1.Lit}
		$$.SetPosition($1.Position())
	}
	| NIL
	{
		$$ = &ast.ConstExpr{Value: $1.Lit}
		$$.SetPosition($1.Position())
	}
	| expr '?' expr ':' expr
	{
		$$ = &ast.TernaryOpExpr{Expr: $1, Lhs: $3, Rhs: $5}
		$$.SetPosition($1.Position())
	}
	| expr '.' IDENT
	{
		$$ = &ast.MemberExpr{Expr: $1, Name: $3.Lit}
		$$.SetPosition($1.Position())
	}
	| FUNC '(' expr_idents ')' '{' compstmt '}'
	{
		$$ = &ast.FuncExpr{Args: $3, Stmts: $6}
		$$.SetPosition($1.Position())
	}
	| FUNC '(' IDENT VARARG ')' '{' compstmt '}'
	{
		$$ = &ast.FuncExpr{Args: []string{$3.Lit}, Stmts: $7, VarArg: true}
		$$.SetPosition($1.Position())
	}
	| FUNC IDENT '(' expr_idents ')' '{' compstmt '}'
	{
		$$ = &ast.FuncExpr{Name: $2.Lit, Args: $4, Stmts: $7}
		$$.SetPosition($1.Position())
	}
	| FUNC IDENT '(' IDENT VARARG ')' '{' compstmt '}'
	{
		$$ = &ast.FuncExpr{Name: $2.Lit, Args: []string{$4.Lit}, Stmts: $8, VarArg: true}
		$$.SetPosition($1.Position())
	}
	| '[' opt_terms exprs opt_terms ']'
	{
		$$ = &ast.ArrayExpr{Exprs: $3}
		if l, ok := yylex.(*Lexer); ok { $$.SetPosition(l.pos) }
	}
	| '[' opt_terms exprs ',' opt_terms ']'
	{
		$$ = &ast.ArrayExpr{Exprs: $3}
		if l, ok := yylex.(*Lexer); ok { $$.SetPosition(l.pos) }
	}
	| '{' opt_terms expr_pairs opt_terms '}'
	{
		mapExpr := make(map[string]ast.Expr)
		for _, v := range $3 {
			mapExpr[v.(*ast.PairExpr).Key] = v.(*ast.PairExpr).Value
		}
		$$ = &ast.MapExpr{MapExpr: mapExpr}
		if l, ok := yylex.(*Lexer); ok { $$.SetPosition(l.pos) }
	}
	| '{' opt_terms expr_pairs ',' opt_terms '}'
	{
		mapExpr := make(map[string]ast.Expr)
		for _, v := range $3 {
			mapExpr[v.(*ast.PairExpr).Key] = v.(*ast.PairExpr).Value
		}
		$$ = &ast.MapExpr{MapExpr: mapExpr}
		if l, ok := yylex.(*Lexer); ok { $$.SetPosition(l.pos) }
	}
	| '(' expr ')'
	{
		$$ = &ast.ParenExpr{SubExpr: $2}
		if l, ok := yylex.(*Lexer); ok { $$.SetPosition(l.pos) }
	}
	| NEW '(' typ ')'
	{
		$$ = &ast.NewExpr{Type: $3.Name}
		$$.SetPosition($1.Position())
	}
	| expr '+' expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "+", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr '-' expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "-", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr '*' expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "*", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr '/' expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "/", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr '%' expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "%", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr POW expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "**", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr SHIFTLEFT expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "<<", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr SHIFTRIGHT expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: ">>", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr EQEQ expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "==", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr NEQ expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "!=", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr '>' expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: ">", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr GE expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: ">=", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr '<' expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "<", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr LE expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "<=", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr PLUSEQ expr
	{
		$$ = &ast.AssocExpr{Lhs: $1, Operator: "+=", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr MINUSEQ expr
	{
		$$ = &ast.AssocExpr{Lhs: $1, Operator: "-=", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr MULEQ expr
	{
		$$ = &ast.AssocExpr{Lhs: $1, Operator: "*=", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr DIVEQ expr
	{
		$$ = &ast.AssocExpr{Lhs: $1, Operator: "/=", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr ANDEQ expr
	{
		$$ = &ast.AssocExpr{Lhs: $1, Operator: "&=", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr OREQ expr
	{
		$$ = &ast.AssocExpr{Lhs: $1, Operator: "|=", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr PLUSPLUS
	{
		$$ = &ast.AssocExpr{Lhs: $1, Operator: "++"}
		$$.SetPosition($1.Position())
	}
	| expr MINUSMINUS
	{
		$$ = &ast.AssocExpr{Lhs: $1, Operator: "--"}
		$$.SetPosition($1.Position())
	}
	| expr '|' expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "|", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr OROR expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "||", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr '&' expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "&", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| expr ANDAND expr
	{
		$$ = &ast.BinOpExpr{Lhs: $1, Operator: "&&", Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| IDENT '(' exprs VARARG ')'
	{
		$$ = &ast.CallExpr{Name: $1.Lit, SubExprs: $3, VarArg: true}
		$$.SetPosition($1.Position())
	}
	| IDENT '(' exprs ')'
	{
		$$ = &ast.CallExpr{Name: $1.Lit, SubExprs: $3}
		$$.SetPosition($1.Position())
	}
	| GO IDENT '(' exprs VARARG ')'
	{
		$$ = &ast.CallExpr{Name: $2.Lit, SubExprs: $4, VarArg: true, Go: true}
		$$.SetPosition($2.Position())
	}
	| GO IDENT '(' exprs ')'
	{
		$$ = &ast.CallExpr{Name: $2.Lit, SubExprs: $4, Go: true}
		$$.SetPosition($2.Position())
	}
	| expr '(' exprs VARARG ')'
	{
		$$ = &ast.AnonCallExpr{Expr: $1, SubExprs: $3, VarArg: true}
		$$.SetPosition($1.Position())
	}
	| expr '(' exprs ')'
	{
		$$ = &ast.AnonCallExpr{Expr: $1, SubExprs: $3}
		$$.SetPosition($1.Position())
	}
	| GO expr '(' exprs VARARG ')'
	{
		$$ = &ast.AnonCallExpr{Expr: $2, SubExprs: $4, VarArg: true, Go: true}
		$$.SetPosition($2.Position())
	}
	| GO expr '(' exprs ')'
	{
		$$ = &ast.AnonCallExpr{Expr: $2, SubExprs: $4, Go: true}
		$$.SetPosition($1.Position())
	}
	| IDENT '[' expr ']'
	{
		$$ = &ast.ItemExpr{Value: &ast.IdentExpr{Lit: $1.Lit}, Index: $3}
		$$.SetPosition($1.Position())
	}
	| expr '[' expr ']'
	{
		$$ = &ast.ItemExpr{Value: $1, Index: $3}
		$$.SetPosition($1.Position())
	}
	| IDENT '[' expr ':' expr ']'
	{
		$$ = &ast.SliceExpr{Value: &ast.IdentExpr{Lit: $1.Lit}, Begin: $3, End: $5}
		$$.SetPosition($1.Position())
	}
	| expr '[' expr ':' expr ']'
	{
		$$ = &ast.SliceExpr{Value: $1, Begin: $3, End: $5}
		$$.SetPosition($1.Position())
	}
	| MAKE '(' typ ')'
	{
		$$ = &ast.MakeExpr{Type: $3.Name}
		$$.SetPosition($1.Position())
	}
	| MAKE '(' CHAN typ ')'
	{
		$$ = &ast.MakeChanExpr{Type: $4.Name, SizeExpr: nil}
		$$.SetPosition($1.Position())
	}
	| MAKE '(' CHAN typ ',' expr ')'
	{
		$$ = &ast.MakeChanExpr{Type: $4.Name, SizeExpr: $6}
		$$.SetPosition($1.Position())
	}
	| MAKE '(' ARRAYLIT typ ',' expr ')'
	{
		$$ = &ast.MakeArrayExpr{Type: $4.Name, LenExpr: $6}
		$$.SetPosition($1.Position())
	}
	| MAKE '(' ARRAYLIT typ ',' expr ',' expr ')'
	{
		$$ = &ast.MakeArrayExpr{Type: $4.Name, LenExpr: $6, CapExpr: $8}
		$$.SetPosition($1.Position())
	}
	| expr OPCHAN expr
	{
		$$ = &ast.ChanExpr{Lhs: $1, Rhs: $3}
		$$.SetPosition($1.Position())
	}
	| OPCHAN expr
	{
		$$ = &ast.ChanExpr{Rhs: $2}
		$$.SetPosition($2.Position())
	}

opt_terms : /* none */
	| terms
	;


terms : term
	{
	}
	| terms term
	{
	}
	;

term : ';'
	{
	}
	| '\n'
	{
	}
	;

%%
