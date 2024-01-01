package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	games, err := Parse(string(buf))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part1:", Part1(games))
	fmt.Println("Part2:", Part2(games))
}

func Parse(s string) ([]Game, error) {
	return parse.Lines(parse.Struct[Game](
		split.On(": "),
		parse.Prefix("Game ", parse.Signed[int]),
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
	))(s)
}

type Game struct {
	Index int
	Cubes [][]Cube
}

type Cube struct {
	N     int
	Color string
}

func Part1(games []Game) int {
	var n int
	for _, game := range games {
		if err := IsPossible(game, 12, 13, 14); err == nil {
			n += game.Index
		}
	}
	return n
}

func Part2(games []Game) int {
	var n int
	for _, game := range games {
		r, g, b := MinCubes(game)
		n += r * g * b
	}
	return n
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
