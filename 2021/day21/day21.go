// Solves only part 2

//go:build ignore

package main

import (
	"fmt"
	"sort"
)

func main() {
	nwin := PlayDirac(State{
		Pos:    [2]int{4, 7},
		Player: 1,
	})
	fmt.Printf("Nwins after playing Dirac Dice: %v\n", nwin)
}

type State struct {
	Pos    [2]int
	Score  [2]int
	Player int
}

var rolls [27]int

func init() {
	for i := range rolls {
		rolls[i] = i%3 + (i/3)%3 + (i/9)%3 + 3
	}
	sort.Ints(rolls[:])
	fmt.Println(rolls)
}

var (
	n   = 0
	mem = make(map[State][2]int)
)

func PlayDirac(s State) [2]int {
	if nw, ok := mem[s]; ok {
		return nw
	}
	var nw [2]int
	for _, r := range rolls[:] {
		s2 := s
		s2.Player = (s2.Player + 1) % 2
		s2.Pos[s2.Player] += r
		for s2.Pos[s2.Player] > 10 {
			s2.Pos[s2.Player] -= 10
		}
		s2.Score[s2.Player] += s2.Pos[s2.Player]
		if s2.Score[s2.Player] >= 21 {
			nw[s2.Player] += 1
			continue
		}
		nw2 := PlayDirac(s2)
		nw[0] += nw2[0]
		nw[1] += nw2[1]
	}
	mem[s] = nw
	return nw
}
