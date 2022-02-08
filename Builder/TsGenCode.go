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

type TsBuilder struct {
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

func NewTsBuilder(w *parser.Walker) *TsBuilder {
	return &TsBuilder{
		HeaderPart: `/*Generator Code , do not modify*/\n`,
		vnode:      w.VistorNode.(*parser.RootVistor),
	}
}

func (b *TsBuilder) buildConstPart() {
	b.ConstPart = "// const part \n"
	for _, identifier := range b.vnode.GetIdsymtabl() {
		if identifier.IDTyp == parser.TERMID &&
			!parser.TestPrefix(identifier.Name) {
			b.ConstPart += fmt.Sprintf("const %s = %d\n", identifier.Name, identifier.Value)
		}
	}
	b.ConstPart += fmt.Sprintf("const ERROR_ACTION = 0 \nconst ACCEPT_ACTION = %d\n", b.vnode.GenAcceptCode())
}

func (b *TsBuilder) buildUionAndCode() {
	str := `
var StateSymStack :StateSym[] = [];
var StackPointer = 0;
class ValType {
    %s
};
class StateSym  {
	Yystate :number; // state
	//sym val
	 YySymIndex :number; 
	//other
	ValType :ValType;
    constructor(Yystate :number, YySymIndex :number) {
        this.Yystate = Yystate;
        this.YySymIndex = YySymIndex;
       
    }
    // getter action code
    Action(a :number) :number {
        return StateActionArray[this.Yystate][a]
    }
    
};`
	str = fmt.Sprintf(str, b.vnode.GetUion())
	b.UnionPart = str
	b.CodeHeader = b.vnode.GetCode()
	b.CodeLast = b.vnode.GetCodeCopy()
}

//  make AnalyTable
func (b *TsBuilder) buildAnalyTable() {
	AnalyTable := `
var StateActionArray :number[][] =[
	%s 
]
`
	s := "/*     "
	for _, sy := range b.vnode.G.Symbols {
		s += fmt.Sprintf("%s\t", sy.Name)
	}
	s += "*/\n"
	for index, row := range b.vnode.GTable {
		s += fmt.Sprintf("/* %d */ [", index)
		for _, val := range row {
			s += fmt.Sprintf("%d,\t", val)
		}
		s += "],\n"
	}
	b.AnalyTable = fmt.Sprintf(AnalyTable, s)
}

func (b *TsBuilder) buildStateFunc() {
	b.StateFunc = `
function PushStateSym(state:StateSym) {
	if (StackPointer >= StateSymStack.length) {
		StateSymStack.push(state);
	} else {
		StateSymStack[StackPointer] = state;
	}
	StackPointer++;
}

// Pop num StateSym
function PopStateSym(num :number) {
	StackPointer -= num
}
function initialize() {
    StateSymStack = [new StateSym(0,1)];
    StackPointer = 1;
}

function Parser(input :string) :ValType {
	var currentPos :number = 0
	var val :ValType
	const model = {ValType :val, pos :currentPos}
	var lookAhead = fetchLookAhead(input, model)
	while (true) {
		if (StackPointer == 0) {
			break
		}
		if (StackPointer > StateSymStack.length) {
			break
		}
		let state = StateSymStack[StackPointer - 1]
		let action = state.Action(lookAhead)
		if (action == ERROR_ACTION) {
			console.error("Grammer error")
			break
		}else if (action == ACCEPT_ACTION) {
			return state.ValType
		}else {
			if (action > 0) {
				// shift
				let sym = new StateSym(action, lookAhead)
				sym.ValType = model.ValType
				PushStateSym(sym)
				lookAhead = fetchLookAhead(input, model)
			}else {
				// reduce
				let SymTy = ReduceFunc(-action)
				state = StateSymStack[StackPointer-1]
				let gotoState = state.Action(SymTy.YySymIndex)
				SymTy.Yystate = gotoState
				PushStateSym(SymTy)
			}
		}

	}
    return null;
}
function fetchLookAhead(input :string, 
	model:{ValType :ValType, pos :number})  {
	let token = GetToken(input, model)
	 return translate(token)   
 
}
`
}

func (b *TsBuilder) buildTranslate() {
	str := `
function translate(c :number) :number {
	var conv :number = 0
	switch (c) {
%s
	}
	return conv;
}
`
	caseCodes := ""
	for _, sy := range b.vnode.G.Symbols {
		if !sy.IsNonTerminator {
			caseCodes += fmt.Sprintf("\tcase %d:\n \tconv = %d;\nbreak;\n", sy.Value, sy.ID)
		}
	}
	b.Translate = fmt.Sprintf(str, caseCodes)
}

// make ReduceFunc
func (b *TsBuilder) buildReduceFunc() {
	str := `
function ReduceFunc(reduceIndex :number) :StateSym {
	let dollarDolar = new StateSym(-1,-1)
	dollarDolar.ValType = new ValType()
	let topIndex = StackPointer - 1
	switch (reduceIndex) {
%s
	}
	return dollarDolar;
}
`
	caseCode := ""
	for i := 1; i < len(b.vnode.G.ProductoinRules); i++ {
		productionRule := b.vnode.G.ProductoinRules[i]
		caseCode += fmt.Sprintf("case %d: {\n", i)
		rightPartlen := len(productionRule.RighPart)
		caseCode += fmt.Sprintf("\tdollarDolar.YySymIndex = %d\n", productionRule.LeftPart.ID)
		caseCode += fmt.Sprintf("\tlet Dollar = StateSymStack.slice(topIndex-%d , StackPointer);\n", rightPartlen)
		//fetch the action code here
		caseCode += actionCodeReplaceTs(b.vnode, i, productionRule)
		caseCode += fmt.Sprintf("\tPopStateSym(%d);\n\tbreak;\n}\n", rightPartlen)
	}
	b.ReduceFunc = fmt.Sprintf(str, caseCode) + "\ninitialize();\n"
}

func actionCodeReplaceTs(vnode *parser.RootVistor,
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
		fmt.Sprintf("dollarDolar.ValType.%s", pr.LeftPart.Tag))

	// find the $ and digits
	reg := regexp.MustCompile(`\$[0-9]+`)
	str = reg.ReplaceAllStringFunc(str, func(s string) string {
		index := s[1:]
		i, _ := strconv.Atoi(index)
		return fmt.Sprintf("Dollar[%s].ValType.%s", index, pr.RighPart[i-1].Tag)
	})
	return strComment + str + "\n"
}

func TsGenFromString(input string, file string) error {
	w, err := parser.ParseAndBuild(input)
	if err != nil {
		return fmt.Errorf("parse error: %s", err)
	}
	b := NewTsBuilder(w)
	b.buildConstPart()
	b.buildUionAndCode()
	b.buildAnalyTable()
	b.buildStateFunc()
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
	f.WriteString(b.StateFunc)
	f.WriteString(b.ReduceFunc)
	f.WriteString(b.Translate)
	f.WriteString(b.CodeLast)
	f.Close()
	return nil
}
