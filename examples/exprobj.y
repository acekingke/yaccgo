// language: go
%{
package main
import (
	"fmt"
)
%}

%union {
	val int
}

%type	<val>	E
%token '+'  '*'   '(' ')' '/' '>'
%nonassoc '>'
%left '+'  
%left '*'  '/'
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
	| E '/' E {
		$$	=	$1 / $3
	}
	| E '>' E {
		if 	$1 > $3 {
			$$	=	1
		} else {
			$$	=	0
		}
	}
	| '(' E ')' {
		$$	=	$2
	}
	| NUM {
		$$	=	$1
	}
	
%%
const EOF = -1
// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
func GetToken(input string, valTy *ValType, pos *int) int {
	if *pos >= len(input) {
		return -1
	} else {
		*valTy = ValType{0}
	loop:
		if *pos >= len(input) {
			return EOF
		}
		c := input[*pos]
		*pos++
		switch c {
		case '+':
			fallthrough
		case '(':
			fallthrough
		case ')':
			fallthrough

		case '*':
			return int(c)
		case '/':
			return int(c)
		case '>':
			return int(c)
		default:
			if c >= '0' && c <= '9' { // is digit
				valTy.val = (valTy.val)*10 + int(c) - '0'
				// next is digit
				if *pos < len(input) && input[*pos] >= '0' && input[*pos] <= '9' {
					goto loop
				}
				return NUM
			}

		}
		return 0
	}
}

func main() {
	c := MakeParserContext()

	v := c.Parser("1+2*31>2").val
	fmt.Println(v)
}
