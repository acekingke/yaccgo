/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package grammar

import (
	"fmt"

	graph "github.com/acekingke/yaccgo/Graph"
	item "github.com/acekingke/yaccgo/Items"
	utils "github.com/acekingke/yaccgo/Utils"
)

func (g *Grammar) ItemToStr(It *item.Item) string {
	res := ""
	r := g.ProductoinRules[It.RuleIndex]
	res += r.LeftPart.Name + "-\\>"
	if len(g.ProductoinRules[It.RuleIndex].RighPart) == 0 {
		res += "\u03b5"
		return res
	}

	for index, sym := range r.RighPart {
		if index == It.Dot {
			res += "\u2022"
		}
		res += fmt.Sprintf(" %s", utils.EscapeDotGraph(utils.RemoveTempName(sym.Name)))
	}
	if len(r.RighPart) == It.Dot {
		res += "\u2022"
	}
	return res
}

func (g *Grammar) StateGraphNode(IC *item.ItemCloure) *graph.GraghNode {
	Node := &graph.GraghNode{
		StateNumber: IC.Index,
	}
	for _, item := range IC.Items {
		Node.Children = append(Node.Children, g.ItemToStr(item))
	}
	return Node
}
