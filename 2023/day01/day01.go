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
	var part1 = map[string]int{
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
	if v, err := sumValues(lines, part1); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Sum of calibration values (part1):", v)
	}

	var part2 = map[string]int{
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
	if v, err := sumValues(lines, part2); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Sum of calibration values (part2):", v)
	}
}

func sumValues(lines []string, m map[string]int) (int, error) {
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
			return 0, err
		}
		n += v
	}
	return n, nil
}
