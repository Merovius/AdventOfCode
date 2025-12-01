package container

import (
	"iter"
	"slices"
	"strings"
	"testing"

	"gonih.org/AdventOfCode/internal/xiter"
)

func TestRadixSet(t *testing.T) {
	var r RadixSet
	t.Run("Add", func(t *testing.T) {
		tcs := []struct {
			s       string
			want    bool
			wantAll []string
		}{
			{"bar", true, []string{"bar"}},
			{"bar", false, []string{"bar"}},
			{"foo", true, []string{"bar", "foo"}},
			{"baz", true, []string{"bar", "baz", "foo"}},
			{"", true, []string{"", "bar", "baz", "foo"}},
			{"foobar", true, []string{"", "bar", "baz", "foo", "foobar"}},
		}
		for _, tc := range tcs {
			if got := r.Add(tc.s); got != tc.want {
				t.Errorf("Add(%q) = %v, want %v", tc.s, got, tc.want)
			}
			gotAll := collect(r.All())
			if !slices.Equal(gotAll, tc.wantAll) {
				t.Errorf("Add(%q) = %q, want %q", tc.s, gotAll, tc.wantAll)
			}
			if got := r.Len(); got != len(gotAll) {
				t.Errorf("Len() = %d, want %d", got, len(gotAll))
			}
		}
	})
	if got, want := r.Len(), 5; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}

	t.Run("Contains", func(t *testing.T) {
		tcs := []struct {
			key  string
			want bool
		}{
			{"", true},
			{"foo", true},
			{"bar", true},
			{"baz", true},
			{"foobar", true},
			{"f", false},
			{"foob", false},
			{"foobarbaz", false},
			{"spam", false},
		}
		for _, tc := range tcs {
			if got := r.Contains(tc.key); got != tc.want {
				t.Errorf("Get(%q) = %v, want %v", tc.key, got, tc.want)
			}
		}
	})
	t.Run("WithPrefix", func(t *testing.T) {
		tcs := []struct {
			prefix string
			want   []string
		}{
			{"", []string{"", "bar", "baz", "foo", "foobar"}},
			{"ba", []string{"bar", "baz"}},
			{"bar", []string{"bar"}},
			{"baz", []string{"baz"}},
			{"barr", nil},
			{"fo", []string{"foo", "foobar"}},
			{"foo", []string{"foo", "foobar"}},
		}
		for _, tc := range tcs {
			if got := collect(r.WithPrefix(tc.prefix)); !slices.Equal(got, tc.want) {
				t.Errorf("WithPrefix(%q) = %q, want %q", tc.prefix, got, tc.want)
			}
		}
	})
	t.Run("PrefixesOf", func(t *testing.T) {
		tcs := []struct {
			s    string
			want []string
		}{
			{"", []string{""}},
			{"foo", []string{"", "foo"}},
			{"foob", []string{"", "foo"}},
			{"foobar", []string{"", "foo", "foobar"}},
			{"foobarbaz", []string{"", "foo", "foobar"}},
			{"bazaar", []string{"", "baz"}},
		}
		for _, tc := range tcs {
			if got := collect(r.PrefixesOf(tc.s)); !slices.Equal(got, tc.want) {
				t.Errorf("PrefixesOf(%q) = %q, want %q", tc.s, got, tc.want)
			}
		}
	})
	t.Run("Delete", func(t *testing.T) {
		tcs := []struct {
			s       string
			want    bool
			wantAll []string
		}{
			{"", true, []string{"bar", "baz", "foo", "foobar"}},
			{"", false, []string{"bar", "baz", "foo", "foobar"}},
			{"foo", true, []string{"bar", "baz", "foobar"}},
			{"banana", false, []string{"bar", "baz", "foobar"}},
			{"bar", true, []string{"baz", "foobar"}},
			{"bar", false, []string{"baz", "foobar"}},
			{"baz", true, []string{"foobar"}},
			{"foobar", true, nil},
			{"", false, nil},
		}
		for _, tc := range tcs {
			if got := r.Delete(tc.s); got != tc.want {
				t.Errorf("Delete(%q) = %v, want %v", tc.s, got, tc.want)
			}
			gotAll := collect(r.All())
			if !slices.Equal(gotAll, tc.wantAll) {
				t.Errorf("Delete(%q) = %q, want %q", tc.s, gotAll, tc.wantAll)
			}
			if got := r.Len(); got != len(gotAll) {
				t.Errorf("Len() = %d, want %d", got, len(gotAll))
			}
		}
	})

}

func collect(s iter.Seq[string]) []string {
	return slices.Collect(xiter.Map(s, strings.Clone))
}
