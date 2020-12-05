package main

import "testing"

func TestSeatID(t *testing.T) {
	tcs := []struct {
		in   string
		want SeatID
	}{
		{"FBFBBFFRLR", 357},
		{"BFFFBBFRRR", 567},
		{"FFFBBBFRRR", 119},
		{"BBFFBBFRLL", 820},
	}
	for _, tc := range tcs {
		if got := ParseSeatID(tc.in); got != tc.want {
			t.Errorf("SeatID(%q) = %v (%.16b), want %v (%.16b)", tc.in, got, got, tc.want, tc.want)
		}
	}
}
