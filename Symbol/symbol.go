/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package symbol

import "fmt"

type E_Precedence int

const (
	LEFT E_Precedence = iota
	RIGHT
	NONE // no associate
)

type Symbol struct {
	ID               uint
	Value            int
	Name             string
	Tag              string // Tag indicate the type of symbol
	IsNonTerminator  bool
	IsEpsilonClosure bool
	PrecType         E_Precedence
	Prec             int
}

func NewSymbol(id uint, name string) *Symbol {
	return &Symbol{ID: id, Name: name, PrecType: NONE, Prec: -1}
}

func (s *Symbol) SetValue(val int) {
	s.Value = val
}

func (s *Symbol) SetTag(tag string) {
	s.Tag = tag
}

func (s *Symbol) GetTag() string {
	return s.Tag
}

func (s *Symbol) SetNT() {
	s.IsNonTerminator = true
}

func (s *Symbol) SetEpsilon() {
	s.IsEpsilonClosure = true
}

func (s *Symbol) SetPrecType(ty E_Precedence) {
	s.PrecType = ty
}

func (s *Symbol) SetPrec(Prec int) {
	s.Prec = Prec
}

func (s *Symbol) Show() {
	fmt.Println("=========")
	fmt.Println("ID:", s.ID)
	fmt.Println("Name", s.Name)
	fmt.Println("IsNonTerminator:", s.IsNonTerminator)
	fmt.Println("isEpsilon:", s.IsEpsilonClosure)
	fmt.Println("============")
}
