/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package parser

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	str :=
		`// Copyright 2013 The Go Authors. All rights reserved.
	// Use of this source code is governed by a BSD-style
	// license that can be found in the LICENSE file.
	
	// This is an example of a goyacc program.
	// To build it:
	// goyacc -p "expr" expr.y (produces y.go)
	// go build -o expr y.go
	// expr
	// > <type an expression>
	
	%{
	
	package main
	
	import (
		"bufio"
		"bytes"
		"fmt"
		"io"
		"log"
		"math/big"
		"os"
		"unicode/utf8"
	)
	
	%}
	
	%union {
		num *big.Rat
	}
	
	%type	<num>	expr expr1 expr2 expr3
	
	%token '+' '-' '*' '/' '(' ')' MINUS
	%left '+' '-'
	%left '*' '/'
	%right MINUS
	%token	<num>	NUM
	%token NUM 100
	%start top
	%%
	
	top:
	expr
	{
		if $1.IsInt() {
			fmt.Println($1.Num().String())
		} else {
			fmt.Println($1.String())
		}
	}

expr:
	expr1
|	'+' expr
	{
		$$ = $2
	}
|	'-' expr %prec MINUS
	{
		$$ = $2.Neg($2)
	}

expr1:
	expr2
|	expr1 '+' expr2
	{
		$$ = $1.Add($1, $3)
	}
|	expr1 '-' expr2
	{
		$$ = $1.Sub($1, $3)
	}

expr2:
	expr3
|	expr2 '*' expr3
	{
		$$ = $1.Mul($1, $3)
	}
|	expr2 '/' expr3
	{
		$$ = $1.Quo($1, $3)
	}

expr3:
	NUM
|	'(' expr ')'
	{
		$$ = $2
	}

	
	
	
	%%
	
	// The parser expects the lexer to return 0 on EOF.  Give it a name
	// for clarity.
	const eof = 0
	
	// The parser uses the type <prefix>Lex as a lexer. It must provide
	// the methods Lex(*<prefix>SymType) int and Error(string).
	type exprLex struct {
		line []byte
		peek rune
	}	
	`
	if tr, err := Parse(str); err != nil {
		t.Error(err)
	} else {
		// work in test
		var node Node = tr
		w := DoWalker(&node, &RootVistor{})
		lalr := w.BuildLALR1()
		fmt.Println(lalr)
		root := w.VistorNode.(*RootVistor)
		root.LALR1 = lalr
	}
}

func TestParser2(t *testing.T) {
	str :=
		`// Copyright 2013 The Go Authors. All rights reserved.
	// Use of this source code is governed by a BSD-style
	// license that can be found in the LICENSE file.
	
	// This is an example of a goyacc program.
	// To build it:
	// goyacc -p "expr" expr.y (produces y.go)
	// go build -o expr y.go
	// expr
	// > <type an expression>
	
	%{
	
	package main
	
	import (
		"bufio"
		"bytes"
		"fmt"
		"io"
		"log"
		"math/big"
		"os"
		"unicode/utf8"
	)
	
	%}
	
%left <tga> A B 
%left C
%start s
%%
s : A
`
	if tr, err := Parse(str); err != nil {
		t.Error(err)
	} else {
		// work in test
		var node Node = tr
		w := DoWalker(&node, &RootVistor{})
		lalr := w.BuildLALR1()
		fmt.Println(lalr)
		root := w.VistorNode.(*RootVistor)
		root.LALR1 = lalr
	}
}

func TestParser3(t *testing.T) {
	str := `
	%{
		package main
		%}
		
		%union{
		String string
		Expr expr 
		}
		
		
		%token<String>  IDENTIFIER
		%token<String> NUMBER 100 
		%type <Expr> expr assignment
		
		%left '+' '-'
		%left '*' '/'
		%%
		start: expr {yylex.(*interpreter).parseResult = &astRoot{$1}} 
			 | assignment {yylex.(*interpreter).parseResult = $1}
			 ;
		
		expr:
			  NUMBER {$$ = &number{$1} }
			| IDENTIFIER { $$ = &variable{$1}}
			| expr '+' expr { $$ = &binaryExpr{Op: '+', lhs: $1, rhs: $3} }
			| expr '-' expr { $$ = &binaryExpr{Op: '-', lhs: $1, rhs: $3} }
			| expr '*' expr { $$ = &binaryExpr{Op: '*', lhs: $1, rhs: $3} }
			| expr '/' expr { $$ = &binaryExpr{Op: '/', lhs: $1, rhs: $3} }
			| '(' expr ')'  { $$ = &parenExpr{$2}}
			| '-' expr %prec '*' { $$ = &unaryExpr{$2} }
			;
			
		
		assignment:
				  IDENTIFIER '=' expr {$$ = &assignment{$1, $3}};
		%%
`
	if tr, err := Parse(str); err != nil {
		t.Error(err)
	} else {
		// work in test
		var node Node = tr
		w := DoWalker(&node, &RootVistor{})
		lalr := w.BuildLALR1()
		fmt.Println(lalr)
		root := w.VistorNode.(*RootVistor)
		root.LALR1 = lalr
	}
}
