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
