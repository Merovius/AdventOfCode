package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"gonih.org/AdventOfCode/internal/input/parse"
	"gonih.org/AdventOfCode/internal/input/split"
)

func main() {
	log.SetFlags(0)

	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	in, err := Parse(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Part1(in))
	fmt.Println(Part2(in))
}

func Parse(in []byte) ([]string, error) {
	return parse.Slice(split.Lines, parse.String[string])(string(in))
}

func Part1(in []string) int {
	var n int
	for _, s := range in {
		var (
			vowels     int
			double     bool
			hasSpecial bool
		)
		for i := range len(s) {
			if strings.IndexByte("aeiou", s[i]) >= 0 {
				vowels++
			}
			if i > 0 {
				double = double || (s[i-1] == s[i])
				if x := s[i-1 : i+1]; x == "ab" || x == "cd" || x == "pq" || x == "xy" {
					hasSpecial = true
					break
				}
			}
		}
		if vowels > 2 && double && !hasSpecial {
			n++
		}
	}
	return n
}

func Part2(in []string) int {
	var n int
	for _, s := range in {
		var rule1 bool
		for i := range len(s) - 1 {
			if j := strings.LastIndex(s, s[i:i+2]); j >= i+2 {
				rule1 = true
			}
		}
		if !rule1 {
			continue
		}
		for i := range len(s) - 2 {
			if s[i] == s[i+2] {
				n++
				break
			}
		}
	}
	return n
}
