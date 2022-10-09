package matrix

import (
	"github.com/itsubaki/neu/activation"
	"github.com/itsubaki/neu/loss"
)

type Matrix [][]float64

func New(v ...[]float64) Matrix {
	out := make(Matrix, len(v))
	copy(out, v)
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

func (m Matrix) Addf64(w float64) Matrix {
	p, q := m.Dimension()

	out := make(Matrix, 0, p)
	for i := 0; i < p; i++ {
		v := make([]float64, 0, q)

		for j := 0; j < q; j++ {
			v = append(v, m[i][j]+w)
		}

		out = append(out, v)
	}

	return out
}

func (m Matrix) Mulf64(w float64) Matrix {
	p, q := m.Dimension()

	out := make(Matrix, 0, p)
	for i := 0; i < p; i++ {
		v := make([]float64, 0, q)

		for j := 0; j < q; j++ {
			v = append(v, m[i][j]*w)
		}

		out = append(out, v)
	}

	return out
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

func Dot(m, n Matrix) Matrix {
	return m.Dot(n)
}

func Sigmoid(m Matrix) Matrix {
	out := make(Matrix, 0)
	p, q := m.Shape()

	for i := 0; i < p; i++ {
		v := make([]float64, 0)
		for j := 0; j < q; j++ {
			v = append(v, activation.Sigmoid(m[i][j]))
		}

		out = append(out, v)
	}

	return out
}

func CrossEntropyError(y, t Matrix) []float64 {
	out := make([]float64, 0)
	for i := range y {
		out = append(out, loss.CrossEntropyError(y[i], t[i]))
	}

	return out
}

func Identity(m Matrix) Matrix {
	return m
}

func SumAxis1(m Matrix) []float64 {
	p, q := m.Shape()

	out := make([]float64, q)
	for i := 0; i < q; i++ {
		for j := 0; j < p; j++ {
			out[i] = out[i] + m[j][i]
		}
	}

	return out
}

func Mask(m Matrix, mask [][]bool) Matrix {
	out := make(Matrix, 0)
	for i := range m {
		v := make([]float64, 0)
		for j := range m[i] {
			if mask[i][j] {
				v = append(v, 0)
				continue
			}

			v = append(v, m[i][j])
		}

		out = append(out, v)
	}

	return out
}
