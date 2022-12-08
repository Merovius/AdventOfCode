package main

import "testing"

func TestDay1(t *testing.T) {
	tcs := []struct {
		in         []int
		wantPair   [2]int
		wantTriple [3]int
	}{
		{[]int{1721, 979, 366, 299, 675, 1456}, [2]int{1721, 299}, [3]int{979, 366, 675}},
	}
	for _, tc := range tcs {
		if got, ok := Pair2020(tc.in); !ok || got != tc.wantPair {
			t.Fatalf("Pair2020(%v) = %v, %v, want %v, true", tc.in, got, ok, tc.wantPair)
		}
		if got, ok := Triple2020(tc.in); !ok || got != tc.wantTriple {
			t.Fatalf("ProductOfTriple(%v) = %v, %v, want %v, true", tc.in, got, ok, tc.wantTriple)
		}
	}
}
