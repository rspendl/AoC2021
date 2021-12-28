package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

const DIM = 1500

type Topo struct {
	dots [DIM][DIM]bool
	mx   int
	my   int
}

// func (t Topo) Adjacent(i, j int) ([]int, [][2]int) {
// 	var adj []int
// 	var adjCoords [][2]int
// 	if i > 0 {
// 		adj = append(adj, t[j][i-1])
// 		adjCoords = append(adjCoords, [2]int{i - 1, j})
// 		if j > 0 {
// 			adj = append(adj, t[j-1][i-1])
// 			adjCoords = append(adjCoords, [2]int{i - 1, j - 1})
// 		}
// 		if j < len(t)-1 {
// 			adj = append(adj, t[j+1][i-1])
// 			adjCoords = append(adjCoords, [2]int{i - 1, j + 1})
// 		}
// 	}
// 	if i < len(t[j])-1 {
// 		adj = append(adj, t[j][i+1])
// 		adjCoords = append(adjCoords, [2]int{i + 1, j})
// 		if j > 0 {
// 			adj = append(adj, t[j-1][i+1])
// 			adjCoords = append(adjCoords, [2]int{i + 1, j - 1})
// 		}
// 		if j < len(t)-1 {
// 			adj = append(adj, t[j+1][i+1])
// 			adjCoords = append(adjCoords, [2]int{i + 1, j + 1})
// 		}
// 	}
// 	if j > 0 {
// 		adj = append(adj, t[j-1][i])
// 		adjCoords = append(adjCoords, [2]int{i, j - 1})
// 	}
// 	if j < len(t)-1 {
// 		adj = append(adj, t[j+1][i])
// 		adjCoords = append(adjCoords, [2]int{i, j + 1})
// 	}
// 	return adj, adjCoords
// }

// func (t *Topo) Step() int {
// 	var flashes int
// 	for j := range *t {
// 		for i := range (*t)[j] {
// 			(*t)[j][i]++
// 		}
// 	}

// 	for {
// 		var hasFlashed bool
// 		for j := range *t {
// 			for i, v := range (*t)[j] {
// 				if v > 9 {
// 					_, adjCoords := t.Adjacent(i, j)
// 					for _, c := range adjCoords {
// 						(*t)[c[1]][c[0]]++
// 					}
// 					(*t)[j][i] = -100
// 					hasFlashed = true
// 					flashes++
// 				}
// 			}
// 		}
// 		if !hasFlashed {
// 			break
// 		}
// 	}

// 	for j := range *t {
// 		for i, v := range (*t)[j] {
// 			if v < 0 {
// 				(*t)[j][i] = 0
// 			}
// 		}
// 	}
// 	return flashes
// }

func (t *Topo) Fold(f [2]int) {
	if f[0] > 0 {
		fx := f[0]
		for i := fx + 1; i <= t.mx; i++ {
			for j := 0; j <= t.my; j++ {
				t.dots[j][2*fx-i] = t.dots[j][2*fx-i] || t.dots[j][i]
			}
		}
		t.mx = fx - 1
	}
	if f[1] > 0 {
		fy := f[1]
		for j := fy + 1; j <= t.my; j++ {
			for i := 0; i <= t.mx; i++ {
				t.dots[2*fy-j][i] = t.dots[2*fy-j][i] || t.dots[j][i]
			}
		}
		t.my = fy - 1
	}
}

func (t *Topo) Print() {
	for j := 0; j <= t.my; j++ {
		var s string
		for i := 0; i <= t.mx; i++ {
			if t.dots[j][i] {
				s = s + "#"
			} else {
				s = s + "."
			}
		}
		log.Println(s)
	}
}

func (t *Topo) NDots() int {
	var n int
	for j := 0; j <= t.my; j++ {
		for i := 0; i <= t.mx; i++ {
			if t.dots[j][i] {
				n++
			}
		}
	}
	return n
}

func readTopo(fileName string) Topo {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	var maxX, maxY int
	var topo Topo
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		x, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal(err)
		}
		y, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatal(err)
		}
		topo.dots[y][x] = true
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}
	topo.mx = maxX
	topo.my = maxY
	return topo
}

func main() {

	// fold along y=7
	// fold along x=5
	var Fold = [][2]int{{0, 7}, {5, 0}}

	// fold along x=655
	// fold along y=447
	// fold along x=327
	// fold along y=223
	// fold along x=163
	// fold along y=111
	// fold along x=81
	// fold along y=55
	// fold along x=40
	// fold along y=27
	// fold along y=13
	// fold along y=6
	Fold = [][2]int{
		{655, 0},
		{0, 447},
		{327, 0},
		{0, 223},
		{163, 0},
		{0, 111},
		{81, 0},
		{0, 55},
		{40, 0},
		{0, 27},
		{0, 13},
		{0, 6},
	}
	topo := readTopo("input")
	log.Print(topo.mx, topo.my)
	// topo.Print()
	for i, f := range Fold {
		topo.Fold(f)
		if i == 0 {
			log.Printf("dots after 1st fold: %d", topo.NDots())
		}
	}
	topo.Print()
}
