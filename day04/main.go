package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

var Draw = []int{26, 38, 2, 15, 36, 8, 12, 46, 88, 72, 32, 35, 64, 19, 5, 66, 20, 52, 74, 3, 59, 94, 45, 56, 0, 6, 67, 24, 97, 50, 92, 93, 84, 65, 71, 90, 96, 21, 87, 75, 58, 82, 14, 53, 95, 27, 49, 69, 16, 89, 37, 13, 1, 81, 60, 79, 51, 18, 48, 33, 42, 63, 39, 34, 62, 55, 47, 54, 23, 83, 77, 9, 70, 68, 85, 86, 91, 41, 4, 61, 78, 31, 22, 76, 40, 17, 30, 98, 44, 25, 80, 73, 11, 28, 7, 99, 29, 57, 43, 10}

const W = 5

type Board [][]int
type Boards []Board

func (bs *Boards) Mark(n int) {
	for i, b := range *bs {
		for j, bl := range b {
			for k, v := range bl {
				if v == n {
					(*bs)[i][j][k] = 0
				}
			}
		}
	}
}

func (bs *Boards) Bingo() (Boards, []int) {
	var bingo Boards
	var ix []int
	for i, b := range *bs {
		if b.IsBingo() {
			bingo = append(bingo, b)
			ix = append(ix, i)
		}
	}
	return bingo, ix
}

func (b *Board) IsBingo() bool {
	for i, bl := range *b {
		sl := 0
		sv := 0
		for j, v := range bl {
			sl += v
			sv += (*b)[j][i]
		}
		if sl == 0 || sv == 0 {
			return true
		}
	}
	return false
}

func (b *Board) Sum() int {
	s := 0
	for _, bl := range *b {
		for _, v := range bl {
			s += v
		}
	}
	return s
}

func readBoards(fileName string) Boards {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	r.Comma = ' '
	var (
		boardList Boards
	)
	var board Board
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var line []int
		for i := 0; i < W; i++ {
			n, err := strconv.ParseInt(record[i], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			line = append(line, int(n))
		}
		board = append(board, line)
		if len(board) == W {
			boardList = append(boardList, board)
			board = Board{}
		}
	}
	return boardList
}

func main() {
	boards := readBoards("input")
	for _, n := range Draw {
		boards.Mark(n)
		bingo, ix := boards.Bingo()
		if len(bingo) > 0 {
			nd := 0
			for j, board := range bingo {
				i := ix[j] - nd
				log.Printf("N: %d, sum: %d, last draw: %d, result: %d", len(bingo), board.Sum(), n, board.Sum()*n)
				boards = append(boards[:i], boards[i+1:]...)
				nd++
			}
		}
	}
}
