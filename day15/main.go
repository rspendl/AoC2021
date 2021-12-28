package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"time"
)

type Path []int // 1=left, 2=down

type Topo struct {
	risk    [][]int
	cost    [][]int
	minPath [][]Path
}

func (t *Topo) Dist(r int) [][2]int {
	var nodes [][2]int
	mx := len(t.risk[0])
	my := len(t.risk)
	if r > mx+my {
		log.Panicf("radius %d larger than %d+%d", r, mx, my)
	}
	for j := r; j >= 0; j-- {
		if j < my {
			k := r - j
			if k < mx {
				nodes = append(nodes, [2]int{k, j})
			}
		}
	}

	return nodes
}

func (t *Topo) Costs(radius int) {
	mx := len(t.risk[0])
	my := len(t.risk)
	changed := false
	r := radius
	decreasing := true
	for !changed {
		changed = false
		for r > 0 && r <= radius {
			nodes := t.Dist(r)
			for _, n := range nodes {
				var paths []Path
				var costs []int
				x := n[0]
				y := n[1]

				if x > 0 {
					var p Path
					p = append(p, t.minPath[y][x-1]...)
					p = append(p, 1)
					paths = append(paths, p)
					cost := t.risk[y][x] + t.cost[y][x-1]
					tc := t.Cost(p)
					if cost != tc {
						log.Panicf("cost of path: %v == %d != %d", p, tc, cost)
					}
					costs = append(costs, cost)
				}
				if y > 0 {
					var p Path
					p = append(p, t.minPath[y-1][x]...)
					p = append(p, 2)
					paths = append(paths, p)
					cost := t.risk[y][x] + t.cost[y-1][x]
					tc := t.Cost(p)
					if cost != tc {
						log.Panicf("cost of path: %v == %d != %d", p, tc, cost)
					}
					costs = append(costs, cost)
				}
				if x < mx-1 {
					if t.cost[y][x+1] > 0 {
						var p Path
						p = append(p, t.minPath[y][x+1]...)
						p = append(p, 3)
						paths = append(paths, p)
						cost := t.risk[y][x] + t.cost[y][x+1]
						tc := t.Cost(p)
						if cost != tc {
							log.Panicf("cost of path: %v == %d != %d", p, tc, cost)
						}
						costs = append(costs, cost)
					} else {
					}
				}
				if y < my-1 {
					if t.cost[y+1][x] > 0 {
						var p Path
						p = append(p, t.minPath[y+1][x]...)
						p = append(p, 4)
						paths = append(paths, p)
						cost := t.risk[y][x] + t.cost[y+1][x]
						tc := t.Cost(p)
						if cost != tc {
							log.Panicf("cost of path: %v == %d != %d", p, tc, cost)
						}
						costs = append(costs, cost)
					} else {
					}
				}

				var i int
				mc := costs[0]
				for j, c := range costs {
					if c < mc {
						i = j
						mc = c
					}
				}
				if t.cost[y][x] == 0 || costs[i] < t.cost[y][x] {
					t.minPath[y][x] = paths[i]
					t.cost[y][x] = costs[i]
					// if (x == 8 && y == 8) || (x == 7 && y == 8) {
					// 	log.Printf("changed %d (%d,%d): %d (%d) %v", r, x, y, t.cost[y][x], t.Cost(t.minPath[y][x]), t.minPath[y][x])
					// 	// log.Printf("(%d,%d): %d %v", x, y, t.cost[y][x], t.minPath[y][x])
					// }
					changed = true
				}
				// if (x == 8 && y == 8) || (x == 7 && y == 8) {
				// 	log.Printf("%d (%d,%d): %d (%d) %v", r, x, y, t.cost[y][x], t.Cost(t.minPath[y][x]), t.minPath[y][x])
				// 	// log.Printf("(%d,%d): %d %v", x, y, t.cost[y][x], t.minPath[y][x])
				// }
			}
			if decreasing {
				r--
			} else {
				r++
			}
			if r == 0 && decreasing {
				decreasing = false
				r = 1
			}
		}
		if !decreasing && r > radius {
			break
		}
	}
}

func (t *Topo) Cost(path Path) int {
	x := 0
	y := 0
	var cost int
	for _, p := range path {
		switch p {
		case 1:
			x++
		case 2:
			y++
		case 3:
			x--
		case 4:
			y--
		}
		cost += t.risk[y][x]
	}
	return cost
}

func (t *Topo) Extend(xx, xy int) {
	mx := len(t.risk[0])
	my := len(t.risk)

	newT := Topo{}
	for j := 0; j < xy; j++ {
		for i := 0; i < xx; i++ {
			for l := 0; l < my; l++ {
				if i == 0 {
					line := make([]int, mx*xx)
					newT.risk = append(newT.risk, line)
				}
				for k := 0; k < mx; k++ {
					newT.risk[j*my+l][i*mx+k] = inc9(t.risk[l][k], i+j)
				}
				c0 := make([]int, mx*xx)
				newT.cost = append(newT.cost, c0)
				p0 := make([]Path, mx*xx)
				newT.minPath = append(newT.minPath, p0)
			}
		}
	}
	*t = newT
}

func inc9(x, y int) int {
	v := ((x - 1) + y) % 9
	v++
	return v
}

func readTopo(fileName string) Topo {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	var topo Topo
	for {
		var (
			line  []int
			cost  []int
			path  []Path
			final []bool
		)
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range record[0] {
			line = append(line, int(c-'0'))
			cost = append(cost, 0)
			path = append(path, Path{})
			final = append(final, false)
		}
		topo.risk = append(topo.risk, line)
		topo.cost = append(topo.cost, cost)
		topo.minPath = append(topo.minPath, path)
	}
	return topo
}

func main() {
	topo := readTopo("input")
	topo.Extend(5, 5)
	r := 1
	mx := len(topo.risk[0])
	my := len(topo.risk)
	t0 := time.Now()
	for {
		topo.Costs(r)
		log.Print(r, ": ", time.Since(t0))
		if r == mx+my {
			log.Printf("cost: %d %v", topo.cost[my-1][mx-1], topo.minPath[my-1][mx-1])
			break
		}
		r++
	}
}
