package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)`)
	ms := re.FindAll(buf, -1)
	var total int
	for _, m := range ms {
		a, b, _ := strings.Cut(string(m[4:len(m)-1]), ",")
		ai, err := strconv.Atoi(a)
		if err != nil {
			log.Fatal(err)
		}
		bi, err := strconv.Atoi(b)
		if err != nil {
			log.Fatal(err)
		}
		total += ai * bi
	}
	fmt.Println(total)

	re = regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)|do\(\)|don't\(\)`)
	ms = re.FindAll(buf, -1)
	total = 0
	do := true
	for _, m := range ms {
		if string(m) == "do()" {
			do = true
			continue
		} else if string(m) == "don't()" {
			do = false
			continue
		}
		if !do {
			continue
		}
		a, b, _ := strings.Cut(string(m[4:len(m)-1]), ",")
		ai, err := strconv.Atoi(a)
		if err != nil {
			log.Fatal(err)
		}
		bi, err := strconv.Atoi(b)
		if err != nil {
			log.Fatal(err)
		}
		total += ai * bi
	}
	fmt.Println(total)
}
