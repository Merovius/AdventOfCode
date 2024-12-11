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
		m  int
		b  int
		n  int
		x  int
		M  int
		ok bool
	}{
		{2, 3, 3, 5, 8, 15, true},
		{8, 15, 2, 7, 23, 105, true},
	}
	for _, tc := range tcs {
		if x, M, ok := ChineseRemainder(tc.a, tc.m, tc.b, tc.n); x != tc.x || M != tc.M || ok != tc.ok {
			t.Errorf("ChineseRemainder(%d, %d, %d, %d) = %d, %d, %v, want %d, %d, %v", tc.a, tc.m, tc.b, tc.n, x, M, ok, tc.x, tc.M, tc.ok)
		}
	}

	t.Run("exhaustive", func(t *testing.T) {
		for m := 1; m < 10; m++ {
			for n := 1; n < 10; n++ {
				for a := 0; a < m; a++ {
					for b := 0; b < n; b++ {
						x1, M1, ok1 := ChineseRemainder(a, m, b, n)
						x2, M2, ok2 := crtSlow(a, m, b, n)
						if x1 != x2 || ok1 != ok2 {
							t.Errorf("ChineseRemainder(%d, %d, %d, %d) = %d, %d, %v, want %d, %d %v", a, b, m, n, x1, M1, ok1, x2, M2, ok2)
						}
					}
				}
			}
		}
	})
}

func crtSlow(a, m, b, n int) (int, int, bool) {
	M := LCM(m, n)
	for x := 0; x < M; x++ {
		if x%m == a && x%n == b {
			return x, M, true
		}
	}
	return a, m, false
}

func TestDivMod(t *testing.T) {
	for a := -100; a < 100; a++ {
		for b := -100; b < 100; b++ {
			if b == 0 {
				continue
			}
			q, r := DivMod(a, b)
			if b*q+r != a || r < 0 || r >= Abs(b) {
				t.Errorf("DivMod(%d, %d) = %d, %d and b*q+r = %d, want %d", a, b, q, r, b*q+r, a)
			}
		}
	}
}

func TestChineseRemaniderBig(t *testing.T) {
	for a := int64(-10); a <= 10; a++ {
		for m := int64(1); m <= 50; m++ {
			for b := int64(-10); b <= 10; b++ {
				for n := int64(1); n <= 50; n++ {
					x1, m1, ok1 := ChineseRemainder(a, m, b, n)
					a2, m2, b2, n2 := big.NewInt(a), big.NewInt(m), big.NewInt(b), big.NewInt(n)
					ok2 := ChineseRemainderBig(a2, m2, b2, n2)
					if ok1 != ok2 || a2.Int64() != x1 || m2.Int64() != m1 {
						t.Fatalf("ChineseRemainderBig(%v, %v, %v, %v) = %v, %v, %v want %v, %v, %v", a, m, b, n, a2, m2, ok2, x1, m1, ok1)
					}
				}
			}
		}
	}
}

func TestLog10(t *testing.T) {
	for i := 1; i < 20000; i++ {
		p := Log10(i)
		if i >= 10000 {
			if p != 4 {
				t.Errorf("Log10(%d) = %d, want %d", i, p, 4)
			}
		} else if i >= 1000 {
			if p != 3 {
				t.Errorf("Log10(%d) = %d, want %d", i, p, 3)
			}
		} else if i >= 100 {
			if p != 2 {
				t.Errorf("Log10(%d) = %d, want %d", i, p, 2)
			}
		} else if i >= 10 {
			if p != 1 {
				t.Errorf("Log10(%d) = %d, want %d", i, p, 1)
			}
		} else if p != 0 {
			t.Errorf("Log10(%d) = %d, want %d", i, p, 0)
		}
	}
}
