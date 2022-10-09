package matrix

type Matrix [][]float64

func New(v ...[]float64) Matrix {
	out := make(Matrix, len(v))
	copy(out, v)
	return out
}

func Fill(w float64, n, m int) Matrix {
	out := make(Matrix, 0)
	for i := 0; i < n; i++ {
		v := make([]float64, 0)
		for j := 0; j < m; j++ {
			v = append(v, w)
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Shape() (int, int) {
	return m.Dimension()
}

func (m Matrix) Dimension() (int, int) {
	if len(m) == 0 {
		return 0, 0
	}

	return len(m), len(m[0])
}

func (m Matrix) Apply(n Matrix) Matrix {
	a, b := n.Dimension()
	_, p := m.Dimension()

	out := Matrix{}
	for i := 0; i < a; i++ {
		v := make([]float64, 0)

		for j := 0; j < p; j++ {
			var c float64
			for k := 0; k < b; k++ {
				c = c + n[i][k]*m[k][j]
			}

			v = append(v, c)
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Dot(n Matrix) Matrix {
	return n.Apply(m)
}

func (m Matrix) Add(n Matrix) Matrix {
	p, q := m.Dimension()

	out := make(Matrix, 0, p)
	for i := 0; i < p; i++ {
		v := make([]float64, 0, q)

		for j := 0; j < q; j++ {
			v = append(v, m[i][j]+n[i][j])
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Sub(n Matrix) Matrix {
	p, q := m.Dimension()

	out := make(Matrix, 0, p)
	for i := 0; i < p; i++ {
		v := make([]float64, 0, q)

		for j := 0; j < q; j++ {
			v = append(v, m[i][j]-n[i][j])
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Mul(n Matrix) Matrix {
	p, q := m.Dimension()

	out := make(Matrix, 0, p)
	for i := 0; i < p; i++ {
		v := make([]float64, 0, q)

		for j := 0; j < q; j++ {
			v = append(v, m[i][j]*n[i][j])
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Transpose() Matrix {
	p, q := m.Dimension()

	out := make(Matrix, 0, p)
	for i := 0; i < p; i++ {
		v := make([]float64, 0, q)
		for j := 0; j < q; j++ {
			v = append(v, m[j][i])
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) T() Matrix {
	return m.Transpose()
}

func Dot(m, n Matrix) Matrix {
	return m.Dot(n)
}
