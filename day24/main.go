package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Command int

const (
	INP Command = iota
	ADD
	MUL
	DIV
	MOD
	EQL
)

type Argument string

type Op struct {
	cmd Command
	a   Argument
	b   Argument
}

type State struct {
	w, x, y, z int
}

func readProgram(fileName string) []Op {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	var opList []Op
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		record := strings.Split(rec[0], " ")
		op := Op{}
		switch record[0] {
		case "inp":
			op.cmd = INP
		case "add":
			op.cmd = ADD
		case "mul":
			op.cmd = MUL
		case "div":
			op.cmd = DIV
		case "mod":
			op.cmd = MOD
		case "eql":
			op.cmd = EQL
		default:
			log.Fatal("invalid operation: ", record[0])
		}
		op.a = Argument(record[1])
		if len(record) > 2 {
			op.b = Argument(record[2])
		}
		opList = append(opList, op)
	}
	return opList
}

func (s *State) Arg(arg Argument) *int {
	switch arg {
	case "w":
		return &s.w
	case "x":
		return &s.x
	case "y":
		return &s.y
	case "z":
		return &s.z
	default:
		log.Fatal("invalid argument: ", arg)
		return nil
	}
}

func (s *State) Val(arg Argument) int {
	switch arg {
	case "w":
		return s.w
	case "x":
		return s.x
	case "y":
		return s.y
	case "z":
		return s.z
	default:
		val, err := strconv.Atoi(string(arg))
		if err != nil {
			log.Fatal("invalid argument: ", arg)
		}
		return val
	}
}

func (s *State) Store(arg Argument, val int) {
	*(s.Arg(arg)) = val
}

func (s *State) Add(arg Argument, b Argument) {
	v1 := s.Val(arg)
	v2 := s.Val(b)
	s.Store(arg, v1+v2)
}

func (s *State) Mult(arg Argument, b Argument) {
	v1 := s.Val(arg)
	v2 := s.Val(b)
	s.Store(arg, v1*v2)
}

func (s *State) Div(arg Argument, b Argument) {
	v1 := s.Val(arg)
	v2 := s.Val(b)
	s.Store(arg, v1/v2)
}

func (s *State) Mod(arg Argument, b Argument) {
	v1 := s.Val(arg)
	v2 := s.Val(b)
	s.Store(arg, v1%v2)
}

func (s *State) Equal(arg Argument, b Argument) {
	v1 := s.Val(arg)
	v2 := s.Val(b)
	var eql int
	if v1 == v2 {
		eql = 1
	}
	s.Store(arg, eql)
}

func run(program []Op, input string) int {
	state := State{}
	inp := 0
	for line, op := range program {
		// if line == 96 {
		// 	state.x = state.w
		// }
		switch op.cmd {
		case INP:
			val := int(input[inp] - '0')
			inp++
			state.Store(op.a, val)
		case ADD:
			state.Add(op.a, op.b)
		case MUL:
			state.Mult(op.a, op.b)
		case DIV:
			state.Div(op.a, op.b)
		case MOD:
			state.Mod(op.a, op.b)
		case EQL:
			state.Equal(op.a, op.b)
		default:
			log.Fatal("invalid command: ", op.cmd, " in line: ", line+1)
		}
		if ((line+1)%18 == 0 || (line+1)%18 == 0) && line < 350 { //line == 17 || line == 35 || line == 54 || line == 89 || line == 107 {
			log.Printf("%d: %d,%d,%d,%d z%%26=%d", line+1, state.w, state.x, state.y, state.z, state.z%26)
			zm := state.z
			for zm > 0 {
				log.Printf("%d %c", zm%26, 'A'+zm%26)
				zm = zm / 26
			}
		}
	}
	return state.z
}

func randStr(l int) string {
	var s string
	for i := 0; i < l; i++ {
		x := rand.Intn(8)
		s = fmt.Sprintf("%s%d", s, x+1)
	}
	return s
}

func main() {
	program := readProgram("input")
	input := "13579246899999"
	rand.Seed(time.Now().UnixMilli())
	for i := 0; i < 1; i++ {
		input = "92969593497992" // max
		input = "81514171161381"
		z := run(program, input)
		log.Print(input, " , ", z)
	}
}
