// Occupiable positions are numbered top-to-bottom, left-to-right:
// •••••••••••••••••••
// •0 1. 2. 3. 4. 5 6•
// ••••7 •8 •9 •10••••
//    •11•12•13•14•
//    •15•16•17•18•
//    •19•20•21•22•
//    •••••••••••••

package main

import (
	"errors"
	"fmt"

	"gonih.org/AdventOfCode/2021/go/priority_queue"
)

func main() {
	var example = State{
		7: B, 8: C, 9: B, 10: D,
		11: D, 12: C, 13: B, 14: A,
		15: D, 16: B, 17: A, 18: C,
		19: A, 20: D, 21: C, 22: A,
	}
	var input = State{
		7: D, 8: A, 9: C, 10: C,
		11: D, 12: C, 13: B, 14: A,
		15: D, 16: B, 17: A, 18: C,
		19: D, 20: A, 21: B, 22: B,
	}
	_, _ = example, input
	fmt.Println("Cost to organize:", Organize(input))
}

// Cell is a cell content
type Cell uint8

// Possible cell types
const (
	None Cell = iota
	A
	B
	C
	D
)

// State is a full state of the burrow.
type State [23]Cell

var End = State{
	7: A, 8: B, 9: C, 10: D,
	11: A, 12: B, 13: C, 14: D,
	15: A, 16: B, 17: C, 18: D,
	19: A, 20: B, 21: C, 22: D,
}

// Edge is an edge between two cells
type Edge [2]int

// EdgeCost is the cost associated with the given edge. If an edge is not in
// EdgeCost, it doesn't exist.
var EdgeCost = map[Edge]int{
	{0, 1}: 1,
	{1, 0}: 1, {1, 2}: 2, {1, 7}: 2,
	{2, 1}: 2, {2, 3}: 2, {2, 7}: 2, {2, 8}: 2,
	{3, 2}: 2, {3, 4}: 2, {3, 8}: 2, {3, 9}: 2,
	{4, 3}: 2, {4, 5}: 2, {4, 9}: 2, {4, 10}: 2,
	{5, 4}: 2, {5, 10}: 2, {5, 6}: 1,
	{6, 5}: 1,
	{7, 1}: 2, {7, 2}: 2, {7, 11}: 1,
	{8, 2}: 2, {8, 3}: 2, {8, 12}: 1,
	{9, 3}: 2, {9, 4}: 2, {9, 13}: 1,
	{10, 4}: 2, {10, 5}: 2, {10, 14}: 1,
	{11, 7}: 1, {11, 15}: 1,
	{12, 8}: 1, {12, 16}: 1,
	{13, 9}: 1, {13, 17}: 1,
	{14, 10}: 1, {14, 18}: 1,
	{15, 11}: 1, {15, 19}: 1,
	{16, 12}: 1, {16, 20}: 1,
	{17, 13}: 1, {17, 21}: 1,
	{18, 14}: 1, {18, 22}: 1,
	{19, 15}: 1,
	{20, 16}: 1,
	{21, 17}: 1,
	{22, 18}: 1,
}

type Path struct {
	Length int
	Nodes  []int
}

var Paths [23][23]Path

func (p Path) Append(node int) Path {
	last := p.Nodes[len(p.Nodes)-1]
	c, ok := EdgeCost[Edge{last, node}]
	if !ok {
		panic(fmt.Sprintf("no edge from %d to %d", last, node))
	}
	return Path{
		Length: p.Length + c,
		Nodes:  append(p.Nodes[0:len(p.Nodes):len(p.Nodes)], node),
	}
}

func findShortestPaths(src int) {
	visited := make(map[int]int)
	type QEntry struct {
		cost int
		from int
		to   int
	}
	q := priority_queue.NewFunc(func(a, b QEntry) bool {
		return a.cost < b.cost
	})
	q.Push(QEntry{0, src, src})
	for q.Len() > 0 && len(visited) < 23 {
		e := q.Pop()
		if _, ok := visited[e.to]; ok {
			continue
		}
		visited[e.to] = e.from
		for i := 0; i < 23; i++ {
			if _, ok := visited[i]; ok {
				continue
			}
			c, ok := EdgeCost[Edge{e.to, i}]
			if !ok {
				continue
			}
			q.Push(QEntry{e.cost + c, e.to, i})
		}
	}
	Paths[src][src] = Path{
		Length: 0,
		Nodes:  []int{src},
	}
	delete(visited, src)
	for len(visited) > 0 {
		for to, from := range visited {
			if len(Paths[src][from].Nodes) > 0 {
				Paths[src][to] = Paths[src][from].Append(to)
				delete(visited, to)
			}
		}
	}
}

