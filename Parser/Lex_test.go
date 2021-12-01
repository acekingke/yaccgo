/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package parser

import (
	"fmt"
	"testing"
)

func TestComment(t *testing.T) {
	//1,
	l := Lex("//zzggag\n")
	fmt.Println(l.nextToken())

	//2. /**/
	l = Lex("/*xagaga g*/")
	fmt.Println(l.nextToken())

	txt := `// Copyright 2013 The Go Authors. All rights reserved.
	// Use of this source code is governed by a BSD-style
	// license that can be found in the LICENSE file.
	
	// This is an example of a goyacc program.
	// To build it:
	// goyacc -p "expr" expr.y (produces y.go)
	// go build -o expr y.go
	// expr
	// > <type an expression>`
	if l = Lex(txt); l.nextToken().Kind == tokenError {
		t.Error("err")
	}

}

func TestDirective(t *testing.T) {
	fmt.Println("TestDirective")
	l := Lex(`%% %{ vewry nice %} %type %token `)
	for {
		tok := l.nextToken()
		if tok.Kind == tokenError {
			t.Error(tok.Value)
			break
		} else if tok.Kind == EOF {
			fmt.Println(tok)
			break
		} else {
			fmt.Println(tok)
		}

	}
}

func TestIdentify(t *testing.T) {
	fmt.Println("TestIdentify")
	l := Lex("thello-")
	if tok := l.nextToken(); tok.Kind != tokenError {
		fmt.Println(tok)
	} else {
		t.Error(tok.Value)
	}
	if tok := l.nextToken(); tok.Kind != tokenError {
		fmt.Println(tok)
	} else {
		fmt.Printf("error toke %s\n", tok.Value)
	}
	l = Lex("thello_")
	if tok := l.nextToken(); tok.Kind != tokenError {
		fmt.Println(tok)
	} else {
		t.Error(tok.Value)
	}
	if tok := l.nextToken(); tok.Kind != tokenError {
		fmt.Println(tok)
	} else {
		t.Error(tok.Value)
	}
}

func TestToken(t *testing.T) {
	fmt.Println("===TestToken===")
	s := "%type	<num>	expr expr1 expr2 expr3"
	l := Lex(s)
	for {
		if tok := l.nextToken(); tok.Kind != tokenError {
			fmt.Println(tok)
			if tok.Kind == EOF {
				break
			}
		} else {
			t.Error(tok.Value)
		}
	}

	s = "%token '+' '-' '*' '/' '(' ')'"
	l = Lex(s)
	for {
		if tok := l.nextToken(); tok.Kind != tokenError {
			fmt.Println(tok)
			if tok.Kind == EOF {
				break
			}
		} else {
			t.Error(tok.Value)
		}
	}
	s = "%union { xx } %type"
	l = Lex(s)
	for {
		if tok := l.nextToken(); tok.Kind != tokenError {
			fmt.Println(tok)
			if tok.Kind == EOF {
				break
			}
		} else {
			t.Error(tok.Value)
		}
	}
	s = `expr:
	expr1
|	'+' expr`
	l = Lex(s)
	for {
		if tok := l.nextToken(); tok.Kind != tokenError {
			fmt.Println(tok)
			if tok.Kind == EOF {
				break
			}
		} else {
			t.Error(tok.Value)
		}
	}
	s = `"jellag"`
	l = Lex(s)
	for {
		if tok := l.nextToken(); tok.Kind != tokenError {
			fmt.Println(tok)
			if tok.Kind == EOF {
				break
			}
		} else {
			t.Error(tok.Value)
		}
	}
}
