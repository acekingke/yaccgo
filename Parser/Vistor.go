/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package parser

import (
	"fmt"

	grammar "github.com/acekingke/yaccgo/Grammar"
	item "github.com/acekingke/yaccgo/Items"
	lalr "github.com/acekingke/yaccgo/LALR"
	rule "github.com/acekingke/yaccgo/Rules"
	symbol "github.com/acekingke/yaccgo/Symbol"
)

type Vistor interface {
	Process(node *Node)
}

type astDeclareVistor struct {
	precIndex  int
	idsymtabl  map[string]*Idendity
	idMaxValue int
	preIdList  []precId
	startSym   *Idendity
	union      string
	code       string
}

type precId struct {
	AssocType PrecAssocType
	Prec      int
	Id        *Idendity
}

type oneRule struct {
	LineNo     int
	LeftPart   *Idendity
	PrecIdSym  *precId
	RighPart   []*Idendity
	ActionCode string
}

type RootVistor struct {
	*RuleVistor
	CodeCpy string
	*lalr.LALR1
}

type RuleVistor struct {
	*astDeclareVistor
	rules  []*oneRule
	preMap map[string]*precId
}

type Walker struct {
	VistorNode Vistor
}

func (v *astDeclareVistor) Process(node *Node) {
	if n, ok := (*node).(*DeclareNode); ok {
		// 1. Token
		for _, tokdef := range n.TokenDefList {
			for index, id := range tokdef.IdentifyList {
				if id.Value > v.idMaxValue {
					v.idMaxValue = id.Value
				}
				if in := v.idsymtabl[id.Name]; in != nil {
					if id.Alias != "" {
						in.Alias = id.Alias
					}
					if id.Tag != "" {
						in.Tag = id.Tag
					}
					if id.Value != 0 {
						in.Value = id.Value
					}
				} else {
					v.idsymtabl[id.Name] = &tokdef.IdentifyList[index]
				}
			}
		}
		// 2. Type
		for _, typedef := range n.TypeDefList {
			if in := v.idsymtabl[typedef.IdName]; in != nil {
				in.Tag = typedef.Tag
			} else {
				//Append new name
				id := &Idendity{
					Name:  typedef.IdName,
					Tag:   typedef.Tag,
					IDTyp: NONTERMID,
					Value: 0,
				}
				v.idsymtabl[typedef.IdName] = id
			}

		}
		// 3. Prec
		for _, predefSlice := range n.PrecDefList {
			v.precIndex++
			for _, predef := range predefSlice {
				if v.idsymtabl[predef.IdName] == nil {
					errStr := fmt.Sprintf("prec symbol %s not found, please check token is exist", predef.IdName)
					panic(errStr)
				}
				v.preIdList = append(v.preIdList, precId{
					Prec:      v.precIndex,
					AssocType: predef.AssocType,
					Id:        v.idsymtabl[predef.IdName],
				})
			}
		}
		// 4. start symbol
		if len(n.StartSym) != 0 && v.idsymtabl[n.StartSym] == nil {
			//Append new name
			id := &Idendity{
				Name:  n.StartSym,
				Tag:   "",
				IDTyp: NONTERMID,
				Value: 0,
			}
			v.idsymtabl[n.StartSym] = id
			v.startSym = v.idsymtabl[n.StartSym]
		}

		//set other value
		v.code = n.CodeList
		v.union = n.Union
		for key, id := range v.idsymtabl {
			if id.Value == 0 {
				v.idMaxValue++
				v.idsymtabl[key].Value = v.idMaxValue
			}
		}
	}
}

func (v *RuleVistor) Process(node *Node) {
	if n, ok := (*node).(*RuleDefNode); ok {
		//1. create precsym map
		for index, preid := range v.preIdList {
			v.preMap[preid.Id.Name] = &v.preIdList[index]
		}
		//2. create left part id
		for _, ruledef := range n.RuleDefList {
			//fmt.Println(ruledef)
			leftpart := ruledef.LeftPart
			if v.idsymtabl[leftpart] == nil {
				//Append new name
				v.idMaxValue++
				id := &Idendity{
					Name:  leftpart,
					Tag:   "",
					IDTyp: NONTERMID,
					Value: v.idMaxValue,
				}
				v.idsymtabl[leftpart] = id
			}
		}
		//3. create the rule define
		for _, ruledef := range n.RuleDefList {
			// create ruledef
			left := v.idsymtabl[ruledef.LeftPart]
			r := &oneRule{
				LineNo:     ruledef.LineNo,
				LeftPart:   left,
				PrecIdSym:  nil,
				RighPart:   make([]*Idendity, 0),
				ActionCode: "",
			}
			for _, right := range ruledef.RightPart {
				switch right.ElemType {
				case RightActionType:
					r.ActionCode = right.Element
				case RightSyType:
					if v.idsymtabl[right.Element] == nil {
						panic("It's not define symbol ")
					}
					id := v.idsymtabl[right.Element]
					if precIdsym := v.preMap[id.Name]; precIdsym != nil {
						r.PrecIdSym = precIdsym
					}
					r.RighPart = append(r.RighPart, id)
				}
			}
			// prec sym
			if ruledef.PrecSym != "" {
				precIdsym := v.preMap[ruledef.PrecSym]
				r.PrecIdSym = precIdsym
			}
			v.rules = append(v.rules, r)
		}
	} else {
		panic("not RuleDefNode")
	}

}

