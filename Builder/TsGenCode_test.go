/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package builder

import "testing"

func TestTsCodeGen(t *testing.T) {
	str :=
		`// Language: typescript

	
	%{
		"use strict";
	
	%}
	
	%union {
		val :number;
	}
	
	%type	<val>	E
	%token '+'  '*'   '(' ')' 
	%left '+'  
	%left '*'  
	%token	<val>	NUM
	%token NUM 100
	%start E
%%

E:
	E '+' E {
		$$	=	$1 + $3
	}	
	| E '*' E {
		$$	=	$1 * $3
	}
	| '(' E ')' {
		$$	=	$2
	}
	| NUM {
		$$	=	$1
	}
	
%%
function GetToken(input :string, model:{ValType :ValType, pos :number}) :number {
	if (model.pos >= input.length) {
		return -1
	} else {
        model.ValType = new ValType()
        model.ValType.val = 0
		while (true) {
			if (model.pos >= input.length) {
				return -1
			}
			let c = input.charCodeAt(model.pos)	
			model.pos++
			switch (c) {
				case '$'.charCodeAt(0):
				case '+'.charCodeAt(0):
				case '('.charCodeAt(0):
				case ')'.charCodeAt(0):
				case '*'.charCodeAt(0):
					return c
				default:
					if (c >= '0'.charCodeAt(0) && c <= '9'.charCodeAt(0)) {
						model.ValType.val = model.ValType.val*10 + c - '0'.charCodeAt(0)
						if (model.pos < input.length &&
							 input.charCodeAt(model.pos) >= '0'.charCodeAt(0) && 
							 input.charCodeAt(model.pos) <= '9'.charCodeAt(0)) {
							continue;
						}
						return NUM
					}
			}
		}
		return 0;
	}
}
try {
	console.log(Parser("1+20*31").val);
}catch(e) {
	console.log(e.stack)
	console.error(e)
	
}
	`
	err := TsGenFromString(str, "../../sample.ts")
	if err != nil {
		t.Error(err)
	}
}
