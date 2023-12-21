//go:build goexperiment.rangefunc

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/container"
	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/set"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	g, start, err := Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Part1(g, start))
	fmt.Println("Part 2:", Part2(g, start))
}

func Parse(s string) (*grid.Grid[Cell], grid.Pos, error) {
	g, err := grid.Read(strings.NewReader(s), grid.Enum[Cell]('.', '#', 'S'))
	if err != nil {
		return nil, grid.Pos{}, err
	}
	start := grid.Pos{-1, -1}
	for p := range grid.Find(g, Start) {
		if start != (grid.Pos{-1, -1}) {
			return nil, grid.Pos{}, errors.New("multiple start positions")
		}
		start = p
	}
	g.Set(start, Garden)
	return g, start, nil
}

func Part1(g *grid.Grid[Cell], start grid.Pos) int {
	return NReachable(g, start, 64)
}

func Part2(g *grid.Grid[Cell], start grid.Pos) int {
	// Due to the structure of the input, the reachable nodes form a diamond
	// every w•i+w/2 steps (where w is the width of the grid and i is a natural
	// number) - and in particular, n%w == w/2 for the interesting number of
	// steps.
	//
	// The area of reachable nodes for those diamonds is proportional to the
	// radius of the diamond (i.e. number of used steps), so it is a second
	// degree polynomial.
	//
	// The algebra is easier, if we assume the input is i in the w•i+w/2. We
	// map the input from the step number in the final polynomial.
	//
	// This logic only works for the actual input, not for the example.
	n := 26501365
	w := g.W
	if w/2 != n%w {
		panic(fmt.Errorf("%d/2 = %d but %d%%%d = %d", w, w/2, n, g.W, n%w))
	}
	x := []int{
		n%w + w*0,
		n%w + w*1,
		n%w + w*2,
	}
	y := []int{
		NReachable(&cover{g}, start, x[0]),
		NReachable(&cover{g}, start, x[1]),
		NReachable(&cover{g}, start, x[2]),
	}
	// three points are enough to interpolate a polynomial
	// we have p(i) = y[i] = a•x[i]²+b•x[i]+c

	// p(0) = a•0²+b•0+c
	c := y[0]
	// p(1) = a+b+c
	// p(2) = 4a+2b+c
	// ⇒ 4•p(1)-p(2) = 2•b+3•c
	// ⇒ b = (4•p(1)-p(2)-3•c)/2
	b := (4*y[1] - y[2] - 3*c) / 2
	// p(1) = a+b+c
	// ⇒ a = p(1)-b-c
	a := y[1] - b - c

	p := func(v int) int {
		return a*v*v + b*v + c
	}
	// map actual step numbers into the amount we are interested in.
	// This isn't correct for step numbers not of the w•i+n%w form, but
	// whatever.
	q := func(v int) int {
		return p((v - n%w) / w)
	}
	for i := 0; i < 3; i++ {
		if z := q(x[i]); z != y[i] {
			panic(fmt.Errorf("p(%d) = %d, want %d", x[i], z, y[i]))
		}
	}
	return q(n)
}

type Cell uint8

const (
	Garden Cell = iota
	Rocks
	Start
	Reachable
)

type Grid interface {
	At(grid.Pos) Cell
	Neigh4(grid.Pos) []grid.Pos
}

// universal cover of Grid - presents (r,c) as (r%g.H, c%g.C) of the underlying
// grid.
type cover struct {
	*grid.Grid[Cell]
}

func (w *cover) At(p grid.Pos) Cell {
	p.Row = p.Row % w.H
	if p.Row < 0 {
		p.Row += w.H
	}
	p.Col = p.Col % w.W
	if p.Col < 0 {
		p.Col += w.W
	}
	return w.Grid.At(p)
}

func (w *cover) Neigh4(p grid.Pos) []grid.Pos {
	out := make([]grid.Pos, 0, 4)
	for _, δ := range []grid.Pos{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} {
		q := p.Add(δ)
		if w.At(q) != Rocks {
			out = append(out, q)
		}
	}
	return out
}

