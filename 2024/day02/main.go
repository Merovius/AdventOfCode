package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
	"gonih.org/AdventOfCode/internal/math"
	"golang.org/x/exp/constraints"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	reports, err := parse.Lines(parse.Slice(split.Fields, parse.Signed[int]))(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	var N int
	for _, r := range reports {
		if safe(r) {
			N++
		}
	}
	fmt.Println(N)
	N = 0
	for _, r := range reports {
		if safe(r) {
			N++
			continue
		}
		r2 := slices.Clone(r)
		for i := len(r) - 1; i >= 0; i-- {
			r2 = append(r2[:i], r[i+1:]...)
			if safe(r2) {
				N++
				break
			}
		}
	}
	fmt.Println(N)
}

func abs[T constraints.Signed](a T) T {
	if a < 0 {
		return -a
	}
	return a
}

func safe(r []int) bool {
	sgn := math.Sgn(r[1] - r[0])
	for i := 1; i < len(r); i++ {
		δ := r[i] - r[i-1]
		if a := math.Abs(δ); math.Sgn(δ) != sgn || a < 1 || a > 3 {
			return false
		}
	}
	return true
}
