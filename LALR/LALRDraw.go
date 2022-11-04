/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package lalr

import (
	"fmt"

	graph "github.com/acekingke/yaccgo/Graph"
	utils "github.com/acekingke/yaccgo/Utils"

	"github.com/awalterschulze/gographviz"
)

func (lalr *LALR1) DrawGrammar(tab [][]int) *gographviz.Graph {
	graphInst := graph.NewGraph()
	for _, iC := range lalr.G.LR0.LR0Closure {
		node := lalr.G.StateGraphNode(iC)
		graphInst = node.GenDotGraph(graphInst)
	}

	for stateNum, r := range tab {
		for SymNum, d := range r {
			if d != lalr.GenErrorCode() &&
				d != lalr.GenAcceptCode() && d >= 0 {
				from := stateNum
				to := d
				edage := fmt.Sprintf("\"%s\"",
					utils.EscapeDotGraph(utils.RemoveTempName(lalr.G.Symbols[SymNum].Name)))
				graphInst = graph.AddEdge(graphInst, from, to, edage)
			} else if d == lalr.GenAcceptCode() {
				n := graphInst.Nodes.Lookup[fmt.Sprintf("state_%d", stateNum)]
				n.Attrs.Add("style", "filled")
				n.Attrs.Add("fillcolor", "\"yellow:green\"")
				n.Attrs.Add("gradientangle", "315")
			}

		}
	}
	return graphInst
}
