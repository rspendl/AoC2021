package main

import (
	"container/heap"
	"log"
	"sort"
	"sync"
)

type Location int

const PodDepth = 4

type Pod [PodDepth]rune
type Hallway [11]rune

type Burrow struct {
	A, B, C, D Pod     // Location 1,2/3,4/5,6/7,8
	Hall       Hallway // Location 11-21
	Cost       int
}

type Creatures struct {
	creatures []Creature
	lastMoved int
	history   []Creatures
	cost      int
}

type Creature struct {
	Type   rune
	Cur    Location
	SpentE int
	NMoves int
}

func (c Creature) IsInPod() bool {
	for pod := 1; pod <= 4; pod++ {
		for d := 1; d <= PodDepth; d++ {
			if c.Cur == Location(pod*10+d) {
				return true
			}
		}
	}
	return false
}

func (c Creature) AtBottom() bool {
	return c.Cur%10 == 1
}

func (c Creature) AbovePod() Location {
	return 101 + (c.Cur/10)*2
}

func (c Creature) InOwnPod() bool {
	switch c.Cur / 10 {
	case 1:
		return c.Type == 'A'
	case 2:
		return c.Type == 'B'
	case 3:
		return c.Type == 'C'
	case 4:
		return c.Type == 'D'
	}
	return false
}

func (c Creature) Cost() int {
	switch c.Type {
	case 'A':
		return 1
	case 'B':
		return 10
	case 'C':
		return 100
	case 'D':
		return 1000
	}
	log.Print("invalid type for cost: ", c.Type)
	return 0
}

// FreeLocations - locations in Hallway where a creature can move
func (b Burrow) FreeLocations(initLocation Location) []Location {
	var l []Location
	if initLocation < 101 || initLocation > 111 {
		log.Fatal("invalid init location: ", initLocation)
		return nil
	}
	for direction := 1; direction >= -1; direction -= 2 {
		curLoc := initLocation
		for {
			curLoc += Location(direction)
			if b.IsPodTop(curLoc) != 0 {
				curLoc += Location(direction)
			}
			if curLoc > 111 || curLoc < 101 || b.Occupant(curLoc) != 0 {
				break
			}
			l = append(l, curLoc)
		}
	}
	return l
}

func (c Creature) FreeAbove(b *Burrow) bool {
	if !c.IsInPod() {
		log.Fatal("creature not in pod")
	}
	op := int(c.Cur/10) * 10
	for i := int(c.Cur%10) + 1; i <= PodDepth; i++ {
		if b.Occupant(Location(op+i)) != 0 {
			return false
		}
	}
	return true
}

// Moves returns possible moves for a creature in a burrow
func (c Creature) Moves(b *Burrow) []Creature {
	// if b.Hall[10] == 'D' && b.D[2] == 0 &&
	// 	b.Hall[0] == 'A' && b.Hall[7] == 'B' && b.Hall[9] == 'B' &&
	// 	b.Hall[1] == 'A' && b.Hall[3] == 0 {
	// 	b.Print()
	// }

	if c.InOwnPod() {
		if c.FinishedPod(b) {
			return nil
		}
	}
	if c.IsInPod() {
		var cs []Creature
		if c.FreeAbove(b) { // can move up and anywhere out
			depth := PodDepth - (c.Cur % 10) + 1
			ap := c.AbovePod()
			fl := b.FreeLocations(ap)
			for _, l := range fl {
				var dist int
				if l > ap {
					dist = int(l - ap)
				} else {
					dist = int(ap - l)
				}
				dist += int(depth)
				// Check if creature can move directly from the current pod to it's destination pod.
				pt := c.OwnPodTop()
				if l == pt-1 || l == pt+1 {
					pf := b.PodFree(c)
					if pf > 0 {
						dist += PodDepth - (int(pf) % 10) + 2
						return []Creature{{c.Type, pf, c.SpentE + dist*c.Cost(), c.NMoves + 1}}
					}
				}
				cs = append(cs, Creature{c.Type, l, c.SpentE + dist*c.Cost(), c.NMoves + 1})
			}
		}
		return cs
	}
	return b.HallMoves(c)
}

