package main

import (
	"encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Beacon [3]int

type Scanner struct {
	id      int
	pos     Beacon
	pos0    []Beacon
	orient  int
	orient0 []int
	beacons []Beacon
	hasPos  bool
}

func readScanners(fileName string) []Scanner {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)

	var ss []Scanner
	var beacons []Beacon
	var n int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasPrefix(record[0], "---") {
			s := strings.Split(record[0], " ")
			nn, err := strconv.Atoi(s[2])
			if err != nil {
				log.Fatal("invalid number", s[2])
			}
			if nn > 0 {
				s := Scanner{
					id:      n,
					beacons: beacons,
				}
				ss = append(ss, s)
			}
			n = nn
			beacons = []Beacon{}
		} else {
			if record[0] != "" {
				x, _ := strconv.Atoi(record[0])
				y, _ := strconv.Atoi(record[1])
				z, _ := strconv.Atoi(record[2])
				beacons = append(beacons, [3]int{x, y, z})
			}
		}
	}
	s := Scanner{
		id:      n,
		beacons: beacons,
	}
	ss = append(ss, s)
	return ss
}

func inverseRotate(o int) int {
	var aa, ab, ag int // rz, ry, rx
	og := o % 4
	ag = og * 90
	if ag == 90 || ag == 270 {
		ag += 180
		if ag >= 360 {
			ag -= 360
		}
	}
	ob := (o / 4) % 4
	ab = ob * 90
	if ab == 90 || ab == 270 {
		ab += 180
		if ab >= 360 {
			ab -= 360
		}
	}

	oa := (o / 16) % 4
	aa = oa * 90
	if aa == 90 || ag == 270 {
		aa += 180
		if aa >= 360 {
			aa -= 360
		}
	}
	return (aa/90)*16 + (ab/90)*4 + (ag / 90)
}

func addRotate(o1, o2 int) int {
	o1g := o1 % 4
	o1b := (o1 / 4) % 4
	o1a := (o1 / 16) % 4

	o2g := o2 % 4
	o2b := (o2 / 4) % 4
	o2a := (o2 / 16) % 4

	og := o1g + o2g
	ob := o1b + o2b
	oa := o1a + o2a

	if og >= 4 {
		og -= 4
	}
	if ob >= 4 {
		ob -= 4
	}
	if oa >= 4 {
		oa -= 4
	}
	o := oa*16 + ob*4 + og
	// log.Printf("orient: %d + %d = %d", o1, o2, o)
	return o
}

func (s Scanner) RotateBeacons(o int) []Beacon {
	var a [3][3]int
	var sa, sb, sg, ca, cb, cg int
	var aa, ab, ag int // rz, ry, rx

	og := o % 4
	ag = og * 90

	ob := (o / 4) % 4
	ab = ob * 90

	oa := (o / 16) % 4
	aa = oa * 90

	ca = int(math.Cos(float64(aa) * math.Pi / 180))
	sa = int(math.Sin(float64(aa) * math.Pi / 180))
	cb = int(math.Cos(float64(ab) * math.Pi / 180))
	sb = int(math.Sin(float64(ab) * math.Pi / 180))
	cg = int(math.Cos(float64(ag) * math.Pi / 180))
	sg = int(math.Sin(float64(ag) * math.Pi / 180))
	a = [3][3]int{
		{ca * cb, ca*sb*sg - sa*cg, ca*sb*cg + sa*sg},
		{sa * cb, sa*sb*sg + ca*cg, sa*sb*cg - ca*sg},
		{-sb, cb * sg, cb * cg},
	}
	var rb []Beacon
	for _, b := range s.beacons {
		b1 := [3]int{
			b[0]*a[0][0] + b[1]*a[0][1] + b[2]*a[0][2],
			b[0]*a[1][0] + b[1]*a[1][1] + b[2]*a[1][2],
			b[0]*a[2][0] + b[1]*a[2][1] + b[2]*a[2][2],
		}
		rb = append(rb, b1)
	}
	return rb
}

func rotatedIdentity(o int) Beacon {
	is := Scanner{beacons: []Beacon{{1, 1, 1}}}
	return is.RotateBeacons(o)[0]
}

