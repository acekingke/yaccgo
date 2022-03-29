/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package parser

/* inspired by Rob Pikes' video Lexical Scanning
In Go and golang's 'template' package.
*/
import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	tokenError Kind = "Error"
	Identifier Kind = "Identifier"
	Number     Kind = "Number"
	Section    Kind = "Section"

	CodeQuote   Kind = "CodeQuote"
	ActionQuote Kind = "ActionQuote"
	EOF         Kind = "EOF"

	TypeDirective  Kind = "TypeDirective"
	TokenDirective Kind = "TokenDirective"
	UnionDirective Kind = "UnionDirective"
	LeftAssoc      Kind = "LeftAssoc"
	RightAssoc     Kind = "RightAssoc"
	NoneAssoc      Kind = "NoneAssoc"
	PrecDirective  Kind = "PrecDirective"
	Precedence     Kind = "Precedence"
	StartDirective Kind = "StartDirective"

	ActionSelf   Kind = "ActionSelf"
	ActionN      Kind = "ActionN"
	ActionAccept Kind = "ActionAccept"
	ActionEnd    Kind = "ActionEnd"

	RuleOR     Kind = "RuleOR"     // |
	RuleDefine Kind = "RuleDefine" // :

	LeftAngleBracket  Kind = "LeftAngleBracket"  // Aangle brackets <
	RightAngleBracket Kind = "RightAngleBracket" // RightAngleBracket >

	Charater   Kind = "Charater"
	StringKind Kind = "StringKind"
)

type stateFn func(*lexer) stateFn
type Kind string
type Token struct {
	Location
	Kind  Kind
	Value string
	EndAt int
}

type Location struct {
	Line   int // The 1-based line of the location.
	Column int // The 0-based column number of the location.
}

type lexer struct {
	input      string
	tokens     chan Token // chan for tokens
	start, end int        // current position in input
	width      int        // last rune width
	startLoc   Location   // start location
	prev, loc  Location   // prev location of end location, end location

}

func rootState(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.end:], "//") || strings.HasPrefix(l.input[l.end:], "/*") {
		return CommentState
	}

	switch r := l.next(); {
	case r == eof:
		l.emitEOF()
		return nil
	case r == '%':
		return DirectiveState
	case r == '$':
		return ActionState
	case r == '|':
		l.emit(RuleOR)
	case r == ':':
		l.emit(RuleDefine)
	case r == ' ', r == '\t', r == '\n': //skip the space
		l.ignore()
	// case  alpha , identify
	case r == '\'':
		return charaterState
	case r == '"':
		return stringKindState
	case unicode.IsLetter(r), r == '_':
		return IdentifyState
	case r == '<':
		l.emit(LeftAngleBracket)
	case r == '>':
		l.emit(RightAngleBracket)
	case unicode.IsDigit(r):
		l.backup()
		l.acceptRun("0123456789")
		l.emit(Number)
	case r == '-':
		l.acceptRun("0123456789")
		l.emit(Number)
	case r == '{':
		return ActionQuoteState

	default:
		l.error("not correct lexer")
		return nil
	}
	return rootState
}

func CommentState(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.end:], "//") {
		for {
			r := l.next()
			l.ignore()
			if r == '\n' || r == eof {
				break
			}
		}
	} else {
		//start with /*
		for {
			r := l.next()
			if r == '*' {
				r = l.next()
				if r == '/' {
					l.ignore()
					break
				}
			}
			if r == eof {
				l.error("comment do not has */")
			}
			l.ignore()
		}
	}
	return rootState
}

// { xxx }
func ActionQuoteState(l *lexer) stateFn {
	depth := 1

	for {
		r := l.next()
		if r == '{' {
			depth++
		}
		if r == '}' {
			depth--
		}
		if depth == 0 {
			break
		}
		if r == eof {
			l.error("`{` and`}` not match")
			return nil
		}
	}
	l.emit(ActionQuote)
	return rootState
}

// 'a' '\'' and other ,'ab' is error
func charaterState(l *lexer) stateFn {
	value := ""
	if r := l.next(); r != '\\' {
		value += string(r)
		r = l.next()
		if r != '\'' {
			l.error("just can quote single char")
			return nil
		}
		l.emitValue(Charater, value)
	} else {
		//transalte symbol
		r = l.next()
		switch r {
		case '\'':
			value += "'"
			l.emitValue(Charater, value)
		default:
			l.error("not correct translate")
			return nil
		}
	}
	return rootState
}

// " "  "aaa\""
func stringKindState(l *lexer) stateFn {
	value := ""
	r := l.next()
	for r != '"' {

		if r == '\\' && l.next() == '"' {
			value += string('"')
		} else if r != eof {
			value += string(r)

		} else {
			l.error("do not has \" end")
			return nil
		}
		r = l.next()
	}
	l.emitValue(StringKind, value)
	return rootState
}

func IdentifyState(l *lexer) stateFn {
	//alpha, number, underscore
	r := l.next()
	for ; unicode.IsLetter(r) ||
		unicode.IsDigit(r) ||
		r == '_'; r = l.next() {
	}
	l.backup()
	l.emit(Identifier)
	return rootState
}

