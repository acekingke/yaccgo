/*
Copyright (c) 2021 Ke Yuchang(aceking.ke@gmail.com). All rights reserved.
Use of this source code is governed by MIT license that can be found in the LICENSE file.
*/
package graph

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/awalterschulze/gographviz"
)

type GraghNode struct {
	StateNumber int
	Children    []string
}

func (node *GraghNode) GenDotGraph(graphInst *gographviz.Graph) *gographviz.Graph {
	labels := fmt.Sprintf("\"<f0> state %d|", node.StateNumber)
	if len(node.Children) != 0 {
		labels += "{"
		for _, val := range node.Children {
			labels += val + "|"
		}
		labels = labels[0 : len(labels)-1]
		labels += "}\""
	}

	graphInst.AddNode("G", fmt.Sprintf("state_%d", node.StateNumber), map[string]string{
		"shape": "record",
		"label": labels,
	})
	return graphInst
}

func AddEdge(graphInst *gographviz.Graph, from, to int, Label string) *gographviz.Graph {
	fromNode, toNode := fmt.Sprintf("state_%d", from), fmt.Sprintf("state_%d", to)
	graphInst.AddEdge(fromNode, toNode, true, map[string]string{
		"label": Label,
	})
	return graphInst
}

func NewGraph() *gographviz.Graph {
	graphAst, _ := gographviz.ParseString(`digraph G {}`)
	graphInst := gographviz.NewGraph()
	//graphInst.AddAttr("", "rankdir", "LR")
	if err := gographviz.Analyse(graphAst, graphInst); err != nil {
		panic(err)
	}
	return graphInst
}
func SaveGraph(path string, graph *gographviz.Graph) error {

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	output := graph.String()
	fmt.Println(output)
	cmd := exec.Command("dot", "-Tpng")
	cmd.Stdin = strings.NewReader(output)
	cmd.Stdout = file
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	cmd.Wait()

	return nil
}
func DrawSimple() {
	graphInst := NewGraph()
	graphInst.AddNode("G", "a", map[string]string{
		"shape": "record",
		"label": "\"<f0> state 0|{hello|world |nice}\"",
	})
	graphInst.AddNode("G", "b", nil)
	graphInst.AddEdge("a", "b", true, nil)
	//SaveGraph("./one.png", graphInst)
	fmt.Println(graphInst.String())

}
