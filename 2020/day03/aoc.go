package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	var grid []string
	for s.Scan() {
		grid = append(grid, s.Text())
	}
	slopes := [][2]int{
		{1, 1}, {3, 1}, {5, 1}, {7, 1}, {1, 2},
	}
	product := 1
	for _, s := range slopes {
		v := CountTrees(grid, s[0], s[1])
		fmt.Printf("(%d,%d) -> %d\n", s[0], s[1], v)
		product *= v
	}
	fmt.Println(product)
}

func CountTrees(grid []string, slopeX, slopeY int) int {
	var (
		N int
		X int
	)
	for Y := 0; Y < len(grid); Y += slopeY {
		l := grid[Y]
		if l[X%len(l)] == '#' {
			N++
		}
		X += slopeX
	}
	return N
}
