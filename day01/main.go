package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

func readData(fileName string) []int {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)

	var dataList []int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		v, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal(err)
		}
		dataList = append(dataList, v)
	}
	return dataList
}

func main() {
	data := readData("input")
	increases, i3 := 0, 0
	for i, v := range data {
		if i == 0 {
			continue
		}
		if v > data[i-1] {
			increases++
		}
		if i < 3 {
			continue
		}
		if data[i-3]+data[i-2]+data[i-1] < data[i-2]+data[i-1]+data[i] {
			i3++
		}
	}
	log.Printf("%d increases, %d 3-measeurement increases", increases, i3)
}
