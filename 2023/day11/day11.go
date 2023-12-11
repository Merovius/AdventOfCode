package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/grid"
	"github.com/Merovius/AdventOfCode/internal/slices"
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(string(b))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Part1(in))
	fmt.Println("Part 2:", Part2(in))
}

// Parse returns a list of galaxy positions from the input.
func Parse(in string) ([]grid.Pos, error) {
	g := make([]grid.Pos, 0, 1024)
	for r := 0; len(in) > 0; r++ {
		i := strings.IndexByte(in, '\n')
		if i < 0 {
			i = len(in)
		}
		l := in[:i]
		var lastC int
		for len(l) > 0 {
			c := strings.IndexByte(l, '#')
			if c < 0 {
				break
			}
			g = append(g, grid.Pos{r, c + lastC})
			lastC += c + 1
			l = l[c+1:]
		}
		in = in[i+1:]
	}
	return g, nil
}

func Part1(in []grid.Pos) int {
	return SumDistances(in, 1)
}

func Part2(in []grid.Pos) int {
	return SumDistances(in, 999999)
}

// SumDistances returns the sum of the minimal distance between all pairs of galaxies.
func SumDistances(in []grid.Pos, age int) int {
	// The minimal distance is the manhatten distance.
	// We can calculate the distances for rows and columns separately, sorting
	// them individually, as addition is commutative.
	// That makes the overall complexity n•log(n)+n²: First, we sort. Then we
	// can walk the "gaps" together with the inner loop of the pair - we don't
	// have to consider gaps that are before the first element of the pair or
	// the last second element we considered.
	// We also can run the calculation for rows/columns in parallel.
	ch := make(chan int, 2)
	go func() {
		rows := make([]int, 0, 1024)
		for _, g := range in {
			rows = append(rows, g.Row)
		}
		slices.Sort(rows)
		ch <- sumDistances(rows, age)
	}()
	go func() {
		cols := make([]int, 0, 1024)
		for _, g := range in {
			cols = append(cols, g.Col)
		}
		slices.Sort(cols)
		ch <- sumDistances(cols, age)
	}()
	return <-ch + <-ch
}

// sumDistances returns the sum of distances between all pairs in vs. Gaps in
// vs contribute age extra to the distance.
//
// vs must be sorted.
func sumDistances(vs []int, age int) int {
	empty := make([]int, 0, vs[len(vs)-1])
loop:
	for w := range vs[len(vs)-1] {
		for _, v := range vs {
			if v == w {
				continue loop
			}
		}
		empty = append(empty, w)
	}
	var (
		sum int
		e   int
	)
	for i, v1 := range vs {
		e += search(empty[e:], v1)
		last := v1 // last value considered for inner loop
		δ := 0     // distance between v1 and last
		e2 := e    // offset into empty for last
		for _, v2 := range vs[i+1:] {
			// accumulate extra distances to v1
			δe := search(empty[e2:], v2)
			δ += v2 - last + δe*age
			e2 += δe
			sum += δ
			last = v2
		}
	}
	return sum
}

// search returns the smallest index i of s where s[i] > v.
func search(s []int, v int) int {
	// our slices are short, so linear search is faster than binary search.
	for i, w := range s {
		if w > v {
			return i
		}
	}
	return len(s)
}