func (c Creature) FinishedPod(b *Burrow) bool {
	if !c.InOwnPod() {
		return false
	}
	op := c.OwnPod() - 1
	for i := 1; i < int(c.Cur%10); i++ {
		if b.Occupant(op+Location(i)) != c.Type {
			return false
		}
	}
	if c.Cur%10 < PodDepth {
		emptyAbove := b.Occupant(c.Cur+1) == 0
		for i := int(c.Cur%10) + 1; i <= PodDepth; i++ {
			oc := b.Occupant(op + Location(i))
			if emptyAbove {
				if oc != 0 {
					return false
				}
			} else {
				if oc == 0 {
					emptyAbove = true
				} else {
					if oc != c.Type {
						return false
					}
				}
			}
		}
	}
	return true
}

func (b *Burrow) Put(loc rune, num int, t rune) {
	switch loc {
	case 'A':
		b.A[num] = t
	case 'B':
		b.B[num] = t
	case 'C':
		b.C[num] = t
	case 'D':
		b.D[num] = t
	case 'H':
		b.Hall[num] = t
	default:
		log.Fatal("invalid location: ", loc)
	}
}

func (b Burrow) Place(pod rune, i int) rune {
	switch pod {
	case 'A':
		return b.A[i]
	case 'B':
		return b.B[i]
	case 'C':
		return b.C[i]
	case 'D':
		return b.D[i]
	case 'H':
		return b.Hall[i]
	}
	log.Fatal("invalid place name: ", pod)
	return 0
}

func (b Burrow) Occupant(loc Location) rune {
	switch loc / 10 {
	case 1:
		return b.A[loc-11]
	case 2:
		return b.B[loc-21]
	case 3:
		return b.C[loc-31]
	case 4:
		return b.D[loc-41]
	case 10, 11:
		return b.Hall[loc-101]
	default:
		log.Panic("invalid location for occupant: ", loc)
		return 0
	}
}

func (b Burrow) IsPodTop(loc Location) rune {
	if loc < 101 || loc > 111 {
		// log.Fatal("location not in hallway: ", loc)
		return 0
	}
	switch loc {
	case 103:
		return 'A'
	case 105:
		return 'B'
	case 107:
		return 'C'
	case 109:
		return 'D'
	}
	return 0
}

func (c Creature) OwnPod() Location {
	switch c.Type {
	case 'A':
		return 11
	case 'B':
		return 21
	case 'C':
		return 31
	case 'D':
		return 41
	}
	log.Fatal("wrong creature type: ", c.Type)
	return 0
}

func (c Creature) OwnPodTop() Location {
	switch c.Type {
	case 'A':
		return 103
	case 'B':
		return 105
	case 'C':
		return 107
	case 'D':
		return 109
	}
	log.Fatal("wrong creature type: ", c.Type)
	return 0
}

func (c Creature) AboveOwnPod() bool {
	pt := c.OwnPodTop()
	return c.Cur == pt-1 || c.Cur == pt+1
}

func (cl Creatures) OwnPodOccupied(t rune) []Creature {
	var oc []Creature
	c := Creature{Type: t}
	op := c.OwnPod()
	for _, c := range cl.creatures {
		if c.Type != t && c.Cur/10 == op/10 {
			oc = append(oc, c)
		}
	}
	return oc
}

func (b Burrow) PodFree(c Creature) Location {
	op := c.OwnPod()
	l := op + PodDepth - 1
	for l >= op {
		if b.Occupant(l) != 0 {
			break
		}
		l--
	}
	if l%10 >= PodDepth { // Pod full
		return 0
	}
	if l%10 == 0 { // Pod empty
		return op
	}
	cx := Creature{Type: c.Type, Cur: l + 1}
	if cx.FinishedPod(&b) {
		return l + 1
	}
	return 0
}

