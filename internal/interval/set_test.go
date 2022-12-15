package interval

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"gonih.org/rnd"
)

func FuzzSet(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		calls, err := decodeCalls(b)
		if err != nil {
			return
		}
		type Interval = CO[int16]
		s := new(Set[Interval, int16])
		for _, c := range calls {
			switch c.Op {
			case opAdd:
				i := Interval{c.Arg1, c.Arg2}
				s.Add(i)
				s.set.check()
				if i.Empty() {
					continue
				}
				t.Log(i.Max, i.Min, i.Max-i.Min)
				v := int16(rnd.Intn(int(i.Max)-int(i.Min))) + i.Min
				if !s.Contains(v) {
					t.Fatalf("s.Contains(%v) = false, want true", v)
				}
			case opIntersect:
				i := Interval{c.Arg1, c.Arg2}
				s.Intersect(i)
				s.set.check()
				for _, j := range s.Intervals() {
					if j.Min < i.Min || j.Max > i.Max {
						t.Fatalf("s.Intervals contains %v, which is not a subset of %v", j, i)
					}
				}
			default:
				return
			}
		}
	})
}

func decodeCalls(b []byte) ([]call, error) {
	var out []call
	r := bytes.NewReader(b)
	for r.Len() > 0 {
		var c call
		if err := binary.Read(r, binary.BigEndian, &c); err != nil {
			return nil, err
		}
		if c.Arg1 > c.Arg2 {
			return nil, errors.New("invalid interval in fuzzing input")
		}
		out = append(out, c)
	}
	return out, nil
}

const (
	opAdd = iota
	opIntersect
)

type call struct {
	Op   int16
	Arg1 int16
	Arg2 int16
}
