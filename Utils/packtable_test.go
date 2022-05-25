/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPackTable1(t *testing.T) {
	table := [][]int{
		{0, 0, 1, 2, 0, 0},
		{0, 0, 1, 2, 0, 0},
		{9, 0, 0, 0, 0, 8},
		{0, 10, 0, 0, 0, 0},
		{9, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 11, 0},
		{0, 0, 1, 2, 0, 0},
		{0, 0, 1, 2, 11, 0},
		{0, 10, 0, 0, 0, 0},
	}
	fmt.Println(table)
	T, D, C := PackTable(table)
	fmt.Println(T)
	fmt.Println(D)
	fmt.Println(C)
	R := UnPackTable(9, 6, T, D, C)
	if !reflect.DeepEqual(table, R) {
		t.Error("PackTable1 failed")
	}

}

func TestPackTable2(t *testing.T) {
	table := [][]int{
		{0, 0, 1},
		{0, 1, 0},
		{0, 0, 1},
	}
	fmt.Println(table)
	T, D, C := PackTable(table)
	fmt.Println(T)
	fmt.Println(D)
	fmt.Println(C)
	R := UnPackTable(3, 3, T, D, C)
	if !reflect.DeepEqual(table, R) {
		t.Error("PackTable2 failed")
	}
}

func TestPackTable3(t *testing.T) {
	table := [][]int{
		{0, 0, 2, 0, 3},
		{0, 206, 0, 0, 0},
		{0, 0, 2, 0, 3},
		{0, -2, 0, -2, 0},
		{0, 0, 0, 5, 0},
		{0, -1, 0, -1, 0},
	}
	fmt.Println(table)
	T, D, C := PackTable(table)
	fmt.Println(T)
	fmt.Println(D)
	fmt.Println(C)
	R := UnPackTable(6, 5, T, D, C)
	if !reflect.DeepEqual(table, R) {
		t.Error("PackTable2 failed")
	}
}
