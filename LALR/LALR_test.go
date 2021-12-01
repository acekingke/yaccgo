/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package lalr

import (
	"fmt"
	"testing"

	grammar "github.com/acekingke/yaccgo/Grammar"
	item "github.com/acekingke/yaccgo/Items"
	rule "github.com/acekingke/yaccgo/Rules"
	symbol "github.com/acekingke/yaccgo/Symbol"
)

func TestUnion(t *testing.T) {
	type args struct {
		a []int
		b []int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "--empty--", args: args{a: []int{1, 2}, b: []int{2}}},
		{name: "--222--", args: args{a: []int{1, 2}, b: []int{3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Union(tt.args.a, tt.args.b)
			fmt.Println(tt.name)
			fmt.Println(a)
		})
	}
}

var FP_X map[int][]int = map[int][]int{
	1: {100},
	2: {200},
	3: {300},
}

var F_X map[int][]int = map[int][]int{
	1: {},
	2: {},
	3: {},
}

func TestDigraph(t *testing.T) {
	type args struct {
		X  []int
		R  []Relation
		Fp map[int][]int
		F  *map[int][]int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "dig",
			args: args{
				X:  []int{1, 2, 3},
				R:  []Relation{{0, 1, 2}, {1, 2, 3}, {2, 3, 1}},
				Fp: FP_X,
				F:  &F_X,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Digraph(tt.args.X, tt.args.R, tt.args.Fp, tt.args.F)
			fmt.Println("end")
		})
	}
}

func TestLALR1_BuildTrans(t *testing.T) {
	g := grammar.NewGrammar()
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
	c := symbol.NewSymbol(7, "c")

	g.InsertNewSymbol(dollar)
	g.InsertNewSymbol(C)
	g.InsertNewSymbol(A)
	g.InsertNewSymbol(B)
	g.InsertNewSymbol(a)
	g.InsertNewSymbol(b)
	g.InsertNewSymbol(c)

	// 0 start -> C
	g.InsertNewRules(rule.NewProductoinRule(g.StartSymbol, []*symbol.Symbol{
		C, dollar}))
	//1 C → A  B c
	g.InsertNewRules(rule.NewProductoinRule(C, []*symbol.Symbol{A, B}))

	//(2) A ->
	g.InsertNewRules(rule.NewProductoinRule(A, []*symbol.Symbol{}))
	//(3) B -> b
	g.InsertNewRules(rule.NewProductoinRule(B, []*symbol.Symbol{b}))
	// (4) B ->
	g.InsertNewRules(rule.NewProductoinRule(B, []*symbol.Symbol{}))
	// (5) B-> a
	g.InsertNewRules(rule.NewProductoinRule(B, []*symbol.Symbol{a}))
	g.ResolveSymbols()
	g.CalculateEpsilonClosure()

	fmt.Println("................")
	item_var := item.NewItem(0, 0)
	Icloures := item.NewItemCloure()
	Icloures.InsertItem(item_var)
	g.ComputeIClosure(Icloures)
	g.LR0.InsertItemClosure(Icloures, true)
	g.ComputeGotoItemRecursive(Icloures)
	lalr := NewLALR(&g)
	g.Show()
	lalr.BuildTrans()
	fmt.Println("=======show trans=======")
	for k := range lalr.trans {
		fmt.Println(lalr.showTrans(k))
	}
	//fmt.Println(lalr)
	lalr.CalcDR()

	fmt.Println("====DrSet====")
	lalr.ShowDrSet()
	fmt.Println("====ReadsReliatin====")
	ReadRel := lalr.CalcAllReadRelations()
	for _, r := range ReadRel {
		lalr.ShowReads(r)
	}

	fmt.Println("==========")
	lalr.CalcReadSet()
	lalr.ShowReadSet()
	fmt.Println("=====includes=====")
	res := lalr.CaclIncludes()
	for _, r := range res {
		lalr.ShowIncludes(r)
	}
	fmt.Println("====lookbacks======")
	res = lalr.CalcLookbacks()
	for _, r := range res {
		lalr.ShowLookbacks(r)
	}
	fmt.Println("=======follow set=====")
	lalr.CalcFollowSet()
	lalr.ShowFollowSet()
	fmt.Println("=======lookAheadSet==")
	lalr.CalcLookAheadSet()
	lalr.ShowLookAheadSet()
}

