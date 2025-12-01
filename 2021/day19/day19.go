package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"gonih.org/AdventOfCode/2021/go/set"
	"gonih.org/AdventOfCode/2021/go/vec"
)

func main() {
	log.SetFlags(log.Lshortfile)
	scanners, err := Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	done := set.Make(scanners[0])
	open := set.Make(scanners[1:]...)
	for len(open) > 0 {
		for s1 := range open {
			for s2 := range done {
				if !s1.Intersect(s2) {
					continue
				}
				s1.Transform = s2.Transform.Compose(s1.Transform)
				done.Add(s1)
				open.Delete(s1)
				break
			}
		}
	}

	beacons := make(set.Set[vec.V])
	for s := range done {
		for b := range s.Beacons {
			beacons.Add(s.Transform.Apply(b))
		}
	}
	log.Printf("There are %d beacons in total", len(beacons))

	m := math.MinInt
	for _, s1 := range scanners {
		for _, s2 := range scanners {
			d := vec.DistL1(s1.Pos(), s2.Pos())
			if d > m {
				m = d
			}
		}
	}
	log.Printf("Maximum distance between scanners is %d", m)
}

func Parse(r io.Reader) ([]*Scanner, error) {
	var (
		beacons  = make(set.Set[vec.V])
		scanners []*Scanner
	)
	s := bufio.NewScanner(r)
	for s.Scan() {
		if s.Text() == "" {
			scanners = append(scanners, NewScanner(len(scanners), beacons))
			beacons = make(set.Set[vec.V])
			continue
		}
		if strings.HasPrefix(s.Text(), "---") {
			continue
		}
		sp := strings.Split(s.Text(), ",")
		if len(sp) != 3 {
			return nil, errors.New("wrong number of components")
		}
		x, err := strconv.Atoi(sp[0])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(sp[1])
		if err != nil {
			return nil, err
		}
		z, err := strconv.Atoi(sp[2])
		if err != nil {
			return nil, err
		}
		beacons.Add(vec.Vec(x, y, z))
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(beacons) > 0 {
		scanners = append(scanners, NewScanner(len(scanners), beacons))
	}
	return scanners, nil
}

var Rotations = set.Make(
	vec.Mat(1, 0, 0, 0, 0, -1, 0, 1, 0),
	vec.Mat(0, 0, 1, 0, 1, 0, -1, 0, 0),
	vec.Mat(0, -1, 0, 1, 0, 0, 0, 0, 1),
	vec.Mat(1, 0, 0, 0, -1, 0, 0, 0, -1),
	vec.Mat(0, 0, 1, 1, 0, 0, 0, 1, 0),
	vec.Mat(0, -1, 0, 0, 0, -1, 1, 0, 0),
	vec.Mat(0, 1, 0, 0, 0, -1, -1, 0, 0),
	vec.Mat(-1, 0, 0, 0, 1, 0, 0, 0, -1),
	vec.Mat(0, 0, -1, 0, -1, 0, -1, 0, 0),
	vec.Mat(0, 1, 0, 1, 0, 0, 0, 0, -1),
	vec.Mat(0, -1, 0, 0, 0, 1, -1, 0, 0),
	vec.Mat(-1, 0, 0, 0, -1, 0, 0, 0, 1),
	vec.Mat(-1, 0, 0, 0, 0, 1, 0, 1, 0),
	vec.Mat(0, 0, 1, 0, -1, 0, 1, 0, 0),
	vec.Mat(0, -1, 0, -1, 0, 0, 0, 0, -1),
	vec.Mat(1, 0, 0, 0, 0, 1, 0, -1, 0),
	vec.Mat(0, 0, -1, 1, 0, 0, 0, -1, 0),
	vec.Mat(0, 1, 0, 0, 0, 1, 1, 0, 0),
	vec.Mat(-1, 0, 0, 0, 0, -1, 0, -1, 0),
	vec.Mat(0, 0, 1, -1, 0, 0, 0, -1, 0),
	vec.Mat(0, 0, -1, 0, 1, 0, 1, 0, 0),
	vec.Mat(0, 0, -1, -1, 0, 0, 0, 1, 0),
	vec.Mat(0, 1, 0, -1, 0, 0, 0, 0, 1),
	vec.Mat(1, 0, 0, 0, 1, 0, 0, 0, 1),
)

type Affine struct {
	A vec.M
	B vec.V
}

func (a Affine) Apply(v vec.V) vec.V {
	return a.A.MulV(v).Add(a.B)
}

func (a Affine) Compose(b Affine) Affine {
	return Affine{A: a.A.MulM(b.A), B: a.A.MulV(b.B).Add(a.B)}
}

func (a Affine) String() string {
	return fmt.Sprintf("%v•v+%v", a.A, a.B)
}

type Scanner struct {
	N         int
	Beacons   set.Set[vec.V]
	Transform Affine
}

func NewScanner(n int, beacons set.Set[vec.V]) *Scanner {
	return &Scanner{
		N:         n,
		Beacons:   beacons,
		Transform: Affine{A: vec.ID()},
	}
}

func (s1 *Scanner) Intersect(s2 *Scanner) bool {
	for r := range Rotations {
		rb := set.Map(s1.Beacons, r.MulV)
		Δ := CountDeltas(rb, s2.Beacons)
		for δ, n := range Δ {
			if n >= 12 {
				s1.Transform = Affine{
					A: r,
					B: δ,
				}
				return true
			}
		}
	}
	return false
}

func (s1 *Scanner) Pos() vec.V {
	return s1.Transform.B
}

func CountDeltas(b1, b2 set.Set[vec.V]) map[vec.V]int {
	Δ := make(map[vec.V]int)
	for v := range b1 {
		for w := range b2 {
			Δ[w.Sub(v)] += 1
		}
	}
	return Δ
}
