package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
)

func readChunks(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	var chunks []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		chunks = append(chunks, record[0])
	}
	return chunks
}

func closes(a, b rune) bool {
	switch b {
	case '}':
		return a == '{'
	case '>':
		return a == '<'
	case ']':
		return a == '['
	case ')':
		return a == '('
	}
	return false
}

func main() {
	chunks := readChunks("input")
	var sum int
	var scores []int
	for _, chunk := range chunks {
		var (
			cur []rune
		)
		count := make(map[rune]int, 4)
		var errFound bool
		for _, c := range chunk {
			errFound = false
			switch c {
			case '{', '<', '[', '(':
				cur = append(cur, c)
				if _, ok := count[c]; !ok {
					count[c] = 0
				}
				count[c]++
			case '}', '>', ']', ')':
				if !closes(cur[len(cur)-1], c) {
					log.Printf("chunk: %s, error: %c", chunk, c)
					switch c {
					case ')':
						sum += 3
					case ']':
						sum += 57
					case '}':
						sum += 1197
					case '>':
						sum += 25137
					}
					errFound = true
				}
				cur = cur[:len(cur)-1]
			}
			if errFound {
				break
			}
		}
		if !errFound {
			// Complete the line
			log.Printf("complete: %v", string(cur))
			var score int
			for i := range cur {
				c := cur[len(cur)-1-i]
				switch c {
				case '(':
					score = 5*score + 1
				case '[':
					score = 5*score + 2
				case '{':
					score = 5*score + 3
				case '<':
					score = 5*score + 4
				default:
					log.Panic(cur)
				}
			}
			log.Printf("score: %d", score)
			scores = append(scores, score)
		}
	}
	sort.Slice(scores, func(i, j int) bool { return scores[i] < scores[j] })
	log.Printf("sum: %d, score: %d", sum, scores[len(scores)/2])
}
