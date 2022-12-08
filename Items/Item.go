/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package items

import symbol "github.com/acekingke/yaccgo/Symbol"

type Item struct {
	RuleIndex int // The index in Grammar  production rules
	Dot       int // such as A->.b  Dot is zero, A->b. Dit is 1
}
type GoToCloure struct {
	Sym    *symbol.Symbol
	ItemCl int
	ICref  *ItemCloure
}

type ItemCloure struct {
	Index   int
	Items   []*Item // the item clouser
	itemMap map[Item]bool
	GoTo    []*GoToCloure //ItemClure Index in LR0
	GoToMap map[*symbol.Symbol]*GoToCloure
}

func NewItem(r_index int, dot int) *Item {
	return &Item{RuleIndex: r_index, Dot: dot}
}

func NewItemCloure() *ItemCloure {
	return &ItemCloure{Items: make([]*Item, 0),
		GoTo:    make([]*GoToCloure, 0),
		itemMap: make(map[Item]bool),
		GoToMap: make(map[*symbol.Symbol]*GoToCloure)}
}

func (IC *ItemCloure) InsertItem(It *Item) int {
	if IC.itemMap[*It] {
		return 0
	}
	IC.Items = append(IC.Items, It)
	IC.itemMap[*It] = true
	return 1
}

//one Core may be has many core Item
func (IC *ItemCloure) InsertGoTO(Goto *GoToCloure) int {
	if IC.GoToMap[Goto.Sym] != nil {
		return 0
	}
	IC.GoTo = append(IC.GoTo, Goto)
	IC.GoToMap[Goto.Sym] = Goto
	return 1
}

func (IC *ItemCloure) FindItemClosure(sy *symbol.Symbol) *GoToCloure {
	return IC.GoToMap[sy]
}
