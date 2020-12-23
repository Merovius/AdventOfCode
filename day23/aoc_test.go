package main

import (
	"reflect"
	"testing"
)

func TestDoMove(t *testing.T) {
	states := []Game{
		{0, ring{3, 8, 9, 1, 2, 5, 4, 6, 7}},
		{1, ring{3, 2, 8, 9, 1, 5, 4, 6, 7}},
		{2, ring{3, 2, 5, 4, 6, 7, 8, 9, 1}},
		{3, ring{7, 2, 5, 8, 9, 1, 3, 4, 6}},
		{4, ring{3, 2, 5, 8, 4, 6, 7, 9, 1}},
		{5, ring{9, 2, 5, 8, 4, 1, 3, 6, 7}},
		{6, ring{7, 2, 5, 8, 4, 1, 9, 3, 6}},
		{7, ring{8, 3, 6, 7, 4, 1, 9, 2, 5}},
		{8, ring{7, 4, 1, 5, 8, 3, 9, 2, 6}},
		{0, ring{5, 7, 4, 1, 8, 3, 9, 2, 6}},
		{1, ring{5, 8, 3, 7, 4, 1, 9, 2, 6}},
	}
	for i := 0; i < len(states)-1; i++ {
		got := doMove(states[i])
		if !reflect.DeepEqual(got, states[i+1]) {
			t.Errorf("doMove(%v) = %v, want %v (step %d)", states[i], got, states[i+1], i)
		}
	}
}
