# Command
If you run yaccgo , it will show as follow:
```
Understandable yacc generator , it can generate go/js/rust code

Usage:
  yaccgo [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  debug       open debug mode
  generate    generate filetype input.y output.go
  help        Help about any command

Flags:
  -h, --help   help for yaccgo

Use "yaccgo [command] --help" for more information about a command.
```
It has sub command like `generate, debug, help`
## sub command
### generate
The command use as follow:
```
yaccgo generate {filetype} {inputfile} {outputfile} 
```
filetype it has `go` `typescript`, and `rust` will support in feature.

**example:**
```
yaccgo generate go examples/expr.y examples/e.go
```
it will generate golang file in `examples/e.go` , and the input file is `examples/expr.y`, file type is `go`

**example**
```
yaccgo generate typescript examples/exprts.y examples/e.ts
```
It will create a typescript file in `examples/e.ts`, input file is `examples/exprts.y`

### Debug
command format as follow:
```
yaccgo debug {inputfile}
```
it will generate yacc debug information.

# grammar file specification
grammar file is '.y' as postfix files, use as input file for yaccgo.
 `examples/expr.y` and  `examples/exprts.y` above is grammar files
## Sections
grammar file consists of the following sections:

- Declarations
- Rules
- Programs

and sections are seperated by `%%` each other, so a complete grammar file is like this:
```
Declarations
%%
Rules
%%
Programs
```

### Declarations
they are consist of like that:
```
%{
Import codes
const variables define
global variables define
other header codes
%}
start define
sematics value type define
terminal symbol define
precedence define
```
all parts will persent as follow:

* header codes

 header codes are wrap by `%{ %}`,Usually they contain import codes, const defines , global variables define and so on.

  **example**
  ```
  %{
	
	package main
	
	import (

		"fmt"

	)
	
	%}
  ```
* start define 

 start define specify nonterminal symbol start in grammar. 
  ```
  %start nonterminal
  ```

 **example**
```
%start E
```
indicate E is the start symbol in grammar.

* sematics value type define

Every terminal symbol or nonterminal symbol has a sematics value type,  use `%union ` to define the value type,  
use `%token` to specify a value type to terminal symbol, 
or use `%type` to specify a value type to nonterminal symbol.

`%token`'s form is:

```

%token <tag> name

```
`%type`'s form is:

```
%type <tag> name
```
`%union` form is:

```
%union {
		XXX
	}
```
**examples**

```
%union {
    //it is go code form
		val int
	}
	
%type	<val>	E
%token	<val>	NUM
```

* terminal symbol define

define symbol is terminal symbol, the forms are:

```

%token name
%token name interger-value
```

if `%token` is not following with interger-value, the system specify identification number automatically, or following with interger-value, specify identification number by user.

* precedence define

  Precedence is a method to resolve shift/reduce conflict. for example, if the rule is

  ```
  ...
  %%
  expr	   : expr '=' expr    
         	| expr '+' expr 
         	| expr '-' expr 
         	| expr '*' expr 
           | expr '/' expr 
         	| NAME 
         	
  ```

  When `expr + expr  ` ,and next symbol is '\*', Parser do not know should shift symbol  '*' , or reduce `expr + expr  ` to expr . this grammar is ambiguous, but you do not need to rewrite the grammar, just use precedence can resolve the confict.precedence has keywords:` %left`, `%right`,  `%nonassoc` or `%precedence` , followed by  tokens. 
  
  Thus:
  
  ```
  %left '+' '-' 
  %left '*' '/'
  ```
  
  in define above,  '+', '-' is left associative, and lower precedence than '*' and '/' .  if the operators is not associative , use  `%nonassoc` or `%precedence` 

### Rules

rules form is :

```
A : Body {Actions}
```

`A` represents a nonterminal symbol, `Body `represents a list of zero or more symbol,(containing terminal symbols or nonterminal symbols) .

 rules (or production rules ) has three parts. 

* Left part : At the beginning of the rule, such as `A` above, followed by`:`
* Body or Right part: after `:` , a list of symbols.
* Actions: after the Body，one or more statements enclosed in { and }.  `$` specify  pseudo-variable for action operate the symbol sematics value. `$$` represents the left part sematics value, `$n`( n is number starting from 1) specify the nth symbol's  sematic value .

**for examples**

```
expr		   : expr '+' expr 
        	{ 
           		$$ = node( '+', $1, $3 ); 
         }
```



### Programs
Other subroutines.