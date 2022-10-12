package matrix

import (
	"math/rand"
)

type Matrix [][]float64

func New(v ...[]float64) Matrix {
	out := make(Matrix, len(v))
	copy(out, v)
	return out
}

func Zero(n, m int) Matrix {
	out := make(Matrix, n)
	for i := 0; i < n; i++ {
		out[i] = make([]float64, m)
	}

	return out
}

func Randn(n, m int) Matrix {
	out := make(Matrix, 0)
	for i := 0; i < n; i++ {
		v := make([]float64, 0)
		for j := 0; j < m; j++ {
			v = append(v, rand.NormFloat64())
		}

		out = append(out, v)
	}

	return out
}

func Rand(n, m int) Matrix {
	out := make(Matrix, 0)
	for i := 0; i < n; i++ {
		v := make([]float64, 0)
		for j := 0; j < m; j++ {
			v = append(v, rand.Float64())
		}

		out = append(out, v)
	}

	return out
}

func Mask(x Matrix, f func(x float64) bool) Matrix {
	out := make(Matrix, 0)
	for i := range x {
		v := make([]float64, 0)
		for j := range x[i] {
			if f(x[i][j]) {
				v = append(v, 0)
				continue
			}

			v = append(v, 1)
		}

		out = append(out, v)
	}

	return out
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
	return m.FuncWith(n, func(a, b float64) float64 { return a + b })
}

func (m Matrix) Mul(n Matrix) Matrix {
	return m.FuncWith(n, func(a, b float64) float64 { return a * b })
}

func (m Matrix) Transpose() Matrix {
	p, q := m.Dimension()

	out := make(Matrix, q)
	for i := range out {
		out[i] = make([]float64, p)
	}

	for i := 0; i < q; i++ {
		for j := 0; j < p; j++ {
			out[i][j] = m[j][i]
		}
	}

	return out
}

func (m Matrix) T() Matrix {
	return m.Transpose()
}

func (m Matrix) Func(f func(v float64) float64) Matrix {
	p, q := m.Dimension()

	out := make(Matrix, 0, p)
	for i := 0; i < p; i++ {
		v := make([]float64, 0, q)

		for j := 0; j < q; j++ {
			v = append(v, f(m[i][j]))
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) FuncWith(n Matrix, f func(a, b float64) float64) Matrix {
	p, q := m.Dimension()

	out := make(Matrix, 0, p)
	for i := 0; i < p; i++ {
		v := make([]float64, 0, q)

		for j := 0; j < q; j++ {
			v = append(v, f(m[i][j], n[i][j]))
		}

		out = append(out, v)
	}

	return out
}

func Dot(m, n Matrix) Matrix {
	return m.Dot(n)
}

func Func(m Matrix, f func(a float64) float64) Matrix {
	return m.Func(f)
}

func FuncWith(m, n Matrix, f func(a, b float64) float64) Matrix {
	return m.FuncWith(n, f)
}
