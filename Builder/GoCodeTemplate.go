package builder

var goCodeTemplateStr string = `
/*Generator Code , do not modify*/
// Code header part 
{{.CodeHeader}}
import 	"strings"
// const part
{{.ConstPart}}
{{ if .NeedPacked }}
// Terminal Size
const NTERMINALS = {{.NTerminals}}
{{end}}

var IsTrace bool = false
var StateSymStack = []StateSym{}
var StackPointer = 0
type Context struct {
	StackSym []StateSym
	Stackpos int
}

var globalContext = []Context{}

type ValType struct {
	// Union part
	{{.UnionPart}}
}
type StateSym struct {
	Yystate int // state
	
 	//sym val
	YySymIndex int
	//other
	ValType
}

{{ if .NeedPacked }}
 // It is NeedPacked  
 {{.PackAnalyTable}}
 func (s *StateSym) Action(a int) int {
	if StatePackOffset[s.Yystate]+a  < 0 {
		 return ERROR_ACTION
	}
	if StatePackOffset[s.Yystate]+a >= len(StackPackCheck) || 
		StackPackCheck[StatePackOffset[s.Yystate]+a] != s.Yystate {
		 if a > NTERMINALS {
			 return StackPackGotoDef[a - NTERMINALS - 1]
		 }else {
			 return StackPackActDef[s.Yystate]
		 }
	}else{
		return StatePackAction[StatePackOffset[s.Yystate]+a]
	}
}
 {{else}} 
 // It is not packed
 var StateActionArray = [][]int{
	{{.AnalyTable}}
}
func (s *StateSym) Action(a int) int {
	return StateActionArray[s.Yystate][a]
}

 {{ end }}

 func TraceShift(s *StateSym) {
	if IsTrace {
	fmt.Printf("Shift %s, push state %d\n", TraceTranslate(s.YySymIndex), s.Yystate)
	}
}

// Reduce function
func ReduceFunc(reduceIndex int) *StateSym {
	dollarDolar := &StateSym{}
	topIndex := StackPointer - 1
	switch reduceIndex {
		{{.ReduceFunc}}
	}
	return dollarDolar
}

// Push StateSym
func PushStateSym(state *StateSym) {
	TraceShift(state)
	if StackPointer >= len(StateSymStack) {
		StateSymStack = append(StateSymStack, *state)
	} else {
		StateSymStack[StackPointer] = *state
	}
	StackPointer++
}

// Pop num StateSym
func PopStateSym(num int) {
	StackPointer -= num
}


func init() {
	ParserInit()
}
func PushContex() {
	globalContext = append(globalContext, Context{
		StackSym: StateSymStack,
		Stackpos: StackPointer,
	})
}
func PopContex() {
	StackPointer = globalContext[len(globalContext)-1].Stackpos
	StateSymStack = globalContext[len(globalContext)-1].StackSym
	globalContext = globalContext[:len(globalContext)-1]
}
func ParserInit() {
	StateSymStack = []StateSym{
		{
			Yystate:    0,
			YySymIndex: 1, //$
		},
	}
	StackPointer = 1
}

func Parser(input string) *ValType {
	var currentPos int = 0
	var val ValType 
	lookAhead := fetchLookAhead(input, &val, &currentPos)
	for {

		if StackPointer == 0 {
			break
		}
		if StackPointer > len(StateSymStack) {
			break
		}
		s := &StateSymStack[StackPointer-1]
		a := s.Action(lookAhead)
		if a == ERROR_ACTION {
			lines := strings.Split(input[:currentPos], "\n")
			panic("Grammar parse error near :" + lines[len(lines)-1])
		} else if a == ACCEPT_ACTION {
			return &s.ValType
		} else {
			if a > 0 {
				// shift
				PushStateSym(&StateSym{
					Yystate:    a,
					YySymIndex: lookAhead,
					ValType:      val,
				})
				lookAhead = fetchLookAhead(input, &val, &currentPos)
			} else {
				reduceIndex := -a
				SymTy := ReduceFunc(reduceIndex)
				s := &StateSymStack[StackPointer-1]
				gotoState := s.Action(SymTy.YySymIndex)
				SymTy.Yystate = gotoState
				TraceReduce(reduceIndex, gotoState, TraceTranslate(lookAhead))
				PushStateSym(SymTy)
			}
		}
	}
	return nil
}
func fetchLookAhead(input string, val *ValType, pos *int) int {
	token := GetToken(input, val, pos)
	return translate(token)
}
func translate(c int) int {
	var conv int = 0
	switch c {
		{{.Translate}}
	}
	return conv
}
// Trace function for translate
func TraceTranslate(c int) string {
	var conv string = ""
	switch c {
{{.TranslateTrace}}
	}
	return conv
}
// Trace function for reduce
func TraceReduce(reduceIndex, s int, look string) {
	if IsTrace {
		switch reduceIndex {
{{.ReduceTrace}}
		}
	}
}
// Code Last part
{{.CodeLast}}`
