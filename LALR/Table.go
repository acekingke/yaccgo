/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package lalr

import (
	"fmt"

	symbol "github.com/acekingke/yaccgo/Symbol"
	utils "github.com/acekingke/yaccgo/Utils"
)

type E_ActionType int

const (
	SHIFT E_ActionType = iota
	REDUCE
	ERROR
	SHIFT_REDUCE_CONFLICT
	REDUCE_REDUCE_CONFLICT
)

type Action struct {
	Sym         *symbol.Symbol
	ActionType  E_ActionType
	ActionIndex int
	PrecType    symbol.E_Precedence
	Prec        int
}

func (lalr *LALR1) ShowAndCheckConflict(state int, tranlist []Transistor) {
	fmt.Printf("=====%d==\n", state)
	for _, tr := range tranlist {
		if tr.sym_or_rule&CheckMask != 0 {
			//reduce
			fmt.Println("")
			for _, sy := range lalr.LookAheadSet[tr.Index] {
				fmt.Printf("%s ", lalr.G.Symbols[sy].Name)
			}
			fmt.Printf(":reduce by %s", lalr.showTrans(tr.Index))
		} else {
			fmt.Printf("\n%s shift to %d ", lalr.G.Symbols[tr.sym_or_rule].Name, tr.to)
		}
	}
	fmt.Println("\n========")
}

//  Resolve the conflict by Prec
func (lalr *LALR1) ResolveConflict(act01, act02 *Action) (*Action, error) {
	act_first := act01
	act_second := act02
	// Reduce first, Shift Second
	if act02.ActionType == REDUCE && act01.ActionType == SHIFT {
		act_first = act02
		act_second = act01
	}
	if act_first.Prec == -1 || act_second.Prec == -1 {
		return nil, fmt.Errorf("cannot resolve conflict")
	}
	if act_first.Prec > act_second.Prec {
		return act_first, nil
	} else if act_first.Prec == act_second.Prec {
		if act_first.PrecType == symbol.NONE || act_second.PrecType == symbol.NONE {
			actionError := &Action{
				Sym:         act_first.Sym,
				ActionType:  ERROR,
				ActionIndex: 0,
				PrecType:    symbol.NONE,
				Prec:        act_first.Prec,
			}
			return actionError, nil
		}
		if act_first.PrecType == symbol.LEFT {
			return act_first, nil
		}
		if act_first.PrecType == symbol.RIGHT {
			return act_second, nil
		}
	} else {
		return act_second, nil
	}
	return nil, fmt.Errorf("other error")
}

// check weath it has shift/reduce conflict or reduce/reduce conflict
func (lalr *LALR1) CheckAndResolveConflict(state int, tranlist []Transistor) (map[int][]*Action, error) {
	// key is symbol, slice is action list, minus is the reduce, postive is shift
	action_set := make(map[int][]*Action)
	for _, tr := range tranlist {
		if tr.sym_or_rule&CheckMask != 0 {
			// Reduce
			for _, sy := range lalr.LookAheadSet[tr.Index] {
				r := lalr.G.ProductoinRules[int(tr.sym_or_rule&Mask)]
				PreTy := symbol.NONE
				Pre := -1
				if r.PrecSymbol != nil {
					PreTy = r.PrecSymbol.PrecType
					Pre = r.PrecSymbol.Prec
				}
				action_set[sy] = append(action_set[sy], &Action{
					Sym:         lalr.G.Symbols[sy],
					ActionType:  REDUCE,
					ActionIndex: -int(tr.sym_or_rule & Mask),
					PrecType:    PreTy,
					Prec:        Pre,
				})
			}
		} else {
			// Shift
			sym := lalr.G.Symbols[tr.sym_or_rule]
			action_set[int(tr.sym_or_rule)] = append(action_set[int(tr.sym_or_rule)],
				&Action{
					Sym:         sym,
					ActionType:  SHIFT,
					ActionIndex: tr.to,
					PrecType:    sym.PrecType,
					Prec:        sym.Prec,
				})
		}
	}
	//Check conflict
	for syIndex, actionlist := range action_set {
		if len(actionlist) > 1 {
			// Conflict
			res := actionlist

			var err error
			var act *Action
			for {
				if len(res) == 1 {
					action_set[syIndex] = res
					break
				}
				if act, err = lalr.ResolveConflict(res[0], res[1]); err != nil {
					return action_set, fmt.Errorf("cannot resolve the conflic %d, sym %d, conflict Type %s, %s ",
						state, syIndex, getActionTypeName(res[0].ActionType), getActionTypeName(res[1].ActionType))
				}
				res[1] = act
				res = res[1:]
			}

		}
	}
	return action_set, nil
}

// generate Table
func (lalr *LALR1) GenTable() ([][]int, error) {
	//group transist
	trans_set := make(map[int][]Transistor)
	for q := range lalr.G.LR0.LR0Closure {
		for _, tr := range lalr.trans {
			if tr.q == q {
				trans_set[q] = append(trans_set[q], tr)
			}
		}

	}
	tableGen := [][]int{}
	for q := 0; q < len(lalr.G.LR0.LR0Closure); q++ {
		if set, err := lalr.CheckAndResolveConflict(q, trans_set[q]); err != nil {
			lalr.ShowAndCheckConflict(q, trans_set[q])
			return nil, err
		} else {
			row := make([]int, len(lalr.G.Symbols))
			for i := 0; i < len(row); i++ {
				row[i] = lalr.GenErrorCode()
			}
			for syInd, act := range set {
				if act[0].ActionType != ERROR {
					if act[0].ActionIndex != 0 {
						row[syInd] = act[0].ActionIndex
					} else {
						row[syInd] = lalr.GenAcceptCode()
					}
				} else {
					if utils.DebugFlags {
						fmt.Println("it is nonassoc, should error")
					}

				}
			}
			tableGen = append(tableGen, row)
		}
	}
	return tableGen, nil
}

func (lalr *LALR1) SplitActionAndGotoTable(tab [][]int) ([][]int, [][]int) {
	nTerminals := len(lalr.G.VtSet)
	nNonTerminals := len(lalr.G.VnSet)
	fmt.Println("nTerminals", nTerminals, "nNonTerminals", nNonTerminals)
	actionTable := [][]int{}
	for i := 0; i < len(tab); i++ {
		// Copy start symbol and all terminal symbol
		row := make([]int, nTerminals+1)
		copy(row, tab[i][:nTerminals+1])
		actionTable = append(actionTable, row)
	}
	// goto table , make row index is symbol ,column index is state.
	gotoTable := [][]int{}
	//skip start symbol
	for i := 0; i < nNonTerminals; i++ {
		row := make([]int, len(tab))
		for j := 0; j < len(row); j++ {
			row[j] = tab[j][nTerminals+i+1]
		}
		gotoTable = append(gotoTable, row)
	}
	return actionTable, gotoTable

}

func getActionTypeName(act E_ActionType) string {
	switch act {
	case REDUCE:
		return "reduce"
	case SHIFT:
		return "shift"
	}
	return "error"
}