func NReachable(g Grid, start grid.Pos, n int) int {
	type node struct {
		pos   grid.Pos
		steps int
	}
	var (
		q     container.FIFO[node]
		seen  = make(set.Set[grid.Pos])
		total int
	)

	q.Push(node{start, 0})
	seen.Add(start)
	for q.Len() > 0 {
		st := q.Pop()
		if st.steps > n {
			continue
		}
		if (st.steps % 2) == (n % 2) {
			total++
		}
		for _, neigh := range g.Neigh4(st.pos) {
			if g.At(neigh) != Garden {
				continue
			}
			if seen.Contains(neigh) {
				continue
			}
			seen.Add(neigh)
			q.Push(node{neigh, st.steps + 1})
		}
	}
	return total
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func Print(g *grid.Grid[Cell]) {
	for r := 0; r < g.H; r++ {
		for c := 0; c < g.H; c++ {
			switch g.At(grid.Pos{r, c}) {
			case Garden:
				fmt.Print(" ")
			case Rocks:
				fmt.Print("█")
			case Start:
				fmt.Print("S")
			case Reachable:
				fmt.Print("•")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

//func NReachable2(g Grid, start grid.Pos, n int) int {
//	type visit struct {
//		pos   grid.Pos // grid-copy, relative to the start grid
//		start grid.Pos // where the copy is entered
//		steps int      // how many steps we have left
//	}
//
//	var (
//		q          container.FIFO[visit]
//		seen       = make(set.Set[grid.Pos])
//		total      int
//		reachables = make(map[grid.Pos][]int)
//		neighbors  = make(map[grid.Pos][]visit)
//	)
//	var minR, maxR, minC, maxC int
//	push := func(v visit) {
//		if seen.Contains(v.pos) {
//			return
//		}
//		seen.Add(v.pos)
//		q.Push(v)
//		if v.pos.Col < minC || v.pos.Col > maxC || v.pos.Row < minR || v.pos.Row > maxR {
//			minC = min(v.pos.Col, minC)
//			maxC = max(v.pos.Col, maxC)
//			minR = min(v.pos.Row, minR)
//			maxR = max(v.pos.Row, maxR)
//		}
//		if len(seen)%1000 == 0 {
//			log.Println(len(seen), minR, minC, maxR, maxC)
//		}
//	}
//	findNeighbors := func(v visit) []visit {
//		var (
//			vs []visit
//			ok bool
//		)
//		if vs, ok = neighbors[v.start]; !ok {
//			s, p := FindRow(g, v.start, 0)
//			vs = append(vs, visit{
//				grid.Pos{-1, 0},
//				grid.Pos{g.H - 1, p.Col},
//				s + 1,
//			})
//			s, p = FindRow(g, v.start, g.H-1)
//			vs = append(vs, visit{
//				grid.Pos{1, 0},
//				grid.Pos{0, p.Col},
//				s + 1,
//			})
//			s, p = FindCol(g, v.start, 0)
//			vs = append(vs, visit{
//				grid.Pos{0, -1},
//				grid.Pos{p.Row, g.W - 1},
//				s + 1,
//			})
//			s, p = FindCol(g, v.start, g.W-1)
//			vs = append(vs, visit{
//				grid.Pos{0, 1},
//				grid.Pos{p.Row, 0},
//				s + 1,
//			})
//			s = FindPos(g, v.start, grid.Pos{0, 0})
//			vs = append(vs, visit{
//				grid.Pos{-1, -1},
//				grid.Pos{g.H - 1, g.W - 1},
//				s + 2,
//			})
//			s = FindPos(g, v.start, grid.Pos{0, g.W - 1})
//			vs = append(vs, visit{
//				grid.Pos{-1, 1},
//				grid.Pos{g.H - 1, 0},
//				s + 2,
//			})
//			s = FindPos(g, v.start, grid.Pos{g.H - 1, g.W - 1})
//			vs = append(vs, visit{
//				grid.Pos{1, 1},
//				grid.Pos{0, 0},
//				s + 2,
//			})
//			s = FindPos(g, v.start, grid.Pos{g.H - 1, 0})
//			vs = append(vs, visit{
//				grid.Pos{1, -1},
//				grid.Pos{0, g.W - 1},
//				s + 2,
//			})
//			neighbors[v.start] = vs
//		}
//		vs = slices.Clone(vs)
//		for i := range vs {
//			vs[i].pos = vs[i].pos.Add(v.pos)
//			vs[i].steps = v.steps - vs[i].steps
//		}
//		return vs
//	}
//
//	push(visit{grid.Pos{0, 0}, start, n})
//	for q.Len() > 0 {
//		v := q.Pop()
//		if v.steps < 0 {
//			continue
//		}
//		r, ok := reachables[v.start]
//		if !ok {
//			r = Reachables(g, v.start)
//			reachables[v.start] = r
//		}
//		var m int
//		if v.steps >= len(r) {
//			m = r[len(r)-1]
//		} else {
//			m = r[v.steps]
//		}
//		total += m
//		for _, n := range findNeighbors(v) {
//			push(n)
//		}
//	}
//	return total
//}
//
//func Reachables(g *Grid, start grid.Pos) []int {
//	type node struct {
//		p grid.Pos
//		n int
//	}
//	var (
//		q    container.FIFO[node]
//		seen = make(set.Set[grid.Pos])
//		out  []int
//	)
//	q.Push(node{start, 0})
//	for q.Len() > 0 {
//		n := q.Pop()
//		if len(out) <= n.n {
//			if len(out) >= 2 {
//				out = append(out, out[len(out)-2])
//			} else {
//				out = append(out, 0)
//			}
//		}
//		out[n.n]++
//		for _, p := range g.Neigh4(n.p) {
//			if !seen.Contains(p) {
//				q.Push(node{p, n.n + 1})
//				seen.Add(p)
//			}
//		}
//	}
//	return out
//}
//
//// FindRow finds the shortest path from start to row r. It returns the number
//// of steps that takes and the position where row r is entered.
//func FindRow(g Grid, start grid.Pos, r int) (int, grid.Pos) {
//	goal := func(p grid.Pos) bool {
//		return p.Row == r
//	}
//	h := func(p grid.Pos) float64 {
//		v := math.Abs(float64(p.Row - r))
//		// Give a slight preference towards cells in the middle column.
//		// This increases the chance of cache-hits.
//		v -= math.Abs(float64(p.Col)-float64(g.W)/2) / float64(g.W)
//		return v
//	}
//	path := graph.AStar(&Graph{g}, start, goal, h)
//	if len(path) == 0 {
//		if !goal(start) {
//			panic(fmt.Errorf("row %d is unreachable from %v", r, start))
//		}
//		return 0, start
//	}
//	return len(path), path[len(path)-1][1]
//}
//
//// FindCol finds the shortest path from start to column c. It returns the
//// number of steps that takes and the position where column c is entered.
//func FindCol(g Grid, start grid.Pos, c int) (int, grid.Pos) {
//	goal := func(p grid.Pos) bool {
//		return p.Col == c
//	}
//	h := func(p grid.Pos) float64 {
//		v := math.Abs(float64(p.Col - c))
//		// Give a slight preference towards cells in the middle row.
//		// This increases the chance of cache-hits.
//		v -= math.Abs(float64(p.Row)-float64(g.H)/2) / float64(g.H)
//		return v
//	}
//	path := graph.AStar(&Graph{g}, start, goal, h)
//	if len(path) == 0 {
//		if !goal(start) {
//			panic(fmt.Errorf("column %d is unreachable from %v", c, start))
//		}
//		return 0, start
//	}
//	return len(path), path[len(path)-1][1]
//}
//
//// FindPos  finds the shortest path from start to end. It returns the number of
//// steps that takes.
//func FindPos(g Grid, start, end grid.Pos) int {
//	if start == end {
//		return 0
//	}
//	goal := func(p grid.Pos) bool {
//		return p == end
//	}
//	h := func(p grid.Pos) float64 {
//		return float64(end.Sub(p).Length())
//	}
//	path := graph.AStar(&Graph{g}, start, goal, h)
//	if len(path) == 0 {
//		panic(fmt.Errorf("%v is unreachable from %v", end, start))
//	}
//	return len(path)
//}
//
//type Graph struct{ g Grid }
//
//func (g *Graph) Edges(p grid.Pos) [][2]grid.Pos {
//	var e [][2]grid.Pos
//	for _, q := range g.g.Neigh4(p) {
//		if g.g.At(q) != Garden {
//			continue
//		}
//		e = append(e, [2]grid.Pos{p, q})
//	}
//	return e
//}
//
//func (g *Graph) From(e [2]grid.Pos) grid.Pos {
//	return e[0]
//}
//
//func (g *Graph) To(e [2]grid.Pos) grid.Pos {
//	return e[1]
//}
//
//func (g *Graph) Weight(e [2]grid.Pos) float64 {
//	return 1.0
//}
//
//func Copy(g *Grid, scale int) *Grid {
//	h := grid.New[Cell](scale*g.W, scale*g.H)
//	for r := 0; r < h.H; r++ {
//		for c := 0; c < h.W; c++ {
//			h.Set(grid.Pos{r, c}, g.At(grid.Pos{r % g.H, c % g.W}))
//		}
//	}
//	return h
//}
