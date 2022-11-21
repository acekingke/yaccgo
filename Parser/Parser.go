/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package parser

import (
	"fmt"
	"strconv"
)

type parser struct {
	lex       *lexer
	current   Token
	tokenArr  [3]Token
	peekCount int
	pos       int
	err       error

	// map Def
	TokenDefMap map[string]bool
}

type Node interface {
	Location() Location
	SetLocation(Location)
}

type Base struct {
	loc Location
}

func (b *Base) Location() Location {
	return b.loc
}

func (b *Base) SetLocation(l Location) {
	b.loc = l
}

// Nil Node

type NilNode struct {
	Base
}
type IDType int

const (
	TERMID = iota + 1
	NONTERMID
)

type Idendity struct {
	Name  string // ID
	IDTyp IDType
	Value int //value
	Tag   string
	Alias string
}

// token define
type TokenDef struct {
	IdentifyList []Idendity
}

type PrecAssocType int

const (
	LeftAssocType = iota + 1
	RightAssocype
	NonAssocType
)

// prec
type PrecDef struct {
	Prec      int
	AssocType PrecAssocType
	IdName    string
}

// type
//%type	<num>	expr expr1 expr2 expr3
type TypeDef struct {
	Tag    string
	IdName string
}
type DeclareNode struct {
	Base
	CodeList     string
	TokenDefList []TokenDef
	PrecDefList  [][]PrecDef
	TypeDefList  []TypeDef
	Union        string
	StartSym     string
}
type RightSymActType int

const (
	RightSyType = iota + 1
	RightActionType
)

type RightSymOrAction struct {
	ElemType RightSymActType
	Element  string
}

// Rules
type RuleDef struct {
	LineNo    int
	LeftPart  string
	RightPart []RightSymOrAction
	PrecSym   string // which prec symbol
}

//RuleDefNode
type RuleDefNode struct {
	Base
	RuleDefList []RuleDef
}

// Root
type RootNode struct {
	Base
	Declare Node
	Rules   Node
	rest    string
}

func Parse(input string) (*RootNode, error) {

	p := &parser{
		lex: Lex(input),
		pos: 0,
	}
	p.TokenDefMap = make(map[string]bool)
	var nodeDeclare Node
	if nodeDeclare = p.parseDeclare(); nodeDeclare == nil {
		return nil, fmt.Errorf("do not has declare %s", p.err)
	}
	decl, _ := nodeDeclare.(*DeclareNode)
	if !p.current.Is(Section) {
		return nil, fmt.Errorf(
			fmt.Sprintf("parser error! %s", p.current.Value),
		)
	}
	RuDlist := make([]RuleDef, 0)
	p.next() // get the first identify
	for {
		if ruleslice := p.parseRule(&decl.TokenDefList); ruleslice == nil {
			break
		} else {
			RuDlist = append(RuDlist, ruleslice...)
		}
	}

	if !p.current.Is(Section) && !p.current.Is(EOF) {
		return nil, fmt.Errorf(fmt.Sprintf("parser err :%s", p.current.Value))
	}
	restcode := p.lex.input[p.current.EndAt:]
	return &RootNode{
			Declare: nodeDeclare,
			Rules:   &RuleDefNode{RuleDefList: RuDlist},
			rest:    restcode,
		},
		nil

}

func (p *parser) next() {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.tokenArr[0] = p.lex.nextToken()
	}
	p.current = p.tokenArr[p.peekCount]
}

// backup backs the input stream up one token.
func (p *parser) backup() {
	p.peekCount++
}

// backup2 backs the input stream up two tokens.
// The zeroth token is already there.
func (p *parser) backup2(t1 Token) {
	p.tokenArr[1] = t1
	p.peekCount = 2
}

func (p *parser) error(format string, args ...interface{}) {
	format = fmt.Sprintf("at line %d, column %d:%s",
		p.current.Line, p.current.Column, format)
	p.err = fmt.Errorf(format, args...)
}

// expect token.
func (p *parser) expect(kind Kind, values ...string) {
	if p.current.Is(kind, values...) {
		p.next()
		return
	}
	p.error("unexpected token %v", p.current)
}

