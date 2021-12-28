package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type TID int

const (
	T_VERSION TID = iota
	T_TYPE
	T_LENTYPE
	T_LITERAL
	T_LENGTH
	T_SUB
	T_END
)

var counter int

type State struct {
	currentEl   TID // 0-versin, 1=type, 2=
	version     int
	typeID      int
	bitPosition int
	// lastGroup   bool
	lengthType  bool // false =15bit len; true = 11bit N-subpackets
	length      int
	packetNum   int
	packetLen   int
	packetStart int
	packetEnd   int
	literal     int
	lastLiteral bool
	literals    []int
	parent      *State
	id          int
}

func (s State) Value(ss []*State) int {
	log.Printf("value of %d", s.id)
	if s.typeID == 4 {
		var v int
		for _, l := range s.literals {
			v = v<<4 | l
		}
		return v
	}
	var childValues []int
	for i, cs := range ss {
		if cs.parent != nil && cs.parent.id == s.id {
			cv := cs.Value(ss[i+1:])
			log.Printf("getting value of child %d of parent %d = %d", cs.id, cs.parent.id, cv)
			childValues = append(childValues, cv)
		}
	}
	log.Printf("parent %d has %d children", s.id, len(childValues))
	var str string
	switch s.typeID {
	case 0:
		var v int
		for _, cv := range childValues {
			v += cv
			if str == "" {
				str = fmt.Sprint(cv)
			} else {
				str = fmt.Sprintf("%s+%d", str, cv)
			}
		}
		str = fmt.Sprintf("%s=%d", str, v)
		log.Print(str)
		return v
	case 1:
		v := 1
		for _, cv := range childValues {
			v *= cv
			if str == "" {
				str = fmt.Sprint(cv)
			} else {
				str = fmt.Sprintf("%s*%d", str, cv)
			}
		}
		str = fmt.Sprintf("%s=%d", str, v)
		log.Print(str)
		return v
	case 2:
		noMin := true
		var v int
		for _, cv := range childValues {
			if noMin || v > cv {
				v = cv
				if str == "" {
					str = fmt.Sprintf("min(%d", cv)
				} else {
					str = fmt.Sprintf("%s,%d", str, cv)
				}
				noMin = false
			}
		}
		str = fmt.Sprintf("%s)=%d", str, v)
		log.Print(str)
		return v
	case 3:
		v := 0
		for _, cv := range childValues {
			if v < cv {
				v = cv
				if str == "" {
					str = fmt.Sprintf("max(%d", cv)
				} else {
					str = fmt.Sprintf("%s,%d", str, cv)
				}
			}
		}
		str = fmt.Sprintf("%s)=%d", str, v)
		log.Print(str)
		return v
	case 5:
		if childValues[0] > childValues[1] {
			return 1
		}
		return 0
	case 6:
		if childValues[0] < childValues[1] {
			return 1
		}
		return 0
	case 7:
		if childValues[0] == childValues[1] {
			return 1
		}
		return 0
	}

	log.Fatal("unknown type", s.typeID)
	return 0
}

func readPacket(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	r := csv.NewReader(file)
	record, err := r.Read()
	// if err == io.EOF {
	// 	break
	// }
	if err != nil {
		log.Fatal(err)
	}
	var s string
	for _, c := range record[0] {
		if c >= '0' && c <= '9' {
			s = s + fmt.Sprintf("%04b", int(c-'0'))
		} else {
			s = s + fmt.Sprintf("%04b", int(c-'A')+10)
		}
	}
	return s
}

