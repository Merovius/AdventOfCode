package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	stars := readStars()

	prevBounds := bounds(stars)
	step(stars, 1)
	i := 1
	for ; ; i++ {
		step(stars, 1)
		b := bounds(stars)
		if b.Dx()*b.Dy() > prevBounds.Dx()*prevBounds.Dy() {
			step(stars, -1)
			break
		}
		prevBounds = b
	}
	b := prevBounds
	fmt.Println(b, i)

	w, h := b.Dx(), b.Dy()
	marks := make([]bool, w*h)
	for _, st := range stars {
		marks[(st.P.Y-b.Min.Y)*w+st.P.X-b.Min.X] = true
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if marks[y*w+x] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func bounds(stars []star) image.Rectangle {
	var bounds image.Rectangle
	for _, s := range stars {
		bounds = bounds.Union(image.Rectangle{s.P, s.P.Add(image.Pt(1, 1))})
	}
	return bounds
}

func step(stars []star, i int) {
	for j, st := range stars {
		stars[j].P = st.P.Add(st.V.Mul(i))
	}
}

type star struct {
	P image.Point
	V image.Point
}

func readStars() []star {
	re := regexp.MustCompile(`position=<\s*(-?\d+),\s*(-?\d+)> velocity=<\s*(-?\d+),\s*(-?\d+)>`)
	var out []star
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		var st star
		m := re.FindStringSubmatch(s.Text())
		if len(m) != 5 {
			log.Fatal("Could not match")
		}
		st.P.X, _ = strconv.Atoi(m[1])
		st.P.Y, _ = strconv.Atoi(m[2])
		st.V.X, _ = strconv.Atoi(m[3])
		st.V.Y, _ = strconv.Atoi(m[4])
		out = append(out, st)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	return out
}
