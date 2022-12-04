package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := ReadInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("There are %d pairs where one contains the other\n", len(Filter(input, EitherContains)))
	fmt.Printf("There are %d overlapping pairs\n", len(Filter(input, Overlap)))
}

type Range struct {
	Min int
	Max int
}

func ReadInput(r io.Reader) ([][2]Range, error) {
	var out [][2]Range
	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		a, b, ok := strings.Cut(l, ",")
		if !ok {
			return nil, fmt.Errorf("invalid input line %q", l)
		}
		mina, maxa, ok := strings.Cut(a, "-")
		if !ok {
			return nil, fmt.Errorf("invalid input line %q", l)
		}
		minb, maxb, ok := strings.Cut(b, "-")
		if !ok {
			return nil, fmt.Errorf("invalid input line %q", l)
		}
		var (
			pair [2]Range
			err  error
		)
		pair[0].Min, err = strconv.Atoi(mina)
		if err != nil {
			return nil, fmt.Errorf("invalid input line %q: %w", l, err)
		}
		pair[0].Max, err = strconv.Atoi(maxa)
		if err != nil {
			return nil, fmt.Errorf("invalid input line %q: %w", l, err)
		}
		pair[1].Min, err = strconv.Atoi(minb)
		if err != nil {
			return nil, fmt.Errorf("invalid input line %q: %w", l, err)
		}
		pair[1].Max, err = strconv.Atoi(maxb)
		if err != nil {
			return nil, fmt.Errorf("invalid input line %q: %w", l, err)
		}
		out = append(out, pair)
	}
	return out, s.Err()
}

func Contains(a, b Range) bool {
	return a.Min <= b.Min && a.Max >= b.Max
}

func EitherContains(a, b Range) bool {
	return Contains(a, b) || Contains(b, a)
}

func Overlap(a, b Range) bool {
	return !(a.Max < b.Min || b.Max < a.Min)
}

func Filter(p [][2]Range, include func(a, b Range) bool) [][2]Range {
	var out [][2]Range
	for _, p := range p {
		if include(p[0], p[1]) {
			out = append(out, p)
		}
	}
	return out
}
