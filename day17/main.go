package main

import (
	"log"
)

func hits(target [][2]int, v [2]int) (bool, int) {
	var x, y, mh int
	for {
		if x >= target[0][0] && x <= target[0][1] &&
			y >= target[1][0] && y <= target[1][1] {
			return true, mh
		}
		if x > target[0][1] || y < target[1][0] {
			return false, mh
		}
		x += v[0]
		y += v[1]
		if y > mh {
			mh = y
		}
		if v[0] > 0 {
			v[0]--
		}
		v[1]--
	}
}

func main() {
	target := [][2]int{{20, 30}, {-10, -5}}
	target = [][2]int{{153, 199}, {-114, -75}}
	minH := 0
	maxV := 200
	for {
		if minH*(minH+1)/2 >= target[0][0] {
			break
		}
		minH++
	}
	maxHeight := 0
	nHits := 0
	var mhV [2]int
	for vv := target[1][0]; vv <= maxV; vv++ {
		for hv := minH; hv <= target[0][1]; hv++ {
			if hits, mh := hits(target, [2]int{hv, vv}); hits {
				nHits++
				if mh > maxHeight {
					maxHeight = mh
					mhV = [2]int{hv, vv}
				}
			}
		}
	}
	log.Print("max height: ", maxHeight, mhV, nHits)
}