func (b Burrow) HallMoves(c Creature) []Creature {
	fl := b.FreeLocations(c.Cur)
	pt := c.OwnPodTop()
	pf := b.PodFree(c)
	if pf == 0 {
		return nil
	}
	for _, l := range fl { // check if the creature can jump into its own pod - this is the only move, if it can be done
		if pt-1 == l || pt+1 == l {
			var dist int
			if c.Cur > pt {
				dist = int(c.Cur - pt)
			} else {
				dist = int(pt - c.Cur)
			}
			dist += PodDepth - (int(pf) % 10) + 1
			return []Creature{{c.Type, pf, c.SpentE + dist*c.Cost(), c.NMoves + 1}}
		}
	}
	// try to fall in immediately if on the edge of own pod
	l := c.Cur
	if pt-1 == l || pt+1 == l {
		var dist int
		if c.Cur > pt {
			dist = int(c.Cur - pt)
		} else {
			dist = int(pt - c.Cur)
		}
		dist += PodDepth - (int(pf) % 10) + 1
		return []Creature{{c.Type, pf, c.SpentE + dist*c.Cost(), c.NMoves + 1}}
	}
	// no moving once in hallway
	// for _, l := range fl {
	// 	var dist int
	// 	if c.Cur > l {
	// 		dist = int(c.Cur - l)
	// 	} else {
	// 		dist = int(l - c.Cur)
	// 	}
	// 	cl = append(cl, Creature{c.Type, l, c.SpentE + dist*c.Cost(), c.NMoves + 1})
	// }
	return nil
}

func (cl Creatures) MakeBurrow() Burrow {
	var b Burrow
	for _, c := range cl.creatures {
		switch c.Cur / 10 {
		case 1:
			b.Put('A', int(c.Cur%10)-1, c.Type)
		case 2:
			b.Put('B', int(c.Cur%10)-1, c.Type)
		case 3:
			b.Put('C', int(c.Cur%10)-1, c.Type)
		case 4:
			b.Put('D', int(c.Cur%10)-1, c.Type)
		default:
			if c.Cur >= 101 && c.Cur <= 111 {
				b.Put('H', int(c.Cur)-101, c.Type)
			} else {
				log.Fatal("invalid current location: ", c.Cur)
			}
		}
	}
	return b
}

func (cl Creatures) IsEnd() bool {
	burrow := cl.MakeBurrow()
	for pod := 1; pod <= 4; pod++ {
		t := rune('A' + (pod - 1))
		for d := 1; d <= PodDepth; d++ {
			loc := Location(pod*10 + d)
			if burrow.Occupant(loc) != t {
				return false
			}
		}
	}
	return true
}

func (cl Creatures) SpentEnergy() int {
	var e int
	for _, c := range cl.creatures {
		e += c.SpentE
	}
	return e
}

func (cl Creatures) SpentMoves() int {
	var sm int
	for _, c := range cl.creatures {
		sm += c.NMoves
	}
	return sm
}

func (c Creature) EndLocation(b *Burrow) []Creature {
	var obstacles []Creature // whole pod, including the very creature is an obstacle

	el := c.OwnPod()
	for ; el%10 <= PodDepth; el++ {
		if b.Occupant(el) != c.Type {
			break
		}
	}

	if c.InOwnPod() {
		atBottom := true
		for p := el - 1; p%10 > 0; p-- {
			if b.Occupant(p) != c.Type {
				atBottom = false
				break
			}
		}
		if atBottom {
			return nil // creature already in its final place
		}
		op := c.OwnPod()
		for i := 0; i < PodDepth; i++ {
			l := op + Location(i)
			occ := b.Occupant(l)
			if occ != 0 {
				obstacles = append(obstacles, Creature{Type: occ, Cur: l})
			}
		}
		return obstacles
	}

	var ap Location
	if c.IsInPod() {
		l := c.Cur + 1
		for ; l%10 <= PodDepth; l++ {
			occ := b.Occupant(l)
			if occ != 0 {
				obstacles = append(obstacles, Creature{Type: occ, Cur: l})
			}
		}
		ap = c.AbovePod()
	} else {
		ap = c.Cur
	}
	cop := Creature{Type: c.Type, Cur: c.OwnPod()}
	aop := cop.AbovePod()
	direction := 1
	if aop < ap {
		direction = -1
	}
	for l := ap; l != aop; l += Location(direction) {
		occ := b.Occupant(l)
		if occ != 0 {
			obstacles = append(obstacles, Creature{Type: occ, Cur: l})
		}
	}

	op := c.OwnPod()
	for i := PodDepth - 1; i >= 0; i-- {
		l := op + Location(i)
		occ := b.Occupant(l)
		if occ != 0 {
			obstacles = append(obstacles, Creature{Type: occ, Cur: l})
		}
		if l == el {
			break
		}
	}
	return obstacles
}

