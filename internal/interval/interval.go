// Package interval implements intervals of integers.
//
// Interval types are named by whether or not either of their ends is
// open/closed. For example, the type CO is closed on the left and open on the
// right, so it describes the interval type [a,b).
//
// Only signed integers are supported, to make arithmetic simpler. In general,
// the package will note behave well in the presence of overflows.
//
// The package also includes a Set type, which can efficiently store sparse
// integer intervals. It is parameterized, so can be used with any of the
// interval types in this package.
//
// The easiest way to use it is to define a local alias for the interval kind
// appropriate for the problem, e.g.
//
//	type Interval = interval.CO[int]
package interval

import (
	"fmt"

	"github.com/Merovius/AdventOfCode/internal/math"
	"golang.org/x/exp/constraints"
)

// Interval is a constraint for any of the interval types in this package
type Interval[T constraints.Signed] interface {
	CO[T] | OC[T] | OO[T] | CC[T]
	toCO() CO[T]
}

// Convert between interval kinds.
func Convert[J, I Interval[T], T constraints.Signed](i I) J {
	switch any(*new(J)).(type) {
	case CO[T]:
		return any(i.toCO()).(J)
	case OC[T]:
		return any(i.toCO().toOC()).(J)
	case OO[T]:
		return any(i.toCO().toOO()).(J)
	case CC[T]:
		return any(i.toCO().toCC()).(J)
	default:
		panic("unreachable")
	}
}

// CO is a right half-open interval of T, that is [Min, Max). Well-formed
// intervals have Min≤Max, the empty interval is {0,0}.
type CO[T constraints.Signed] struct {
	Min T
	Max T
}

func MakeCO[T constraints.Signed](a, b T) CO[T] {
	return CO[T]{min(a, b), max(a, b)}
}

func (i CO[T]) toCO() CO[T] {
	return i
}

func (i CO[T]) toOC() OC[T] {
	return OC[T]{i.Min - 1, i.Max - 1}
}

func (i CO[T]) toOO() OO[T] {
	return OO[T]{i.Min - 1, i.Max}
}

func (i CO[T]) toCC() CC[T] {
	return CC[T]{i.Min, i.Max - 1}
}

// Empty returns if i is the empty interval.
func (i CO[T]) Empty() bool {
	return i.Len() <= 0
}

// Len returns the size of the interval.
func (i CO[T]) Len() int {
	return int(i.Max - i.Min)
}

func (i CO[T]) valid() bool {
	return i.Max >= i.Min
}

// Contains returns whether i contains v.
func (i CO[T]) Contains(v T) bool {
	return v >= i.Min && v < i.Max
}

// Intersect returns the intersection of i and j.
func (i CO[T]) Intersect(j CO[T]) CO[T] {
	if i.Min >= j.Max || j.Min >= i.Max {
		return CO[T]{}
	}
	return CO[T]{math.Max(i.Min, j.Min), math.Min(i.Max, j.Max)}
}

// Intersects returns whether i and j intersect.
func (i CO[T]) Intersects(j CO[T]) bool {
	return i.Min < j.Max && j.Min < i.Max
}

// Union returns the union of i and j, if it is an interval. Otherwise, it
// returns the empty interval and false.
func (i CO[T]) Union(j CO[T]) (CO[T], bool) {
	if i.Min == j.Max {
		return CO[T]{j.Min, i.Max}, true
	}
	if j.Min == i.Max {
		return CO[T]{i.Min, j.Max}, true
	}
	if i.Intersects(j) {
		return CO[T]{math.Min(i.Min, j.Min), math.Max(i.Max, j.Max)}, true
	}
	return CO[T]{}, false
}

// String implements fmt.Stringer.
func (i CO[T]) String() string {
	return fmt.Sprintf("[%d,%d)", i.Min, i.Max)
}

// OC is a left half-open interval of T, that is (Min, Max]. Well-formed intervals have
// Min≤Max, the empty interval is {0,0}.
type OC[T constraints.Signed] struct {
	Min T
	Max T
}

func MakeOC[T constraints.Signed](a, b T) OC[T] {
	return OC[T]{min(a, b), max(a, b)}
}

func (i OC[T]) toCO() CO[T] {
	return CO[T]{i.Min + 1, i.Max + 1}
}

// Empty returns if i is the empty interval.
func (i OC[T]) Empty() bool {
	return i.toCO().Empty()
}

