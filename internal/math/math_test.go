package math

import (
	"math/big"
	"testing"

	"golang.org/x/exp/constraints"
)

func TestGCD(t *testing.T) {
	tcs := []struct {
		name string
		f    func(*testing.T)
	}{
		{"int", testGCDSigned[int]},
		{"int8", testGCDSigned[int8]},
		{"int16", testGCDSigned[int16]},
		{"int32", testGCDSigned[int32]},
		{"int64", testGCDSigned[int64]},
		{"uint", testGCDUnsigned[uint]},
		{"uint8", testGCDUnsigned[uint8]},
		{"uint16", testGCDUnsigned[uint16]},
		{"uint32", testGCDUnsigned[uint32]},
		{"uint64", testGCDUnsigned[uint64]},
		{"uintptr", testGCDUnsigned[uintptr]},
	}
	for _, tc := range tcs {
		t.Run(tc.name, tc.f)
	}
}

func testGCDSigned[T constraints.Signed](t *testing.T) {
	for a := T(-100); a <= 100; a++ {
		for b := T(-100); b <= 100; b++ {
			z, x, y := GCD[T](a, b)
			zz, xx, yy := gcdBig[T](a, b)
			if z != zz || x != xx || y != yy {
				t.Errorf("GCD[%T](%d, %d) = %d, %d, %d, want %d, %d, %d", a, a, b, z, x, y, zz, xx, yy)
			}
		}
	}
}

func testGCDUnsigned[T constraints.Unsigned](t *testing.T) {
	for a := T(0); a <= 100; a++ {
		for b := T(0); b <= 100; b++ {
			z, x, y := GCD[T](a, b)
			zz, xx, yy := gcdBig[T](a, b)
			if z != zz || x != xx || y != yy {
				t.Errorf("GCD[%T](%d, %d) = %d, %d, %d, want %d, %d, %d", a, a, b, z, x, y, zz, xx, yy)
			}
		}
	}
}

func gcdBig[T constraints.Integer](a, b T) (z, x, y T) {
	ab, bb := big.NewInt(int64(a)), big.NewInt(int64(b))
	xb, yb, zb := new(big.Int), new(big.Int), new(big.Int)
	zb.GCD(xb, yb, ab, bb)
	return T(zb.Int64()), T(xb.Int64()), T(yb.Int64())
}

func TestChineseRemainder(t *testing.T) {
	tcs := []struct {
		a  int
		b  int
		m  int
		n  int
		x  int
		ok bool
	}{
		{2, 3, 3, 5, 8, true},
		{8, 2, 15, 7, 23, true},
	}
	for _, tc := range tcs {
		if x, ok := ChineseRemainder(tc.a, tc.b, tc.m, tc.n); x != tc.x || ok != tc.ok {
			t.Errorf("ChineseRemainder(%d, %d, %d, %d) = %d, %v, want %d, %v", tc.a, tc.b, tc.m, tc.n, x, ok, tc.x, tc.ok)
		}
	}

	t.Run("exhaustive", func(t *testing.T) {
		for m := 1; m < 10; m++ {
			for n := 1; n < 10; n++ {
				for a := 0; a < m; a++ {
					for b := 0; b < n; b++ {
						x1, ok1 := ChineseRemainder(a, b, m, n)
						x2, ok2 := crtSlow(a, b, m, n)
						if x1 != x2 || ok1 != ok2 {
							t.Errorf("ChineseRemainder(%d, %d, %d, %d) = %d, %v, want %d, %v", a, b, m, n, x1, ok1, x2, ok2)
						}
					}
				}
			}
		}
	})
}

func crtSlow(a, b, m, n int) (int, bool) {
	for x := 0; x < m*n; x++ {
		if x%m == a && x%n == b {
			return x, true
		}
	}
	return 0, false
}

func TestDiv(t *testing.T) {
	for a := -100; a < 100; a++ {
		for b := -100; b < 100; b++ {
			if b == 0 {
				continue
			}
			q, r := Div(a, b)
			if b*q+r != a || r < 0 || r >= Abs(b) {
				t.Errorf("Div(%d, %d) = %d, %d and b*q+r = %d, want %d", a, b, q, r, b*q+r, a)
			}
		}
	}
}
