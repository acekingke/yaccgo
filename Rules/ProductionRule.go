/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package rules

import (
	symbol "github.com/acekingke/yaccgo/Symbol"
)

type ProductoinRule struct {
	LeftPart *symbol.Symbol
	RighPart []*symbol.Symbol
	//precedance symbol
	PrecSymbol *symbol.Symbol
}

func NewProductoinRule(LeftPart *symbol.Symbol, RighPart []*symbol.Symbol) *ProductoinRule {
	return &ProductoinRule{LeftPart: LeftPart, RighPart: RighPart}
}

func (p *ProductoinRule) SetPrecSymbol(s *symbol.Symbol) {
	p.PrecSymbol = s
}
