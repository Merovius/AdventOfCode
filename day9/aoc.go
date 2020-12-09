package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	vs, err := slurpNumbers(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	v, ok := findInvalidNumber(vs)
	if !ok {
		log.Fatal("No invalid number found")
	}
	fmt.Println("Invalid number:", v)
	sum := findSum(vs, v)
	fmt.Println("Range:", sum)
	min, max := bounds(sum)
	fmt.Printf("%d + %d = %d\n", min, max, min+max)
}

const bufSize = 25

func slurpNumbers(r io.Reader) ([]int, error) {
	var out []int
	s := bufio.NewScanner(r)
	for s.Scan() {
		n, err := strconv.Atoi(s.Text())
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	if len(out) < bufSize {
		return nil, errors.New("not enough numbers")
	}
	return out, s.Err()
}

func findInvalidNumber(vs []int) (v int, ok bool) {
loop:
	for i := bufSize; i < len(vs); i++ {
		for j := i - bufSize; j < i; j++ {
			for k := j + 1; k < i; k++ {
				if vs[j]+vs[k] == vs[i] {
					continue loop
				}
			}
		}
		return vs[i], true
	}
	return 0, false
}

func findSum(vs []int, v int) []int {
	for i := 0; i < len(vs); i++ {
		sum := 0
		for j := i; j < len(vs); j++ {
			sum += vs[j]
			if sum == v {
				return vs[i : j+1]
			}
		}
	}
	return nil
}

func bounds(vs []int) (min, max int) {
	min, max = math.MaxInt64, math.MinInt64
	for _, v := range vs {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}
