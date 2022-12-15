package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/input"
	"github.com/Merovius/AdventOfCode/internal/interval"
	"github.com/Merovius/AdventOfCode/internal/math"
	"github.com/Merovius/AdventOfCode/internal/set"
)

func main() {
	example := flag.Bool("example", false, "input is example data")
	flag.Parse()

	data, err := input.Lines(func(s string) (Scanner, error) {
		var sc Scanner
		i, err := fmt.Sscanf(s, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sc.Pos.X, &sc.Pos.Y, &sc.Beacon.X, &sc.Beacon.Y)
		if err != nil || i != 4 {
			return sc, fmt.Errorf("invalid line %q", s)
		}
		sc.R = sc.Pos.Dist(sc.Beacon)
		return sc, nil
	}).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	row, max := 2000000, 4000000
	if *example {
		row, max = 10, 20
	}

	is := CoveredIntervals(data, row)
	fmt.Printf("There are %v cells covered in row %d\n", is.Len(), row)

	p := Find(data, max)
	fmt.Printf("Beacon is at %v and has tuning frequency %v\n", p, p.X*4000000+p.Y)
}

type Scanner struct {
	Pos    Point
	Beacon Point
	R      int
}

type Point struct {
	X int
	Y int
}

func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Len() int {
	return math.Abs(p.X) + math.Abs(p.Y)
}

func (p Point) Dist(q Point) int {
	return p.Sub(q).Len()
}

func (p Point) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

func (s Scanner) IntersectRow(y int) (interval.I[int], bool) {
	R := s.Pos.Dist(s.Beacon)
	δy := math.Abs(s.Pos.Y - y)
	if δy > R {
		return interval.I[int]{}, false
	}
	return interval.I[int]{Min: s.Pos.X - (R - δy), Max: s.Pos.X + (R - δy)}, true
}

func CoveredIntervals(s []Scanner, row int) *interval.Set[int] {
	covered := new(interval.Set[int])
	for _, s := range s {
		i, ok := s.IntersectRow(row)
		if !ok {
			continue
		}
		covered.Add(i)
	}
	return covered
}

func Find(s []Scanner, bound int) Point {
	occupied := make(set.Set[Point])
	for _, s := range s {
		occupied.Add(s.Pos)
		occupied.Add(s.Beacon)
	}
	valid := interval.I[int]{Min: 0, Max: bound}
	for y := 0; y <= bound; y++ {
		excluded := CoveredIntervals(s, y)
		excluded.Intersect(valid)
		is := excluded.Intervals()
		if len(is) <= 1 {
			continue
		}
		return Point{is[0].Max + 1, y}
	}
	panic("not found")
}
