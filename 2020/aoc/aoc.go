package aoc

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

func SlurpNumbers(r io.Reader) ([]int, error) {
	line := 0
	var out []int
	s := bufio.NewScanner(r)
	for s.Scan() {
		line++
		n, err := strconv.Atoi(s.Text())
		if err != nil {
			return nil, fmt.Errorf("%d: %w", line, err)
		}
		out = append(out, n)
	}
	return out, s.Err()
}
