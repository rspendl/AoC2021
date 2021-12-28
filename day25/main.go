package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type Cucumber int

type Topo struct {
	dots [][]Cucumber
	mx   int
	my   int
}

func (t *Topo) Print() {
	my := len(t.dots)
	mx := len(t.dots[0])
	for j := 0; j < my; j++ {
		var s string
		for i := 0; i < mx; i++ {
			switch t.dots[j][i] {
			case 0:
				s = s + "."
			case 1:
				s = s + ">"
			case 2:
				s = s + "v"
			default:
				log.Fatal("invalid cucumber value ", t.dots[j][i])
			}
		}
		log.Println(s)
	}
}

func (c Cucumber) Move(t *Topo, i, j int) (int, int) {
	var x, y int
	switch c {
	case 0:
		log.Fatal("can't move empty on ", i, ",", j)
	case 1:
		x = i + 1
		if x >= t.mx {
			x = 0
		}
		y = j
	case 2:
		y = j + 1
		if y >= t.my {
			y = 0
		}
		x = i
	}
	return x, y
}

func (t *Topo) Move() (Topo, bool) {
	nd := make([][]Cucumber, t.my)
	for i := range nd {
		nd[i] = make([]Cucumber, t.mx)
	}
	nd1 := make([][]Cucumber, t.my)
	for i := range nd1 {
		nd1[i] = make([]Cucumber, t.mx)
	}
	var moved bool
	// move > first
	for j, l := range t.dots {
		for i, c := range l {
			if c == 1 {
				i2, j2 := c.Move(t, i, j)
				// log.Printf("%d: %d,%d->%d,%d", c, i, j, i2, j2)
				if t.dots[j2][i2] == 0 {
					nd1[j2][i2] = c
					nd1[j][i] = 0
					moved = true
				} else {
					nd1[j][i] = c
				}
			}
			if c == 2 {
				nd1[j][i] = c
			}
		}
	}
	t1 := Topo{
		dots: nd1,
		mx:   t.mx,
		my:   t.my,
	}
	// move v after
	for j, l := range nd1 {
		for i, c := range l {
			if c == 2 {
				i2, j2 := c.Move(&t1, i, j)
				// log.Printf("%d: %d,%d->%d,%d", c, i, j, i2, j2)
				if nd1[j2][i2] == 0 {
					nd[j2][i2] = c
					nd[j][i] = 0
					moved = true
				} else {
					nd[j][i] = c
				}
			}
			if c == 1 {
				nd[j][i] = c
			}
		}
	}

	return Topo{
		dots: nd,
		mx:   t.mx,
		my:   t.my,
	}, moved
}

// func (t *Topo) NDots() int {
// 	var n int
// 	for _, l := range t.dots {
// 		for _, c := range l {
// 			if c {
// 				n++
// 			}
// 		}
// 	}
// 	return n
// }

// func (t *Topo) Pixel(i, j int, blink bool) bool {
// 	my := len(t.dots)
// 	mx := len(t.dots[0])
// 	if i < 0 || i >= mx || j < 0 || j >= my {
// 		return blink
// 	}
// 	return t.dots[j][i]
// }

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

		var line []Cucumber
		for _, c := range record[0] {
			var b Cucumber
			switch c {
			case '.':
				b = 0
			case '>':
				b = 1
			case 'v':
				b = 2
			default:
				log.Print("invalid input char: ", c)
			}
			line = append(line, b)
		}
		topo.dots = append(topo.dots, line)
	}
	topo.mx = len(topo.dots[0])
	topo.my = len(topo.dots)
	return topo
}

func main() {
	topo := readTopo("input")
	// topo.Print()
	// log.Print("---")
	nm := 0
	var moved bool
	for {
		topo, moved = topo.Move()
		nm++
		if !moved {
			break
		}
		// log.Print(nm, ": ---")
	}
	log.Print("steps: ", nm)
}
