package main

import (
	"fmt"
	"image"
	"io"
	"log"
)

func main() {
	pts := readPoints()
	var bounds image.Rectangle
	for _, pt := range pts {
		bounds = bounds.Union(image.Rectangle{pt, pt.Add(image.Pt(1, 1))})
	}
	for i := range pts {
		pts[i] = pts[i].Sub(bounds.Min)
	}
	bounds = bounds.Sub(bounds.Min)
	w, h := bounds.Dx(), bounds.Dy()

	counts := make(map[int]int)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var (
				tie  = true
				dmin = int(^uint(0) >> 1)
				imin int
			)
			q := image.Pt(x, y)
			for i, p := range pts {
				δ := dist(p, q)
				if δ == dmin {
					tie = true
				}
				if δ < dmin {
					imin, dmin, tie = i, δ, false
				}
			}
			if tie {
				continue
			}
			if y == 0 || y == h-1 || x == 0 || x == w-1 {
				counts[imin] = -int(^uint(0) >> 1)
			} else {
				counts[imin]++
			}
		}
	}

	var max int
	for _, c := range counts {
		if c > max {
			max = c
		}
	}
	fmt.Println(max)
}

func readPoints() []image.Point {
	var out []image.Point
	for {
		var pt image.Point
		_, err := fmt.Scanf("%d, %d", &pt.X, &pt.Y)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		out = append(out, pt)
	}
	return out
}

func dist(a, b image.Point) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}
