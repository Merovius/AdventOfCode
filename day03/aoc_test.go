package main

import "testing"

func TestCountTrees(t *testing.T) {
	grid := []string{
		"..##.......",
		"#...#...#..",
		".#....#..#.",
		"..#.#...#.#",
		".#...##..#.",
		"..#.##.....",
		".#.#.#....#",
		".#........#",
		"#.##...#...",
		"#...##....#",
		".#..#...#.#",
	}
	tcs := []struct {
		slopeX int
		slopeY int
		want   int
	}{
		{1, 1, 2},
		{3, 1, 7},
		{5, 1, 3},
		{7, 1, 4},
		{1, 2, 2},
	}
	for _, tc := range tcs {
		if got := CountTrees(grid, tc.slopeX, tc.slopeY); got != tc.want {
			t.Errorf("CountTrees(…, %d, %d) = %v, want %v", tc.slopeX, tc.slopeY, got, tc.want)
		}
	}
}
