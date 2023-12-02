package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	buf = bytes.TrimSpace(buf)
	games, err := parse.Lines(parse.Struct[Game](
		split.On(": "),
		func(s string) (int, error) {
			idx, ok := strings.CutPrefix(s, "Game ")
			if !ok {
				return 0, fmt.Errorf(`expected "Game <idx>, got %q"`, s)
			}
			return strconv.Atoi(idx)
		},
		parse.Slice(
			split.On("; "),
			parse.Slice(
				split.On(", "),
				parse.Struct[Cube](
					split.Fields,
					parse.Signed[int],
					parse.Enum("red", "green", "blue"),
				),
			),
		),
	))(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	var part1 int
	for _, game := range games {
		if err := IsPossible(game, 12, 13, 14); err == nil {
			part1 += game.Index
		}
	}
	fmt.Printf("Sum of index of possible games: %d\n", part1)
	var part2 int
	for _, game := range games {
		r, g, b := MinCubes(game)
		part2 += r * g * b
	}
	fmt.Printf("Power of all games: %d\n", part2)
}

type Game struct {
	Index int
	Cubes [][]Cube
}

type Cube struct {
	N     int
	Color string
}

func IsPossible(game Game, r, g, b int) error {
	for _, cs := range game.Cubes {
		for _, c := range cs {
			switch c.Color {
			case "red":
				if c.N > r {
					return fmt.Errorf("%d red cubes", c.N)
				}
			case "green":
				if c.N > g {
					return fmt.Errorf("%d green cubes", c.N)
				}
			case "blue":
				if c.N > b {
					return fmt.Errorf("%d blue cubes", c.N)
				}
			}
		}
	}
	return nil
}

func MinCubes(game Game) (r, g, b int) {
	r, g, b = math.MinInt, math.MinInt, math.MinInt
	for _, cs := range game.Cubes {
		for _, c := range cs {
			switch c.Color {
			case "red":
				r = max(r, c.N)
			case "green":
				g = max(g, c.N)
			case "blue":
				b = max(b, c.N)
			}
		}
	}
	return r, g, b
}
