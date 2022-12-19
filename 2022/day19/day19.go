package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Merovius/AdventOfCode/internal/input"
	"github.com/Merovius/AdventOfCode/internal/math"
	"github.com/Merovius/AdventOfCode/internal/set"
)

var _ = math.MaxInt

func main() {
	log.SetFlags(log.Lshortfile)

	data, err := input.Lines(func(s string) (Blueprint, error) {
		var (
			b Blueprint
			n int
		)
		i, err := fmt.Sscanf(s, "Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.", &n, &b[Ore][Ore], &b[Clay][Ore], &b[Obsidian][Ore], &b[Obsidian][Clay], &b[Geode][Ore], &b[Geode][Obsidian])
		if err != nil || i != 7 {
			return b, errors.New("invalid input line")
		}
		return b, nil
	}).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Part1(data))
	fmt.Println("Part 2:", Part2(data[:3]))
}

func Part1(b []Blueprint) int {
	wg := new(sync.WaitGroup)
	wg.Add(len(b))
	result := make([]int, len(b))
	for i, b := range b {
		go func(i int, b Blueprint) {
			defer wg.Done()
			result[i] = MaximizeGeodes(b, 24)
		}(i, b)
	}
	wg.Wait()
	var total int
	for i, g := range result {
		total += (i + 1) * g
	}
	return total
}

func Part2(b []Blueprint) int {
	wg := new(sync.WaitGroup)
	wg.Add(len(b))
	result := make([]int, len(b))
	for i, b := range b {
		go func(i int, b Blueprint) {
			defer wg.Done()
			result[i] = MaximizeGeodes(b, 32)
		}(i, b)
	}
	wg.Wait()
	total := 1
	for _, g := range result {
		total *= g
	}
	return total

}

func MaximizeGeodes(b Blueprint, time int) int {
	var (
		max = b.MaxNeeded()
		cur = set.Make(State{
			Bots: Materials{Ore: 1},
		})
		next set.Set[State]
	)
	add := func(s State) {
		for ss := range next {
			switch Compare(ss, s) {
			case -1:
				delete(next, ss)
			case 1:
				return
			}
		}
		next.Add(s)
	}
	for t := 0; t < time; t++ {
		next = make(set.Set[State])
		for s := range cur {
			for _, n := range b.Next(s, max) {
				add(n)
			}
		}
		cur = next
	}
	var best int
	for s := range cur {
		best = math.Max(best, s.Raw[Geode])
	}
	return best
}

type Material int

const (
	Ore Material = iota
	Clay
	Obsidian
	Geode

	None
)

func (m Material) Vec() Materials {
	var v Materials
	if m != None {
		v[m] = 1
	}
	return v
}

// Ore, clay, obsidian, geodes
type Materials [4]int

// Cost for Ore, clay, obsidian, geode bots respectively
type Blueprint [4]Materials

type State struct {
	Raw  Materials
	Bots Materials
}

func (s State) CanBuild(r Materials) bool {
	for i, v := range s.Raw {
		if v < r[i] {
			return false
		}
	}
	return true
}

func (b Blueprint) MaxNeeded() Materials {
	var maxNeeded Materials
	for _, r := range b {
		for m, v := range r {
			maxNeeded[m] = math.Max(maxNeeded[m], v)
		}
	}
	maxNeeded[Geode] = 1000
	return maxNeeded
}

func (b Blueprint) Next(s State, max Materials) []State {
	var out []State
	for m := Ore; m < None; m++ {
		if s.Bots[m] >= max[m] {
			continue
		}
		if s.CanBuild(b[m]) {
			out = append(out, State{
				Raw:  Sub(Add(s.Raw, s.Bots), b[m]),
				Bots: Add(s.Bots, m.Vec()),
			})
		}
	}
	return append(out, State{
		Raw:  Add(s.Raw, s.Bots),
		Bots: s.Bots,
	})
}

func Add(a, b Materials) Materials {
	return Materials{
		Ore:      b[Ore] + a[Ore],
		Clay:     b[Clay] + a[Clay],
		Obsidian: b[Obsidian] + a[Obsidian],
		Geode:    b[Geode] + a[Geode],
	}
}

func Sub(a, b Materials) Materials {
	return Materials{
		Ore:      a[Ore] - b[Ore],
		Clay:     a[Clay] - b[Clay],
		Obsidian: a[Obsidian] - b[Obsidian],
		Geode:    a[Geode] - b[Geode],
	}
}

// Compare return -1 if s is definitely better than t, 1 if t is definitely
// better than s and 0 otherwise.
func Compare(s, t State) int {
	var cmp int
	for m := Ore; m < None; m++ {
		switch math.Cmp(s.Bots[m], t.Bots[m]) {
		case -1:
			if cmp > 0 {
				return 0
			}
			cmp = -1
		case 1:
			if cmp < 0 {
				return 0
			}
			cmp = 1
		}
		switch math.Cmp(s.Raw[m], t.Raw[m]) {
		case -1:
			if cmp > 0 {
				return 0
			}
			cmp = -1
		case 1:
			if cmp < 0 {
				return 0
			}
			cmp = 1
		}
	}
	return cmp
}
