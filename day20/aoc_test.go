package main

import (
	"reflect"
	"testing"
)

func TestRotateTile(t *testing.T) {
	tcs := []struct {
		input []string
		edge  int
		want  []string
	}{
		{[]string{"123", "456", "789"}, 0, []string{"123", "456", "789"}},
		{[]string{"123", "456", "789"}, 1, []string{"369", "258", "147"}},
		{[]string{"123", "456", "789"}, 2, []string{"987", "654", "321"}},
		{[]string{"123", "456", "789"}, 3, []string{"741", "852", "963"}},
	}
	for _, tc := range tcs {
		in := Tile{
			Index: 42,
			Tile:  tc.input,
		}
		got := rotateTile(in, tc.edge)
		want := Tile{
			Index: 42,
			Tile:  tc.want,
		}
		if got.Index != want.Index {
			t.Errorf("rotateTile(%q, %d) mutated index", tc.input, tc.edge)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("rotateTile(%q, %d) = %q, want %q", tc.input, tc.edge, got.Tile, tc.want)
		}
	}
}

func TestFlipTile(t *testing.T) {
	tcs := []struct {
		input []string
		want  []string
	}{
		{[]string{"123", "456", "789"}, []string{"789", "456", "123"}},
	}
	for _, tc := range tcs {
		in := Tile{
			Index: 42,
			Tile:  tc.input,
		}
		got := flipTile(in)
		want := Tile{
			Index: 42,
			Tile:  tc.want,
		}
		if got.Index != want.Index {
			t.Errorf("flipTile(%q) mutated index", tc.input)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("flipTile(%q) = %q, want %q", tc.input, got.Tile, tc.want)
		}
	}
}