func (s Scanner) Orient(s0 Scanner) (Beacon, int, int) {
	type BCount struct {
		beacon Beacon
		count  int
	}
	var mo, mc int
	var mb Beacon
	for o := 0; o < 4*4*4; o++ {
		tb := s.RotateBeacons(o)
		var db []Beacon
		for _, b := range tb {
			for _, sb := range s0.beacons {
				db = append(db, Beacon{
					sb[0] - b[0],
					sb[1] - b[1],
					sb[2] - b[2],
				})
			}
		}
		var dbc []BCount
		for _, b := range db {
			found := false
			for i, dc := range dbc {
				if dc.beacon[0] == b[0] && dc.beacon[1] == b[1] && dc.beacon[2] == b[2] {
					dbc[i].count++
					found = true
					break
				}
			}
			if !found {
				dbc = append(dbc, BCount{beacon: b, count: 1})
			}
		}
		sort.Slice(dbc, func(i, j int) bool { return dbc[i].count > dbc[j].count })
		if dbc[0].count > mc {
			mc = dbc[0].count
			mo = o
			mb = dbc[0].beacon
		}
	}
	log.Printf("o: %d, s: %d/%d, maxC: %d", mo, s.id, s0.id, mc)
	return mb, mo, mc
}

func (s Scanner) TransBeacons(s0 Scanner) ([]Beacon, Beacon, []Beacon, int, []int) {
	pos, o, count := s.Orient(s0)
	if count < 12 {
		return []Beacon{}, Beacon{}, []Beacon{}, 0, []int{}
	}

	rb := s.RotateBeacons(o)

	// iv := rotatedIdentity(o)
	// dx := iv[0]
	// dy := iv[1]
	// dz := iv[2]
	// dx := -1
	// if o%4 > 1 {
	// 	dx = 1
	// }
	// dy := -1
	// if (o/4)%4 > 1 {
	// 	dy = 1
	// }
	// dz := -1
	// if (o/16)%4 > 1 {
	// 	dz = 1
	// }
	var trB []Beacon
	for _, b := range rb {
		trB = append(trB, Beacon{
			b[0] + pos[0],
			b[1] + pos[1],
			b[2] + pos[2],
		})
	}

	sb := Scanner{
		beacons: trB,
	}

	// dx = -1
	// if s0.orient%4 > 1 {
	// 	dx = 1
	// }
	// dy = -1
	// if (s0.orient/4)%4 > 1 {
	// 	dy = 1
	// }
	// dz = -1
	// if (s0.orient/16)%4 > 1 {
	// 	dz = 1
	// }

	for k := 0; k < len(s0.orient0); k++ {
		sb.beacons = trB
		rb = sb.RotateBeacons(s0.orient0[k])
		// iv = rotatedIdentity(s0.orient)
		// dx = iv[0]
		// dy = iv[1]
		// dz = iv[2]

		trB = []Beacon{}
		for _, b := range rb {
			trB = append(trB, Beacon{
				b[0] + s0.pos0[k][0],
				b[1] + s0.pos0[k][1],
				b[2] + s0.pos0[k][2],
			})
		}
	}

	// sb.beacons = trB
	// io = inverseRotate(s0.orient)
	// rb = sb.RotateBeacons(io)

	// trB = []Beacon{}
	// for _, b := range rb {
	// 	trB = append(trB, Beacon{
	// 		b[0] + s0.pos[0],
	// 		b[1] + s0.pos[1],
	// 		b[2] + s0.pos[2],
	// 	})
	// }
	p0 := []Beacon{pos}
	for _, ho := range s0.orient0 {
		sp := Scanner{
			beacons: p0,
		}
		p0 = sp.RotateBeacons(ho)
	}
	return trB, Beacon{
		s0.pos[0] + p0[0][0],
		s0.pos[1] + p0[0][1],
		s0.pos[2] + p0[0][2],
	}, append([]Beacon{pos}, s0.pos0...), addRotate(o, s0.orient), append([]int{o}, s0.orient0...)
}