func (cl Creatures) Evaluate() int {
	var score int
	b := cl.MakeBurrow()
	for _, c := range cl.creatures {
		endLocList := c.EndLocation(&b)
		for _, obstacle := range endLocList {
			score -= 5 * obstacle.Cost()
		}
	}
	score -= cl.SpentEnergy()
	return score
}

func (cl Creatures) Equals(c Creatures) bool {
	if len(cl.creatures) != len(c.creatures) {
		return false
	}
	for i := range cl.creatures {
		if cl.creatures[i].Type != c.creatures[i].Type || cl.creatures[i].Cur != c.creatures[i].Cur {
			return false
		}
	}
	return true
}

func (cl Creatures) IntoPod() bool {
	return cl.creatures[cl.lastMoved].Cur/10 >= 1 && cl.creatures[cl.lastMoved].Cur/10 <= 4
}

func (b Burrow) Equals(b2 *Burrow) bool {
	for pod := 1; pod <= 4; pod++ {
		for d := 1; d <= PodDepth; d++ {
			l := Location(pod*10 + d)
			if b.Occupant(l) != b2.Occupant(l) {
				return false
			}
		}
	}
	for l := Location(101); l <= Location(111); l++ {
		if b.Occupant(l) != b2.Occupant(l) {
			return false
		}
	}
	return true
}

func (b Burrow) Hash() int {
	var h int
	for pod := 1; pod <= 4; pod++ {
		var v int
		for d := 1; d <= PodDepth; d++ {
			l := Location(pod*10 + d)
			c := b.Occupant(l)
			if c != 0 {
				v = 4 * (int(c-'A') + 1)
				v = 5 * v
			}
		}
		h = 625*h + v
	}
	return h
}

func (b Burrow) Print() {
	var l string
	for i := 0; i < 13; i++ {
		l = l + "#"
	}
	log.Print(l)
	l = "#"
	for i := 101; i <= 111; i++ {
		oc := b.Occupant(Location(i))
		if oc == 0 {
			l = l + "."
		} else {
			l = l + string(oc)
		}
	}
	l = l + "#"
	log.Print(l)
	for d := PodDepth; d > 0; d-- {
		if d == PodDepth {
			l = "###"
		} else {
			l = "  #"
		}
		for i := 1; i <= 4; i++ {
			oc := b.Occupant(Location(i*10 + d))
			if oc == 0 {
				l = l + "."
			} else {
				l = l + string(oc)
			}
			l = l + "#"
		}
		if d == PodDepth {
			l = l + "##"
		}
		log.Print(l)
	}
	l = "  "
	for i := 0; i < 9; i++ {
		l = l + "#"
	}
	log.Print(l)
}

var bestBurrows map[int][]Burrow
var mubMap map[int]*sync.Mutex
var muBB sync.Mutex

