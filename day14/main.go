package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"time"
)

type Rule struct {
	rule    [2]string
	applied int
}

type Rules []Rule

func (rs *Rules) From(fs [][2]string) {
	var ar Rules
	for _, r := range fs {
		ar = append(ar, Rule{
			rule:    r,
			applied: 0,
		})
	}
	(*rs) = ar
}

func (rs *Rules) Set(s string) {
	for i := range s[:len(s)-1] {
		ri := rs.Index(s[i : i+2])
		(*rs)[ri].applied++
	}
}

func (rs *Rules) Index(s string) int {
	for i, r := range *rs {
		if r.rule[0] == s {
			return i
		}
	}
	return -1
}

func (rs *Rules) Step() {
	var newRules Rules
	for i, r := range *rs {
		if r.applied > 0 {
			r1, r2 := rs.Spawn(r)
			newRules = append(newRules, r1, r2)
			(*rs)[i].applied = 0
		}
	}
	for _, nr := range newRules {
		for i, r := range *rs {
			if r.rule[0] == nr.rule[0] {
				(*rs)[i].applied += nr.applied
			}
		}
	}
}

func (rs *Rules) Spawn(r0 Rule) (Rule, Rule) {
	s1 := string(r0.rule[0][0]) + r0.rule[1]
	s2 := r0.rule[1] + string(r0.rule[0][1])
	var r1, r2 Rule
	for _, r := range *rs {
		if s1 == r.rule[0] {
			r1 = Rule{
				rule:    [2]string{s1, r.rule[1]},
				applied: r0.applied,
			}
		}
		if s2 == r.rule[0] {
			r2 = Rule{
				rule:    [2]string{s2, r.rule[1]},
				applied: r0.applied,
			}
		}
	}
	return r1, r2
}

func (rs *Rules) Count(last rune) map[rune]int {
	n := make(map[rune]int)
	for _, r := range *rs {
		for i, c := range r.rule[0] {
			if i == 0 {
				n[c] += r.applied
			}
		}
	}
	n[last]++
	return n
}

func readRules(fileName string) [][2]string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	r.Comma = ' '
	var ins [][2]string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		ins = append(ins, [2]string{record[0], record[2]})
	}
	return ins
}

func insert(s string, rules [][2]string) string {
	for _, r := range rules {
		if s == r[0] {
			return string(s[0]) + r[1] + string(s[1])
		}
	}
	return s
}

func replace(s string, rules [][2]string) string {
	var out string
	for i := range s[:len(s)-1] {
		if len(out) > 0 {
			out = out[:len(out)-1] + insert(s[i:i+2], rules)
		} else {
			out = out + insert(s[i:i+2], rules)
		}
	}
	return out
}

func main() {
	rules := readRules("input")
	var ar Rules
	ar.From(rules)
	// ar.Set("NNCB")
	ar.Set("PSVVKKCNBPNBBHNSFKBO")
	// poly := "NNCB"
	poly := "PSVVKKCNBPNBBHNSFKBO"
	var last rune
	for i, c := range poly {
		if i == len(poly)-1 {
			last = c
		}
	}
	t0 := time.Now()
	for i := 0; i < 40; i++ {
		// poly = replace(poly, rules)
		ar.Step()
		log.Print(i, ": ", time.Since(t0))
		// log.Print(poly)
		// log.Print(ar)
		// log.Print(ar.Count(last))
	}
	// sort.Slice(poly, func(i, j int) bool { return poly[i] < poly[j] })
	// n := make(map[rune]int)
	n := ar.Count(last)
	var first rune
	for i, c := range ar[0].rule[0] {
		if i == 0 {
			first = c
		}
		// n[c]++
	}
	min := n[first]
	max := min
	minC := first
	maxC := first
	for j := range n {
		if n[j] < min {
			min = n[j]
			minC = j
		}
		if n[j] > max {
			max = n[j]
			maxC = j
		}
	}
	log.Printf("(%c, %d) - (%c, %d) = %d", maxC, max, minC, min, max-min)
}
