// Solves only part 2.
// In an attempt to calculate to higher scores than the exercise demands, this
// go solution is parallelized. I managed to run it to a maximum score of 600 -
// after that my workstation runs out of memory.
// Memory usage is roughly cubic in the maximum score:
// - The number of possible states to memoize is quadratic: nScores²•nPositions²•nPlayers
// - The number of total games to play (so the number of possible results) is
//   exponential in the score, so the storage required to represent them is linear
// It might be possible to optimize this down to quadratic memory by using an
// LRU cache or similar - if we assume that all goroutines have roughly the
// same speed, we only need to memoize a small window of actual results, as the
// actual scores of players grow monotonically.
// Doing that is left as an exercise to the reader :)

//go:build ignore

package main

import (
	"fmt"
	"math/big"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"unsafe"
)

const maxScore = 650

func main() {
	var (
		nwin [2]*big.Int
		o    sync.Once
		wg   sync.WaitGroup
	)
	wg.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			defer wg.Done()
			nw := PlayDirac(State{
				Pos:    [2]uint8{4, 7},
				Player: 1,
			})
			o.Do(func() {
				nwin[0] = &nw[0]
				nwin[1] = &nw[1]
			})
		}()
	}
	wg.Wait()
	fmt.Printf("Nwins after playing Dirac Dice: %v\n", nwin)
}

type State struct {
	Score  [2]int
	Pos    [2]uint8
	Player uint8
}

var rolls [27]uint8

func init() {
	for i := range rolls {
		i := uint8(i)
		rolls[i] = i%3 + (i/3)%3 + (i/9)%3 + 3
	}
	sort.Slice(rolls[:], func(i, j int) bool {
		return rolls[i] < rolls[j]
	})
}

var (
	mem    = NewStateMap(maxScore)
	bigOne = big.NewInt(1)
)

func PlayDirac(s State) [2]big.Int {
	if nw, ok := mem.Load(s); ok {
		return nw
	}

	var nw [2]big.Int
	for _, r := range rolls[:] {
		s2 := s
		s2.Player = (s2.Player + 1) % 2
		s2.Pos[s2.Player] += r
		for s2.Pos[s2.Player] > 10 {
			s2.Pos[s2.Player] -= 10
		}
		s2.Score[s2.Player] += int(s2.Pos[s2.Player])
		if s2.Score[s2.Player] >= maxScore {
			nw[s2.Player].Add(&nw[s2.Player], bigOne)
			continue
		}
		nw2 := PlayDirac(s2)
		nw[0].Add(&nw[0], &nw2[0])
		nw[1].Add(&nw[1], &nw2[1])
	}
	mem.Store(s, nw)
	return nw
}

type StateMap struct {
	score int
	els   []unsafe.Pointer
}

func NewStateMap(score int) *StateMap {
	nEls := 2 * 10 * 10 * score * score
	return &StateMap{score: score, els: make([]unsafe.Pointer, nEls)}
}

func (m *StateMap) idx(s State) int {
	i := s.Score[0]
	i *= m.score
	i += s.Score[1]
	i *= 10
	i += int(s.Pos[0] - 1)
	i *= 10
	i += int(s.Pos[1] - 1)
	i *= 2
	i += int(s.Player)
	return i
}

func (m *StateMap) Load(s State) ([2]big.Int, bool) {
	i := m.idx(s)
	p := (*[2]big.Int)(atomic.LoadPointer(&m.els[i]))
	if p == nil {
		return [2]big.Int{}, false
	}
	return *p, true
}

func (m *StateMap) Store(s State, n [2]big.Int) {
	i := m.idx(s)
	atomic.StorePointer(&m.els[i], unsafe.Pointer(&n))
}
