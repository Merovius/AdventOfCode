package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	data, err := parse.Lines(
		parse.Array[[2]int](
			split.Fields,
			parse.MapParser(
				parse.Rune,
				func(r rune) (int, error) {
					if r >= 'A' && r <= 'C' {
						return int(r - 'A'), nil
					}
					if r >= 'X' && r <= 'Z' {
						return int(r - 'X'), nil
					}
					return 0, fmt.Errorf("invalid character %q", r)
				},
			),
		),
	).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Predicted score: %d\n", Score(data))
	fmt.Printf("Real score: %d\n", RealScore(data))
}

func Score(guide [][2]int) int {
	var score int
	for _, r := range guide {
		score += RoundScore(r[0], r[1])
	}
	return score
}

func RealScore(guide [][2]int) int {
	var score int
	for _, r := range guide {
		play := (3 + r[0] + (r[1] - 1)) % 3
		score += RoundScore(r[0], play)
	}
	return score
}

func RoundScore(they, we int) int {
	// 0: draw, 1: we lost, 2: we won
	s := (3 + they - we) % 3
	switch s {
	case 0:
		s = 3
	case 1:
		s = 0
	case 2:
		s = 6
	default:
		panic("unreachable")
	}
	return s + 1 + we
}
