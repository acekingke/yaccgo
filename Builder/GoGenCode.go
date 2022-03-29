/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package builder

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	parser "github.com/acekingke/yaccgo/Parser"
	rules "github.com/acekingke/yaccgo/Rules"
)

type GoBuilder struct {
	vnode      *parser.RootVistor
	HeaderPart string
	CodeHeader string
	ConstPart  string
	UnionPart  string
	AnalyTable string
	CodeLast   string
	StateFunc  string
	ReduceFunc string
	Translate  string
}

func NewGoBuilder(w *parser.Walker) *GoBuilder {
	return &GoBuilder{
		HeaderPart: `/*Generator Code , do not modify*/\n`,
		vnode:      w.VistorNode.(*parser.RootVistor),
	}
}

func (b *GoBuilder) buildConstPart() {
	b.ConstPart = "// const part \n"
	for _, identifier := range b.vnode.GetIdsymtabl() {
		if identifier.IDTyp == parser.TERMID &&
			!parser.TestPrefix(identifier.Name) {
			b.ConstPart += fmt.Sprintf("const %s = %d\n", identifier.Name, identifier.Value)
		}
	}
	b.ConstPart += fmt.Sprintf("const ERROR_ACTION = 0\nconst ACCEPT_ACTION = %d\n", b.vnode.GenAcceptCode())
}

func (b *GoBuilder) buildUionAndCode() {
	str := `
var StateSymStack = []StateSym{}
var StackPointer = 0
type ValType struct {
	%s
}
type StateSym struct {
	Yystate int // state
	
 	//sym val
	YySymIndex int
	//other
	ValType
}`
	str = fmt.Sprintf(str, b.vnode.GetUion())
	b.UnionPart = str
	b.CodeHeader = b.vnode.GetCode()
	b.CodeLast = b.vnode.GetCodeCopy()
}

//  make AnalyTable
func (b *GoBuilder) buildAnalyTable() {
	AnalyTable := `
var StateActionArray = [][]int {
	%s 
}
`
	s := "/*     "
	for _, sy := range b.vnode.G.Symbols {
		s += fmt.Sprintf("%s\t", sy.Name)
	}
	s += "*/\n"
	for index, row := range b.vnode.GTable {
		s += fmt.Sprintf("/* %d */ {", index)
		for _, val := range row {
			s += fmt.Sprintf("%d,\t", val)
		}
		s += "},\n"
	}
	b.AnalyTable = fmt.Sprintf(AnalyTable, s)
}

func (b *GoBuilder) buildStateFunc() {
	b.StateFunc = `
// Push StateSym
func PushStateSym(state *StateSym) {
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
	return StateActionArray[s.Yystate][a]
}

func init() {
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
			panic("Grammar parse error")
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

func (b *GoBuilder) buildTranslate() {
	str := `
func translate(c int) int {
	var conv int = 0
	switch c {
%s	
	}
	return conv
}
`
	caseCodes := ""
	for _, sy := range b.vnode.G.Symbols {
		if !sy.IsNonTerminator {
			caseCodes += fmt.Sprintf("\tcase %d:\n \tconv = %d\n", sy.Value, sy.ID)
		}
	}
	b.Translate = fmt.Sprintf(str, caseCodes)
}

// make ReduceFunc
func (b *GoBuilder) buildReduceFunc() {
	str := `
func ReduceFunc(reduceIndex int) *StateSym {
	dollarDolar := &StateSym{}
	topIndex := StackPointer - 1
	switch reduceIndex {
		%s
	}
	return dollarDolar
}
`
	caseCode := ""
	for i := 1; i < len(b.vnode.G.ProductoinRules); i++ {
		productionRule := b.vnode.G.ProductoinRules[i]
		caseCode += fmt.Sprintf("case %d: \n", i)
		rightPartlen := len(productionRule.RighPart)
		caseCode += fmt.Sprintf("\tdollarDolar.YySymIndex = %d\n", productionRule.LeftPart.ID)
		caseCode += fmt.Sprintf("\tDollar := StateSymStack[topIndex-%d : StackPointer]\n\t_ = Dollar\n", rightPartlen)
		//fetch the action code here
		caseCode += actionCodeReplace(b.vnode, i, productionRule)
		caseCode += fmt.Sprintf("\tPopStateSym(%d)\n", rightPartlen)
	}
	b.ReduceFunc = fmt.Sprintf(str, caseCode)
}

func actionCodeReplace(vnode *parser.RootVistor,
	index int, pr *rules.ProductoinRule) string {
	oneRule := vnode.GetRules(index - 1)
	//  generate the comments.
	strComment := "\n/*\n%s*/\n"
	leftPartString := fmt.Sprint("\nLineNo:", oneRule.LineNo, "\n") + parser.RemoveTempName(oneRule.LeftPart.Name)
	var rightPartString string = ""
	for _, rightPart := range oneRule.RighPart {
		rightPartString += parser.RemoveTempName(rightPart.Name) + " "
	}
	strComment = fmt.Sprintf(strComment,
		fmt.Sprintf("%s -> %s\n %s\n",
			leftPartString, rightPartString, oneRule.ActionCode))

	str := oneRule.ActionCode
	str = strings.ReplaceAll(str, "$$",
		fmt.Sprintf("dollarDolar.%s", pr.LeftPart.Tag))

	// find the $ and digits
	reg := regexp.MustCompile(`\$[0-9]+`)
	str = reg.ReplaceAllStringFunc(str, func(s string) string {
		index := s[1:]
		i, _ := strconv.Atoi(index)
		return fmt.Sprintf("Dollar[%s].%s", index, pr.RighPart[i-1].Tag)
	})
	return strComment + str + "\n"
}

func GoGenFromString(input string, file string) error {
	w, err := parser.ParseAndBuild(input)
	if err != nil {
		return fmt.Errorf("parse error: %s", err)
	}
	b := NewGoBuilder(w)
	b.buildConstPart()
	b.buildUionAndCode()
	if b.vnode.NeedPacked {
		b.buildAnalyPackTable()
		b.buildPackStateFunc()
	} else {
		b.buildAnalyTable()
		b.buildStateFunc()
	}

	b.buildReduceFunc()
	b.buildTranslate()
	// Create file and write to it
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("create file error: %s", err)
	}
	f.WriteString(b.CodeHeader)
	f.WriteString(b.ConstPart)
	f.WriteString(b.UnionPart)
	f.WriteString(b.AnalyTable)
	f.WriteString(b.CodeLast)
	f.WriteString(b.StateFunc)
	f.WriteString(b.ReduceFunc)
	f.WriteString(b.Translate)
	f.Close()
	return nil
}
