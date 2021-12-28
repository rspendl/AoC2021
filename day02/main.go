package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type Command struct {
	Cmd  string
	Unit int
}

func readData(fileName string) []Command {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	r.Comma = ' '
	var dataList []Command
	for {
		var p Command
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		p.Unit, err = strconv.Atoi(record[1])
		if err != nil {
			log.Fatal(err)
		}
		p.Cmd = record[0]
		dataList = append(dataList, p)
	}
	return dataList
}

func main() {
	data := readData("input")
	depth := 0
	aim := 0
	hor := 0
	for _, p := range data {
		switch p.Cmd {
		case "down":
			aim += p.Unit
		case "up":
			aim -= p.Unit
		case "forward":
			hor += p.Unit
			depth += p.Unit * aim
		}
	}
	log.Printf("hor: %d, depth1: %d, %d", hor, aim, hor*aim)
	log.Printf("hor: %d, depth: %d, %d", hor, depth, hor*depth)
}
