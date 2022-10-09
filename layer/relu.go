package layer

import "github.com/itsubaki/neu/math/matrix"

type ReLU struct {
	mask [][]bool
}

func (l *ReLU) Forward(x, _ matrix.Matrix) matrix.Matrix {
	l.mask = mask(x)
	return matrix.Mask(x, l.mask)
}

func (l *ReLU) Backward(dout matrix.Matrix) (matrix.Matrix, matrix.Matrix) {
	dx := matrix.Mask(dout, l.mask)
	return dx, matrix.New()
}

func mask(x matrix.Matrix) [][]bool {
	out := make([][]bool, 0)
	for i := range x {
		v := make([]bool, 0)
		for j := range x[i] {
			v = append(v, x[i][j] <= 0)
		}

		out = append(out, v)
	}

	return out
}
