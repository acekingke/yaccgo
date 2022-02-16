/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package lalr

import (
	"fmt"
	"sort"

	grammar "github.com/acekingke/yaccgo/Grammar"

	utils "github.com/acekingke/yaccgo/Utils"
)

type LALR1 struct {
	G      *grammar.Grammar
	trans  []Transistor
	GTable [][]int
	DRSet, // trans index --> DrSet map
	ReadSet, // trans index --> DrSet map
	FollowSet,
	LookAheadSet map[int][]int
	NumOfStates int //States number
	// The last is for pack table
	NeedPacked  bool
	ActionTable []int
	OffsetTable []int
	CheckTable  []int
}

// q --t--> p
type Transistor struct {
	Index int
	q     int // state index

	// sym index, or reduce rule index
	// if it's rule, hight bit is 1
	sym_or_rule uint

	//state index to, if it a rule, to will be maxInt
	to int
}

// NewLALR
func NewLALR(g *grammar.Grammar) *LALR1 {
	return &LALR1{
		G:            g,
		trans:        make([]Transistor, 0),
		DRSet:        make(map[int][]int),
		ReadSet:      make(map[int][]int),
		FollowSet:    make(map[int][]int),
		LookAheadSet: make(map[int][]int),
	}
}

func (lalr *LALR1) BuildTrans() {
	g := lalr.G
	lr0 := g.LR0
	states := lr0.LR0Closure
	//build the symbol trans
	//(q, A)
	for _, iC := range states {
		q := iC.Index
		for _, gt := range iC.GoTo {
			t := gt.Sym.ID
			p := gt.ItemCl
			lalr.trans = append(lalr.trans,
				Transistor{
					q:           q,
					sym_or_rule: t,
					to:          p,
				})
		}
	}
	//build the reduce symbol
	//(q, A-->w)
	for _, iC := range states {
		for _, it := range iC.Items {
			r := g.ProductoinRules[it.RuleIndex]
			if it.Dot == len(r.RighPart) {
				q := iC.Index
				t := uint(it.RuleIndex)
				t = t | (1 << 32)
				//p := gt.ItemCl
				lalr.trans = append(lalr.trans,
					Transistor{
						q:           q,
						sym_or_rule: t,
						to:          MaxInt,
					})
			}
		}
	}
	//sort
	sort.SliceStable(lalr.trans, func(i, j int) bool {
		return lalr.trans[i].q < lalr.trans[j].q
	})
	//fill index
	for index := range lalr.trans {
		lalr.trans[index].Index = index
	}
}

func (lalr *LALR1) fetchOneDr(tr Transistor) []int {
	sy := lalr.fetchSymbol(int(tr.sym_or_rule))
	var res []int = nil
	if sy.IsNonTerminator {
		nextState := tr.to
		//traverse and fetch terminator t
		for _, tr := range lalr.trans {
			if tr.q == nextState &&
				lalr.isTermSymIndex(tr.sym_or_rule) {
				res = append(res, int(tr.sym_or_rule))
			}
		}
	}
	return res
}

//calc reads Relation
func (lalr *LALR1) calcReadsRelation(transIndex int) []Relation {
	tr := lalr.trans[transIndex]
	var res []Relation = nil
	// traverse
	nextState := tr.to
	//traverse and fetch terminator t
	for _, tr := range lalr.trans {
		if tr.q == nextState &&
			lalr.isNonAndEpsilonSymIndex(tr.sym_or_rule) {
			//index need use
			res = append(res, Relation{x: transIndex, y: tr.Index})
		}
	}
	return res
}

//func calc all reads relations
func (lalr *LALR1) CalcAllReadRelations() []Relation {
	var res []Relation = nil
	for key := range lalr.DRSet {
		tmp := lalr.calcReadsRelation(key)
		res = append(res, tmp...)
	}
	return res
}

//calculate the Reads Set
func (lalr *LALR1) CalcReadSet() {
	X := []int{}
	for key := range lalr.DRSet {
		X = append(X, key)
		lalr.ReadSet[key] = []int{}
	}
	R := lalr.CalcAllReadRelations()
	Digraph(X, R, lalr.DRSet, &lalr.ReadSet)
}

//calc the DR
func (lalr *LALR1) CalcDR() {

	for _, tr := range lalr.trans {
		Index := tr.sym_or_rule
		if lalr.isNonSymIndex(Index) {
			lalr.DRSet[tr.Index] = lalr.fetchOneDr(tr)
		}
	}
	//set the start -> S' Dr is $ sym
	// 0 index transistor is I0--S'-->
	// 1 is dollor symbol
	lalr.DRSet[0] = append(lalr.DRSet[0], []int{1}...)
}

