package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	points, folds, err := parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	points.Apply(folds[0])
	fmt.Printf("After first fold, there are %d points\n", len(points))
	for _, f := range folds[1:] {
		points.Apply(f)
	}
	fmt.Println("After all folds the paper looks like:")
	points.PrettyPrint()
}

type Axis int

const (
	AxisX Axis = iota
	AxisY
)

type Fold struct {
	Axis  Axis
	Value int
}

type Point struct {
	X int
	Y int
}

func (p Point) Apply(f Fold) Point {
	switch f.Axis {
	case AxisX:
		if p.X > f.Value {
			return Point{2*f.Value - p.X, p.Y}
		}
	case AxisY:
		if p.Y > f.Value {
			return Point{p.X, 2*f.Value - p.Y}
		}
	}
	return p
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

type PointSet map[Point]struct{}

func (s PointSet) Contains(p Point) bool {
	_, ok := s[p]
	return ok
}

func (s PointSet) Add(p Point) {
	s[p] = struct{}{}
}

func (s PointSet) Remove(p Point) {
	delete(s, p)
}

func (s PointSet) Apply(f Fold) {
	var additions []Point
	for p := range s {
		q := p.Apply(f)
		if p != q {
			delete(s, p)
			additions = append(additions, q)
		}
	}
	for _, p := range additions {
		s.Add(p)
	}
}

func (s PointSet) PrettyPrint() {
	mx, my := -1, -1
	for p := range s {
		if p.X > mx {
			mx = p.X
		}
		if p.Y > my {
			my = p.Y
		}
	}
	for y := 0; y <= my; y++ {
		for x := 0; x <= mx; x++ {
			if s.Contains(Point{x, y}) {
				fmt.Print("•")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func parse(r io.Reader) (PointSet, []Fold, error) {
	points := make(PointSet)
	var folds []Fold

	s := bufio.NewScanner(r)
	for s.Scan() {
		if s.Text() == "" {
			break
		}
		sp := strings.Split(s.Text(), ",")
		if len(sp) != 2 {
			return nil, nil, errors.New("invalid point")
		}
		x, err := strconv.Atoi(sp[0])
		if err != nil {
			return nil, nil, err
		}
		y, err := strconv.Atoi(sp[1])
		if err != nil {
			return nil, nil, err
		}
		p := Point{x, y}
		if points.Contains(p) {
			return nil, nil, fmt.Errorf("duplicate point %v", p)
		}
		points.Add(p)
	}
	for s.Scan() {
		l := s.Text()
		if !strings.HasPrefix(l, "fold along ") {
			return nil, nil, errors.New("invalid fold instruction")
		}
		l = strings.TrimPrefix(l, "fold along ")
		sp := strings.Split(l, "=")
		val, err := strconv.Atoi(sp[1])
		if err != nil {
			return nil, nil, err
		}
		switch sp[0] {
		case "x":
			folds = append(folds, Fold{AxisX, val})
		case "y":
			folds = append(folds, Fold{AxisY, val})
		default:
			return nil, nil, fmt.Errorf("invalid axis %q", sp[0])
		}
	}
	return points, folds, s.Err()
}
