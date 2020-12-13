package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	n, err := ParseNotes(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	wait := func(id int) int {
		return id - n.Depart%id
	}

	var next int = math.MaxInt64
	for _, id := range n.IDs {
		if id == 0 {
			continue
		}
		if wait(id) < wait(next) {
			next = id
		}
	}
	fmt.Printf("Next bus: %d (in %d minutes)\n", next, wait(next))
	fmt.Println("Product:", next*wait(next))

	var (
		t = 0
		m = 1
	)
	for i, id := range n.IDs {
		if id == 0 {
			continue
		}
		for (t+i)%id != 0 {
			t += m
		}
		m = lcm(m, id)
	}
	fmt.Printf("Content-winning timestamp: %d\n", t)
}

type Notes struct {
	Depart int
	IDs    []int
}

func ParseNotes(r io.Reader) (Notes, error) {
	var n Notes
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return n, err
	}
	s := string(bytes.TrimSpace(buf))
	sp := strings.Split(s, "\n")
	if len(sp) != 2 {
		return n, errors.New("need exactly two lines of input")
	}
	n.Depart, err = strconv.Atoi(sp[0])
	if err != nil {
		return n, err
	}
	sp = strings.Split(sp[1], ",")
	for _, s := range sp {
		var id int
		if s != "x" {
			id, err = strconv.Atoi(s)
			if err != nil {
				return n, err
			}
		}
		n.IDs = append(n.IDs, id)
	}
	return n, nil
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}

func lcm(a, b int) int {
	return abs(a * (b / gcd(a, b)))
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
