/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package grammar

import (
	"fmt"
	"sort"

	item "github.com/acekingke/yaccgo/Items"
	lr "github.com/acekingke/yaccgo/LR"
	rule "github.com/acekingke/yaccgo/Rules"
	symbol "github.com/acekingke/yaccgo/Symbol"
)

/* G =(Vt, Vn, S, P) */
type Grammar struct {
	//Vt Set
	VtSet       map[*symbol.Symbol]bool
	VnSet       map[*symbol.Symbol]bool
	Symbols     []*symbol.Symbol
	SymbolsMap  map[string]*symbol.Symbol
	StartSymbol *symbol.Symbol
	//Productoin Rules
	ProductoinRules []*rule.ProductoinRule
	LR0             *lr.LR0
}

func NewGrammar() Grammar {
	return Grammar{
		VtSet:      make(map[*symbol.Symbol]bool),
		VnSet:      make(map[*symbol.Symbol]bool),
		SymbolsMap: make(map[string]*symbol.Symbol),
		LR0:        lr.NewLR0(),
	}
}

func (g *Grammar) InsertNewSymbol(s *symbol.Symbol) {
	g.Symbols = append(g.Symbols, s)
	g.SymbolsMap[s.Name] = s
}

func (g *Grammar) FindSymbolByName(name string) *symbol.Symbol {
	return g.SymbolsMap[name]
}

func (g *Grammar) GenStartSymbol() {
	s := symbol.NewSymbol(0, "start")
	s.SetNT()
	g.InsertNewSymbol(s)
	g.SetStartSymbol(s)
}

func (g *Grammar) InsertNewRules(r *rule.ProductoinRule) {
	g.ProductoinRules = append(g.ProductoinRules, r)
	r.LeftPart.SetNT()
	g.VnSet[r.LeftPart] = true
}

func (g *Grammar) SetStartSymbol(s *symbol.Symbol) {
	g.StartSymbol = s
}

//set NonTerminator set
func (g *Grammar) ResolveSymbols() {
	for _, s := range g.Symbols {
		if !s.IsNonTerminator {
			g.VtSet[s] = true
		}
	}
}

//calculate all epsilon
func (g *Grammar) CalculateEpsilonClosure() {
	change := 0
	for {
		change = 0
		for _, r := range g.ProductoinRules {
			//for every sym in RightPart, is IsEpsilonClosure, Is EpsilonClosure
			//empty is true
			every_isEpsilon := true
			for _, every_sy := range r.RighPart {
				every_isEpsilon = every_isEpsilon && every_sy.IsNonTerminator && every_sy.IsEpsilonClosure
			}
			if every_isEpsilon {
				if !(r.LeftPart.IsEpsilonClosure) {
					r.LeftPart.IsEpsilonClosure = true
					// epsilon is terminate too
					r.LeftPart.CanTerminate = true
					change++
				}
			}
		}
		if change == 0 {
			break
		}
	}
}

//calculate all epsilon
func (g *Grammar) CalculateCanTerminate() []*symbol.Symbol {
	change := 0
	for {
		change = 0
		for _, r := range g.ProductoinRules {
			//for every sym in RightPart, is IsEpsilonClosure, Is EpsilonClosure
			//empty is true
			every_CanTerm := true
			for _, every_sy := range r.RighPart {
				every_CanTerm = every_CanTerm && every_sy.CanTerminate
			}
			if every_CanTerm {
				if !(r.LeftPart.CanTerminate) {
					r.LeftPart.CanTerminate = true
					change++
				}
			}
		}
		if change == 0 {
			break
		}
	}
	var inf_cycles []*symbol.Symbol
	for sy := range g.VnSet {
		if !sy.CanTerminate {
			inf_cycles = append(inf_cycles, sy)
		}
	}
	return inf_cycles
}

func (g Grammar) PrintInfLoop(some []*symbol.Symbol) {
	fmt.Println("Error:")
	for _, s := range some {
		fmt.Printf("%s ", s.Name)
	}
	fmt.Println(" have Infinite recursion loop")
}

func (g Grammar) ShowAllSymbols() {
	for _, s := range g.Symbols {
		s.Show()
	}
}

func (g *Grammar) getItemCloure(It *item.Item) []*item.Item {
	if It.Dot == len(g.ProductoinRules[It.RuleIndex].RighPart) {
		return []*item.Item{}
	}
	dotsym := g.ProductoinRules[It.RuleIndex].RighPart[It.Dot]
	if !dotsym.IsNonTerminator {
		return []*item.Item{}
	} else {
		items := []*item.Item{}
		for index, r := range g.ProductoinRules {
			if r.LeftPart.ID == dotsym.ID {
				items = append(items, item.NewItem(index, 0))
			}
		}
		return items
	}
}

