/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package builder

import "fmt"

//  make AnalyTable
func (b *GoBuilder) buildAnalyPackTable() {
	AnalyTable := `
var StatePackAction = []int {
	%s 
}
var StatePackOffset = []int {
	%s
}
var StackPackCheck = []int {
	%s
}	
`
	saction := ""
	soffset := ""
	scheck := ""
	for _, val := range b.vnode.ActionTable {
		saction += fmt.Sprintf("%d,\t", val)
	}
	for _, val := range b.vnode.OffsetTable {
		soffset += fmt.Sprintf("%d,\t", val)
	}
	for _, val := range b.vnode.CheckTable {
		scheck += fmt.Sprintf("%d,\t", val)
	}
	b.AnalyTable = fmt.Sprintf(AnalyTable, saction, soffset, scheck)
}

func (b *GoBuilder) buildPackStateFunc() {
	b.StateFunc = `
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

func (s *StateSym) Action(a int) int {
	if StatePackOffset[s.Yystate]+a  < 0|| StatePackOffset[s.Yystate]+a >= len(StackPackCheck) || StackPackCheck[StatePackOffset[s.Yystate]+a] != s.Yystate {
		return 0
	}else{
		return StatePackAction[StatePackOffset[s.Yystate]+a]
	}
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
func init() {
	ParserInit()
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
`
}
