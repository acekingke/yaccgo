/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package utils

import "sort"

type Pair struct {
	a, b interface{}
}

// Use Tarjan and Yao method compress the two-dimensional array
func PackTable(table [][]int) ( /*T*/ []int /*D*/, []int /*Check*/, []int) {
	var row []int
	entry := make([]bool, len(table)*len(table[0]))
	//step 1 count every row non-zero element
	rowCount := []Pair{}
	for i := 0; i < len(table); i++ {
		rowCount = append(rowCount, Pair{i, 0})
		for j := 0; j < len(table[i]); j++ {
			if table[i][j] != 0 {
				rowCount[i].b = rowCount[i].b.(int) + 1
			}
		}
		row = append(row, 0)
	}
	//step 2 fetch all non-zero element position
	nonZeroPos := make(map[int][]int) //  list(i)
	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[i]); j++ {
			if table[i][j] != 0 {
				nonZeroPos[i] = append(nonZeroPos[i], j)
			}
		}
	}
	//step 3 compress
	//sort the count
	sort.SliceStable(rowCount, func(i, j int) bool {
		return rowCount[i].b.(int) > rowCount[j].b.(int)
	}) // rowCount is equal the bucket
	maxIndex := 0
	// from the largest to the smallest
	for _, p := range rowCount {
		i := p.a.(int)
		row[i] = 0
		//check overlap
	checkoverlap:
		for _, j := range nonZeroPos[i] {
			if entry[row[i]+j] {
				row[i]++
				goto checkoverlap
			}
		}
		for _, k := range nonZeroPos[i] {
			entry[row[i]+k] = true
			if maxIndex < row[i]+k {
				maxIndex = row[i] + k
			}
		}
	}
	var ret []int = make([]int, maxIndex+1)
	var check []int = make([]int, maxIndex+1)
	//init check with -1
	for i := 0; i < maxIndex+1; i++ {
		check[i] = -1
	}
	//step 4 output
	for i, js := range nonZeroPos {
		for _, k := range js {
			ret[row[i]+k] = table[i][k]
			check[row[i]+k] = i
		}
	}
	//Trim the zero element at the begin
	for i := 0; i < len(ret); i++ {
		if ret[i] != 0 {
			break
		}
		ret = ret[1:]
		check = check[1:]
		for j := 0; j < len(row); j++ {
			row[j]--
		}
	}
	return ret, row, check
}

func UnPackTable(rows int, cols int, T []int, D []int, C []int) [][]int {
	var table [][]int = make([][]int, rows)
	// step 1 find the maxIndex

	// step 2 allocate the table
	for i := 0; i < len(table); i++ {
		table[i] = make([]int, cols)
	}
	// step 3 fill the table
	for i := 0; i < len(D); i++ {
		for j := 0; j < cols; j++ {
			if D[i]+j < 0 || D[i]+j >= len(C) || C[D[i]+j] != i {
				table[i][j] = 0
			} else {
				table[i][j] = T[D[i]+j]
			}
		}
	}
	return table
}