func main() {
	ss := readScanners("input")
	var beacons []Beacon
	beacons = append(beacons, ss[0].beacons...)
	log.Print("added: ", len(beacons))
	// beacons := []Beacon{
	// 	{-892, 524, 684},
	// 	{-876, 649, 763},
	// 	{-838, 591, 734},
	// 	{-789, 900, -551},
	// 	{-739, -1745, 668},
	// 	{-706, -3180, -659},
	// 	{-697, -3072, -689},
	// 	{-689, 845, -530},
	// 	{-687, -1600, 576},
	// 	{-661, -816, -575},
	// 	{-654, -3158, -753},
	// 	{-635, -1737, 486},
	// 	{-631, -672, 1502},
	// 	{-624, -1620, 1868},
	// 	{-620, -3212, 371},
	// 	{-618, -824, -621},
	// 	{-612, -1695, 1788},
	// 	{-601, -1648, -643},
	// 	{-584, 868, -557},
	// 	{-537, -823, -458},
	// 	{-532, -1715, 1894},
	// 	{-518, -1681, -600},
	// 	{-499, -1607, -770},
	// 	{-485, -357, 347},
	// 	{-470, -3283, 303},
	// 	{-456, -621, 1527},
	// 	{-447, -329, 318},
	// 	{-430, -3130, 366},
	// 	{-413, -627, 1469},
	// 	{-345, -311, 381},
	// 	{-36, -1284, 1171},
	// 	{-27, -1108, -65},
	// 	{7, -33, -71},
	// 	{12, -2351, -103},
	// 	{26, -1119, 1091},
	// 	{346, -2985, 342},
	// 	{366, -3059, 397},
	// 	{377, -2827, 367},
	// 	{390, -675, -793},
	// 	{396, -1931, -563},
	// 	{404, -588, -901},
	// 	{408, -1815, 803},
	// 	{423, -701, 434},
	// 	{432, -2009, 850},
	// 	{443, 580, 662},
	// 	{455, 729, 728},
	// 	{456, -540, 1869},
	// 	{459, -707, 401},
	// 	{465, -695, 1988},
	// 	{474, 580, 667},
	// 	{496, -1584, 1900},
	// 	{497, -1838, -617},
	// 	{527, -524, 1933},
	// 	{528, -643, 409},
	// 	{534, -1912, 768},
	// 	{544, -627, -890},
	// 	{553, 345, -567},
	// 	{564, 392, -477},
	// 	{568, -2007, -577},
	// 	{605, -1665, 1952},
	// 	{612, -1593, 1893},
	// 	{630, 319, -379},
	// 	{686, -3108, -505},
	// 	{776, -3184, -501},
	// 	{846, -3110, -434},
	// 	{1135, -1161, 1235},
	// 	{1243, -1093, 1063},
	// 	{1660, -552, 429},
	// 	{1693, -557, 386},
	// 	{1735, -437, 1738},
	// 	{1749, -1800, 1813},
	// 	{1772, -405, 1572},
	// 	{1776, -675, 371},
	// 	{1779, -442, 1789},
	// 	{1780, -1548, 337},
	// 	{1786, -1538, 337},
	// 	{1847, -1591, 415},
	// 	{1889, -1729, 1762},
	// 	{1994, -1805, 1792},
	// }
	ss[0].hasPos = true
	changed := false
	for {
		for i := range ss {
			// if i == len(ss)-1 {
			// 	break
			// }
			// if !s1.hasPos {
			// 	continue
			// }
			for {
				addedPos := false
				j := 0
				for j < len(ss) {
					if i == j {
						j++
						continue
					}

					if ss[j].id == ss[i].id {
						j++
						continue
					}
					var (
						tb      []Beacon
						pos     Beacon
						pos0    []Beacon
						orient  int
						orient0 []int
					)
					if ss[i].hasPos && !ss[j].hasPos {
						tb, pos, pos0, orient, orient0 = ss[j].TransBeacons(ss[i])
						if len(tb) > 0 {
							ss[j].pos = pos
							ss[j].pos0 = pos0
							ss[j].orient = orient
							ss[j].orient0 = orient0
							ss[j].hasPos = true
							addedPos = true
							changed = true
							log.Printf("added position %d/%d: [%d %d %d]", ss[j].id, ss[i].id, pos[0], pos[1], pos[2])
						}
					}
					if !ss[i].hasPos && ss[j].hasPos {
						tb, pos, pos0, orient, orient0 = ss[i].TransBeacons(ss[j])
						if len(tb) > 0 {
							ss[i].pos = pos
							ss[i].pos0 = pos0
							ss[i].orient = orient
							ss[i].orient0 = orient0
							ss[i].hasPos = true
							addedPos = true
							changed = true
							log.Printf("added position %d/%d: [%d %d %d]", ss[i].id, ss[j].id, pos[0], pos[1], pos[2])
						}
					}

					l0 := len(beacons)
					for _, b := range tb {
						found := false
						for _, xb := range beacons {
							if b[0] == xb[0] && b[1] == xb[1] && b[2] == xb[2] {
								// log.Printf("dupl %d,%d,%d", b[0], b[1], b[2])
								found = true
								break
								// d++
							}
						}
						if !found {
							beacons = append(beacons, b)
						}
					}
					if len(tb) > 0 {
						log.Print("added: ", len(beacons)-l0)
					}
					j++
				}
				if !addedPos {
					break
				}
			}
		}
		if !changed {
			break
		} else {
			changed = false
		}
	}
	log.Print("beacons: ", len(beacons))
	mdd := 0
	for _, s1 := range ss {
		for _, s2 := range ss {
			b1 := s1.pos
			b2 := s2.pos
			d := int(math.Abs(float64(b1[0]-b2[0])) + math.Abs(float64(b1[1]-b2[1])) + math.Abs(float64(b1[2]-b2[2])))
			if d > mdd {
				mdd = d
			}
		}
	}
	log.Print("max dist: ", mdd)
}
