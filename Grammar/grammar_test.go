/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package grammar

import (
	"fmt"
	"testing"

	item "github.com/acekingke/yaccgo/Items"
	rule "github.com/acekingke/yaccgo/Rules"
	symbol "github.com/acekingke/yaccgo/Symbol"
)

func TestNewGrammar(t *testing.T) {
	g := NewGrammar()
	if g.StartSymbol != nil {
		t.Errorf("value name %v ", g)
	}

}
func TestInsertNewSymbol(t *testing.T) {
	g := NewGrammar()
	fmt.Println("jel;l;p")
	g.InsertNewSymbol(symbol.NewSymbol(2, "A"))
	for r := range g.Symbols {
		fmt.Println(r)
	}
}

func TestGrammar_GenStartSymbol(t *testing.T) {
	g := NewGrammar()
	g.GenStartSymbol()
	/*
		(1) E → E * B
		(2) E → E + B
		(3) E → B
		(4) B → NUMBER
	*/
	dollar := symbol.NewSymbol(1, "$")
	E := symbol.NewSymbol(2, "E")
	PLUS := symbol.NewSymbol(3, "PLUS")
	B := symbol.NewSymbol(4, "B")
	NUMBER := symbol.NewSymbol(5, "NUMBER")
	MULTI := symbol.NewSymbol(6, "*")
	g.InsertNewSymbol(dollar)
	g.InsertNewSymbol(E)
	g.InsertNewSymbol(PLUS)
	g.InsertNewSymbol(B)
	g.InsertNewSymbol(NUMBER)
	g.InsertNewSymbol(MULTI)

	// 0 start -> E
	g.InsertNewRules(rule.NewProductoinRule(g.StartSymbol, []*symbol.Symbol{
		E}))
	//1 E → E * B
	g.InsertNewRules(rule.NewProductoinRule(E, []*symbol.Symbol{E, MULTI, B}))
	//(2) E → E + B
	g.InsertNewRules(rule.NewProductoinRule(E, []*symbol.Symbol{E, PLUS, B}))
	//(3) E → B
	g.InsertNewRules(rule.NewProductoinRule(E, []*symbol.Symbol{B}))
	//(4) B -> NUMBER
	g.InsertNewRules(rule.NewProductoinRule(B, []*symbol.Symbol{NUMBER}))
	g.ResolveSymbols()
	for sy := range g.VnSet {
		sy.Show()
	}
	fmt.Println("................")
	item_var := item.NewItem(0, 0)
	fmt.Println(g.ItemToStr(item_var))
	Icloures := item.NewItemCloure()
	Icloures.InsertItem(item_var)
	g.ComputeIClosure(Icloures)
	g.LR0.InsertItemClosure(Icloures, true)
	g.ComputeGotoItemRecursive(Icloures)
	g.Show()
}
func TestGrammar_ComputeGotoItemRecursive(t *testing.T) {
	g := NewGrammar()
	g.GenStartSymbol()
	/*
		(1) E → E * B
		(2) E → E + B
		(3) E → B
		(4) B → NUMBER
	*/
	dollar := symbol.NewSymbol(1, "$")
	E := symbol.NewSymbol(2, "E")
	PLUS := symbol.NewSymbol(3, "PLUS")
	B := symbol.NewSymbol(4, "B")
	ONE := symbol.NewSymbol(5, "1")
	ZERO := symbol.NewSymbol(6, "0")
	MULTI := symbol.NewSymbol(7, "*")
	g.InsertNewSymbol(dollar)
	g.InsertNewSymbol(E)
	g.InsertNewSymbol(PLUS)
	g.InsertNewSymbol(B)
	g.InsertNewSymbol(ONE)
	g.InsertNewSymbol(ZERO)
	g.InsertNewSymbol(MULTI)

	// 0 start -> E
	g.InsertNewRules(rule.NewProductoinRule(g.StartSymbol, []*symbol.Symbol{
		E}))
	//1 E → E * B
	g.InsertNewRules(rule.NewProductoinRule(E, []*symbol.Symbol{E, MULTI, B}))
	//(2) E → E + B
	g.InsertNewRules(rule.NewProductoinRule(E, []*symbol.Symbol{E, PLUS, B}))
	//(3) E → B
	g.InsertNewRules(rule.NewProductoinRule(E, []*symbol.Symbol{B}))
	//(4) B -> 1
	g.InsertNewRules(rule.NewProductoinRule(B, []*symbol.Symbol{ONE}))
	// (5) B -> 0
	g.InsertNewRules(rule.NewProductoinRule(B, []*symbol.Symbol{ZERO}))
	g.ResolveSymbols()
	g.CalculateEpsilonClosure()
	for sy := range g.VnSet {
		sy.Show()
	}
	fmt.Println("................")
	item_var := item.NewItem(0, 0)
	Icloures := item.NewItemCloure()
	Icloures.InsertItem(item_var)
	g.ComputeIClosure(Icloures)
	g.LR0.InsertItemClosure(Icloures, true)
	g.ComputeGotoItemRecursive(Icloures)
	g.Show()
}

func TestGrammar_CalculateEpsilonClosure(t *testing.T) {
	g := NewGrammar()
	g.GenStartSymbol()
	/*
		(1) C → A B
		(2) A ->
		(3) B -> b
		(4) B->
	*/
	dollar := symbol.NewSymbol(1, "$")
	C := symbol.NewSymbol(2, "C")
	A := symbol.NewSymbol(3, "A")
	B := symbol.NewSymbol(4, "B")
	a := symbol.NewSymbol(5, "a")
	b := symbol.NewSymbol(6, "b")

	g.InsertNewSymbol(dollar)
	g.InsertNewSymbol(C)
	g.InsertNewSymbol(A)
	g.InsertNewSymbol(B)
	g.InsertNewSymbol(a)
	g.InsertNewSymbol(b)

	// 0 start -> C
	g.InsertNewRules(rule.NewProductoinRule(g.StartSymbol, []*symbol.Symbol{
		C}))
	//1 C → A  B
	g.InsertNewRules(rule.NewProductoinRule(C, []*symbol.Symbol{A, B}))

	//(2) A ->
	g.InsertNewRules(rule.NewProductoinRule(A, []*symbol.Symbol{}))
	//(3) B -> b
	g.InsertNewRules(rule.NewProductoinRule(B, []*symbol.Symbol{b}))
	// (4) B ->
	g.InsertNewRules(rule.NewProductoinRule(B, []*symbol.Symbol{}))
	g.ResolveSymbols()
	g.CalculateEpsilonClosure()
	fmt.Println("................")
	if !A.IsEpsilonClosure {
		t.Errorf("%s should be nullable", A.Name)
	}
	if !C.IsEpsilonClosure {
		t.Errorf("%s should be nullable", C.Name)
	}
	if !B.IsEpsilonClosure {
		t.Errorf("%s should be nullable", B.Name)
	}

}
