package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"gonih.org/AdventOfCode/internal/grid"
	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
	"gonih.org/AdventOfCode/internal/math"
)

func main() {
	log.SetFlags(0)
	w := flag.Int("w", 0, "width of the grid")
	h := flag.Int("h", 0, "height of the grid")
	flag.Parse()
	if *w <= 0 || *h <= 0 {
		log.Fatal("-w and -h must be given and non-negative")
	}

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(in, *w, *h))
	fmt.Println(Part2(in, *w, *h))
}

func Parse(in []byte) ([]Robot, error) {
	return parse.Slice(
		split.Lines,
		parse.Struct[Robot](
			split.On(" "),
			parse.Prefix(
				"p=",
				parse.Array[Vec](
					split.On(","),
					parse.Signed[int],
				),
			),
			parse.Prefix(
				"v=",
				parse.Array[Vec](
					split.On(","),
					parse.Signed[int],
				),
			),
		),
	)(string(in))
}

type Vec = [2]int

type Robot struct {
	P Vec
	V Vec
}

func Part1(in []Robot, w, h int) int {
	var t0, t1, t2, t3 int
	for _, r := range in {
		r.P[0] = math.Mod(r.P[0]+r.V[0]*100, w)
		r.P[1] = math.Mod(r.P[1]+r.V[1]*100, h)
		if r.P[0] < w/2 {
			if r.P[1] < h/2 {
				t0++
			} else if r.P[1] > h/2 {
				t1++
			}
		} else if r.P[0] > w/2 {
			if r.P[1] < h/2 {
				t2++
			} else if r.P[1] > h/2 {
				t3++
			}
		}
	}
	return t0 * t1 * t2 * t3
}

func Part2(in []Robot, w, h int) int {
	// We find a timestep where no two robots overlap. It turns out, there
	// is exactly one of those and it's the solution.
	//
	// After w steps, a robot's x-coordinate will be (px+vx*w)%w = px, so
	// it will be back where it started. Likewise for h and its
	// y-coordinate.
	//
	// So after w*h steps, every robot must be back where it started. So we
	// only need to check w*h steps.
	g := grid.New[bool](w, h)
times:
	for i := range w * h {
		for _, r := range in {
			r.P[0] = math.Mod(r.P[0]+r.V[0]*i, w)
			r.P[1] = math.Mod(r.P[1]+r.V[1]*i, h)
			p := grid.Pos{r.P[1], r.P[0]}
			if g.At(p) {
				clear(g.G)
				continue times
			}
			g.Set(p, true)
		}
		return i
	}
	return -1
}
