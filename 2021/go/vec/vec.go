package vec

import "strconv"

type V struct {
	X, Y, Z int
}

func Vec(x, y, z int) V {
	return V{x, y, z}
}

func (v V) Add(w V) V {
	return V{v.X + w.X, v.Y + w.Y, v.Z + w.Z}
}

func (v V) Sub(w V) V {
	return V{v.X - w.X, v.Y - w.Y, v.Z - w.Z}
}

func (v V) Mul(w V) int {
	return v.X*w.X + v.Y*w.Y + v.Z*w.Z
}

func (v V) String() string {
	return "[" + strconv.Itoa(v.X) + "," + strconv.Itoa(v.Y) + "," + strconv.Itoa(v.Z) + "]"
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func DistL1(v, w V) int {
	return abs(v.X-w.X) + abs(v.Y-w.Y) + abs(v.Z-w.Z)
}

type M struct {
	A11, A12, A13 int
	A21, A22, A23 int
	A31, A32, A33 int
}

func Mat(a11, a12, a13, a21, a22, a23, a31, a32, a33 int) M {
	return M{a11, a12, a13, a21, a22, a23, a31, a32, a33}
}

func ID() M {
	return M{1, 0, 0, 0, 1, 0, 0, 0, 1}
}

func (m M) MulV(v V) V {
	return V{
		m.A11*v.X + m.A12*v.Y + m.A13*v.Z,
		m.A21*v.X + m.A22*v.Y + m.A23*v.Z,
		m.A31*v.X + m.A32*v.Y + m.A33*v.Z,
	}
}

func (m M) MulM(n M) M {
	return M{
		m.A11*n.A11 + m.A12*n.A21 + m.A13*n.A31, m.A11*n.A12 + m.A12*n.A22 + m.A13*n.A32, m.A11*n.A13 + m.A12*n.A23 + m.A13*n.A33,
		m.A21*n.A11 + m.A22*n.A21 + m.A23*n.A31, m.A21*n.A12 + m.A22*n.A22 + m.A23*n.A32, m.A21*n.A13 + m.A22*n.A23 + m.A23*n.A33,
		m.A31*n.A11 + m.A32*n.A21 + m.A33*n.A31, m.A31*n.A12 + m.A32*n.A22 + m.A33*n.A32, m.A31*n.A13 + m.A32*n.A23 + m.A33*n.A33,
	}
}

func (m M) String() string {
	return "[[" + strconv.Itoa(m.A11) + "," + strconv.Itoa(m.A12) + "," + strconv.Itoa(m.A13) + "]," +
		"[" + strconv.Itoa(m.A21) + "," + strconv.Itoa(m.A22) + "," + strconv.Itoa(m.A23) + "]," +
		"[" + strconv.Itoa(m.A31) + "," + strconv.Itoa(m.A32) + "," + strconv.Itoa(m.A33) + "]]"
}
