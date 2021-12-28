package main

import "log"

type Dice int
type Player struct {
	pos   int
	score int
	rolls int
}

const SCORE = 21

func (d *Dice) Roll() int {
	s := 0
	for i := 0; i < 3; i++ {
		*d++
		if *d > 100 {
			*d = 1
		}
		s += int(*d)
	}
	return s
}

func nPlays(d int) int { // plays with this result
	switch d {
	case 3:
		return 1
	case 4:
		return 3
	case 5:
		return 6
	case 6:
		return 7
	case 7:
		return 6
	case 8:
		return 3
	case 9:
		return 1
	default:
		log.Fatal("dice throw invalid: ", d)
		return 0
	}
}

func (p Player) Throw(d int) Player {
	pos := p.pos + d
	pos = pos % 10
	if pos == 0 {
		pos = 10
	}
	return Player{
		pos:   pos,
		score: p.score + pos,
		rolls: p.rolls + 3,
	}
}
func (p Player) Wins() bool {
	return p.score >= SCORE
}

func QPlay(p1, p2 Player) (int, int) {
	var w1, w2 int
	for d31 := 3; d31 <= 9; d31++ {
		plays1 := nPlays(d31)
		nP1 := p1.Throw(d31)
		if nP1.Wins() {
			w1 += plays1
			// log.Printf("p1 wins with %d after %d rolls in %d plays", nP1.score, nP1.rolls, plays1)
		} else {
			for d32 := 3; d32 <= 9; d32++ {
				plays2 := nPlays(d32)
				nP2 := p2.Throw(d32)
				if nP2.Wins() {
					w2 += (plays2 * plays1)
					// log.Printf("p2 wins with %d after %d rolls in %d plays", nP2.score, nP2.rolls, plays2)
				} else {
					nw1, nw2 := QPlay(nP1, nP2)
					w1 += (nw1 * plays2 * plays1)
					w2 += (nw2 * plays2 * plays1)
				}
			}
		}
	}
	return w1, w2
}

func main() {
	const (
		P1 = 8
		P2 = 4
	)
	p1 := P1
	p2 := P2
	s1 := 0
	s2 := 0
	dice := Dice(100)
	rolls := 0
	for {
		p1 += dice.Roll()
		rolls++
		p1 = p1 % 10
		if p1 == 0 {
			p1 = 10
		}
		s1 += p1
		if s1 >= 1000 {
			log.Printf("%d rolls, p1 score: %d, position: %d, p2 score: %d, pos: %d, mult: %d", rolls, s1, p1, s2, p2, rolls*s2*3)
			break
		}
		// log.Printf("%d: p1 rolls %d, pos %d score %d", rolls, int(dice), p1, s1)
		p2 += dice.Roll()
		rolls++
		p2 = p2 % 10
		if p2 == 0 {
			p2 = 10
		}
		s2 += p2
		if s2 >= 1000 {
			log.Printf("%d rolls, p1 score: %d, position: %d, p2 score: %d, pos: %d, mult: %d", rolls, s1, p1, s2, p2, rolls*s1*3)
			break
		}
		// log.Printf("%d: p2 rolls %d, pos %d score %d", rolls, int(dice), p2, s2)
	}

	w1, w2 := QPlay(Player{P1, 0, 0}, Player{P2, 0, 0})
	log.Printf("p1 wins: %d, p2 wins: %d, p1 best: %v", w1, w2, w1 > w2)
}
