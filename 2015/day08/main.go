package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
	return split.Lines(string(in))
}

func Part1(in []string) int {
	var total int
	for _, s := range in {
		total += len(s) - UnquotedLen(s)
	}
	return total
}

func UnquotedLen(s string) int {
	if !strings.HasPrefix(s, `"`) || !strings.HasSuffix(s, `"`) {
		panic("invalid string")
	}
	s = s[1 : len(s)-1]

	var n int
	for i := 0; i < len(s); i++ {
		n++
		if s[i] != '\\' {
			continue
		}
		if s[i+1] == 'x' {
			i += 3
		} else {
			i++
		}
	}
	return n
}

func Part2(in []string) int {
	var total int
	for _, s := range in {
		total += QuotedLen(s) - len(s)
	}
	return total
}

func QuotedLen(s string) int {
	var n int
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"', '\\':
			n += 2
		default:
			n++
		}
	}
	return n + 2
}