//-----------------Caculate the Includes relation-------------
// B-> beta A gama if gama can episilon, p' --- beta --> p
// (p, A) includes (p', B)
//
//calculate the include Relation.
func (lalr *LALR1) CaclIncludeRelation(tr int) []Relation {
	res := []Relation{}
	sy := lalr.G.Symbols[lalr.trans[tr].sym_or_rule]
	for index, r := range lalr.G.ProductoinRules {
		LeftSy := r.LeftPart
		for Dot, sycheck := range r.RighPart {
			if sy == sycheck && lalr.seqenceCanEpsilon(r.RighPart[Dot+1:]) {
				for _, q := range lalr.fechStateNumber(index) {
					if to_index, err := lalr.fetchTransIndex(q, int(LeftSy.ID)); err == nil {
						res = append(res, Relation{x: tr, y: to_index})
					}
				}

			}
		}
	}
	return res
}

//calc all the relation
func (lalr *LALR1) CaclIncludes() []Relation {
	var res []Relation = nil
	for key := range lalr.DRSet {
		tmp := lalr.CaclIncludeRelation(key)
		res = append(res, tmp...)
	}
	return res
}

//calculate lookback relation
func (lalr *LALR1) CalcLookbacks() []Relation {
	var res []Relation = nil

	for _, tr := range lalr.fetchReduceTransistor() {
		trIndex := tr.Index
		ruleIndex := tr.sym_or_rule & Mask
		leftPart := lalr.G.ProductoinRules[ruleIndex].LeftPart
		for tr_2 := range lalr.DRSet {
			SyIndex := lalr.trans[tr_2].sym_or_rule
			if SyIndex == leftPart.ID {
				// trIndex lookback tr2
				res = append(res, Relation{x: trIndex, y: tr_2})
			}
		}
	}
	return res
}

// calc the follow set
func (lalr *LALR1) CalcFollowSet() {
	X := []int{}
	for key := range lalr.ReadSet {
		X = append(X, key)
		lalr.FollowSet[key] = []int{}
	}
	R := lalr.CaclIncludes()
	Digraph(X, R, lalr.ReadSet, &lalr.FollowSet)
}

//calc the look ahead Set
func (lalr *LALR1) CalcLookAheadSet() {
	X := []int{}
	Set := make(map[int][]int)
	for _, tr := range lalr.fetchReduceTransistor() {
		X = append(X, tr.Index)
		Set[tr.Index] = []int{}
	}
	R := lalr.CalcLookbacks()
	Digraph(X, R, lalr.FollowSet, &Set)
	for _, tr := range lalr.fetchReduceTransistor() {
		if tr.sym_or_rule&Mask == 0 { //is the start->C. defalut is dollar
			lalr.LookAheadSet[tr.Index] = []int{1}
		} else {
			lalr.LookAheadSet[tr.Index] = Set[tr.Index]
		}
	}
}

func ComputeLALR(g *grammar.Grammar) *LALR1 {
	lalr := NewLALR(g)
	lalr.BuildTrans()
	lalr.CalcDR()
	lalr.CalcReadSet()
	lalr.CalcFollowSet()
	lalr.CalcLookAheadSet()
	if utils.DebugFlags {
		fmt.Println("=========Show State Closure=========")
		lalr.G.Show()
		fmt.Println("===========SHOW TRANS================")
		for k := range lalr.trans {
			fmt.Println(lalr.showTrans(k))
		}
		fmt.Println("==========Show Direct Read SET===============")
		lalr.ShowDrSet()
		fmt.Println("==========Show Reads SET===============")
		lalr.ShowReadSet()
		fmt.Println("==========Show FollowSet SET===============")
		lalr.ShowFollowSet()
		fmt.Println("==========Show LookAhead SET===============")
		lalr.ShowLookAheadSet()
	}

	if tab, err := lalr.GenTable(); err != nil {
		panic(err.Error())
	} else {
		lalr.GTable = tab
		// try to pack the table
		act, off, check := utils.PackTable(lalr.GTable)
		lalr.NeedPacked = false
		lalr.NumOfStates = len(lalr.G.LR0.LR0Closure)
		if len(act)+len(off)+len(check) > lalr.NumOfStates*len(lalr.G.Symbols) {
			if utils.DebugFlags {
				fmt.Println("The table is no need to pack")
			}
		} else {
			if utils.DebugFlags {
				fmt.Println("The table is packed")
			}
			lalr.NeedPacked = true
			lalr.ActionTable, lalr.OffsetTable, lalr.CheckTable = act, off, check
		}

	}
	return lalr
}