func TestLALR1_A(t *testing.T) {
	fmt.Println("========TestLALR1_A==========")
	g := grammar.NewGrammar()
	g.GenStartSymbol()
	/*
		(1) A -> (A)
		(2) A -> a
	*/
	dollar := symbol.NewSymbol(1, "$")
	g.InsertNewSymbol(dollar)

	A := symbol.NewSymbol(2, "A")
	g.InsertNewSymbol(A)

	LP := symbol.NewSymbol(3, "(")
	g.InsertNewSymbol(LP)

	RP := symbol.NewSymbol(4, ")")
	g.InsertNewSymbol(RP)

	a := symbol.NewSymbol(5, "a")
	g.InsertNewSymbol(a)
	// 0 start -> A
	g.InsertNewRules(rule.NewProductoinRule(g.StartSymbol,
		[]*symbol.Symbol{
			A}))
	// 1 A -> ( A )
	g.InsertNewRules(rule.NewProductoinRule(A, []*symbol.Symbol{LP, A, RP}))
	//2 A -> a
	g.InsertNewRules(rule.NewProductoinRule(A, []*symbol.Symbol{a}))
	// symbol
	g.ResolveSymbols()
	g.CalculateEpsilonClosure()
	//calc LR0 closure state
	item_var := item.NewItem(0, 0)
	Icloures := item.NewItemCloure()
	Icloures.InsertItem(item_var)
	g.ComputeIClosure(Icloures)
	g.LR0.InsertItemClosure(Icloures, true)
	g.ComputeGotoItemRecursive(Icloures)
	g.Show()
	// lalr1 look ahead
	lalr := NewLALR(&g)
	lalr.BuildTrans()
	fmt.Println("=======show trans=======")
	for k := range lalr.trans {
		fmt.Println(lalr.showTrans(k))
	}
	lalr.CalcDR()
	fmt.Println("====DrSet====")
	lalr.ShowDrSet()
	fmt.Println("====ReadsReliatin====")
	ReadRel := lalr.CalcAllReadRelations()
	for _, r := range ReadRel {
		lalr.ShowReads(r)
	}

	fmt.Println("==========")
	lalr.CalcReadSet()
	lalr.ShowReadSet()
	fmt.Println("=====includes=====")
	res := lalr.CaclIncludes()
	for _, r := range res {
		lalr.ShowIncludes(r)
	}
	fmt.Println("====lookbacks======")
	res = lalr.CalcLookbacks()
	for _, r := range res {
		lalr.ShowLookbacks(r)
	}
	fmt.Println("=======follow set=====")
	lalr.CalcFollowSet()
	lalr.ShowFollowSet()
	fmt.Println("=======lookAheadSet==")
	lalr.CalcLookAheadSet()
	lalr.ShowLookAheadSet()
	fmt.Println("=======GenTable==")
	if tab, err := lalr.GenTable(); err != nil {
		fmt.Print(err.Error())
	} else {
		for _, row := range tab {
			for _, val := range row {
				fmt.Printf("%d,", val)
			}
			fmt.Println()
		}
	}

}

func TestLALR1_ambiguity(t *testing.T) {
	fmt.Println("========TestLALR1_ambiguity==========")
	g := grammar.NewGrammar()
	g.GenStartSymbol()
	/*
		(1) E-> E + E
		(2) E-> E * E
		(3) E->(E)
		(4) E->i
	*/
	dollar := symbol.NewSymbol(1, "$")
	g.InsertNewSymbol(dollar)

	E := symbol.NewSymbol(2, "E")
	g.InsertNewSymbol(E)

	PLUS := symbol.NewSymbol(3, "+")
	g.InsertNewSymbol(PLUS)

	LP := symbol.NewSymbol(4, "(")
	g.InsertNewSymbol(LP)
	RP := symbol.NewSymbol(5, ")")
	g.InsertNewSymbol(RP)

	MULT := symbol.NewSymbol(6, "*")
	g.InsertNewSymbol(MULT)
	I := symbol.NewSymbol(7, "i")
	g.InsertNewSymbol(I)
	//Set preType and Prec
	PLUS.SetPrecType(symbol.LEFT)
	PLUS.SetPrec(0)
	MULT.SetPrecType(symbol.LEFT)
	MULT.SetPrec(1)
	//produce rules
	// 0 start -> E
	g.InsertNewRules(rule.NewProductoinRule(g.StartSymbol,
		[]*symbol.Symbol{
			E}))
	// 1 E-> E + E
	r1 := rule.NewProductoinRule(E, []*symbol.Symbol{E, PLUS, E})
	r1.SetPrecSymbol(PLUS)
	g.InsertNewRules(r1)

	// 2 E-> E * E
	r2 := rule.NewProductoinRule(E, []*symbol.Symbol{E, MULT, E})
	r2.SetPrecSymbol(MULT)
	g.InsertNewRules(r2)

	// 3 E-> ( E )
	g.InsertNewRules(rule.NewProductoinRule(E, []*symbol.Symbol{LP, E, RP}))
	// 4 E-> i
	g.InsertNewRules(rule.NewProductoinRule(E, []*symbol.Symbol{I}))
	// symbol
	g.ResolveSymbols()
	g.CalculateEpsilonClosure()
	//calc LR0 closure state
	item_var := item.NewItem(0, 0)
	Icloures := item.NewItemCloure()
	Icloures.InsertItem(item_var)
	g.ComputeIClosure(Icloures)
	g.LR0.InsertItemClosure(Icloures, true)
	g.ComputeGotoItemRecursive(Icloures)
	g.Show()
	// lalr1 look ahead
	lalr := NewLALR(&g)
	lalr.BuildTrans()
	fmt.Println("=======show trans=======")
	for k := range lalr.trans {
		fmt.Println(lalr.showTrans(k))
	}
	lalr.CalcDR()
	fmt.Println("====DrSet====")
	lalr.ShowDrSet()
	fmt.Println("====ReadsReliatin====")
	ReadRel := lalr.CalcAllReadRelations()
	for _, r := range ReadRel {
		lalr.ShowReads(r)
	}

	fmt.Println("==========")
	lalr.CalcReadSet()
	lalr.ShowReadSet()
	fmt.Println("=====includes=====")
	res := lalr.CaclIncludes()
	for _, r := range res {
		lalr.ShowIncludes(r)
	}
	fmt.Println("====lookbacks======")
	res = lalr.CalcLookbacks()
	for _, r := range res {
		lalr.ShowLookbacks(r)
	}
	fmt.Println("=======follow set=====")
	lalr.CalcFollowSet()
	lalr.ShowFollowSet()
	fmt.Println("=======lookAheadSet==")
	lalr.CalcLookAheadSet()
	lalr.ShowLookAheadSet()
	fmt.Println("=======GenTable==")
	if tab, err := lalr.GenTable(); err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Print("/*     ")
		for _, sy := range lalr.G.Symbols {
			fmt.Printf("%s\t", sy.Name)
		}
		fmt.Println("*/")
		for index, row := range tab {
			fmt.Printf("/* %d */ {", index)
			for _, val := range row {
				fmt.Printf("%d,\t", val)
			}
			fmt.Println("},")
		}
	}

}
