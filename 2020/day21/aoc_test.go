package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestReadInput(t *testing.T) {
	tcs := []struct {
		in        string
		want      []Food
		wantError bool
	}{
		{``, nil, false},
		{`foo`, []Food{{Ingredients: MakeStringSet("foo")}}, false},
		{`foo bar`, []Food{{Ingredients: MakeStringSet("foo", "bar")}}, false},
		{`foo (contains nuts)`, []Food{
			{Ingredients: MakeStringSet("foo"), Allergens: MakeStringSet("nuts")},
		}, false},
		{`foo (contains nuts, soy)`, []Food{
			{Ingredients: MakeStringSet("foo"), Allergens: MakeStringSet("nuts", "soy")},
		}, false},
		{`foo bar (contains nuts, soy)`, []Food{
			{Ingredients: MakeStringSet("foo", "bar"), Allergens: MakeStringSet("nuts", "soy")},
		}, false},
		{"foo\nbar", []Food{
			{Ingredients: MakeStringSet("foo")},
			{Ingredients: MakeStringSet("bar")},
		}, false},
		{`foo (contains nuts`, nil, true},
	}
	for _, tc := range tcs {
		got, err := ReadInput(strings.NewReader(tc.in))
		if tc.wantError {
			if err == nil {
				t.Errorf("ReadInput(%q) = _, <nil>, want error", tc.in)
			} else {
				continue
			}
		}
		if diff := cmp.Diff(got, tc.want, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("ReadInput(%q) returned wrong output:\n%s", tc.in, diff)
		}
	}
}