func (cl Creatures) Solve(maxE int, maxMoves int, history []Creatures, lastMoved int, depth int) ([]Creatures, int, bool) {
	MaxIntoPod := 10
	// if len(history) > 5 {
	// 	MaxIntoPod = 2
	// }
	en := cl.SpentEnergy()
	mv := cl.SpentMoves()

	curB := cl.MakeBurrow()
	var found bool
	bh := curB.Hash()
	if _, ok := mubMap[bh]; !ok {
		muBB.Lock()
		mubMap[bh] = &sync.Mutex{}
		muBB.Unlock()
	}
	mubMap[bh].Lock()
	bb, ok := bestBurrows[bh]
	if !ok {
		muBB.Lock()
		bb = []Burrow{}
		bestBurrows[bh] = bb
		muBB.Unlock()
	}
	for i, b := range bb {
		if b.Equals(&curB) {
			if b.Cost < en {
				mubMap[bh].Unlock()
				return nil, en, false
			}
			bestBurrows[bh][i].Cost = en
			found = true
			break
		}
	}
	if !found {
		bestBurrows[bh] = append(bestBurrows[bh], curB)
		if len(bestBurrows[bh])%1000 == 0 {
			log.Printf("%d: [%d] %d", depth, bh, len(bestBurrows[bh]))
		}
	}
	mubMap[bh].Unlock()
	var lastPod int
	// if len(history) > 5 {
	// 	log.Print(history)
	// }
	if len(history) >= MaxIntoPod {
		for i, h := range history {
			if h.IntoPod() {
				// if len(history) == 5 {
				// 	log.Print(history)
				// }
				lastPod = i + 1
				break
			}
		}
		if lastPod == 0 { // not found
			lastPod = 99
		}
	}
	if cl.IsEnd() {
		// log.Print("solution: ", en, mv, cl, history)
		return history, en, true
	}
	if en >= maxE {
		// log.Print("energy out: ", en, mv, cl)
		return nil, en, false
	}
	if mv >= maxMoves {
		// log.Print("moves out: ", en, mv, cl)
		return nil, en, false
	}
	allMoves := cl.Moves()
	var possibleMoves []Creatures
	for _, m := range allMoves {
		found := false
		for _, pm := range history {
			if pm.Equals(m) {
				found = true
				// log.Print("loop ", maxMoves)
				break
			}
		}
		if !found {
			if lastPod < MaxIntoPod || m.IntoPod() { // only add moves that put a creature into pod in at most 3 steps
				possibleMoves = append(possibleMoves, m)
			}
		}
	}
	// if len(allMoves) != len(possibleMoves) {
	// 	log.Print("maxMoves: ", maxMoves, " removed ", len(allMoves)-len(possibleMoves), " from possible moves ", len(allMoves), "history len: ", len(history))
	// }
	// sort.Slice(possibleMoves, func(i, j int) bool { return possibleMoves[i].Evaluate() > possibleMoves[j].Evaluate() })

	var (
		men       []int
		solved    []bool
		solutions [][]Creatures
		mu        sync.Mutex
		wg        sync.WaitGroup
	)
	for _, m := range possibleMoves {
		// if i > 10 {
		// 	break
		// }
		wg.Add(1)
		var newH []Creatures
		newH = append(newH, m)
		newH = append(newH, history...)
		// if m.creatures[2].Cur == 16 && m.creatures[7].Cur == 18 && m.creatures[6].Cur == 20 {
		// 	log.Print("good move: ", mv, cl)
		// }
		go func(mx Creatures) {
			defer wg.Done()
			solution, me, ok := mx.Solve(maxE, maxMoves, newH, mx.lastMoved, depth+1)
			if ok {
				mu.Lock()
				men = append(men, me)
				solved = append(solved, ok)
				solutions = append(solutions, solution)
				mu.Unlock()
			}
		}(m)
	}
	wg.Wait()
	if len(men) == 0 {
		return nil, en, false
	}
	best := 0
	for i, me := range men {
		if solved[i] {
			if me < men[best] {
				best = i
			}
			ok = true
		}
	}
	if ok {
		// log.Print("best solution at depth: ", depth, en, mv, cl, history)
		return solutions[best], men[best], true
	}
	return nil, en, false
}

func (cl Creatures) Replace(c Creature, i int) Creatures {
	var left, right []Creature
	if i > 0 {
		left = append(left, cl.creatures[:i]...)
	}
	if i < len(cl.creatures)-1 {
		right = append(right, cl.creatures[i+1:]...)
	}
	left = append(left, c)
	left = append(left, right...)
	return Creatures{left, i, cl.history, 0}
}

