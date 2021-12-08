%{
package main
import "fmt"
%}

%token   'n'
%start L
%%
L :  /*empty*/
    |E L
E: 'n'
%%
func GetToken(input string, valTy *ValType, pos *int) int {
    if *pos >= len(input) {
        return -1
    }
    c := input[*pos]
    *pos++
    switch c {
    case 'n':
        *valTy = ValType{}
        return 'n'
    default:
        return 0
    }
}
func main() {
	v := Parser("nnn")
	fmt.Println(v)
}