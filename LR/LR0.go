/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package lr

import (
	item "github.com/acekingke/yaccgo/Items"
)

type LR0 struct {
	LR0Closure []*item.ItemCloure
	//mapLR0Cl   map[item.Item]*item.ItemCloure
}

func NewLR0() *LR0 {
	return &LR0{LR0Closure: make([]*item.ItemCloure, 0)}
}

//
func (lr0 *LR0) CheckIsExist(IC *item.ItemCloure) (int, bool) {
	var found bool = true
	for index, ic_in := range lr0.LR0Closure {
		if len(ic_in.Items) == len(IC.Items) {
			found = true
			//if ic_in item all equal IC.items, return true
			for i := 0; i < len(ic_in.Items); i++ {
				found = found && (*ic_in.Items[i] == *IC.Items[i])
			}
			if found {
				return index, true
			}

		}
	}
	return -1, false
}
func (lr0 *LR0) InsertItemClosure(IC *item.ItemCloure, needCheck bool) int {
	//Cannot use by mapLR0Cl, should sort and search
	if len(IC.Items) == 0 {
		panic("Error: Items cannot empty")
	}
	var notExist bool = false
	if needCheck {
		if _, exists := lr0.CheckIsExist(IC); !exists {
			notExist = true
		}
	} else {
		notExist = true
	}
	if notExist {
		IC.Index = len(lr0.LR0Closure)
		lr0.LR0Closure = append(lr0.LR0Closure, IC)
		return IC.Index
	}

	return -1
}
