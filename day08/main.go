package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

func readDigits(fileName string) (input [][10]string, output [][4]string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	r.Comma = ' '
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var (
			inL  [10]string
			outL [4]string
		)
		out := false
		for i, r := range record {
			if out {
				outL[i-11] = r
				continue
			}
			if r == "|" {
				out = true
				continue
			}
			inL[i] = r
		}
		input = append(input, inL)
		output = append(output, outL)
	}
	return
}

func includes(bigString, incS string) bool {
	for _, c := range incS {
		if !strings.ContainsRune(bigString, c) {
			return false
		}
	}
	return true
}

func diff(bigString, incS string) string {
	var s string
	for _, c := range bigString {
		if !includes(incS, string(c)) {
			s = s + string(c)
		}
	}
	return s
}

func index(list []string, s string) int {
	for i, l := range list {
		if len(l) == len(s) && includes(l, s) {
			return i
		}
	}
	return -1
}
func main() {
	input, output := readDigits("input")
	unique := 0
	for _, outL := range output {
		for _, o := range outL {
			l := len(o)
			if l == 2 || l == 3 || l == 4 || l == 7 {
				unique++
			}
		}
	}

	var sum int
	for j, inL := range input {
		digits := make([]string, 10)
		for _, i := range inL {
			l := len(i)
			switch l {
			case 2: // 1
				digits[1] = i
			case 3: // 7
				digits[7] = i
			case 4: // 4
				digits[4] = i
			case 7: // 8
				digits[8] = i
			}
		}
		for _, i := range inL {
			l := len(i)
			switch l {
			case 5: // 	2, 3, 5
				// 3 includes 1
				if includes(i, digits[1]) {
					digits[3] = i
				} else {
					d := diff(digits[4], digits[1]) // 2 segments difference between 4 and 1 are included in 5, but not in 2
					if includes(i, d) {
						digits[5] = i
					} else {
						digits[2] = i
					}
				}
			case 6: // 0, 6, 9
				// 6 does not include 1
				if !includes(i, digits[1]) {
					digits[6] = i
				} else {
					if includes(i, digits[4]) { // 9 includes 4, but 0 doesn't
						digits[9] = i
					} else {
						digits[0] = i
					}
				}
			}
		}
		var outN int
		for _, o := range output[j] {
			outN = 10*outN + index(digits, o)
		}
		// log.Printf("output: %d", outN)
		sum += outN
	}
	log.Printf("unique: %d, sum: %d", unique, sum)
}
