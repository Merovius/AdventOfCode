package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	races, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", Part1(races))
	fmt.Println("Part 2:", Part2(races))
}

func Parse(b []byte) ([]Race, error) {
	type parsed struct {
		Times     []int
		Distances []int
	}
	p, err := parse.Struct[parsed](
		split.Lines,
		parse.Prefix("Time:", parse.Slice(split.Fields, parse.Signed[int])),
		parse.Prefix("Distance:", parse.Slice(split.Fields, parse.Signed[int])),
	)(string(bytes.TrimSpace(b)))
	if err != nil {
		return nil, err
	}
	if len(p.Times) != len(p.Distances) {
		return nil, errors.New("different number of times and distances")
	}
	var r []Race
	for i := range p.Times {
		r = append(r, Race{p.Times[i], p.Distances[i]})
	}
	return r, nil
}

type Race struct {
	Time     int
	Distance int
}

func Part1(races []Race) int {
	out := 1
	for _, r := range races {
		var n int
		for i := 0; i <= r.Time; i++ {
			d := (r.Time - i) * i
			if d > r.Distance {
				n++
			}
		}
		out *= n
	}
	return out
}

func Part2(races []Race) int {
	w1, w2 := new(strings.Builder), new(strings.Builder)
	for _, r := range races {
		fmt.Fprint(w1, r.Time)
		fmt.Fprint(w2, r.Distance)
	}
	time, err := strconv.Atoi(w1.String())
	if err != nil {
		panic(err)
	}
	distance, err := strconv.Atoi(w2.String())
	if err != nil {
		panic(err)
	}
	var n int
	for i := 0; i <= time; i++ {
		d := (time - i) * i
		if d > distance {
			n++
		}
	}
	return n
}
