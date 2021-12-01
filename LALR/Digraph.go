/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package lalr

type Relation struct {
	Index int
	x     int // may be transistor index, may be rule
	y     int // may be transistor index, may be rule
}

type stack []int

func (s *stack) Push(v int) {
	*s = append(*s, v)
}

func (s *stack) Pop() int {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func Digraph(X []int, R []Relation,
	Fp map[int][]int, F *map[int][]int) {
	N := make(map[int]int, len(X))
	for _, x := range X {
		N[x] = 0
	}
	var S stack
	for _, x := range X {
		if N[x] == 0 {
			Traverse(x, R, Fp, F, N, &S)
		}

	}
}

func Union(a []int, b []int) []int {
	c := b
	for _, v := range a {
		found := false
		for _, u := range b {
			if v == u {
				found = true
			}
		}
		if !found {
			c = append(c, v)
		}
	}
	return c
}

func Traverse(x int, R []Relation,
	FP map[int][]int, F *map[int][]int,
	N map[int]int, S *stack) {
	(*S).Push(x)
	d := len(*S)
	N[x] = d

	(*F)[x] = FP[x]
	for _, r := range R {
		if r.x == x {
			y := r.y
			if N[y] == 0 {
				Traverse(y, R, FP, F, N, S)
			}
			N[x] = min(N[x], N[y])
			(*F)[x] = Union((*F)[y], (*F)[x])
		}
	}
	if N[x] == d {
		N[x] = MaxInt
		for {
			top := (*S).Pop()
			N[top] = MaxInt

			(*F)[top] = (*F)[x]
			if top == x {
				break
			}
		}
	}

}
