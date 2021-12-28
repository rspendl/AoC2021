package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

type Node struct {
	id      int
	name    string
	visited int
}

func (n *Node) IsStart() bool {
	return n.name == "start"
}

func (n *Node) IsEnd() bool {
	return n.name == "end"
}

func (g *Graph) IsVisited(id int) bool {
	if g.Start() == id {
		return true // don't visit start
	}
	if g.End() == id {
		return false // can visit end
	}
	n := g.Node(id)
	if n.name == strings.ToUpper(n.name) {
		return false
	}
	if n.visited == 0 {
		return false
	}
	if n.visited >= 2 {
		return true
	}
	// if it's the only cave, visited twice, it's OK
	otherVisited := false
	for _, l := range *g {
		for _, n := range l {
			if n.id != id && !n.IsStart() && !n.IsEnd() {
				if n.name == strings.ToLower(n.name) && n.visited > 1 {
					otherVisited = true
				}
			}
		}
	}
	return otherVisited
}

type Graph [][2]Node

func (g *Graph) Start() int {
	for _, l := range *g {
		if l[0].IsStart() {
			return l[0].id
		}
		if l[1].IsStart() {
			return l[1].id
		}
	}
	return 0
}

func (g *Graph) End() int {
	for _, l := range *g {
		if l[0].IsEnd() {
			return l[0].id
		}
		if l[1].IsEnd() {
			return l[1].id
		}
	}
	return 0
}

func (g *Graph) Node(id int) Node {
	for _, l := range *g {
		if l[0].id == id {
			return l[0]
		}
		if l[1].id == id {
			return l[1]
		}
	}
	log.Panicf("Node %d not found", id)
	return Node{}
}

func (g *Graph) Next(id int) []int {
	var nodes []int
	for _, l := range *g {
		for i, n := range l {
			if n.id == id && !n.IsEnd() { // no paths from "end"
				var j int
				if i == 0 {
					j = 1
				}
				if !g.IsVisited(l[j].id) && !l[j].IsStart() {
					nodes = append(nodes, l[j].id)
				}
			}
		}
	}
	return removeDuplicateInt(nodes)
}

func (g *Graph) Visit(id int) {
	for j, l := range *g {
		for i, n := range l {
			if n.id == id {
				if n.name == strings.ToLower(n.name) {
					(*g)[j][i].visited++
				}
			}
		}
	}
}

func (g *Graph) Leave(id int) {
	for j, l := range *g {
		for i, n := range l {
			if n.id == id {
				(*g)[j][i].visited--
			}
		}
	}
}

func (g *Graph) StringPath(path []int) string {
	var s string
	for _, p := range path {
		added := false
		for _, l := range *g {
			for _, n := range l {
				if n.id == p && !added {
					if s != "" {
						s += ","
					}
					s += n.name
					added = true
				}
			}
		}
	}
	return s
}

func (g *Graph) Paths(start int, length int) [][]int {
	g.Visit(start)
	nx := g.Next(start)
	var p [][]int
	if length == 1 || len(nx) == 0 {
		for _, n := range nx {
			if start != n {
				p = append(p, []int{start, n})
			}
		}
	} else {
		for _, n := range nx {
			newPaths := g.Paths(n, length-1)
			npNoDups := removeDuplicateSlice(newPaths)
			for _, np := range npNoDups {
				path := []int{start}
				path = append(path, np...)
				p = append(p, path)
			}
		}
	}
	g.Leave(start)
	return p
}

func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func removeDuplicateSlice(sliceSlice [][]int) [][]int {
	list := [][]int{}
	for i, item := range sliceSlice {
		if i == 0 {
			list = append(list, item)
			continue
		}
		found := false
		for _, s := range list {
			if len(item) != len(s) {
				continue
			}
			found := true
			for j, v := range s {
				if v != item[j] {
					found = false
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			list = append(list, item)
		}
	}
	return list
}

func readGraph(fileName string) Graph {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	r.Comma = '-'
	var (
		graph Graph
		i     int
	)
	ix := make(map[string]int)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		i1, ok := ix[record[0]]
		if !ok {
			i++
			i1 = i
			ix[record[0]] = i
		}
		i2, ok := ix[record[1]]
		if !ok {
			i++
			i2 = i
			ix[record[1]] = i
		}
		graph = append(graph, [2]Node{
			{i1, record[0], 0},
			{i2, record[1], 0},
		})
	}
	return graph
}

func main() {
	graph := readGraph("input")

	nPaths := 0
	l := 1
	for {
		paths := graph.Paths(graph.Start(), l)
		if len(paths) == 0 {
			log.Printf("longest path: %d, numpaths: %d", l-1, nPaths)
			break
		}
		for _, p := range paths {
			if len(p) > 0 && p[len(p)-1] == graph.End() {
				log.Print(graph.StringPath(p))
				nPaths++
			}
		}
		l++
	}
}
