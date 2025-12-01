package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gonih.org/AdventOfCode/internal/input/parse"
)

func main() {
	log.SetFlags(log.Lshortfile)
	data, err := parse.Lines(parse.String[string]).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var sum int
	for _, s := range data {
		if int2snafu(snafu2int(s)) != s {
			log.Printf("int2snafu(%q) = %d, snafu2int(%d) = %q", s, snafu2int(s), snafu2int(s), int2snafu(snafu2int(s)))
		}
		sum += snafu2int(s)
	}
	fmt.Println("decimal:", sum)
	fmt.Println("snafu:", int2snafu(sum))
}

func snafu2int(s string) int {
	var v int
	for _, r := range s {
		v *= 5
		if i := strings.IndexRune("=-0123", r); i >= 0 {
			v += i - 2
		} else {
			panic(fmt.Errorf("invalid input byte %q", r))
		}
	}
	return v
}

func int2snafu(v int) string {
	var b []byte
	for v > 0 {
		d := v % 5
		b = append(b, "012=-"[d])
		if d > 2 {
			d -= 5
		}
		v = (v - d) / 5
	}
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - 1 - i
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
