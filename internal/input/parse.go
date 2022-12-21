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

// Array splits a string using split and calls p for each piece. A must be an
// array type with element type T.
func Array[A, T any](split Splitter, p Parser[T]) Parser[A] {
	tt := reflect.TypeOf(new(T)).Elem()
	at := reflect.TypeOf(new(A)).Elem()
	if at.Kind() != reflect.Array || at.Elem() != tt {
		panic(fmt.Errorf("%v is not an […]%v array type", at, tt))
	}
	n := at.Len()
	return func(s string) (A, error) {
		var a A
		sp, err := split(s)
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
	return func(s string) (B, error) {
		a, err := p(s)
		if err != nil {
			return *new(B), err
		}
		return f(a)
	}
}

// Slice splits a string using split and calls p on each piece.
func Slice[T any](split Splitter, p Parser[T]) Parser[[]T] {
	return func(s string) ([]T, error) {
		var out []T
		sp, err := split(s)
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

var (
	stringType = reflect.TypeOf("")
	errorType  = reflect.TypeOf(new(error)).Elem()
)

// Struct splits a string using split and parses the result into the fields of
// a struct using the given parsers.
//
// S must be a struct type with the same number of exported fields as split
// returns. The variadic fields argument must all be of type Parser[T] and T
// must match the exported field of S in the same sequence.
func Struct[S any](split Splitter, fields ...any) Parser[S] {
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
			panic(fmt.Errorf("too few parsers for number of exported fields of type %v", rt))
		}
		p := reflect.ValueOf(fields[0])
		fields = fields[1:]
		pt := p.Type()
		if pt.Kind() != reflect.Func || pt.NumIn() != 1 || pt.In(0) != stringType || pt.NumOut() != 2 || pt.Out(0) != rf.Type || pt.Out(1) != errorType {
			panic(fmt.Errorf("%v.%s does not match type %v", rt, rf.Name, pt.Out(0)))
		}

		rp := func(s string, rv reflect.Value) error {
			out := p.Call([]reflect.Value{reflect.ValueOf(s)})
			rv.Set(out[0])
			return out[1].Interface().(error)
		}
		fs = append(fs, field{
			idx: i,
			p:   rp,
		})
	}
	return func(s string) (S, error) {
		var v S
		sp, err := split(s)
		if err != nil {
			return v, err
		}
		if len(sp) != len(fields) {
			return v, fmt.Errorf("got %d strings for %d fields", len(sp), len(fields))
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

func refl(rt reflect.Type) func(s string, p reflect.Value) error {
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
	return func(s string, rv reflect.Value) error {
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

// String parses any string as itself.
//
// It is a Parser[T].
func String[T ~string]() Parser[T] {
	return func(s string) (T, error) {
		return T(s), nil
	}
}

// Unsigned parses an unsigned number using strconv.ParseUint.
//
// It is a Parser[T].
func Unsigned[T constraints.Unsigned]() Parser[T] {
	return func(s string) (T, error) {
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
}

// Signed parses a signed number using strconv.ParseInt.
//
// It is a Parser[T].
func Signed[T constraints.Signed]() Parser[T] {
	return func(s string) (T, error) {
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
}

// Rune parses a single UTF-8 codepoint.
//
// It is a Parser[rune].
func Rune() Parser[rune] {
	return func(s string) (rune, error) {
		r, size := utf8.DecodeRuneInString(s)
		if size != len(s) || r == utf8.RuneError {
			return 0, fmt.Errorf("expected single codepoint, got %q", s)
		}
		return r, nil
	}
}

type Splitter func(string) ([]string, error)

// Split into blocks, separated by empty lines.
func Blocks() Splitter {
	return Split("\n\n")
}

// Split into whitespace-separated fields.
func Fields() Splitter {
	return func(s string) ([]string, error) {
		return strings.Fields(s), nil
	}
}

// Split into lines.
func Lines() Splitter {
	return Split("\n")
}

// Split along sep.
func Split(sep string) Splitter {
	return func(s string) ([]string, error) {
		return strings.Split(s, sep), nil
	}
}

// Split into at most n pieces, along sep.
func SplitN(sep string, n int) Splitter {
	return func(s string) ([]string, error) {
		return strings.SplitN(s, sep, n), nil
	}
}
