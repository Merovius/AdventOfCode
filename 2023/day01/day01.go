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
)

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	b = bytes.TrimSpace(b)
	lines, err := parse.Lines(parse.String[string])(string(b))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part1:", Part1(lines))
	fmt.Println("Part2:", Part2(lines))
}

func Parse(s string) ([]string, error) {
	return parse.Lines(parse.String[string])(s)
}

func Part1(lines []string) int {
	var pieces = map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
	}

	return sumValues(lines, pieces)
}

func Part2(lines []string) int {
	var pieces = map[string]int{
		"0":     0,
		"1":     1,
		"2":     2,
		"3":     3,
		"4":     4,
		"5":     5,
		"6":     6,
		"7":     7,
		"8":     8,
		"9":     9,
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
		"six":   6,
		"seven": 7,
		"eight": 8,
		"nine":  9,
	}
	return sumValues(lines, pieces)
}

func sumValues(lines []string, m map[string]int) int {
	var n int
	for _, l := range lines {
		var (
			i, j   int = math.MaxInt, math.MinInt
			vi, vj int
		)
		for s, v := range m {
			if ii := strings.Index(l, s); ii >= 0 {
				if ii < i {
					i, vi = ii, v
				}
			}
			if jj := strings.LastIndex(l, s); jj >= 0 {
				if jj > j {
					j, vj = jj, v
				}
			}
		}
		v, err := strconv.Atoi(fmt.Sprint(vi) + fmt.Sprint(vj))
		if err != nil {
			panic(err)
		}
		n += v
	}
	return n
}
