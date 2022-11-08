/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package graph

import (
	"fmt"
	"testing"

	"github.com/awalterschulze/gographviz"
)

// func TestDrawSimple(t *testing.T) {
// 	DrawSimple()
// }

func TestGraghNode_GenDotGraph(t *testing.T) {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graph := gographviz.NewGraph()
	graph.AddAttr("", "rankdir", "LR")
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		panic(err)
	}
	node := GraghNode{
		StateNumber: 1,
		Children: []string{
			"E-\\>\u2022E",
			"E-\\>\u2022E + E",
		},
	}
	node2 := GraghNode{
		StateNumber: 2,
		Children: []string{
			"E-\\>E\u2022",
		},
	}
	graph = node2.GenDotGraph(graph)
	graph = node.GenDotGraph(graph)
	graph = AddEdge(graph, 1, 2, "E")
	output := graph.String()
	fmt.Println(output)
	// file, err := os.Create("./one.png")
	// if err != nil {
	// 	fmt.Printf("err %s\n", err.Error())
	// }
	// cmd := exec.Command("dot", "-Tpng")
	// cmd.Stdin = strings.NewReader(output)
	// cmd.Stdout = file
	// cmd.Stderr = os.Stderr
	// if err := cmd.Start(); err != nil {
	// 	panic(err)
	// }
	// cmd.Wait()
	// if err := file.Close(); err != nil {
	// 	fmt.Printf("err %s\n", err.Error())
	// }
}
