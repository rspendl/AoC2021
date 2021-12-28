package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type Topo [][]int

func (t Topo) Adjacent(i, j int) ([]int, [][2]int) {
	var adj []int
	var adjCoords [][2]int
	if i > 0 {
		adj = append(adj, t[j][i-1])
		adjCoords = append(adjCoords, [2]int{i - 1, j})
		if j > 0 {
			adj = append(adj, t[j-1][i-1])
			adjCoords = append(adjCoords, [2]int{i - 1, j - 1})
		}
		if j < len(t)-1 {
			adj = append(adj, t[j+1][i-1])
			adjCoords = append(adjCoords, [2]int{i - 1, j + 1})
		}
	}
	if i < len(t[j])-1 {
		adj = append(adj, t[j][i+1])
		adjCoords = append(adjCoords, [2]int{i + 1, j})
		if j > 0 {
			adj = append(adj, t[j-1][i+1])
			adjCoords = append(adjCoords, [2]int{i + 1, j - 1})
		}
		if j < len(t)-1 {
			adj = append(adj, t[j+1][i+1])
			adjCoords = append(adjCoords, [2]int{i + 1, j + 1})
		}
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

func (t *Topo) Step() int {
	var flashes int
	for j := range *t {
		for i := range (*t)[j] {
			(*t)[j][i]++
		}
	}

	for {
		var hasFlashed bool
		for j := range *t {
			for i, v := range (*t)[j] {
				if v > 9 {
					_, adjCoords := t.Adjacent(i, j)
					for _, c := range adjCoords {
						(*t)[c[1]][c[0]]++
					}
					(*t)[j][i] = -100
					hasFlashed = true
					flashes++
				}
			}
		}
		if !hasFlashed {
			break
		}
	}

	for j := range *t {
		for i, v := range (*t)[j] {
			if v < 0 {
				(*t)[j][i] = 0
			}
		}
	}
	return flashes
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
	var flashes int
	var i int
	for {
		fs := topo.Step()
		if fs == 100 {
			log.Printf("all flashed in step: %d", i+1)
			break
		}
		if i < 100 {
			flashes += fs
			if i == 99 {
				log.Printf("flashes in 100 steps: %d", flashes)
			}
		}
		i++
	}
}
