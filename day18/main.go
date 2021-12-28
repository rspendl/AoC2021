package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type Pair [2]int
type SN struct {
	Num   int
	Left  *SN
	Right *SN
	Top   *SN
}

func (sn *SN) IsLeaf() bool {
	return sn.Left == nil && sn.Right == nil
}

func (sn *SN) IsRoot() bool {
	return sn.Top == nil
}

func ParseSN(s string) SN {
	d := 0
	comma := -1
	maxDepth := 0
	for i, c := range s {
		switch c {
		case '[':
			d++
			if maxDepth < d {
				maxDepth = d
			}
		case ']':
			d--
		case ',':
			if d == 1 {
				comma = i
			}
		}
	}
	if d > 0 {
		log.Panic("invalid string ", s, d)
	}
	if maxDepth == 0 {
		n, err := strconv.Atoi(s)
		if err != nil {
			log.Panic("invalid number ", n)
		}
		return SN{Num: n}
	}
	lsn := ParseSN(s[1:comma])
	rsn := ParseSN(s[comma+1 : len(s)-1])
	sn := SN{Left: &lsn, Right: &rsn}
	sn.SetTops()
	return sn
}

func (sn *SN) SetTops() {
	if sn.Left != nil {
		sn.Left.Top = sn
		sn.Left.SetTops()
	}
	if sn.Right != nil {
		sn.Right.Top = sn
		sn.Right.SetTops()
	}
}

func Split(n int) string {
	return fmt.Sprintf("[%d,%d]", n/2, int(math.Round(float64(n)/2.0)))
}