/*
%token <int> INT "integer"
TAG :int, Name : INT, alias :"integer"
%token <int> 'n'
TAG :int, Name : not specify, alias :"n"
%nterm <int> expr
%token <char const *> ID "identifier"
 TAG :char const *, Name : ID, alias :"identifier"
*/
// when Value is 0, need generate
// when Name start with "noname", and alias not , need generate
func (p *parser) parseTokendef() *TokenDef {
	p.next()
	var Tokdef TokenDef
	Tag := ""
	// match <
	if p.current.Is(LeftAngleBracket) {
		// get Tag
		p.next()
		Tag = p.current.Value
		p.next()
		p.expect(RightAngleBracket) // match >
	}

	for {
		if p.current.Is(Identifier) {
			// get Name
			Name := p.current.Value
			id := Idendity{
				Tag:   Tag,
				Name:  Name,
				IDTyp: TERMID,
				Value: 0,
				Alias: "",
			}
			// Get value, 0 is special value
			value := 0
			p.next()
			if p.current.Is(Number) {
				if intVar, err := strconv.Atoi(p.current.Value); err != nil {
					fmt.Println(err)
				} else {
					value = intVar
				}

			} else if p.current.Is(Charater) || p.current.Is(StringKind) { // get alias
				id.Alias = p.current.Value
			} else {
				p.backup()
			}
			id.Value = value
			p.TokenDefMap[id.Name] = true
			Tokdef.IdentifyList = append(Tokdef.IdentifyList, id)
		} else if p.current.Is(Charater) {
			id := Idendity{
				Tag: Tag,
				// noname need do for sepical.
				Name:  genTempName(p.current.Value),
				Value: int(p.current.Value[0]),
				IDTyp: TERMID,
				Alias: p.current.Value,
			}
			p.TokenDefMap[id.Name] = true
			Tokdef.IdentifyList = append(Tokdef.IdentifyList, id)
		} else {
			break
		}
		p.next()
	}

	return &Tokdef
}

// The Same as token,
/*

%left symbols…

%left <type> symbols…
*/
func (p *parser) parsePrecList(Tklist *[]TokenDef) []PrecDef {
	var assocTy PrecAssocType
	var Tokdef TokenDef
	var IdName string
	res := make([]PrecDef, 0)
	if p.current.Is(LeftAssoc) {
		assocTy = LeftAssocType
	} else if p.current.Is(RightAssoc) {
		assocTy = RightAssocype
	} else {
		assocTy = NonAssocType
	}

	Tag := ""
	p.next()
	// match <
	if p.current.Is(LeftAngleBracket) {
		// get Tag
		p.next()
		Tag = p.current.Value
		p.next()
		p.expect(RightAngleBracket) // match >
	}
	p.backup()
	for {
		p.next()
		if p.current.Is(Identifier) || p.current.Is(Charater) {
			// make loop get id or alias
			IdName = p.current.Value
			idvalue := 0
			if p.current.Is(Charater) {
				IdName = genTempName(IdName)
				idvalue = int(p.current.Value[0])
			}
			if !p.TokenDefMap[IdName] {
				id := Idendity{
					Tag: Tag,
					// noname need do for sepical.
					Name:  IdName,
					Value: idvalue,
					IDTyp: TERMID,
					Alias: "",
				}
				p.TokenDefMap[IdName] = true
				Tokdef.IdentifyList = append(Tokdef.IdentifyList, id)
			}
			node := PrecDef{
				IdName:    IdName,
				AssocType: assocTy,
				Prec:      0,
			}
			res = append(res, node)
		} else {
			break
		}
	}
	if len(Tokdef.IdentifyList) != 0 {
		*Tklist = append(*Tklist, Tokdef)
	}
	return res
}

// type
func (p *parser) parseTypeList() []TypeDef {
	p.next()
	var Tag string
	var idName string
	TypedefList := make([]TypeDef, 0)
	// match <
	if p.current.Is(LeftAngleBracket) {
		// get Tag
		p.next()
		Tag = p.current.Value
		p.next()
		p.expect(RightAngleBracket) // match >

	} else {
		p.error("must be has tag")
	}
	for {
		if p.current.Is(Identifier) {
			// pocess the id
			idName = p.current.Value
			TypedefList = append(TypedefList, TypeDef{
				IdName: idName,
				Tag:    Tag,
			})
		} else {
			break
		}
		p.next()
	}
	if len(TypedefList) == 0 {
		p.error("Error , at least has one id")
	}
	return TypedefList
}

// startsymbol indicatation the rules start from
// %start cmds  it means cmds is the start symbols
// or else , the rules start symbol must be `start`
func (p *parser) parseStartSymbol() string {
	p.next()
	if p.current.Is(Identifier) {
		return p.current.Value
	} else {
		p.error("start should follow with identifier")
	}
	return ""
}

