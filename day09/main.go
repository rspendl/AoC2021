package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
)

type Topo [][]int

func (t Topo) LowPoints() ([]int, [][2]int) {
	var lp []int
	var lpCoords [][2]int
	for j, line := range t {
		for i, p := range line {
			adj, _ := t.Adjacent(i, j)
			isLow := true
			for _, a := range adj {
				if a <= p {
					isLow = false
					break
				}
			}
			if isLow {
				lp = append(lp, p)
				lpCoords = append(lpCoords, [2]int{i, j})
			}
		}
	}
	return lp, lpCoords
}

func (t Topo) Adjacent(i, j int) ([]int, [][2]int) {
	var adj []int
	var adjCoords [][2]int
	if i > 0 {
		adj = append(adj, t[j][i-1])
		adjCoords = append(adjCoords, [2]int{i - 1, j})
	}
	if i < len(t[j])-1 {
		adj = append(adj, t[j][i+1])
		adjCoords = append(adjCoords, [2]int{i + 1, j})
	}
	if j > 0 {
		adj = append(adj, t[j-1][i])
		adjCoords = append(adjCoords, [2]int{i, j - 1})
	}
	if j < len(t)-1 {
		adj = append(adj, t[j+1][i])
		adjCoords = append(adjCoords, [2]int{i, j + 1})
	}
	return adj, adjCoords
}

func (t Topo) BasinSize(i, j int) int {
	var visited [][]bool
	for _, l := range t {
		visited = append(visited, make([]bool, len(l)))
	}
	visited[j][i] = true
	size := 1
	for {
		var newVisited bool
		for y, l := range t {
			for x, v := range l {
				if visited[y][x] {
					adj, adjCoords := t.Adjacent(x, y)
					for k, av := range adj {
						ax := adjCoords[k][0]
						ay := adjCoords[k][1]
						if av < 9 && !visited[ay][ax] && v < av {
							visited[ay][ax] = true
							size++
							newVisited = true
						}
					}
				}
			}
		}
		if !newVisited {
			break
		}
	}
	return size
}

func readTopo(fileName string) Topo {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	var topo Topo
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var line []int
		for _, c := range record[0] {
			n := int(c - '0')
			line = append(line, n)
		}
		topo = append(topo, line)
	}
	return topo
}

func main() {
	topo := readTopo("input")
	lp, lpCoords := topo.LowPoints()
	var rlsum int
	for _, p := range lp {
		rlsum += (p + 1)
	}

	var basinSize []int

	for _, lp := range lpCoords {
		bs := topo.BasinSize(lp[0], lp[1])
		basinSize = append(basinSize, bs)
	}
	sort.Slice(basinSize, func(i, j int) bool { return basinSize[i] > basinSize[j] })
	log.Printf("risk level sum: %d, 3 basin sizes: %d, num-basins: %d", rlsum, basinSize[0]*basinSize[1]*basinSize[2], len(basinSize))
}
