// Package math provides generic arithmetic functions
//
// It intentionally does not support floating point types, as that would
// require reflect in generic code, for special handling of NaN, ±∞ and -0. Use
// the standard library math package instead.
package math

import (
	"fmt"
	"math/big"
	"slices"

	"golang.org/x/exp/constraints"
)

// Mathematical constants.
const (
	E   = 2.71828182845904523536028747135266249775724709369995957496696763 // https://oeis.org/A001113
	Pi  = 3.14159265358979323846264338327950288419716939937510582097494459 // https://oeis.org/A000796
	Phi = 1.61803398874989484820458683436563811772030917980576286213544862 // https://oeis.org/A001622

	Sqrt2   = 1.41421356237309504880168872420969807856967187537694807317667974 // https://oeis.org/A002193
	SqrtE   = 1.64872127070012814684865078781416357165377610071014801157507931 // https://oeis.org/A019774
	SqrtPi  = 1.77245385090551602729816748334114518279754945612238712821380779 // https://oeis.org/A002161
	SqrtPhi = 1.27201964951406896425242246173749149171560804184009624861664038 // https://oeis.org/A139339

	Ln2    = 0.693147180559945309417232121458176568075500134360255254120680009 // https://oeis.org/A002162
	Log2E  = 1 / Ln2
	Ln10   = 2.30258509299404568401799145468436420760110148862877297603332790 // https://oeis.org/A002392
	Log10E = 1 / Ln10
)

// Floating-point limit values.
// Max is the largest finite value representable by the type.
// SmallestNonzero is the smallest positive, non-zero value representable by the type.
const (
	MaxFloat32             = 0x1p127 * (1 + (1 - 0x1p-23)) // 3.40282346638528859811704183484516925440e+38
	SmallestNonzeroFloat32 = 0x1p-126 * 0x1p-23            // 1.401298464324817070923729583289916131280e-45

	MaxFloat64             = 0x1p1023 * (1 + (1 - 0x1p-52)) // 1.79769313486231570814527423731704356798070e+308
	SmallestNonzeroFloat64 = 0x1p-1022 * 0x1p-52            // 4.9406564584124654417656879286822137236505980e-324
)

// Integer limit values.
const (
	intSize = 32 << (^uint(0) >> 63) // 32 or 64

	MaxInt    = 1<<(intSize-1) - 1  // MaxInt32 or MaxInt64 depending on intSize.
	MinInt    = -1 << (intSize - 1) // MinInt32 or MinInt64 depending on intSize.
	MaxInt8   = 1<<7 - 1            // 127
	MinInt8   = -1 << 7             // -128
	MaxInt16  = 1<<15 - 1           // 32767
	MinInt16  = -1 << 15            // -32768
	MaxInt32  = 1<<31 - 1           // 2147483647
	MinInt32  = -1 << 31            // -2147483648
	MaxInt64  = 1<<63 - 1           // 9223372036854775807
	MinInt64  = -1 << 63            // -9223372036854775808
	MaxUint   = 1<<intSize - 1      // MaxUint32 or MaxUint64 depending on intSize.
	MaxUint8  = 1<<8 - 1            // 255
	MaxUint16 = 1<<16 - 1           // 65535
	MaxUint32 = 1<<32 - 1           // 4294967295
	MaxUint64 = 1<<64 - 1           // 18446744073709551615
)

type TotallyOrdered interface {
	constraints.Integer | ~string
}

