package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

// type Command struct {
// 	Cmd  string
// 	Unit int
// }

func readData(fileName string) ([]int, int) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	r.Comma = ' '
	var (
		dataList []int
		width    int
	)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		n, err := strconv.ParseInt(record[0], 2, 64)
		width = len(record[0])
		if err != nil {
			log.Fatal(err)
		}
		dataList = append(dataList, int(n))
	}
	return dataList, width
}

func main() {
	data, width := readData("input")
	var (
		gamma, eps int
	)
	for i := width - 1; i >= 0; i-- {
		var b0, b1 int
		for _, n := range data {
			bit := (n >> i) & 1
			if bit > 0 {
				b1++
			} else {
				b0++
			}
		}
		if b1 >= b0 {
			gamma = gamma | (1 << i)
		} else {
			eps = eps | (1 << i)
		}
	}
	// if gamma != eps {
	// 	log.Fatalf("error gamma: %b, eps: %b", gamma, eps)
	// }
	log.Printf("gamma: %012b, eps: %012b, %d", gamma, eps, gamma*eps)

	oxrate := make([]int, len(data))
	co2rate := make([]int, len(data))
	copy(oxrate, data)
	copy(co2rate, data)
	for i := width - 1; i >= 0; i-- {
		log.Print("ox", i, len(oxrate))
		var b0, b1 int
		for _, n := range oxrate {
			bit := (n >> i) & 1
			if bit > 0 {
				b1++
			} else {
				b0++
			}
		}
		log.Print("ox: b0 ", b0, " b1 ", b1)
		j := 0
		for {
			if (b1 >= b0 && oxrate[j]&(1<<i) == 0) || (b1 < b0 && oxrate[j]&(1<<i) > 0) {
				oxrate = append(oxrate[:j], oxrate[j+1:]...)
			} else {
				j++
			}

			if len(oxrate) == 1 || j == len(oxrate) {
				break
			}
		}
		if len(oxrate) == 1 {
			log.Printf("found oxrate: %012b", oxrate[0])
			break
		}
	}

	for i := width - 1; i >= 0; i-- {
		log.Print("co2: ", i, len(co2rate))
		var b0, b1 int
		for _, n := range co2rate {
			bit := (n >> i) & 1
			if bit > 0 {
				b1++
			} else {
				b0++
			}
		}
		log.Print("co2: b0 ", b0, " b1 ", b1)
		j := 0
		for {
			if (b1 >= b0 && co2rate[j]&(1<<i) > 0) || (b1 < b0 && co2rate[j]&(1<<i) == 0) {
				co2rate = append(co2rate[:j], co2rate[j+1:]...)
			} else {
				j++
			}
			if len(co2rate) == 1 || j == len(co2rate) {
				break
			}
		}
		if len(co2rate) == 1 {
			log.Printf("found co2rate: %012b", co2rate[0])
			break
		}

	}

	log.Printf("oxrate: %012b, eps: %012b, %d", oxrate[0], co2rate[0], oxrate[0]*co2rate[0])
}
