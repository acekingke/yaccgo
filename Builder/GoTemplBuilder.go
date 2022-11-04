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
	"text/template"

	parser "github.com/acekingke/yaccgo/Parser"
	rules "github.com/acekingke/yaccgo/Rules"
	utils "github.com/acekingke/yaccgo/Utils"
)

// TODO: add the class model, not use global variables.
type TemplateBuilder struct {
	vnode          *parser.RootVistor
	NeedPacked     bool
	HttpParser     bool
	NTerminals     int
	HeaderPart     string
	CodeHeader     string
	ConstPart      string
	UnionPart      string
	AnalyTable     string
	PackAnalyTable string
	CodeLast       string
	StateFunc      string
	ReduceFunc     string
	Translate      string
	TranslateTrace string
	ReduceTrace    string
}

func TemplateGenFromString(input string, file string) error {
	w, err := parser.ParseAndBuild(input)
	if err != nil {
		return fmt.Errorf("parse error: %s", err)
	}
	b := NewTemplateBuilder(w)
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

	b.WriteFile(f)
	return nil
}

func NewTemplateBuilder(w *parser.Walker) *TemplateBuilder {
	return &TemplateBuilder{
		HeaderPart: `/*Generator Code , do not modify*/\n`,
		vnode:      w.VistorNode.(*parser.RootVistor),
	}
}

func (b *TemplateBuilder) buildConstPart() {
	b.NeedPacked = b.vnode.NeedPacked && utils.PackFlags
	b.HttpParser = utils.HttpDebug
	b.NTerminals = len(b.vnode.G.VtSet)
	b.CodeHeader = b.vnode.GetCode()
	b.CodeLast = b.vnode.GetCodeCopy()
	for _, identifier := range b.vnode.GetIdsymtabl() {
		if identifier.IDTyp == parser.TERMID &&
			!parser.TestPrefix(identifier.Name) {
			b.ConstPart += fmt.Sprintf("const %s = %d\n", identifier.Name, identifier.Value)
		}
	}
	b.ConstPart += fmt.Sprintf("const ERROR_ACTION = %d\nconst ACCEPT_ACTION = %d\n", b.vnode.GenErrorCode(), b.vnode.GenAcceptCode())
}

func (b *TemplateBuilder) buildUionAndCode() {
	b.UnionPart = b.vnode.GetUion()
}

func (b *TemplateBuilder) buildAnalyTable() {
	if !b.NeedPacked {
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
		b.AnalyTable = s
	} else {
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
var StackPackActDef = []int {
	%s
}
var StackPackGotoDef = []int {
	%s
}
`
		saction := ""
		soffset := ""
		scheck := ""
		sactdef := ""
		sgotodef := ""
		for _, val := range b.vnode.ActionTable {
			saction += fmt.Sprintf("%d,\t", val)
		}
		for _, val := range b.vnode.OffsetTable {
			soffset += fmt.Sprintf("%d,\t", val)
		}
		for _, val := range b.vnode.CheckTable {
			scheck += fmt.Sprintf("%d,\t", val)
		}
		for _, val := range b.vnode.ActionDef {
			sactdef += fmt.Sprintf("%d,\t", val)
		}
		for _, val := range b.vnode.GoToDef {
			sgotodef += fmt.Sprintf("%d,\t", val)
		}
		b.PackAnalyTable = fmt.Sprintf(AnalyTable, saction, soffset, scheck, sactdef, sgotodef)
	}
}

func (b *TemplateBuilder) buildStateFunc() {
}

func (b *TemplateBuilder) buildReduceFunc() {
	chooseStack := "StateSymStack"
	chooseStackPos := "StackPointer"
	pre := ""
	if utils.ObjectMode {
		chooseStack = "c.StackSym"
		chooseStackPos = "c.Stackpos"
		pre = "c."
	}
	caseCode := ""
	for i := 1; i < len(b.vnode.G.ProductoinRules); i++ {
		productionRule := b.vnode.G.ProductoinRules[i]
		caseCode += fmt.Sprintf("case %d: \n", i)
		rightPartlen := len(productionRule.RighPart)
		caseCode += fmt.Sprintf("\tdollarDolar.YySymIndex = %d\n", productionRule.LeftPart.ID)
		caseCode += fmt.Sprintf("\tDollar := "+chooseStack+"[topIndex-%d : "+chooseStackPos+"]\n\t_ = Dollar\n", rightPartlen)
		//fetch the action code here
		caseCode += actionCodeReplace(b.vnode, i, productionRule)
		caseCode += fmt.Sprintf("\t"+pre+"PopStateSym(%d)\n", rightPartlen)
	}
	b.ReduceFunc = caseCode
}

func (b *TemplateBuilder) buildTranslate() {
	caseCodes := ""
	for _, sy := range b.vnode.G.Symbols {
		if !sy.IsNonTerminator {
			caseCodes += fmt.Sprintf("\tcase %d:\n \tconv = %d\n", sy.Value, sy.ID)
		}
	}
	b.Translate = caseCodes
	caseCodes = ""
	for _, sy := range b.vnode.G.Symbols {
		caseCodes += fmt.Sprintf("\tcase %d:\n \tconv = \"%s\"\n", sy.ID, parser.RemoveTempName(sy.Name))
	}
	b.TranslateTrace = caseCodes
	caseCode := ""
	for i := 1; i < len(b.vnode.G.ProductoinRules); i++ {
		caseCode += fmt.Sprintf("\t\tcase %d: \n", i)
		oneRule := b.vnode.GetRules(i - 1)
		leftPartString := "use Reduce:" + parser.RemoveTempName(oneRule.LeftPart.Name)
		var rightPartString string = ""
		for _, rightPart := range oneRule.RighPart {
			rightPartString += parser.RemoveTempName(rightPart.Name) + " "
		}
		strTrace := fmt.Sprintf("%s -> %s",
			leftPartString, rightPartString)
		caseCode += fmt.Sprintf("\n\t\tfmt.Printf(\"look ahead %%s, %s, go to state %%d\\n\", look, s)\n", strTrace)
	}
	b.ReduceTrace = caseCode
}

func (b *TemplateBuilder) WriteFile(f *os.File) {
	chooseTemplate := goCodeTemplateStr
	if utils.ObjectMode {
		chooseTemplate = goObjectTemplateStr
	}
	templ, err := template.New("gotemplate").Parse(chooseTemplate)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := templ.Execute(f, b); err != nil {
		panic(err)
	}
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