// Max returns the maximum of a and b.
func Max[T TotallyOrdered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Max returns the minimum of a and b.
func Min[T TotallyOrdered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

// Cmp returns -1, 0 and 1, if a < b, a == b and a > b, respectively.
func Cmp[T TotallyOrdered](a, b T) int {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// Sgn returns -1, 0, 1 if v < 0, v == 0, v > 0, respectively.
func Sgn[T constraints.Signed](v T) T {
	return T(Cmp(v, 0))
}

// Abs returns the absolute value of v.
func Abs[T constraints.Signed](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

// GCD returns z, x, y, such that gcd(a, b) == z == a*x+b*y.
//
// a and b may be positive, zero or negative.
//
// If a == b == 0, GCD sets z = x = y = 0.
//
// If a == 0 and b != 0, GCD sets z = |b|, x = 0, y = sign(b) * 1.
//
// If a != 0 and b == 0, GCD sets z = |a|, x = sign(a) * 1, y = 0.
func GCD[T constraints.Integer](a, b T) (z, x, y T) {
	if a == 0 && b == 0 {
		return 0, 0, 0
	}
	if a == 0 {
		if b < 0 {
			s := -1
			return -b, 0, T(s)
		}
		return b, 0, 1
	}
	if b == 0 {
		if a < 0 {
			s := -1
			return -a, T(s), 0
		}
		return a, 1, 0
	}

	var r0, r T = a, b
	var s0, s T = 1, 0
	var t0, t T = 0, 1
	for r != 0 {
		q := r0 / r
		r0, r = r, r0-q*r
		s0, s = s, s0-q*s
		t0, t = t, t0-q*t
	}
	if r0 < 0 {
		return -r0, -s0, -t0
	}
	return r0, s0, t0
}

// LCM returns the least common multiple of a and b.
func LCM[T constraints.Integer](a, b T) T {
	g, _, _ := GCD(a, b)
	return (a / g) * b
}

// ChineseRemainder solves the Chinese Remainder Theorem. The return values
// satisfy
//
//	x = a (mod m)
//	x = b (mod n)
//	0 <= x < LCM(m, n) = M
//
// If such an x exists. Otherwise it returns a, m, false.
//
// It panics if m==0 or n==0.
func ChineseRemainder[T constraints.Integer](a, m, b, n T) (x, M T, ok bool) {
	if m == 0 || n == 0 {
		panic("invalid inputs to ChineseRemainder")
	}
	aa, ma := a, m
	if m < 0 {
		m = -m
	}
	if n < 0 {
		n = -n
	}
	if a < 0 || a >= m {
		a = Mod(a, m)
	}
	if b < 0 || b >= n {
		b = Mod(b, n)
	}

	g, u, v := GCD(m, n)
	if (a-b)%g != 0 {
		return aa, ma, false
	}
	M = (m / g) * n
	x = (a * v * (n / g))
	x += (b * u * (m / g))
	x = Mod(x, M)
	return x, M, true
}

// ChineseRemainderBig solves the Chinese Remainder Theorem for big.Int. It
// calculates x such that
//
//	x = a (mod m)
//	x = b (mod n)
//	0 <= x < LCM(m, n) = M
//
// If such an x exists, and stores x and M in a and m. Otherwise, it returns
// false and leaves the inputs unmodified.
//
// It panics if m or n is zero.
func ChineseRemainderBig(a, m, b, n *big.Int) bool {
	if m.Sign() == 0 || n.Sign() == 0 {
		panic("invalid inputs to ChineseRemainderBig")
	}

	u, v := new(big.Int), new(big.Int)
	g := new(big.Int).GCD(u, v, m, n)
	if δ := new(big.Int).Sub(a, b); δ.Mod(δ, g).Sign() != 0 {
		return false
	}
	M := new(big.Int).Div(m, g)
	M = M.Mul(M, n)

	// TODO: Can these be optimized? Maybe the definition of u and v already
	// makes a*v*n divisible by g? Or something?
	x := new(big.Int).Mul(a, v)
	x = x.Mul(x, n)
	y := new(big.Int).Mul(b, u)
	y = y.Mul(y, m)

	x = x.Add(x, y)
	x = x.Div(x, g)
	x = x.Mod(x, M)

	a.Set(x)
	m.Set(M)

	return true
}

// DivMod returns p, q such that a = b•q+r, with 0≤r<|b|. It panics if b is 0.
func DivMod[T constraints.Integer](a, b T) (q, r T) {
	q, r = a/b, a%b
	if r < 0 {
		if b < 0 {
			q, r = q+1, r-b
		} else {
			q, r = q-1, r+b
		}
	}
	return q, r
}

// Div returns the quotient a/b for b != 0. Div implements Euclidean division;
// see DivMod for details.
func Div[T constraints.Integer](a, b T) T {
	q, _ := DivMod(a, b)
	return q
}

// Mod returns the modulus a%b for b != 0. Mod implements Euclidean division;
// see DivMod for details.
func Mod[T constraints.Integer](a, b T) T {
	_, r := DivMod(a, b)
	return r
}

// Log10 returns the base-10 logarithm of n, rounded down. If n is not
// positive, Log10 panics.
func Log10[T constraints.Integer](n T) int {
	if n <= 0 {
		panic(fmt.Errorf("invalid argument: Log10(%d)", n))
	}
	i, ok := slices.BinarySearch(pow10[:], uint64(n))
	if !ok {
		i--
	}
	return i
}

// Pow10 returns 10^n. It panics, if n is negative or the result would overflow.
func Pow10(n int) uint64 {
	if n < 0 {
		panic(fmt.Errorf("invalid argument: Pow10(%d)", n))
	}
	if n >= len(pow10) {
		panic(fmt.Errorf("overflow: Pow10(%d)", n))
	}
	return pow10[n]

}

var pow10 = [...]uint64{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19}

// Digits returns the number of decimal digits of the absolute value of n.
func Digits[T constraints.Integer](n T) int {
	switch {
	case n > 0:
		return Log10(n) + 1
	case n < 0:
		return Log10(-n) + 1
	default:
		return 1
	}
}

// Cramer2 solves the linear equation x₀•a+x₁•b=c. It returns [0, 0], false,
// if there is no unique integer solution.
func Cramer2[T constraints.Integer](a, b, c [2]T) (x [2]T, ok bool) {
	// Cramer's rule
	Δ := det2x2(a, b)
	if Δ == 0 {
		return x, false
	}
	x0 := det2x2(c, b)
	if x0%Δ != 0 {
		return x, false
	}
	x1 := det2x2(a, c)
	if x1%Δ != 0 {
		return x, false
	}
	return [2]T{x0 / Δ, x1 / Δ}, true
}

func det2x2[T constraints.Integer](a, b [2]T) T {
	return a[0]*b[1] - b[0]*a[1]
}

// Cramer3 solves the linear equation x₀•a+x₁•b+x₂•c=d. It returns
// [0, 0, 0], false, if there is no unique integer solutions.
func Cramer3[T constraints.Integer](a, b, c, d [3]T) (x [3]T, ok bool) {
	// Cramer's rule
	Δ := det3x3(a, b, c)
	if Δ == 0 {
		return x, false
	}
	x0 := det3x3(d, b, c)
	if x0%Δ != 0 {
		return x, false
	}
	x1 := det3x3(a, d, c)
	if x1%Δ != 0 {
		return x, false
	}
	x2 := det3x3(a, b, d)
	if x2%Δ != 0 {
		return x, false
	}
	return [3]T{x0 / Δ, x1 / Δ, x2 / Δ}, true
}

func det3x3[T constraints.Integer](a, b, c [3]T) T {
	return a[0]*b[1]*c[2] + b[0]*c[1]*a[2] + c[0]*a[1]*b[2] - a[0]*c[1]*b[2] - b[0]*a[1]*c[2] - c[0]*b[1]*a[2]
}
