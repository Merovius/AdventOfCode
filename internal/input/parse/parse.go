package parse

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Merovius/AdventOfCode/internal/input/split"
	"golang.org/x/exp/constraints"
)

// A Parser is a generic, composable parser.
type Parser[T any] func(string) (T, error)

// Parse the reader, returning the result.
func (p Parser[T]) Parse(r io.Reader) (T, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return *new(T), err
	}
	return p(string(bytes.TrimSpace(buf)))
}

// Array splits a string and calls p for each piece. A must be an array type
// with element type T.
func Array[A, T any](s split.Func, p Parser[T]) Parser[A] {
	tt := reflect.TypeOf(new(T)).Elem()
	at := reflect.TypeOf(new(A)).Elem()
	if at.Kind() != reflect.Array || at.Elem() != tt {
		panic(fmt.Errorf("%v is not an […]%v array type", at, tt))
	}
	n := at.Len()
	return func(in string) (A, error) {
		var a A
		sp, err := s(in)
		if err != nil {
			return a, err
		}
		if len(sp) != n {
			return a, fmt.Errorf("expected %d pieces, got %d", len(sp), n)
		}
		av := reflect.ValueOf(&a).Elem()
		for i, s := range sp {
			v, err := p(s)
			if err != nil {
				return a, err
			}
			av.Index(i).Set(reflect.ValueOf(v))
		}
		return a, nil
	}
}

// MapParser converts a Parser[A] into a Parser[B] using f.
func MapParser[A, B any](p Parser[A], f func(A) (B, error)) Parser[B] {
	return func(in string) (B, error) {
		a, err := p(in)
		if err != nil {
			return *new(B), err
		}
		return f(a)
	}
}