// Len returns the size of the interval.
func (i OC[T]) Len() int {
	return i.toCO().Len()
}

func (i OC[T]) valid() bool {
	return i.toCO().valid()
}

// Contains returns whether i contains v.
func (i OC[T]) Contains(v T) bool {
	return i.toCO().Contains(v)
}

// Intersect returns the intersection of i and j.
func (i OC[T]) Intersect(j OC[T]) OC[T] {
	return i.toCO().Intersect(j.toCO()).toOC()
}

// Intersects returns whether i and j intersect.
func (i OC[T]) Intersects(j OC[T]) bool {
	return i.toCO().Intersects(j.toCO())
}

// Union returns the union of i and j, if it is an interval. Otherwise, it
// returns the empty interval and false.
func (i OC[T]) Union(j OC[T]) (OC[T], bool) {
	u, ok := i.toCO().Union(j.toCO())
	return u.toOC(), ok
}

// String implements fmt.Stringer.
func (i OC[T]) String() string {
	return fmt.Sprintf("(%d,%d]", i.Min, i.Max)
}

// OO is an open interval of T, that is (Min, Max). Well-formed intervals have
// Min≤Max+1, the empty interval is {0,1}.
type OO[T constraints.Signed] struct {
	Min T
	Max T
}

func MakeOO[T constraints.Signed](a, b T) OO[T] {
	return OO[T]{min(a, b), max(a, b)}
}

func (i OO[T]) toCO() CO[T] {
	return CO[T]{i.Min + 1, i.Max}
}

// Empty returns if i is the empty interval.
func (i OO[T]) Empty() bool {
	return i.toCO().Empty()
}

// Len returns the size of the interval.
func (i OO[T]) Len() int {
	return i.toCO().Len()
}

func (i OO[T]) valid() bool {
	return i.toCO().valid()
}

// Contains returns whether i contains v.
func (i OO[T]) Contains(v T) bool {
	return i.toCO().Contains(v)
}

// Intersect returns the intersection of i and j.
func (i OO[T]) Intersect(j OO[T]) OO[T] {
	return i.toCO().Intersect(j.toCO()).toOO()
}

// Intersects returns whether i and j intersect.
func (i OO[T]) Intersects(j OO[T]) bool {
	return i.toCO().Intersects(j.toCO())
}

// Union returns the union of i and j, if it is an interval. Otherwise, it
// returns the empty interval and false.
func (i OO[T]) Union(j OO[T]) (OO[T], bool) {
	u, ok := i.toCO().Union(j.toCO())
	return u.toOO(), ok
}

// String implements fmt.Stringer.
func (i OO[T]) String() string {
	return fmt.Sprintf("(%d,%d)", i.Min, i.Max)
}

// CC is a closed interval of T, that is [Min, Max]. Well-formed intervals have
// Min≤Max+1, the empty interval is {0,-1}.
type CC[T constraints.Signed] struct {
	Min T
	Max T
}

func MakeCC[T constraints.Signed](a, b T) CC[T] {
	return CC[T]{min(a, b), max(a, b)}
}

func (i CC[T]) toCO() CO[T] {
	return CO[T]{i.Min, i.Max + 1}
}

// Empty returns if i is the empty interval.
func (i CC[T]) Empty() bool {
	return i.toCO().Empty()
}

// Len returns the size of the interval.
func (i CC[T]) Len() int {
	return i.toCO().Len()
}

func (i CC[T]) valid() bool {
	return i.toCO().valid()
}

// Contains returns whether i contains v.
func (i CC[T]) Contains(v T) bool {
	return i.toCO().Contains(v)
}

// Intersect returns the intersection of i and j.
func (i CC[T]) Intersect(j CC[T]) CC[T] {
	return i.toCO().Intersect(j.toCO()).toCC()
}

// Intersects returns whether i and j intersect.
func (i CC[T]) Intersects(j CC[T]) bool {
	return i.toCO().Intersects(j.toCO())
}

// Union returns the union of i and j, if it is an interval. Otherwise, it
// returns the empty interval and false.
func (i CC[T]) Union(j CC[T]) (CC[T], bool) {
	u, ok := i.toCO().Union(j.toCO())
	return u.toCC(), ok
}

// String implements fmt.Stringer.
func (i CC[T]) String() string {
	return fmt.Sprintf("[%d,%d]", i.Min, i.Max)
}