func (cl Creatures) Moves() []Creatures {
	var moves []Creatures
	var mu sync.Mutex
	b := cl.MakeBurrow()
	var wg sync.WaitGroup
	for i, c := range cl.creatures {
		if i == cl.lastMoved {
			continue
		}
		wg.Add(1)
		func(cx Creature, ix int) {
			defer wg.Done()

			pm := cx.Moves(&b)
			// if ix == 15 && (len(pm) == 1 && pm[0].SpentE == 7000) {
			// 	log.Print(cx, ix)
			// }
			// if cx.Cur == 111 && cx.Type == 'D' && b.D[0] == 0 && b.Hall[9] == 0 && pm[0].SpentE == 7000 {
			// 	log.Print(cx, ix)
			// }

			for _, m := range pm {
				newC := cl.Replace(m, ix)
				newC.lastMoved = ix
				newC.cost = newC.SpentEnergy()
				newH := Creatures{
					cl.creatures,
					cl.lastMoved,
					nil,
					cl.cost,
				}
				newC.history = nil
				newC.history = append(newC.history, cl.history...)
				newC.history = append(newC.history, newH)
				// if the move is final (creature moved to final pod), just return this one move
				// if newC.creatures[i].InOwnPod() {
				// 	return []Creatures{newC}
				// }
				mu.Lock()
				moves = append(moves, newC)
				mu.Unlock()
			}
		}(c, i)
	}
	wg.Wait()
	for _, m := range moves {
		if m.creatures[m.lastMoved].InOwnPod() {
			return []Creatures{m}
		}
	}
	return moves
}

type PQ []*Creatures

var pq PQ

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	ei := pq[i].SpentEnergy()
	ej := pq[j].SpentEnergy()
	return ei < ej || (ei == ej && (pq[i].SpentMoves() < pq[j].SpentMoves()))
}

func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func (pq *PQ) Push(x interface{}) {
	cr := x.(*Creatures)
	*pq = append(*pq, cr)
}

func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

type CStore struct {
	clists map[int]map[int]*Creatures
}

func (cl Creatures) Hash() (int, int) {
	var h, hh int
	for _, c := range cl.creatures {
		v := int(c.Type-'A') + 1
		if c.IsInPod() {
			n := (int(c.Cur/10)-1)*5 + int(c.Cur%10) - 1
			x := 1
			for n > 0 {
				x = x * 5
				n--
			}
			h += x * v
		} else {
			n := int(c.Cur - 101)
			x := 1
			for n > 0 {
				x = x * 5
				n--
			}
			hh += x * v
		}
	}
	return h, hh
}

func (cs *CStore) CostExists(cl *Creatures) (int, bool) {
	h, hh := cl.Hash()
	if l, ok := cs.clists[h]; ok {
		if c, ok := l[hh]; ok {
			return c.cost, true
		}
	}
	return 0, false
}

func (cs *CStore) Add(cl *Creatures) {
	h, hh := cl.Hash()
	if l, ok := cs.clists[h]; ok {
		l[hh] = cl
		// if c, ok := l[hh]; ok {
		// 	return c.cost, true
		// }
	} else {
		cs.clists[h] = make(map[int]*Creatures)
		cs.clists[h][hh] = cl
	}

	// if _, exists := cs.CostExists(cl); exists {
	// 	for i := range (*cs).clists[h] {
	// 		if cs.clists[h][i].Equals(cl) {
	// 			cs.clists[h][i] = &cl
	// 		}
	// 	}
	// } else {
	// 	if _, ok := cs.clists[h]; ok {
	// 		cs.clists[h] = append(cs.clists[h], &cl)
	// 		return
	// 	}
	// 	cs.clists[h] = []*Creatures{&cl}
	// }
}

func NewCStore() CStore {
	cl := make(map[int]map[int]*Creatures)
	return CStore{
		cl,
	}
}

