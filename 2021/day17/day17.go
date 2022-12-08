package main

import (
	"fmt"
	"math"
)

const (
	MinX = 257
	MaxX = 286
	MinY = -101
	MaxY = -57
)

func main() {
	var (
		numHits = 0
		highest = math.MinInt
	)
	for vx := 1; vx < 1000; vx++ {
		for vy := -1000; vy < 1000; vy++ {
			h, hit := Shoot(vx, vy)
			if hit {
				numHits++
			}
			if h > highest {
				highest = h
			}
		}
	}
	fmt.Printf("Highest Y is %d\n", highest)
	fmt.Printf("Number of hits is %d\n", numHits)
}

func Shoot(vx, vy int) (highest int, hit bool) {
	x, y := 0, 0
	highest = math.MinInt
	for x <= MaxX && y >= MinY {
		x += vx
		y += vy
		if vx > 0 {
			vx -= 1
		}
		vy -= 1

		if y > highest {
			highest = y
		}
		if x >= MinX && x <= MaxX && y >= MinY && y <= MaxY {
			return highest, true
		}
	}
	return math.MinInt, false
}
