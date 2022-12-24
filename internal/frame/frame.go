package frame

import (
	"bufio"
	"bytes"
	"io"
	"unicode/utf8"

	"github.com/Merovius/AdventOfCode/internal/math"
)

type Border string

const (
	Simple Border = "┌─┐││└─┘"
)

func New(w io.Writer, b Border) io.WriteCloser {
	return &frame{
		w:      &errWriter{w: w},
		buf:    new(bytes.Buffer),
		border: []rune(b),
	}
}

type frame struct {
	w      *errWriter
	buf    *bytes.Buffer
	border []rune
}

func (w *frame) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *frame) Close() error {
	cols := math.MinInt
	s := bufio.NewScanner(bytes.NewReader(w.buf.Bytes()))
	for s.Scan() {
		cols = math.Max(cols, utf8.RuneCount(s.Bytes()))
	}
	if s.Err() != nil {
		return s.Err()
	}
	io.WriteString(w.w, string(w.border[0]))
	for i := 0; i < cols; i++ {
		io.WriteString(w.w, string(w.border[1]))
	}
	io.WriteString(w.w, string(w.border[2])+"\n")
	s = bufio.NewScanner(w.buf)
	for s.Scan() {
		io.WriteString(w.w, string(w.border[3]))
		w.w.Write(s.Bytes())
		io.WriteString(w.w, string(w.border[4])+"\n")
	}
	io.WriteString(w.w, string(w.border[5]))
	for i := 0; i < cols; i++ {
		io.WriteString(w.w, string(w.border[6]))
	}
	io.WriteString(w.w, string(w.border[7])+"\n")
	return w.w.err
}

type errWriter struct {
	w   io.Writer
	err error
}

func (w *errWriter) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}
	n, w.err = w.w.Write(p)
	return n, w.err
}

func (w *errWriter) WriteString(s string) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}
	if sw, ok := w.w.(io.StringWriter); ok {
		n, w.err = sw.WriteString(s)
	} else {
		n, w.err = w.w.Write([]byte(s))
	}
	return n, w.err
}
