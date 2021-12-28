package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cube struct {
	on                     bool
	x1, x2, y1, y2, z1, z2 int
	// overlap                CubeList
}
type CubeList []Cube

type Univ [101][101][101]bool

func (u *Univ) Set(c Cube) {
	if c.x2 < -50 || c.x1 > 50 {
		return
	}
	if c.y2 < -50 || c.y1 > 50 {
		return
	}
	if c.z2 < -50 || c.z1 > 50 {
		return
	}
	x1 := c.x1
	if x1 < -50 {
		x1 = -50
	}
	x2 := c.x2
	if x2 > 50 {
		x2 = -50
	}
	y1 := c.y1
	if y1 < -50 {
		y1 = -50
	}
	y2 := c.y2
	if y2 > 50 {
		y2 = 50
	}
	z1 := c.z1
	if z1 < -50 {
		z1 = -50
	}
	z2 := c.z2
	if z2 > 50 {
		z2 = 50
	}
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			for z := z1; z <= z2; z++ {
				if x >= -50 && x <= 50 && y >= -50 && y <= 50 && z >= -50 && z <= 50 {
					u[x+50][y+50][z+50] = c.on
				}
			}
		}
	}
}

func (u *Univ) Count() int {
	var c int
	for x := -50; x <= 50; x++ {
		for y := -50; y <= 50; y++ {
			for z := -50; z <= 50; z++ {
				if u[x+50][y+50][z+50] {
					c++
				}
			}
		}
	}
	return c
}

func extract12(s string) (int, int) {
	ss := strings.Split(s, "=")
	ss = strings.Split(ss[1], "..")
	a, err := strconv.Atoi(ss[0])
	if err != nil {
		log.Fatalf("invalid number: %s", ss[0])
	}
	b, err := strconv.Atoi(ss[1])
	if err != nil {
		log.Fatalf("invalid number: %s", ss[1])
	}
	return a, b
}

func readCubes(fileName string) CubeList {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)

	var cubes CubeList
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		r1 := strings.Split(record[0], " ")
		var onOff bool
		if r1[0] == "on" {
			onOff = true
		}
		x1, x2 := extract12(r1[1])
		y1, y2 := extract12(record[1])
		z1, z2 := extract12(record[2])
		cubes = append(cubes, Cube{onOff, x1, x2, y1, y2, z1, z2})
	}
	return cubes
}

func (c Cube) Overlaps(c2 Cube) bool {
	if c.x2 < c2.x1 || c.x1 > c2.x2 {
		return false
	}
	if c.y2 < c2.y1 || c.y1 > c2.y2 {
		return false
	}
	if c.z2 < c2.z1 || c.z1 > c2.z2 {
		return false
	}
	return true
}

// IsWithin - c within c2
func (c Cube) IsWithin(c2 Cube) bool {
	if c.x1 < c2.x1 || c.x2 > c2.x2 {
		return false
	}
	if c.y1 < c2.y1 || c.y2 > c2.y2 {
		return false
	}
	if c.z1 < c2.z1 || c.z2 > c2.z2 {
		return false
	}
	return true
}

func (c Cube) Overlap(c2 Cube) Cube {
	if !c.Overlaps(c) {
		return Cube{}
	}
	x1 := max(c.x1, c2.x1)
	x2 := min(c.x2, c2.x2)
	y1 := max(c.y1, c2.y1)
	y2 := min(c.y2, c2.y2)
	z1 := max(c.z1, c2.z1)
	z2 := min(c.z2, c2.z2)
	return Cube{c2.on, x1, x2, y1, y2, z1, z2}
}

// Add returns true when cube can be removed
// func (c *Cube) Add(c2 Cube) bool {
// 	if c.on == c2.on && len(c.overlap) == 0 {
// 		return false
// 	}
// 	if c.IsWithin(c2) {
// 		c.on = c2.on
// 		c.overlap = c2.overlap
// 		return true
// 	}
// 	var del []int
// 	for i, oc := range c.overlap {
// 		if oc.Overlaps(c2) {
// 			d := c.overlap[i].Add(oc.Overlap(c2))
// 			if d {
// 				del = append(del, i)
// 			}
// 		}
// 	}
// 	c.overlap.Delete(del)
// 	if c.on != c2.on {
// 		c.overlap = append(c.overlap, c2)
// 	}
// 	return false
// }

func (c Cube) Count() int {
	return (c.x2 - c.x1 + 1) * (c.y2 - c.y1 + 1) * (c.z2 - c.z1 + 1)
}

