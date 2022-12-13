package input

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

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

// SplitBlocks splits s into blocks separated by empty lines.
func SplitBlocks(s string) []string {
	return strings.Split(s, "\n\n")
}

// SplitLines splits s into lines.
func SplitLines(s string) []string {
	return strings.Split(s, "\n")
}

// SplitFunc splits a string using split and calls p on each piece.
func SplitFunc[T any](split func(string) []string, p Parser[T]) Parser[[]T] {
	return func(s string) ([]T, error) {
		var out []T
		for _, s := range split(s) {
			v, err := p(s)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	}

}

// Split splits a string on sep and calls p on each piece.
func Split[T any](sep string, p Parser[T]) Parser[[]T] {
	return SplitFunc(func(s string) []string {
		return strings.Split(s, sep)
	}, p)
}

// Blocks splits a string into blocks separated by empty lines and calls p on
// each block.
func Blocks[T any](p Parser[T]) Parser[[]T] {
	return Split("\n\n", p)
}

// Lines splits a string into lines and calls p on each line.
func Lines[T any](p Parser[T]) Parser[[]T] {
	return Split("\n", p)
}

// Fields splits a string according to strings.Fields and calls p on each
// field.
func Fields[T any](p Parser[T]) Parser[[]T] {
	return SplitFunc(strings.Fields, p)
}

// Array splits a string using split and calls p for each piece. A must be an
// array type with element type T.
func Array[A, T any](split func(string) []string, p Parser[T]) Parser[A] {
	tt := reflect.TypeOf(new(T)).Elem()
	at := reflect.TypeOf(new(A)).Elem()
	if at.Kind() != reflect.Array || at.Elem() != tt {
		panic(fmt.Errorf("%v is not an […]%v array type", at, tt))
	}
	n := at.Len()
	return func(s string) (A, error) {
		var a A
		sp := split(s)
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

// Reflect builds a reflection based parser for T. It panics if T can't be
// parsed automatically.
func Reflect[T any]() Parser[T] {
	p := refl(reflect.TypeOf(new(T)).Elem())
	return func(s string) (T, error) {
		var v T
		err := p(reflect.ValueOf(&v).Elem(), s)
		return v, err
	}
}

func refl(rt reflect.Type) func(p reflect.Value, s string) error {
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
	return func(rv reflect.Value, s string) error {
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
			if err := p.Parse(s); err != nil {
				return fmt.Errorf("parsing %q: %w", s, err)
			}
			return nil
		}
		switch rt.Kind() {
		case reflect.Bool:
			switch strings.ToLower(s) {
			case "true":
				rv.SetBool(true)
			case "false":
				rv.SetBool(false)
			default:
				return fmt.Errorf("can not parse %q as %v", s, rt)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(s, 0, 64)
			if err != nil {
				return fmt.Errorf("can not parse %q as %v", s, rt)
			}
			rv.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, err := strconv.ParseUint(s, 0, 64)
			if err != nil {
				return fmt.Errorf("can not parse %q as %v", s, rt)
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

//func Struct[T any]() Parser[T] {
//	rt := reflect.TypeOf(new(T)).Elem()
//	switch rt.Kind() {
//	case reflect.Struct:
//		return structFields[T](rt)
//	case reflect.Array:
//		n := rt.Len()
//		return func(s string) (T, error) {
//			sp := strings.Fields(s)
//			if len(sp) != n {
//				return *new(T), fmt.Errorf("expected %d fields, got %d", n, len(sp))
//			}
//
//		}
//		panic("TODO")
//	case reflect.Slice:
//		panic("TODO")
//	default:
//		panic(fmt.Errorf("Fields called with invalid type %v", rt))
//	}
//}

func structFields[T any](rt reflect.Type) Parser[T] {
	if rt.Kind() != reflect.Struct {
		panic(fmt.Errorf("structFields called with non-struct type %v", rt))
	}
	type field struct {
		i int
		p func(reflect.Value, string) error
	}
	var fields []field
	for i, n := 0, rt.NumField(); i < n; i++ {
		sf := rt.Field(i)
		if !sf.IsExported() {
			continue
		}
		p := refl(sf.Type)
		fields = append(fields, field{i, p})
	}
	return func(s string) (T, error) {
		sp := strings.Fields(s)
		if len(sp) != len(fields) {
			return *new(T), fmt.Errorf("got %d fields, expected %d", len(sp), len(fields))
		}
		var v T
		rv := reflect.ValueOf(&v).Elem()
		for i, s := range sp {
			f := fields[i]
			rf := rv.Field(f.i)
			if err := f.p(rf, s); err != nil {
				return v, err
			}
		}
		return v, nil
	}
}

// String parses any string as itself.
//
// It is a Parser[T].
func String(s string) (string, error) {
	return s, nil
}

// Unsigned parses an unsigned number using strconv.ParseUint.
//
// It is a Parser[T].
func Unsigned[T constraints.Unsigned](s string) (T, error) {
	var v T
	n, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return v, fmt.Errorf("parsing %q as %T: %w", s, v, err)
	}
	if uint64(T(n)) != n {
		return v, fmt.Errorf("%v overflows %T", n, v)
	}
	return T(n), nil
}

// Signed parses a signed number using strconv.ParseInt.
//
// It is a Parser[T].
func Signed[T constraints.Signed](s string) (T, error) {
	var v T
	n, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return v, fmt.Errorf("parsing %q as %T: %w", s, v, err)
	}
	if int64(T(n)) != n {
		return v, fmt.Errorf("%v overflows %T", n, v)
	}
	return T(n), nil
}

// Rune parses a single UTF-8 codepoint.
//
// It is a Parser[rune].
func Rune(s string) (rune, error) {
	r, size := utf8.DecodeRuneInString(s)
	if size != len(s) || r == utf8.RuneError {
		return 0, fmt.Errorf("expected single codepoint, got %q", s)
	}
	return r, nil
}

// Map converts a Parser[A] into a Parser[B] using f.
func Map[A, B any](p Parser[A], f func(A) (B, error)) Parser[B] {
	return func(s string) (B, error) {
		a, err := p(s)
		if err != nil {
			return *new(B), err
		}
		return f(a)
	}
}