func init() {
	for i := 0; i < 23; i++ {
		findShortestPaths(i)
	}
}

// CellCost are the energy used by amphipod type.
var CellCost = [...]int{
	A: 1,
	B: 10,
	C: 100,
	D: 1000,
}

// Homes are the home-cells for each type of amphipod.
var Homes = map[Cell][4]int{
	A: {7, 11, 15, 19},
	B: {8, 12, 16, 20},
	C: {9, 13, 17, 21},
	D: {10, 14, 18, 22},
}

func (c Cell) IsHome(i int) bool {
	for _, h := range Homes[c] {
		if h == i {
			return true
		}
	}
	return false
}

func (c Cell) Cost() int {
	switch c {
	case A:
		return 1
	case B:
		return 10
	case C:
		return 100
	case D:
		return 1000
	default:
		panic("invalid cell")
	}
}

func (c Cell) String() string {
	switch c {
	case A:
		return "A"
	case B:
		return "B"
	case C:
		return "C"
	case D:
		return "D"
	default:
		return " "
	}
}

func IsHallway(i int) bool {
	return i < 7
}

func (s State) Move(src, dst int) (cost int, next State, err error) {
	if s[src] == None {
		return 0, State{}, errors.New("source cell is unoccupied")
	}
	if s[dst] != None {
		return 0, State{}, errors.New("destination cell is occupied")
	}
	a := s[src]
	if !a.IsHome(dst) {
		if IsHallway(src) {
			return 0, State{}, errors.New("amphipod stands in hallway and can only move home")
		}
		if !IsHallway(dst) {
			return 0, State{}, errors.New("amphipod can only move home or into hallway")
		}
	} else {
		for _, h := range Homes[a] {
			if s[h] != None && s[h] != a {
				return 0, State{}, errors.New("amphipods home is occupied by wrong amphipod type")
			}
		}
	}
	p := Paths[src][dst]
	for _, n := range p.Nodes[1:] {
		if s[n] != None {
			return 0, State{}, errors.New("path is blocked")
		}
	}
	s[src], s[dst] = None, a
	return a.Cost() * p.Length, s, nil
}

type Neighbor struct {
	Cost  int
	State State
}

func (s State) Neighbors() []Neighbor {
	var out []Neighbor
	for i := 0; i < 23; i++ {
		for j := 0; j < 23; j++ {
			c, s, err := s.Move(i, j)
			if err != nil {
				continue
			}
			out = append(out, Neighbor{c, s})
		}
	}
	return out
}

func (s State) Dump() {
	fmt.Println("•••••••••••••")
	fmt.Printf("•%v%v %v %v %v %v%v•\n", s[0], s[1], s[2], s[3], s[4], s[5], s[6])
	fmt.Printf("•••%v•%v•%v•%v•••\n", s[7], s[8], s[9], s[10])
	fmt.Printf("  •%v•%v•%v•%v•\n", s[11], s[12], s[13], s[14])
	fmt.Printf("  •%v•%v•%v•%v•\n", s[15], s[16], s[17], s[18])
	fmt.Printf("  •%v•%v•%v•%v•\n", s[19], s[20], s[21], s[22])
	fmt.Println("  •••••••••")
}

func Organize(s State) int {
	type QEntry struct {
		prio int
		cost int
		from State
		to   State
	}
	q := priority_queue.NewFunc(func(a, b QEntry) bool {
		return a.prio < b.prio
	})
	visited := make(map[State]State)
	q.Push(QEntry{0, 0, s, s})
	update := 1000
	for q.Len() > 0 {
		e := q.Pop()
		if _, ok := visited[e.to]; ok {
			continue
		}
		visited[e.to] = e.from
		if e.prio > update {
			e.to.Dump()
			fmt.Println(e.prio)
			update += 1000
		}
		if e.to == End {
			return e.prio
		}
		for _, n := range e.to.Neighbors() {
			if _, ok := visited[n.State]; ok {
				continue
			}
			cost := e.cost + n.Cost
			q.Push(QEntry{cost, cost, e.to, n.State})
		}
	}
	panic("no solution found")
}