func (p *parser) parseDeclare() Node {
	var node Node
	var Unionstr string
	var Codestr string
	var TokDefList []TokenDef
	var PreDefList [][]PrecDef
	var TypeDefList []TypeDef
	var StartSym string = "start"
	p.next()
	for !(p.current.Is(EOF) || p.current.Is(Section)) {
		if p.current.Is(tokenError) {
			p.error("not correct token")
			return nil
		}

		if p.current.Is(UnionDirective) {
			Unionstr = p.current.Value
		}
		if p.current.Is(CodeQuote) {
			Codestr += p.current.Value
		}
		if p.current.Is(TokenDirective) {
			TokDefList = append(TokDefList, *p.parseTokendef())
			continue
		}
		// Prec
		if p.current.Is(LeftAssoc) ||
			p.current.Is(RightAssoc) ||
			p.current.Is(NoneAssoc) ||
			// precDirective is just used to rules
			p.current.Is(Precedence) {
			PreDefList = append(PreDefList, p.parsePrecList(&TokDefList))
			//Do not need call p.next
			continue
		}
		// TYPE
		if p.current.Is(TypeDirective) {
			TypeDefList = append(TypeDefList, p.parseTypeList()...)
			//Do not need call p.next
			continue
		}
		// Start
		if p.current.Is(StartDirective) {
			StartSym = p.parseStartSymbol()
		}
		p.next()
	}

	node = &DeclareNode{
		CodeList:     Codestr,
		Union:        Unionstr,
		TokenDefList: TokDefList,
		TypeDefList:  TypeDefList,
		PrecDefList:  PreDefList,
		StartSym:     StartSym,
	}
	node.SetLocation(p.current.Location)
	return node
}

// parser rules
/*
ID RuleDefine  {id/char %prec terminal-symbol  |ActionQuote}*
	| RuleOR  {id/char |ActionQuote}*
*/
func (p *parser) parseRule(toklst *[]TokenDef) []RuleDef {
	var Leftpart string
	var Tokdef TokenDef
	if p.current.Is(Identifier) {
		Leftpart = p.current.Value
		p.next()
	} else {
		// because it next will p.next, so it must
		p.backup()
		return nil
	}
	p.expect(RuleDefine)
	rightpart := make([]RightSymOrAction, 0)
	res := make([]RuleDef, 0)

	rule := RuleDef{LeftPart: Leftpart, LineNo: p.current.Line}

	for {
		t1 := p.current
		p.next()
		t2 := p.current
		p.backup2(t1) // backup, then need next
		if t1.Is(Identifier) && t2.Is(RuleDefine) {
			res = append(res, rule)
			// get next identify
			p.next()
			break
		}
		p.next()
		switch p.current.Kind {
		case Charater:
			rightpart = append(rightpart, RightSymOrAction{
				ElemType: RightSyType,
				Element:  genTempName(p.current.Value),
			})
			if !p.TokenDefMap[genTempName(p.current.Value)] {
				id := Idendity{
					Tag: "",
					// noname need do for sepical.
					Name:  genTempName(p.current.Value),
					Value: int(p.current.Value[0]),
					IDTyp: TERMID,
					Alias: "",
				}
				p.TokenDefMap[genTempName(p.current.Value)] = true
				Tokdef.IdentifyList = append(Tokdef.IdentifyList, id)
			}
		case Identifier:
			rightpart = append(rightpart, RightSymOrAction{
				ElemType: RightSyType,
				Element:  p.current.Value,
			})
		case ActionQuote:
			rightpart = append(rightpart, RightSymOrAction{
				ElemType: RightActionType,
				Element:  p.current.Value,
			})
		case RuleOR:
			res = append(res, rule)
			rule = RuleDef{LeftPart: Leftpart, LineNo: p.current.Line}
			rightpart = make([]RightSymOrAction, 0)
		case PrecDirective:
			p.next()
			if p.current.Is(Identifier) {
				rule.PrecSym = p.current.Value
			} else if p.current.Is(Charater) {
				rule.PrecSym = genTempName(p.current.Value)
			} else {
				p.error("need symbol or token")
				return nil
			}
		default:
			res = append(res, rule)
			goto out

		}
		rule.RightPart = rightpart
		p.next()
	}
out:
	if len(Tokdef.IdentifyList) != 0 {
		*toklst = append(*toklst, Tokdef)
	}
	return res
}
