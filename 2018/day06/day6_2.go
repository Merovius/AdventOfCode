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

	var num int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var total int
			q := image.Pt(x, y)
			for _, p := range pts {
				total += dist(p, q)
			}
			if total < 10000 {
				num++
			}
		}
	}
	fmt.Println(num)
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
