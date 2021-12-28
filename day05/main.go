package main

import (
	"encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type Line struct {
	x1, y1, x2, y2 int
}

type World [][]int

const WS = 1000

func (w *World) AddLine(line Line, includeDiagonal bool) bool {
	if line.x1 == line.x2 {
		j := line.y1
		for {
			(*w)[line.x1][j]++
			if j == line.y2 {
				break
			}
			if line.y1 < line.y2 {
				j++
			} else {
				j--
			}
		}
		return true
	}
	if line.y1 == line.y2 {
		i := line.x1
		for {
			(*w)[i][line.y1]++
			if i == line.x2 {
				break
			}
			if line.x1 < line.x2 {
				i++
			} else {
				i--
			}
		}
		return true
	}
	if includeDiagonal && math.Abs(float64(line.y1-line.y2)) == math.Abs(float64(line.x1-line.x2)) {
		i := line.x1
		j := line.y1
		for {
			(*w)[i][j]++
			if i == line.x2 {
				break
			}
			if line.x1 < line.x2 {
				i++
			} else {
				i--
			}
			if line.y1 < line.y2 {
				j++
			} else {
				j--
			}
		}
		return true
	}
	return false
}

func (w *World) CountMore(min int) int {
	var c int
	for _, l := range *w {
		for _, v := range l {
			if v > min {
				c++
			}
		}
	}
	return c
}

func readLines(fileName string) []Line {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	var lineList []Line
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		x1, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		y1, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		x2, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		y2, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		lineList = append(lineList, Line{int(x1), int(y1), int(x2), int(y2)})
	}
	return lineList
}

func main() {
	lines := readLines("input")
	world := make(World, WS)
	for i := range world {
		world[i] = make([]int, WS)
	}
	lc := 0
	for _, l := range lines {
		if world.AddLine(l, false) {
			lc++
		}
	}
	log.Printf(">1 vert-hor: %d, line-count: %d", world.CountMore(1), lc)
	for i := range world {
		world[i] = make([]int, WS)
	}
	lc = 0
	for _, l := range lines {
		if world.AddLine(l, true) {
			lc++
		}
	}
	log.Printf(">1 vert-hor-diag: %d, line-count: %d", world.CountMore(1), lc)
}