func (cl Creatures) SolveH() bool {
	cstore := NewCStore()
	heap.Push(&pq, &cl)
	lastCost := 0
	for pq.Len() > 0 {
		cheapest := heap.Pop(&pq).(*Creatures)
		// en := cheapest.SpentEnergy()

		if cheapest.cost == 4288+40+5000 && //+40+0*5000 && cheapest.creatures[5].Cur == 106 &&
			cheapest.creatures[14].Cur == 101 &&
			cheapest.creatures[4].Cur == 104 {
			cheapest.MakeBurrow().Print()
		}
		if cheapest.IsEnd() {
			log.Print("solved: ", cheapest.SpentEnergy())
			// log.Print(cheapest.SpentEnergy(), cheapest)
			curCost := 0
			var en int
			for i, h := range cheapest.history {
				b := h.MakeBurrow()
				en = h.SpentEnergy()
				if h.lastMoved >= 0 {
					prevLoc := cheapest.history[i-1].creatures[h.lastMoved].Cur
					log.Printf("%d = %d + %d (%c: %d->%d)", en, curCost, en-curCost, h.creatures[h.lastMoved].Type, prevLoc, h.creatures[h.lastMoved].Cur)
				} else {
					log.Printf("%d = %d + %d", en, curCost, en-curCost)
				}
				curCost = en
				b.Print()
			}
			b := cheapest.MakeBurrow()
			en = cheapest.SpentEnergy()
			log.Printf("%d = %d + %d", en, curCost, en-curCost)
			curCost = en
			b.Print()
			return true
			log.Print("---------------")
		}
		allMoves := cheapest.Moves()
		if lastCost/1000 != cheapest.cost/1000 {
			log.Print(cheapest.cost)
		}
		lastCost = cheapest.cost
		// log.Print(cheapest.cost, cheapest.SpentMoves(), cheapest.Evaluate(), cheapest.creatures)
		// if (cheapest.creatures[15].Cur == 111 &&
		// 	cheapest.creatures[14].Cur == 101 &&
		// 	cheapest.creatures[11].Cur == 110 &&
		// 	cheapest.creatures[10].Cur == 108 &&
		// 	cheapest.creatures[9].Cur == 102) ||
		// 	(cheapest.creatures[7].Cur == 41) {
		// 	log.Print(cheapest.cost, cheapest.SpentMoves(), cheapest.Evaluate(), cheapest.creatures)
		// 	for i, h := range cheapest.history {
		// 		log.Print(i+1, ": ", h)
		// 	}
		// 	cheapest.MakeBurrow().Print()
		// 	// 	// return true
		// }
		for i, m := range allMoves {
			mc := m.SpentEnergy()
			// isEnd := false
			if m.IsEnd() {
				log.Print("move solves ", m.cost, m.SpentEnergy())
				// isEnd = true
			}
			if cost, exists := cstore.CostExists(&m); exists {
				if cost <= mc {
					continue
				}
			}
			cstore.Add(&allMoves[i])
			heap.Push(&pq, &allMoves[i])
		}
		// log.Print(pq.Len())
	}
	return false
}

func (cl Creatures) SolveD(maxE int, maxMoves int) bool {
	en := cl.SpentEnergy()
	// if en > maxE {
	// 	// log.Print("over max cost:", en)
	// 	return false
	// }

	mv := cl.SpentMoves()
	// if mv > maxMoves {
	// 	// log.Print("over max moves:", mv)
	// 	return false
	// }

	curB := cl.MakeBurrow()
	curB.Cost = en
	bh := curB.Hash()
	if cl.IsEnd() {
		log.Printf("solved: %d (%d) %v %d %v", en, mv, cl, bh, bestBurrows[bh])
		return true
	}

	var found bool

	// muBB.Lock()
	// mub, ok := mubMap[bh]
	// if !ok {
	// 	mubMap[bh] = &sync.Mutex{}
	// 	mub = mubMap[bh]
	// }
	bb, ok := bestBurrows[bh]
	if !ok {
		bb = []Burrow{}
		bestBurrows[bh] = bb
	}
	// mub.Lock()
	for i, b := range bb {
		if b.Equals(&curB) {
			if b.Cost < en {
				// muBB.Unlock()
				return false
			}
			bestBurrows[bh][i].Cost = en
			found = true
			break
		}
	}
	if !found {
		bestBurrows[bh] = append(bestBurrows[bh], curB)
		if len(bestBurrows[bh])%100 == 0 {
			log.Printf("%d: [%d] %d", mv, bh, len(bestBurrows[bh]))
		}
	}
	// muBB.Unlock()

	allMoves := cl.Moves()
	sort.Slice(allMoves, func(i, j int) bool { return allMoves[i].SpentEnergy() < allMoves[j].SpentEnergy() })

	// var wg sync.WaitGroup
	var solved bool
	for _, m := range allMoves {
		solved = solved || m.SolveD(maxE, maxMoves)
		// wg.Add(1)
		// func(mx Creatures) {
		// defer wg.Done()
		// 	mx.SolveD()
		// }(m)
	}
	// wg.Wait()
	return solved
}