func (sn *SN) Reduce() bool {
	s := sn.String()
	d := 0
	for i, c := range s {
		switch c {
		case '[':
			d++
			if d == 5 {
				ln := 0
				rn := 0
				isL := true
				j := i + 1
				closeBracket := false
				for !closeBracket {
					a := s[j]

					switch a {
					case ']':
						closeBracket = true
					case ',':
						isL = false
					default:
						if isL {
							ln = ln*10 + int(s[j]-'0')
						} else {
							rn = rn*10 + int(s[j]-'0')
						}
					}
					j++
				}
				lastBefore := -1
				lastBeforeLast := -1
				for i1 := i; i1 > 0; i1-- {
					if s[i1] >= '0' && s[i1] <= '9' {
						for i2 := i1; i2 > 0; i2-- {
							if s[i2] < '0' || s[i2] > '9' {
								lastBefore = i2 + 1
								lastBeforeLast = i1
								break
							}
						}
						break
					}
				}

				firstAfter := -1
				firstAfterLast := -1
				for j1 := j; j1 < len(s)-1; j1++ {
					if s[j1] >= '0' && s[j1] <= '9' {
						for j2 := j1; j2 < len(s)-1; j2++ {
							if s[j2] < '0' || s[j2] > '9' {
								firstAfter = j1
								firstAfterLast = j2 - 1
								break
							}
						}
						break
					}
				}

				var s1 string
				if lastBefore == -1 {
					s1 = s[:i]
				} else {
					lbV, err := strconv.Atoi(s[lastBefore : lastBeforeLast+1])
					if err != nil {
						log.Fatal("invalid before number ", s[lastBefore:lastBeforeLast+1])
					}
					s1 = s[:lastBefore] + fmt.Sprint(lbV+ln) + s[lastBeforeLast+1:i]
				}

				var s2 string
				if firstAfter == -1 {
					s2 = s[j:]
				} else {
					faV, err := strconv.Atoi(s[firstAfter : firstAfterLast+1])
					if err != nil {
						log.Fatal("invalid after number ", s[firstAfter:firstAfterLast+1])
					}
					s2 = s[j:firstAfter] + fmt.Sprint(faV+rn) + s[firstAfterLast+1:]
				}
				reducedStr := s1 + "0" + s2
				rsn := ParseSN(reducedStr)
				// log.Print("explode: ", reducedStr)
				rsn.Reduce()
				*sn = rsn
				return true
			}
		case ']':
			d--
		}
	}

	if d > 0 {
		log.Fatal("no ending bracket")
	}

	curN := 0
	curNpos := -1
	for i, c := range s {
		switch c {
		case '[':
			if curN > 9 {
				splitStr := s[:curNpos] + Split(curN) + s[i:]
				ssn := ParseSN(splitStr)
				// log.Print("split: ", splitStr)
				ssn.Reduce()
				*sn = ssn
				return true
			}
			curN = 0
			curNpos = -1
			d++
			if d == 5 {
				ln := 0
				rn := 0
				isL := true
				j := i + 1
				closeBracket := false
				for !closeBracket {
					a := s[j]

					switch a {
					case ']':
						closeBracket = true
					case ',':
						isL = false
					default:
						if isL {
							ln = ln*10 + int(s[j]-'0')
						} else {
							rn = rn*10 + int(s[j]-'0')
						}
					}
					j++
				}
				lastBefore := -1
				lastBeforeLast := -1
				for i1 := i; i1 > 0; i1-- {
					if s[i1] >= '0' && s[i1] <= '9' {
						for i2 := i1; i2 > 0; i2-- {
							if s[i2] < '0' || s[i2] > '9' {
								lastBefore = i2 + 1
								lastBeforeLast = i1
								break
							}
						}
						break
					}
				}

				firstAfter := -1
				firstAfterLast := -1
				for j1 := j; j1 < len(s)-1; j1++ {
					if s[j1] >= '0' && s[j1] <= '9' {
						for j2 := j1; j2 < len(s)-1; j2++ {
							if s[j2] < '0' || s[j2] > '9' {
								firstAfter = j1
								firstAfterLast = j2 - 1
								break
							}
						}
						break
					}
				}

				var s1 string
				if lastBefore == -1 {
					s1 = s[:i]
				} else {
					lbV, err := strconv.Atoi(s[lastBefore : lastBeforeLast+1])
					if err != nil {
						log.Fatal("invalid before number ", s[lastBefore:lastBeforeLast+1])
					}
					s1 = s[:lastBefore] + fmt.Sprint(lbV+ln) + s[lastBeforeLast+1:i]
				}

				var s2 string
				if firstAfter == -1 {
					s2 = s[j:]
				} else {
					faV, err := strconv.Atoi(s[firstAfter : firstAfterLast+1])
					if err != nil {
						log.Fatal("invalid after number ", s[firstAfter:firstAfterLast+1])
					}
					s2 = s[j:firstAfter] + fmt.Sprint(faV+rn) + s[firstAfterLast+1:]
				}
				reducedStr := s1 + "0" + s2
				rsn := ParseSN(reducedStr)
				// log.Print("explode: ", reducedStr)
				rsn.Reduce()
				*sn = rsn
				return true
			}
		case ']':
			if curN > 9 {
				splitStr := s[:curNpos] + Split(curN) + s[i:]
				ssn := ParseSN(splitStr)
				// log.Print("split: ", splitStr)
				ssn.Reduce()
				*sn = ssn
				return true
			}
			curN = 0
			curNpos = -1
			d--
		case ',':
			if curN > 9 {
				splitStr := s[:curNpos] + Split(curN) + s[i:]
				ssn := ParseSN(splitStr)
				// log.Print("split: ", splitStr)
				ssn.Reduce()
				*sn = ssn
				return true
			}
			curN = 0
			curNpos = -1
		default:
			if curNpos == -1 {
				curNpos = i
			}
			curN = curN*10 + int(c-'0')
		}
	}

	return false
}

func (sn *SN) Magnitude() int {
	if sn.Left == nil && sn.Right == nil {
		return sn.Num
	}
	return sn.Left.Magnitude()*3 + sn.Right.Magnitude()*2
}

func (sn *SN) String() string {
	if sn.Left == nil && sn.Right == nil {
		return fmt.Sprint(sn.Num)
	}
	return fmt.Sprintf("[%s,%s]", sn.Left.String(), sn.Right.String())
}

func readSN(fileName string) []SN {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	r.Comma = ';'
	var sns []SN
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		sn := ParseSN(record[0])
		sns = append(sns, sn)
	}
	return sns
}

func addSNs(sns []SN) SN {
	var sn SN
	for n, snx := range sns {
		if n == 0 {
			sn = snx
		} else {
			l := sn
			r := snx
			log.Print("adding: ", r.String())
			sn = SN{Left: &l, Right: &r}
			sn.Reduce()
		}
	}
	return sn
}

func main() {
	sns := readSN("input")
	sn := addSNs(sns)
	sn.Reduce()
	log.Print("magnitude: ", sn.Magnitude(), " ", sn.String())

	maxMag := 0
	for i, sn1 := range sns {
		for j, sn2 := range sns {
			if i != j {
				log.Print(i, "+", j)
				snX := SN{Left: &sn1, Right: &sn2}
				snX.Reduce()
				m := snX.Magnitude()
				if m > maxMag {
					maxMag = m
				}
			}
		}
	}
	log.Print("max magnitude: ", maxMag)
}
