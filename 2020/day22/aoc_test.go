package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestReadInput(t *testing.T) {
	tcs := []struct {
		input   string
		wantP1  []int
		wantP2  []int
		wantErr bool
	}{
		{
			"Player 1:\n9\n2\n6\n3\n1\n\nPlayer 2:\n5\n8\n4\n7\n10",
			[]int{9, 2, 6, 3, 1},
			[]int{5, 8, 4, 7, 10},
			false,
		},
	}
	for _, tc := range tcs {
		got1, got2, err := ReadInput(strings.NewReader(tc.input))
		if err != nil {
			if !tc.wantErr {
				t.Errorf("ReadInput(%q) = _, _, %v, want <nil>", tc.input, err)
			}
			continue
		}
		if !reflect.DeepEqual(got1, tc.wantP1) || !reflect.DeepEqual(got2, tc.wantP2) {
			t.Errorf("ReadInput(%q) = %v, %v, <nil>, want %v, %v\n", tc.input, got1, got2, tc.wantP1, tc.wantP2)
		}
	}
}

func TestPlayRound(t *testing.T) {
	tcs := []struct {
		p1    []int
		p2    []int
		want1 []int
		want2 []int
	}{
		{[]int{42}, []int{23}, []int{42, 23}, nil},
		{[]int{23}, []int{42}, nil, []int{42, 23}},
		{[]int{1, 2}, []int{3, 4}, []int{2}, []int{4, 3, 1}},
	}
	for _, tc := range tcs {
		got1, got2 := PlayRound(append([]int(nil), tc.p1...), append([]int(nil), tc.p2...))
		if !cmp.Equal(got1, tc.want1, cmpopts.EquateEmpty()) || !cmp.Equal(got2, tc.want2, cmpopts.EquateEmpty()) {
			t.Errorf("PlayRound(%v, %v) = %v, %v, want %v, %v", tc.p1, tc.p2, got1, got2, tc.want1, tc.want2)
		}
	}
}

func TestPlayGame(t *testing.T) {
	tcs := []struct {
		p1    []int
		p2    []int
		want1 []int
		want2 []int
	}{
		{[]int{42}, []int{23}, []int{42, 23}, nil},
		{[]int{23}, []int{42}, nil, []int{42, 23}},
		{[]int{1, 2}, []int{3, 4}, nil, []int{3, 1, 4, 2}},
		{[]int{9, 2, 6, 3, 1}, []int{5, 8, 4, 7, 10}, nil, []int{3, 2, 10, 6, 8, 5, 9, 4, 7, 1}},
	}
	for _, tc := range tcs {
		got1, got2 := PlayGame(append([]int(nil), tc.p1...), append([]int(nil), tc.p2...))
		if !cmp.Equal(got1, tc.want1, cmpopts.EquateEmpty()) || !cmp.Equal(got2, tc.want2, cmpopts.EquateEmpty()) {
			t.Errorf("PlayRound(%v, %v) = %v, %v, want %v, %v", tc.p1, tc.p2, got1, got2, tc.want1, tc.want2)
		}
	}
}

func TestScoreDeck(t *testing.T) {
	tcs := []struct {
		deck []int
		want int
	}{
		{[]int{3, 2, 10, 6, 8, 5, 9, 4, 7, 1}, 306},
	}
	for _, tc := range tcs {
		got := ScoreDeck(tc.deck)
		if got != tc.want {
			t.Errorf("ScoreDeck(%v) = %v, want %v", tc.deck, got, tc.want)
		}
	}
}
