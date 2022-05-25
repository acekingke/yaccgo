/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package lalr

import (
	"fmt"

	symbol "github.com/acekingke/yaccgo/Symbol"
)

func (lalr *LALR1) fetchSymbol(index int) *symbol.Symbol {
	return lalr.G.Symbols[index]
}

func (lalr *LALR1) isNonSymIndex(in uint) bool {
	if in&(1<<32) != 0 {
		return false
	}
	return lalr.fetchSymbol(int(in)).IsNonTerminator
}

func (lalr *LALR1) isNonAndEpsilonSymIndex(in uint) bool {
	if in&(CheckMask) != 0 {
		return false
	}
	sy := lalr.fetchSymbol(int(in))
	return sy.IsNonTerminator && sy.IsEpsilonClosure
}

func (lalr *LALR1) isTermSymIndex(in uint) bool {
	if in&(CheckMask) != 0 {
		return false
	}
	return !lalr.fetchSymbol(int(in)).IsNonTerminator
}

//show Dr set
func (lalr *LALR1) ShowDrSet() {
	for trIndex, set := range lalr.DRSet {
		q := lalr.trans[trIndex].q
		symName := lalr.fetchSymbol(int(lalr.trans[trIndex].sym_or_rule)).Name
		var str_set string = "["
		for _, v := range set {
			str_set += lalr.fetchSymbol(v).Name + " "
		}
		str_set += " ]"
		fmt.Printf("%d--%s--> %s\n", q, symName, str_set)
	}
}

func (lalr *LALR1) ShowReadSet() {
	for trIndex, set := range lalr.ReadSet {
		q := lalr.trans[trIndex].q
		symName := lalr.fetchSymbol(int(lalr.trans[trIndex].sym_or_rule)).Name
		var str_set string = "["
		for _, v := range set {
			str_set += lalr.fetchSymbol(v).Name + " "
		}
		str_set += " ]"
		fmt.Printf("%d--%s--> %s\n", q, symName, str_set)
	}
}

func (lalr *LALR1) ShowFollowSet() {
	var str_set []string
	for trIndex, set := range lalr.FollowSet {
		str_set = nil
		q := lalr.trans[trIndex].q
		symName := lalr.fetchSymbol(int(lalr.trans[trIndex].sym_or_rule)).Name
		for _, v := range set {
			str_set = append(str_set, lalr.fetchSymbol(v).Name)
		}
		fmt.Printf("%d--%s--> %v\n", q, symName, str_set)
	}
}

// return state number.
func (lalr *LALR1) fechStateNumber(rIndex int) []int {
	ret := []int{}
	for num, s := range lalr.G.LR0.LR0Closure {
		for _, it := range s.Items {
			if it.RuleIndex == rIndex {
				ret = append(ret, num)
			}
		}
	}
	return ret
}

// return  trans index
func (lalr *LALR1) fetchTransIndex(state, sym int) (int, error) {
	for index, tr := range lalr.trans {
		if tr.q == state && sym == int(tr.sym_or_rule) {
			return index, nil
		}
	}
	return MaxInt, fmt.Errorf("not found")
}

func (lalr *LALR1) seqenceCanEpsilon(slice []*symbol.Symbol) bool {
	ret := true
	for _, sy := range slice {
		if !sy.IsEpsilonClosure {
			ret = false
			break
		}
	}
	return ret
}

//func fetch A --> omega
func (lalr *LALR1) fetchReduceTransistor() []Transistor {
	res := []Transistor{}
	for _, tr := range lalr.trans {
		if tr.sym_or_rule&CheckMask != 0 {
			res = append(res, tr)
		}
	}
	return res
}

// Show the trans
func (lalr *LALR1) showTrans(tr int) string {
	trIt := lalr.trans[tr]
	var s string = fmt.Sprintf("%d:", trIt.q)

	if trIt.sym_or_rule&CheckMask != 0 {
		ruleIndex := trIt.sym_or_rule & Mask
		r := lalr.G.ProductoinRules[ruleIndex]
		s += fmt.Sprintf("%s-->", r.LeftPart.Name)
		for _, sy := range r.RighPart {
			s += fmt.Sprintf(" %s ", sy.Name)
		}
	} else {
		s += lalr.fetchSymbol(int(trIt.sym_or_rule)).Name
	}
	return s
}

func (lalr *LALR1) ShowReads(R Relation) {
	tr1 := R.x
	tr2 := R.y
	fmt.Printf("%s....Reads....%s\n", lalr.showTrans(tr1), lalr.showTrans(tr2))
}

// Show the lookbacks
func (lalr *LALR1) ShowLookbacks(R Relation) {
	tr1 := R.x
	tr2 := R.y
	fmt.Printf("%s....lookbacks....%s\n", lalr.showTrans(tr1), lalr.showTrans(tr2))
}

// Show the includes
func (lalr *LALR1) ShowIncludes(R Relation) {
	tr1 := R.x
	tr2 := R.y
	fmt.Printf("%s....includes....%s\n", lalr.showTrans(tr1), lalr.showTrans(tr2))
}

// Show the LookAheadSet
func (lalr *LALR1) ShowLookAheadSet() {
	var str_set string = ""
	for trId, set := range lalr.LookAheadSet {
		str_set = ""
		for _, sId := range set {
			str_set += " "
			str_set += lalr.fetchSymbol(sId).Name
		}
		fmt.Printf("%s : %s\n", lalr.showTrans(trId), str_set)
	}
}

// Get the error action code
// Max State + 100 is the error
func (lalr *LALR1) GenErrorCode() int {
	return len(lalr.G.LR0.LR0Closure) + 100
}

// Get the accept action code
func (lalr *LALR1) GenAcceptCode() int {
	return len(lalr.G.LR0.LR0Closure) + 200
}