func (g *Grammar) ComputeIClosure(IC *item.ItemCloure) {
	var change int = 0
	var items []*item.Item

	for {
		change = 0
		items = []*item.Item{}
		for _, it := range IC.Items {
			items = append(items, g.getItemCloure(it)...)
		}
		for _, i := range items {
			change += IC.InsertItem(i)
		}
		if change == 0 {
			break
		}
	}
	//sort for check
	sort.SliceStable(IC.Items, func(i, j int) bool {
		if IC.Items[i].RuleIndex < IC.Items[j].RuleIndex {
			return true
		} else if IC.Items[i].RuleIndex == IC.Items[j].RuleIndex {
			if IC.Items[i].Dot < IC.Items[j].Dot {
				return true
			}
		}
		return false
	})
}

//compute the closure goto item
func (g *Grammar) ComputeGotoItemRecursive(IC *item.ItemCloure) {
	change := 0
	for _, it := range IC.Items {
		r := g.ProductoinRules[it.RuleIndex]
		if it.Dot < len(r.RighPart) {
			sy := r.RighPart[it.Dot]
			if goToPointer := IC.FindItemClosure(sy); goToPointer != nil {
				ExistIC := g.LR0.LR0Closure[goToPointer.ItemCl]
				change += ExistIC.InsertItem(item.NewItem(it.RuleIndex, it.Dot+1))
				g.ComputeIClosure(ExistIC)

			} else {
				newIC := item.NewItemCloure()
				newIC.InsertItem(item.NewItem(it.RuleIndex, it.Dot+1))
				var index_goto int
				g.ComputeIClosure(newIC)
				if exist_index, exist := g.LR0.CheckIsExist(newIC); exist {
					index_goto = exist_index
				} else {
					//check weath exist IC
					g.LR0.InsertItemClosure(newIC, false)
					index_goto = newIC.Index
				}
				newGoto := &item.GoToCloure{Sym: sy, ItemCl: index_goto}
				change += IC.InsertGoTO(newGoto)
			}

			if change == 0 { //find a cycle, return
				return
			}
		} // if Dot == len, just return
	}
	// then Compute Recursive
	for _, goTo := range IC.GoTo {
		ic := goTo.ItemCl
		g.ComputeGotoItemRecursive(g.LR0.LR0Closure[ic])
	}
}

func (g *Grammar) ComputeAllGoto() {
	var i int = 0
	for i < len(g.LR0.LR0Closure) {
		g.ComputeGotoItemNoneRec(g.LR0.LR0Closure[i])
		if len(g.LR0.LR0Closure) >= 2000 {
			g.Show()
			panic("too manay states!")
		}
		i++
	}
}

//compute the closure goto item
func (g *Grammar) ComputeGotoItemNoneRec(IC *item.ItemCloure) {
	change := 0
	for _, it := range IC.Items {
		r := g.ProductoinRules[it.RuleIndex]
		if it.Dot < len(r.RighPart) {
			sy := r.RighPart[it.Dot]
			if goToPointer := IC.FindItemClosure(sy); goToPointer != nil {
				ExistIC := goToPointer.ICref
				if goToPointer.ItemCl != -1 {
					ExistIC = g.LR0.LR0Closure[goToPointer.ItemCl]
				}
				change += ExistIC.InsertItem(item.NewItem(it.RuleIndex, it.Dot+1))
				g.ComputeIClosure(ExistIC)

			} else {
				newIC := item.NewItemCloure()
				newIC.InsertItem(item.NewItem(it.RuleIndex, it.Dot+1))

				var index_goto int = -1
				g.ComputeIClosure(newIC)

				newGoto := &item.GoToCloure{Sym: sy, ItemCl: index_goto, ICref: newIC}
				change += IC.InsertGoTO(newGoto)
			}

			if change == 0 { //find a cycle, return
				return
			}
		} // if Dot == len, just return
	}

	for _, goItem := range IC.GoToMap {
		IcTemp := goItem.ICref
		var index_goto int = -1
		if exist_index, exist := g.LR0.CheckIsExist(IcTemp); exist {
			index_goto = exist_index
		} else {
			//check weath exist IC
			g.LR0.InsertItemClosure(IcTemp, false)
			index_goto = IcTemp.Index
		}
		goItem.ItemCl = index_goto
		goItem.ICref = nil
	}
}

func (g *Grammar) ShowCloure(IC *item.ItemCloure) {
	fmt.Printf("--------state %d------------\n", IC.Index)
	for _, it := range IC.Items {
		r := g.ProductoinRules[it.RuleIndex]
		fmt.Printf("%s-->", r.LeftPart.Name)
		for _, sy := range r.RighPart[:it.Dot] {
			fmt.Printf(" %s ", sy.Name)
		}
		fmt.Print("@")
		for _, sy := range r.RighPart[it.Dot:] {
			fmt.Printf(" %s ", sy.Name)
		}
		fmt.Print("\n")
	}
	fmt.Println("GOTO:")
	for _, g := range IC.GoTo {
		fmt.Printf("at %s goto %d \n", g.Sym.Name, g.ItemCl)
	}

}

func (g *Grammar) Show() {
	for _, ic := range g.LR0.LR0Closure {
		g.ShowCloure(ic)
	}
}