func ActionState(l *lexer) stateFn {
	switch r := l.next(); {
	case r == '$': //$$
		l.emit(ActionSelf)
	case unicode.IsDigit(r):
		l.acceptRun("0123456789")
		l.emit(ActionN)
	case l.acceptWord("accept"): //$accept
		l.emit(ActionAccept)
	case l.acceptWord("end"): //$end
		l.emit(ActionEnd)
	default:
		l.error("Action lexer error")
	}

	return rootState
}

func DirectiveState(l *lexer) stateFn {
	switch l.next() {
	case '%':
		l.emit(Section)
	case '{':
		return CodeQuoteBegin
	default:
		l.backup()
		return DirectiveOtherState
	}

	return rootState
}

func DirectiveOtherState(l *lexer) stateFn {
	if l.acceptWord("type") {
		l.emit(TypeDirective)
	}
	if l.acceptWord("token") {
		l.emit(TokenDirective)
	}
	if l.acceptWord("union") {
		return DirectiveUnionState
	}
	if l.acceptWord("left") {
		l.emit(LeftAssoc)
	}
	if l.acceptWord("right") {
		l.emit(RightAssoc)
	}
	if l.acceptWord("prec") {
		l.emit(PrecDirective)
	}
	if l.acceptWord("precedence") {
		l.emit(Precedence)
	}
	if l.acceptWord("start") {
		l.emit(StartDirective)
	}
	return rootState
}

func CodeQuoteBegin(l *lexer) stateFn {
	vstart := l.end
	for {
		// Skip spaces (U+0020) if any
		r := l.peek()
		for ; r == '\t' || r == '\n' || r == ' '; r = l.peek() {
			l.next()
		}
		if l.acceptWord("%}") {
			vend := l.end - 2
			l.emitValue(CodeQuote, l.input[vstart:vend])
			break
		}
		if r := l.next(); r == eof {
			l.error("not correct code quote")
			return nil
		}
	}
	return rootState
}

/* %uinon {

}
*/
func DirectiveUnionState(l *lexer) stateFn {
	//skip space
	for {
		r := l.next()
		if r != ' ' && r != '\t' {
			break
		}
	}
	l.backup()
	level := 0
	if !l.acceptWord("{") {
		l.error("union directive need { to start")
		return nil
	}
	vstart := l.end
	level++
Loop:
	for {
		switch l.next() {
		case '{':
			level++
		case '}':
			level--
			if level == 0 {
				break Loop
			}
		case eof:
			l.error("not complete in union Directive")
			return nil
		}
	}
	vend := l.end - 1 //skip }
	l.emitValue(UnionDirective, l.input[vstart:vend])
	return rootState
}

func Lex(source string) *lexer {
	l := &lexer{
		input:  source,
		tokens: make(chan Token),
	}

	l.loc = Location{Line: 1, Column: 0}
	l.prev = l.loc
	l.startLoc = l.loc

	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for state := rootState; state != nil; {
		state = state(l)
	}

	close(l.tokens)
}

const eof rune = -1

func (l *lexer) next() rune {
	if l.end >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.end:])
	l.width = w
	l.end += w

	l.prev = l.loc
	if r == '\n' {
		l.loc.Line++
		l.loc.Column = 0
	} else {
		l.loc.Column++
	}

	return r
}

func (l *lexer) nextToken() Token {
	return <-l.tokens
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.end -= l.width
	l.loc = l.prev
}

func (l *lexer) emit(t Kind) {
	l.emitValue(t, l.word())
}

func (l *lexer) emitValue(t Kind, value string) {
	l.tokens <- Token{
		Location: l.startLoc,
		Kind:     t,
		Value:    value,
		EndAt:    l.end,
	}
	l.start = l.end
	l.startLoc = l.loc
}

func (l *lexer) emitEOF() {
	l.tokens <- Token{
		Location: l.prev, // Point to previous position for better error messages.
		Kind:     EOF,
	}
	l.start = l.end
	l.startLoc = l.loc
}

func (l *lexer) word() string {
	return l.input[l.start:l.end]
}

func (l *lexer) ignore() {
	l.start = l.end
	l.startLoc = l.loc
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) acceptWord(word string) bool {
	pos, loc, prev := l.end, l.loc, l.prev

	// Skip spaces (U+0020) if any
	r := l.peek()
	for ; r == ' '; r = l.peek() {
		l.next()
	}

	for _, ch := range word {
		if l.next() != ch {
			l.end, l.loc, l.prev = pos, loc, prev
			return false
		}
	}
	if r = l.peek(); r != ' ' && r != '\t' && r != '\n' && r != eof {
		l.end, l.loc, l.prev = pos, loc, prev
		return false
	}

	return true
}

func (l *lexer) error(format string, args ...interface{}) stateFn {
	err := fmt.Sprintf("at line %d , pos %d:", l.loc.Line, l.loc.Column+1)

	l.tokens <- Token{
		Location: l.loc,
		Kind:     tokenError,
		Value:    err + fmt.Sprintf(format, args...),
	}
	return nil
}

func (l Location) Empty() bool {
	return l.Column == 0 && l.Line == 0
}

func (t Token) Is(kind Kind, values ...string) bool {
	if len(values) == 0 {
		return kind == t.Kind
	}

	for _, v := range values {
		if v == t.Value {
			goto found
		}
	}
	return false

found:
	return kind == t.Kind
}
