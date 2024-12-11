package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Merovius/AdventOfCode/internal/input/parse"
	"github.com/Merovius/AdventOfCode/internal/input/split"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	input, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(input))
	fmt.Println(Part2(input))
}

func Parse(in []byte) ([]int, error) {
	return parse.Slice(split.Fields, parse.Signed[int])(string(in))
}

func Part1(input []int) int {
	var total int
	for _, v := range input {
		total += rec(v, 25)
	}
	return total
}

func Part2(input []int) int {
	var total int
	for _, v := range input {
		total += rec(v, 75)
	}
	return total
}

var memo = make(map[[2]int]int)

func rec(v, n int) (m int) {
	if m, ok := memo[[2]int{v, n}]; ok {
		return m
	}
	defer func() { memo[[2]int{v, n}] = m }()

	if n == 0 {
		return 1
	}
	if v == 0 {
		return rec(1, n-1)
	} else if s := strconv.Itoa(v); len(s)%2 == 0 {
		l, _ := strconv.Atoi(s[len(s)/2:])
		r, _ := strconv.Atoi(s[:len(s)/2])
		return rec(l, n-1) + rec(r, n-1)
	} else {
		return rec(v*2024, n-1)
	}
}
