package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestExplode(t *testing.T) {
	tcs := []struct {
		input string
		want  string
	}{
		{"[[[[[9,8],1],2],3],4]", "[[[[0,9],2],3],4]"},
		{"[7,[6,[5,[4,[3,2]]]]]", "[7,[6,[5,[7,0]]]]"},
		{"[[6,[5,[4,[3,2]]]],1]", "[[6,[5,[7,0]]],3]"},
		{"[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]", "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]"},
		{"[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]", "[[3,[2,[8,0]]],[9,[5,[7,0]]]]"},
	}
	for _, tc := range tcs {
		input, want := new(Num), new(Num)
		if err := json.Unmarshal([]byte(tc.input), input); err != nil {
			t.Fatalf("json.Unmarshal(%q, _) = %v, want <nil>", tc.input, err)
		}
		if err := json.Unmarshal([]byte(tc.want), want); err != nil {
			t.Fatalf("json.Unmarshal(%q, _) = %v, want <nil>", tc.want, err)
		}
		got := input.copy()
		got.explode()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%v.explode() = %v, want %v", input, got, want)
		}
	}
}

func TestReduce(t *testing.T) {
	tcs := []struct {
		input string
		want  string
	}{
		{"[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]", "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]"},
	}
	for _, tc := range tcs {
		input, want := new(Num), new(Num)
		if err := json.Unmarshal([]byte(tc.input), input); err != nil {
			t.Fatalf("json.Unmarshal(%q, _) = %v, want <nil>", tc.input, err)
		}
		if err := json.Unmarshal([]byte(tc.want), want); err != nil {
			t.Fatalf("json.Unmarshal(%q, _) = %v, want <nil>", tc.want, err)
		}
		got := input.copy()
		got.reduce()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%v.reduce() = %v, want %v", input, got, want)
		}
	}
}

func TestAdd(t *testing.T) {
	tcs := []struct {
		inputA string
		inputB string

		want string
	}{
		{"[1,1]", "[2,2]", "[[1,1],[2,2]]"},
		{"[[1,1],[2,2]]", "[3,3]", "[[[1,1],[2,2]],[3,3]]"},
		{"[[[1,1],[2,2]],[3,3]]", "[4,4]", "[[[[1,1],[2,2]],[3,3]],[4,4]]"},
		{"[[[[1,1],[2,2]],[3,3]],[4,4]]", "[5,5]", "[[[[3,0],[5,3]],[4,4]],[5,5]]"},
		{"[[[[3,0],[5,3]],[4,4]],[5,5]]", "[6,6]", "[[[[5,0],[7,4]],[5,5]],[6,6]]"},
		{"[[[[4,3],4],4],[7,[[8,4],9]]]", "[1,1]", "[[[[0,7],4],[[7,8],[6,0]]],[8,1]]"},
		{"[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]", "[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]", "[[[[4,0],[5,4]],[[7,7],[6,0]]],[[8,[7,7]],[[7,9],[5,0]]]]"},
	}
	for _, tc := range tcs {
		inputA, inputB, want := new(Num), new(Num), new(Num)
		if err := json.Unmarshal([]byte(tc.inputA), inputA); err != nil {
			t.Fatalf("json.Unmarshal(%q, _) = %v, want <nil>", tc.inputA, err)
		}
		if err := json.Unmarshal([]byte(tc.inputB), inputB); err != nil {
			t.Fatalf("json.Unmarshal(%q, _) = %v, want <nil>", tc.inputB, err)
		}
		if err := json.Unmarshal([]byte(tc.want), want); err != nil {
			t.Fatalf("json.Unmarshal(%q, _) = %v, want <nil>", tc.want, err)
		}
		got := Add(inputA, inputB)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("add(%v, %v) = %v, want %v", inputA, inputB, got, want)
		}
	}
}