// Slice splits a string and calls p on each piece.
func Slice[T any](s split.Func, p Parser[T]) Parser[[]T] {
	return func(in string) ([]T, error) {
		var out []T
		sp, err := s(in)
		if err != nil {
			return nil, err
		}
		for _, s := range sp {
			v, err := p(s)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	}
}

// Fields is a shorthand for Slice(split.Fields, …).
func Fields[T any](p Parser[T]) Parser[[]T] {
	return Slice(split.Fields, p)
}

// Lines is a shorthand for Slice(split.Lines, …).
func Lines[T any](p Parser[T]) Parser[[]T] {
	return Slice(split.Lines, p)
}

// Blocks is a shorthand for Slice(split.Blocks, …).
func Blocks[T any](p Parser[T]) Parser[[]T] {
	return Slice(split.Blocks, p)
}

// Map parses a map. It uses split to separate the input into key/value pairs
// and then uses cut to cut them up into a key and a value, which are parsed
// using k and v respectively.
//
// It errors if a key is duplicated.
func Map[K comparable, V any](split, cut split.Func, k Parser[K], v Parser[V]) Parser[map[K]V] {
	return func(in string) (map[K]V, error) {
		sp, err := split(in)
		if err != nil {
			return nil, err
		}
		m := make(map[K]V)
		for _, s := range sp {
			sp, err := cut(s)
			if err != nil || len(sp) != 2 {
				return nil, fmt.Errorf("invalid key/value pair %q", s)
			}
			kk, err := k(sp[0])
			if err != nil {
				return nil, err
			}
			vv, err := v(sp[1])
			if err != nil {
				return nil, err
			}
			if _, ok := m[kk]; ok {
				return nil, fmt.Errorf("duplicate key %q", sp[0])
			}
			m[kk] = vv
		}
		return m, nil
	}
}

var (
	stringType = reflect.TypeOf("")
	errorType  = reflect.TypeOf(new(error)).Elem()
)

// isParser returns T, if t is a Parser[T]. Otherwise, returns nil.
func isParser(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Func {
		return nil
	}
	if t.NumIn() != 1 || t.In(0) != stringType {
		return nil
	}
	if t.NumOut() != 2 || t.Out(1) != errorType {
		return nil
	}
	return t.Out(0)
}

// Any runs a set of parser, one after the other, returning the first succesful result.
// The parsers don't have to return the same type, but they must return
// something assignable to T.
func Any[T any](p ...any) Parser[T] {
	rp := make([]func(string, reflect.Value) error, len(p))
	tt := reflect.TypeOf(new(T)).Elem()
	for i, p := range p {
		rf := reflect.ValueOf(p)
		t := isParser(rf.Type())
		if t == nil {
			panic(fmt.Errorf("%T is not a Parser[T]", p))
		}
		if !t.AssignableTo(tt) {
			panic(fmt.Errorf("%v is not assignable to %v", t, tt))
		}
		rp[i] = func(s string, rv reflect.Value) error {
			out := rf.Call([]reflect.Value{reflect.ValueOf(s)})
			rv.Set(out[0])
			if e := out[1].Interface(); e != nil {
				return e.(error)
			}
			return nil
		}
	}
	return func(in string) (v T, err error) {
		rv := reflect.ValueOf(&v).Elem()
		for _, p := range rp {
			if err := p(in, rv); err == nil {
				return v, err
			}
		}
		return v, fmt.Errorf("can not parse %q as %T", in, v)
	}
}

// Struct splits a string and parses the result into the fields of a struct
// using the given parsers.
//
// S must be a struct type with the same number of exported fields as split
// returns. The variadic fields argument must all be of type Parser[T] and T
// must match the exported field of S in the same sequence.
func Struct[S any](s split.Func, fields ...any) Parser[S] {
	type field struct {
		idx int
		p   func(string, reflect.Value) error
	}
	var fs []field
	rt := reflect.TypeOf(new(S)).Elem()
	for i, n := 0, rt.NumField(); i < n; i++ {
		rf := rt.Field(i)
		if !rf.IsExported() {
			continue
		}
		if len(fields) == 0 {
			panic(fmt.Errorf("Struct[%T]: too few parsers for number of exported fields of type %v", *new(S), rt))
		}
		p := reflect.ValueOf(fields[0])
		fields = fields[1:]
		pt := p.Type()
		tt := isParser(pt)
		if tt == nil {
			panic(fmt.Errorf("Struct[%T]: %v is not a Parser[T]", *new(S), pt))
		}
		if !tt.AssignableTo(rf.Type) {
			panic(fmt.Errorf("Struct[%T]: %v is not assignable to type %v of field %v.%s", *new(S), tt, rf.Type, rt, rf.Name))
		}

		rp := func(in string, rv reflect.Value) error {
			out := p.Call([]reflect.Value{reflect.ValueOf(in)})
			rv.Set(out[0])
			if e := out[1].Interface(); e != nil {
				return e.(error)
			}
			return nil
		}
		fs = append(fs, field{
			idx: i,
			p:   rp,
		})
	}
	return func(in string) (S, error) {
		var v S
		sp, err := s(in)
		if err != nil {
			return v, err
		}
		if len(sp) != len(fs) {
			return v, fmt.Errorf("Struct[%T]: got %q for %d fields", *new(S), sp, len(fs))
		}
		rv := reflect.ValueOf(&v).Elem()
		for i, f := range fs {
			if err := f.p(sp[i], rv.Field(f.idx)); err != nil {
				return v, err
			}
		}
		return v, nil
	}
}

func refl(rt reflect.Type) func(in string, p reflect.Value) error {
	var (
		indir    int
		isParser bool
		ptr      bool
	)
	for rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
		indir++
	}
	isParser = rt.Implements(parseT)
	if !isParser {
		if reflect.PointerTo(rt).Implements(parseT) {
			isParser, ptr = true, true
		}
	}
	return func(in string, rv reflect.Value) error {
		for i := 0; i < indir; i++ {
			rv.Set(reflect.New(rv.Type().Elem()))
			rv = rv.Elem()
		}
		if isParser {
			var p parseIface
			if ptr {
				p = rv.Addr().Interface().(parseIface)
			} else {
				p = rv.Interface().(parseIface)
			}
			if err := p.Parse(in); err != nil {
				return fmt.Errorf("parsing %q: %w", in, err)
			}
			return nil
		}
		switch rt.Kind() {
		case reflect.Bool:
			switch strings.ToLower(in) {
			case "true":
				rv.SetBool(true)
			case "false":
				rv.SetBool(false)
			default:
				return fmt.Errorf("can not parse %q as %v", in, rt)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(in, 0, 64)
			if err != nil {
				return fmt.Errorf("can not parse %q as %v", in, rt)
			}
			rv.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(in, 0, 64)
			if err != nil {
				return fmt.Errorf("can not parse %q as %v", in, rt)
			}
			rv.SetUint(n)
		default:
			return fmt.Errorf("unknown field type %v", rt)
		}
		return nil
	}
}

type parseIface interface {
	Parse(string) error
}

var parseT = reflect.TypeOf(new(parseIface)).Elem()

// String parses any string as itself.
func String[T ~string](in string) (T, error) {
	return T(in), nil
}

// Prefix expects a prefix and parses the rest of the string using p.
func Prefix[T any](prefix string, p Parser[T]) Parser[T] {
	return func(in string) (T, error) {
		rest, ok := strings.CutPrefix(in, prefix)
		if !ok {
			return *new(T), fmt.Errorf("expect %q", prefix)
		}
		return p(rest)
	}
}

// Unsigned parses an unsigned number using strconv.ParseUint.
func Unsigned[T constraints.Unsigned](in string) (T, error) {
	var v T
	n, err := strconv.ParseUint(in, 0, 64)
	if err != nil {
		return v, fmt.Errorf("parsing %q as %T: %w", in, v, err)
	}
	if uint64(T(n)) != n {
		return v, fmt.Errorf("%v overflows %T", n, v)
	}
	return T(n), nil
}

// Signed parses a signed number using strconv.ParseInt.
func Signed[T constraints.Signed](in string) (T, error) {
	var v T
	n, err := strconv.ParseInt(in, 0, 64)
	if err != nil {
		return v, fmt.Errorf("parsing %q as %T: %w", in, v, err)
	}
	if int64(T(n)) != n {
		return v, fmt.Errorf("%v overflows %T", n, v)
	}
	return T(n), nil
}

// Rune parses a single UTF-8 codepoint.
func Rune(in string) (rune, error) {
	r, size := utf8.DecodeRuneInString(in)
	if size != len(in) || r == utf8.RuneError {
		return 0, fmt.Errorf("expected single codepoint, got %q", in)
	}
	return r, nil
}

// Enum parses as any of opts.
func Enum[T ~byte | ~rune | ~string](opts ...T) Parser[T] {
	return func(in string) (v T, err error) {
		for _, o := range opts {
			if string(o) == in {
				return o, nil
			}
		}
		return v, fmt.Errorf("expected one of %q", opts)
	}
}