func main() {
	packet := readPacket("input")
	var stateStack []*State
	var parentList []*State

	log.Print(packet)
	counter++
	stateStack = append(stateStack, &State{
		id:          counter,
		packetStart: 0,
	})
	st := stateStack[len(stateStack)-1]
	for pos, b := range packet {
		var bVal int
		if b == '1' {
			bVal = 1
		}
		st.packetLen++
		switch st.currentEl {
		case T_VERSION: // version
			st.version = st.version<<1 | bVal
			if st.bitPosition < 2 {
				st.bitPosition++
			} else {
				st.bitPosition = 0
				st.currentEl = T_TYPE
			}

		case T_TYPE: //type
			st.typeID = st.typeID<<1 | bVal
			if st.bitPosition < 2 {
				st.bitPosition++
			} else {
				st.bitPosition = 0
				if st.typeID == 4 { // literal
					st.currentEl = T_LITERAL
				} else {
					st.currentEl = T_LENTYPE
				}
			}
		case T_LENTYPE:
			st.lengthType = bVal == 1
			st.bitPosition = 0
			st.currentEl = T_LENGTH
		case T_LENGTH:
			st.length = st.length<<1 | bVal
			if (!st.lengthType && st.bitPosition < 14) || (st.lengthType && st.bitPosition < 10) {
				st.bitPosition++
			} else {
				st.bitPosition = 0
				st.currentEl = T_SUB
				st.packetEnd = pos
				parentList = append(parentList, stateStack[len(stateStack)-1])
				log.Printf("added parent %d parentList %v", parentList[len(parentList)-1].id, parentList)
				counter++
				stateStack = append(stateStack, &State{
					id:          counter,
					packetStart: pos + 1,
					parent:      parentList[len(parentList)-1],
				})
				st = stateStack[len(stateStack)-1]
			}
		case T_LITERAL:
			if st.bitPosition == 0 {
				st.lastLiteral = bVal == 0
				st.bitPosition++
				continue
			}
			st.literal = st.literal<<1 | bVal
			if st.bitPosition < 4 {
				st.bitPosition++
			} else {
				st.literals = append(st.literals, st.literal)
				st.packetEnd = pos
				st.bitPosition = 0
				if st.lastLiteral {
					// stateStack = append(stateStack, st)
					if st.parent == nil {
						st.currentEl = T_END
						continue
						// } else {
						// 	st.parent.packetNum++
					}
					par := st.parent
					for par != nil {
						par.packetLen += st.packetLen
						par.packetNum++
						if /*(!par.lengthType && par.packetLen == par.length) ||*/
						(par.lengthType && par.packetNum == par.length) ||
							(!par.lengthType && par.length == (pos-par.packetEnd)) {
							log.Printf("%d: removed parent %d parentList %v", pos, parentList[len(parentList)-1].id, parentList)
							parentList = parentList[:len(parentList)-1]
							log.Printf("parentList %v", parentList)
							// if st.parent.parent == nil {
							// st = st.parent
							// if st.parent != nil {
							// 	st.parent.packetNum++
							// } else {
							// 	st.currentEl = T_END
							// }
							// continue
							// }
							// st = st.parent.parent
							// counter++
							// stateStack = append(stateStack, &State{
							// 	id:     counter,
							// 	parent: st.parent.parent,
							// })
							// st = stateStack[len(stateStack)-1]
							// continue
						} else {
							break
						}
						par = par.parent
					}
					// parent := st.parent
					if par != nil {
						counter++
						var parent *State
						if len(parentList) > 0 {
							parent = parentList[len(parentList)-1]
						}
						stateStack = append(stateStack, &State{
							id:          counter,
							packetStart: pos + 1,
							parent:      parent,
						})
						st = stateStack[len(stateStack)-1]
					} else {
						st = stateStack[0]
						st.currentEl = T_END
					}
				}
			}
		case T_END:
			st.bitPosition++ // just read till the end
		}
	}
	if stateStack[len(stateStack)-1].typeID == 0 {
		stateStack = stateStack[:len(stateStack)-1]
	}
	ssStr := ""
	for _, ss := range stateStack {
		ssStr += fmt.Sprintf("{%p, %d, %d, (%d-%d) %v}", ss, ss.typeID, ss.id, ss.packetStart+1, ss.packetEnd+1, *ss)
	}
	log.Print(ssStr)
	var vs int
	for _, ss := range stateStack {
		vs += ss.version
		// log.Printf("version: %d, value: %d", ss.version, ss.Value(stateStack))
	}
	log.Print("version sum:", vs, " value: ", stateStack[0].Value(stateStack))
}