func (v *RootVistor) Process(node *Node) {
	if n, ok := (*node).(*RootNode); ok {
		astv := &astDeclareVistor{
			precIndex:  0,
			idsymtabl:  make(map[string]*Idendity),
			idMaxValue: 2,
		}
		DoWalker(&n.Declare, astv)
		rulev := &RuleVistor{
			astDeclareVistor: astv,
			rules:            make([]*oneRule, 0),
			preMap:           make(map[string]*precId),
		}
		DoWalker(&n.Rules, rulev)
		v.CodeCpy = n.rest
		v.RuleVistor = rulev
	}
}

func DoWalker(node *Node, vistor Vistor) *Walker {
	w := Walker{
		VistorNode: vistor,
	}
	w.Walk(node)
	return &w
}

func (w *Walker) Walk(node *Node) {
	w.VistorNode.Process(node)
}

func (w *Walker) BuildLALR1() *lalr.LALR1 {
	g := grammar.NewGrammar()
	g.GenStartSymbol()
	dollar := symbol.NewSymbol(1, "$")
	dollar.Value = -1
	g.InsertNewSymbol(dollar)
	var SymbolFirst *symbol.Symbol
	var Idendities []*Idendity = make([]*Idendity, 0)
	var terminals []*Idendity = make([]*Idendity, 0)
	var NonTerminals []*Idendity = make([]*Idendity, 0)

	if v, ok := w.VistorNode.(*RootVistor); ok {
		//1. create symbo
		index := 1
		// first move the terminal symbol first
		for _, id := range v.idsymtabl {
			if id.IDTyp == TERMID {
				terminals = append(terminals, id)
			}
			if id.IDTyp == NONTERMID {
				NonTerminals = append(NonTerminals, id)
			}
		}
		Idendities = append(terminals, NonTerminals...)
		for _, id := range Idendities {
			if id.Value == -1 {
				continue
			}
			index++
			sy := symbol.NewSymbol(uint(index), id.Name)
			sy.SetValue(id.Value)
			if id.Tag != "" {
				sy.SetTag(id.Tag)
			}
			if id == v.startSym {
				SymbolFirst = sy
			}
			if id.IDTyp == NONTERMID {
				sy.SetNT()
			} else {
				// check prec
				if prec := v.preMap[id.Name]; prec != nil {
					switch prec.AssocType {
					case LeftAssocType:
						sy.SetPrecType(symbol.LEFT)
					case RightActionType:
						sy.SetPrecType(symbol.RIGHT)
					case NonAssocType:
						fallthrough
					default:
						sy.SetPrecType(symbol.NONE)
					}
					sy.SetPrec(prec.Prec)
				}

			}
			g.InsertNewSymbol(sy)
		}
		//2. create rules
		//2.1 insert S'->S
		g.InsertNewRules(rule.NewProductoinRule(g.StartSymbol,
			[]*symbol.Symbol{
				SymbolFirst,
			}))
		for _, onerule := range v.rules {
			// create rule
			leftsym := g.FindSymbolByName(onerule.LeftPart.Name)
			rightsyms := make([]*symbol.Symbol, 0)
			for _, right := range onerule.RighPart {
				rightsyms = append(rightsyms, g.FindSymbolByName(right.Name))
			}
			r := rule.NewProductoinRule(leftsym, rightsyms)
			if onerule.PrecIdSym != nil {
				r.SetPrecSymbol(g.FindSymbolByName(onerule.PrecIdSym.Id.Name))
			}

			g.InsertNewRules(r)

		}
		// resolve symbol
		// check the nonTerminal symbol is at left side
		for _, sym := range g.Symbols {
			if sym.IsNonTerminator {
				if _, ok := g.VnSet[sym]; !ok {
					panic(fmt.Sprintf("Check the nonterminal %s in left part of rules", sym.Name))
				}
			}
		}
		g.ResolveSymbols()
		g.CalculateEpsilonClosure()
		item_var := item.NewItem(0, 0)
		Icloures := item.NewItemCloure()
		Icloures.InsertItem(item_var)
		g.ComputeIClosure(Icloures)
		g.LR0.InsertItemClosure(Icloures, true)
		g.ComputeGotoItemRecursive(Icloures)
	} else {
		panic("not generate root node")
	}

	lalr1 := lalr.ComputeLALR(&g)
	return lalr1
}

func ParseAndBuild(input string) (*Walker, error) {
	if tr, err := Parse(input); err != nil {
		return nil, err
	} else {
		//
		var node Node = tr
		w := DoWalker(&node, &RootVistor{})
		lalr := w.BuildLALR1()
		root := w.VistorNode.(*RootVistor)
		root.LALR1 = lalr
		return w, nil
	}
}

func (v *RootVistor) GetIdsymtabl() map[string]*Idendity {
	return v.idsymtabl
}

func (v *RootVistor) GetUion() string {
	return v.union
}

func (v *RootVistor) GetCode() string {
	return v.code
}

func (v *RootVistor) GetCodeCopy() string {
	return v.CodeCpy
}

func (v *RootVistor) GetRules(index int) *oneRule {
	return v.rules[index]
}
