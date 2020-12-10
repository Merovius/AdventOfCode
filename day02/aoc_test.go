package main

import (
	"strings"
	"testing"
)

func TestValidPasswords(t *testing.T) {
	tcs := []struct {
		in    string
		wantC int
		wantP int
	}{
		{"1-3 a: abcde\n1-3 b: cdefg\n2-9 c: ccccccccc", 2, 1},
	}
	for _, tc := range tcs {
		gotC, gotP, err := ValidPasswords(strings.NewReader(tc.in))
		if err != nil || gotC != tc.wantC || gotP != tc.wantP {
			t.Errorf("ValidPasswords(%q) = %v, %v, %v, want %d, %d, <nil>", tc.in, gotC, gotP, err, tc.wantC, tc.wantP)
		}
	}
}
