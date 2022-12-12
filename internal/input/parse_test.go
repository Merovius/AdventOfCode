package input

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFields(t *testing.T) {
	type X struct {
		A int
		B uint
		c string
		D testParser
		E *int
		F *testParser
	}
	tcs := []struct {
		in   string
		want X
	}{
		{"42 23 foo 13 foo", X{42, 23, "", 1337, ptr(13), ptr[testParser](1337)}},
	}
	p := Fields[X]()
	for _, tc := range tcs {
		if x, err := p(tc.in); err != nil || !cmp.Equal(x, tc.want, cmpopts.IgnoreUnexported(X{})) {
			t.Errorf("Fields[X](%q) = %v, %v, want %v, <nil>", tc.in, x, err, tc.want)
		}
	}
}

type testParser int

func (p *testParser) Parse(s string) error {
	if s == "foo" {
		*p = 1337
		return nil
	}
	return fmt.Errorf("only %q is valid", "foo")
}

func ptr[T any](v T) *T { return &v }
