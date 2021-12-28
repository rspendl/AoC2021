package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type Topo struct {
	dots [][]bool
}

type IE [512]bool

func (t *Topo) Print() {
	my := len(t.dots)
	mx := len(t.dots[0])
	for j := 0; j < my; j++ {
		var s string
		for i := 0; i < mx; i++ {
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
	for _, l := range t.dots {
		for _, c := range l {
			if c {
				n++
			}
		}
	}
	return n
}

func (t *Topo) Pixel(i, j int, blink bool) bool {
	my := len(t.dots)
	mx := len(t.dots[0])
	if i < 0 || i >= mx || j < 0 || j >= my {
		return blink
	}
	return t.dots[j][i]
}

func (t *Topo) Get9(i, j int, blink bool) int {
	var pv int
	for y := j - 1; y <= j+1; y++ {
		for x := i - 1; x <= i+1; x++ {
			var bv int
			if t.Pixel(x, y, blink) {
				bv = 1
			}
			pv = 2*pv + bv
		}
	}
	// log.Printf("%d,%d: %9b", i, j, pv)
	return pv
}

func (t *Topo) Enhance(ie IE, step int) Topo {
	var nt Topo
	for j := -1; j < len(t.dots)+1; j++ {
		var l []bool
		for i := -1; i < len(t.dots[0])+1; i++ {
			var blink bool
			if ie[0] {
				blink = step%2 == 0
			}
			pv := t.Get9(i, j, blink)
			l = append(l, ie[pv])
		}
		nt.dots = append(nt.dots, l)
	}
	return nt
}

func readTopo(fileName string) (IE, Topo) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	var ie IE

	var topo Topo
	readingIE := true
	var ieLine string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if readingIE {
			if record[0] == "X" {
				for i := 0; i < len(ieLine); i++ {
					if ieLine[i] == '#' {
						ie[i] = true
					}
				}
				readingIE = false
				continue
			}
			ieLine += record[0]
			continue
		}

		var line []bool
		for _, c := range record[0] {
			var b bool
			if c == '#' {
				b = true
			}
			line = append(line, b)
		}
		topo.dots = append(topo.dots, line)
	}
	return ie, topo
}

func main() {
	ie, topo := readTopo("input")
	topo.Print()
	log.Print("---")
	for i := 1; i <= 50; i++ {
		topo = topo.Enhance(ie, i)
		topo.Print()
		log.Print("---")
	}
	log.Print(topo.NDots())
}