func (c Cube) Break(c2 Cube) CubeList {
	if !c.Overlaps(c2) {
		return CubeList{c}
	}

	var b CubeList
	// above c2
	az1 := c2.z2 + 1
	az2 := c.z2
	if az2 >= az1 {
		b = append(b, Cube{c.on, c.x1, c.x2, c.y1, c.y2, az1, az2})
	}
	// below c2
	bz2 := c2.z1 - 1
	bz1 := c.z1
	if bz2 >= bz1 {
		b = append(b, Cube{c.on, c.x1, c.x2, c.y1, c.y2, bz1, bz2})
	}
	// add cubes around c2
	// front / back (y)
	z1 := max(c.z1, c2.z1)
	z2 := min(c.z2, c2.z2)
	if z2 >= z1 {
		ay1 := c2.y2 + 1
		ay2 := c.y2
		if ay2 >= ay1 {
			b = append(b, Cube{c.on, c.x1, c.x2, ay1, ay2, z1, z2})
		}
		by2 := c2.y1 - 1
		by1 := c.y1
		if by2 >= by1 {
			b = append(b, Cube{c.on, c.x1, c.x2, by1, by2, z1, z2})
		}

		// left/right (x)
		y1 := max(c.y1, c2.y1)
		y2 := min(c.y2, c2.y2)
		if y2 >= y1 {
			ax1 := c2.x2 + 1
			ax2 := c.x2
			if ax2 >= ax1 {
				b = append(b, Cube{c.on, ax1, ax2, y1, y2, z1, z2})
			}
			bx2 := c2.x1 - 1
			bx1 := c.x1
			if bx2 >= bx1 {
				b = append(b, Cube{c.on, bx1, bx2, y1, y2, z1, z2})
			}
		}
	}
	return b
}

func (cl CubeList) Count() int {
	var v int
	for _, c := range cl {
		v += c.Count()
	}
	return v
}

func (cl *CubeList) Add(c Cube) {
	if len(*cl) == 0 {
		if c.on {
			*cl = append(*cl, c)
			return
		}
		return
	}

	var del []int
	var addCubes CubeList
	for i, ci := range *cl {
		if !c.on {
			if ci.IsWithin(c) {
				del = append(del, i)
			} else {
				if ci.Overlaps(c) {
					addCubes = append(addCubes, ci.Break(c)...)
					del = append(del, i)
				}
			}
		} else {
			if c.IsWithin(ci) {
				continue
			}
			if ci.Overlaps(c) {
				addCubes = append(addCubes, ci.Break(c)...)
				del = append(del, i)
			}
		}
	}
	cl.Delete(del)
	*cl = append(*cl, addCubes...)
	if c.on {
		*cl = append(*cl, c)
	}
}

func (cl *CubeList) Delete(del []int) {
	if len(del) == 0 {
		return
	}

	var ol CubeList
	for i, ix := range del {
		if i > 0 {
			ol = append(ol, (*cl)[del[i-1]+1:ix]...)
		} else {
			ol = append(ol, (*cl)[:ix]...)
		}
	}
	ol = append(ol, (*cl)[del[len(del)-1]+1:]...)
	(*cl) = ol
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	var u Univ
	cubes := readCubes("input")
	var cubelist CubeList
	t0 := time.Now()
	for i, c := range cubes {
		log.Printf("setting cube: %d", i)
		u.Set(c)
		cubelist.Add(c)
		log.Print(cubelist.Count())
	}
	log.Print("-50-50 on: ", u.Count())
	log.Print("total count: ", cubelist.Count(), " cubes: ", len(cubelist))
	log.Print("time: ", time.Since(t0))
}

// func (t *Topo) Print() {
// 	my := len(t.dots)
// 	mx := len(t.dots[0])
// 	for j := 0; j < my; j++ {
// 		var s string
// 		for i := 0; i < mx; i++ {
// 			if t.dots[j][i] {
// 				s = s + "#"
// 			} else {
// 				s = s + "."
// 			}
// 		}
// 		log.Println(s)
// 	}
// }

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

// func (t *Topo) Get9(i, j int, blink bool) int {
// 	var pv int
// 	for y := j - 1; y <= j+1; y++ {
// 		for x := i - 1; x <= i+1; x++ {
// 			var bv int
// 			if t.Pixel(x, y, blink) {
// 				bv = 1
// 			}
// 			pv = 2*pv + bv
// 		}
// 	}
// 	// log.Printf("%d,%d: %9b", i, j, pv)
// 	return pv
// }

// func (t *Topo) Enhance(ie IE, step int) Topo {
// 	var nt Topo
// 	for j := -1; j < len(t.dots)+1; j++ {
// 		var l []bool
// 		for i := -1; i < len(t.dots[0])+1; i++ {
// 			var blink bool
// 			if ie[0] {
// 				blink = step%2 == 0
// 			}
// 			pv := t.Get9(i, j, blink)
// 			l = append(l, ie[pv])
// 		}
// 		nt.dots = append(nt.dots, l)
// 	}
// 	return nt
// }