func main() {
	heap.Init(&pq)
	creatures := Creatures{
		[]Creature{
			{'A', 11, 0, 0},
			{'B', 12, 0, 0},
			{'D', 21, 0, 0},
			{'C', 22, 0, 0},
			{'C', 31, 0, 0},
			{'B', 32, 0, 0},
			{'A', 41, 0, 0},
			{'D', 42, 0, 0},
		},
		-1,
		nil,
		0,
	}

	creatures = Creatures{
		[]Creature{
			{'C', 11, 0, 0},
			{'D', 12, 0, 0},
			{'A', 21, 0, 0},
			{'B', 22, 0, 0},
			{'D', 31, 0, 0},
			{'A', 32, 0, 0},
			{'B', 41, 0, 0},
			{'C', 42, 0, 0},
		},
		-1,
		nil,
		0,
	}

	if PodDepth == 4 {
		creatures = Creatures{
			[]Creature{
				{'C', 11, 0, 0},
				{'D', 12, 0, 0},
				{'D', 13, 0, 0},
				{'D', 14, 0, 0},
				{'A', 21, 0, 0},
				{'B', 22, 0, 0},
				{'C', 23, 0, 0},
				{'B', 24, 0, 0},
				{'D', 31, 0, 0},
				{'A', 32, 0, 0},
				{'B', 33, 0, 0},
				{'A', 34, 0, 0},
				{'B', 41, 0, 0},
				{'C', 42, 0, 0},
				{'A', 43, 0, 0},
				{'C', 44, 0, 0},
			},
			-1,
			nil,
			0,
		}

		// creatures = Creatures{
		// 	[]Creature{
		// 		{'A', 11, 0, 0},
		// 		{'D', 12, 0, 0},
		// 		{'D', 13, 0, 0},
		// 		{'B', 14, 0, 0},
		// 		{'D', 21, 0, 0},
		// 		{'B', 22, 0, 0},
		// 		{'C', 23, 0, 0},
		// 		{'C', 24, 0, 0},
		// 		{'C', 31, 0, 0},
		// 		{'A', 32, 0, 0},
		// 		{'B', 33, 0, 0},
		// 		{'B', 34, 0, 0},
		// 		{'A', 41, 0, 0},
		// 		{'C', 42, 0, 0},
		// 		{'A', 43, 0, 0},
		// 		{'D', 44, 0, 0},
		// 	},
		// 	-1,
		// 	nil,
		// 	0,
		// }
	}

	bestBurrows = make(map[int][]Burrow)
	creatures.SolveH()
	return
	mubMap = make(map[int]*sync.Mutex)
	// maxE := 45000
	// for {
	// 	log.Print("trying to solve for: ", maxE)
	// 	if creatures.SolveD(maxE) {
	// 		break
	// 	}
	// 	maxE += 1000
	// }
	// return
	maxE := 45000
	maxMoves := 22
	for {
		bestBurrows = make(map[int][]Burrow)
		log.Print("trying to solve for: ", maxE, maxMoves)
		if solved := creatures.SolveD(maxE, maxMoves); solved {
			break
		}
		// if solution, spentE, ok := creatures.Solve(maxE, maxMoves, nil, -1, 1); ok {
		// 	log.Print("min energy: ", spentE, " moves: ", maxMoves, " solution (", len(solution), "): ", solution)
		// 	// break
		// }
		maxMoves++

		if maxMoves > 40 {
			maxMoves = 22
			maxE += 5000
			if maxE > 1000000 {
				break
			}
		}
	}
}
