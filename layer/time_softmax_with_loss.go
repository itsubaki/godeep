package layer

import (
	"fmt"

	"github.com/itsubaki/neu/math/matrix"
)

type TimeSoftmaxWithLoss struct {
	ts []matrix.Matrix
	ys []matrix.Matrix
}

func (l *TimeSoftmaxWithLoss) Params() []matrix.Matrix      { return make([]matrix.Matrix, 0) }
func (l *TimeSoftmaxWithLoss) Grads() []matrix.Matrix       { return make([]matrix.Matrix, 0) }
func (l *TimeSoftmaxWithLoss) SetParams(p ...matrix.Matrix) {}
func (l *TimeSoftmaxWithLoss) SetState(h ...matrix.Matrix)  {}
func (l *TimeSoftmaxWithLoss) ResetState()                  {}

func (l *TimeSoftmaxWithLoss) Forward(xs, ts []matrix.Matrix, _ ...Opts) []matrix.Matrix {
	T, V := len(xs), len(xs[0][0])
	l.ys = make([]matrix.Matrix, T)
	l.ts = oneHot(ts, V)

	// naive
	var loss float64
	for t := 0; t < T; t++ {
		l.ys[t] = Softmax(xs[t])
		loss += Loss(l.ys[t], l.ts[t])
	}

	return []matrix.Matrix{{{loss / float64(T)}}}
}

func (l *TimeSoftmaxWithLoss) Backward(dout []matrix.Matrix) []matrix.Matrix {
	T := len(l.ys)
	dx := make([]matrix.Matrix, T)

	// naive
	for t := T - 1; t > -1; t-- {
		size, _ := l.ts[t].Dimension()
		dx[t] = l.ys[t].Sub(l.ts[t]).Mul(dout[0]).MulC(1.0 / float64(size)) // (y - t) * dout / size
	}

	return dx
}

func (l *TimeSoftmaxWithLoss) String() string {
	return fmt.Sprintf("%T", l)
}

func oneHot(ts []matrix.Matrix, size int) []matrix.Matrix {
	out := make([]matrix.Matrix, 0)
	for _, t := range ts {
		m := make(matrix.Matrix, 0)
		for _, r := range t {
			for _, v := range r {
				onehot := make([]float64, size)
				onehot[int(v)] = 1.0
				m = append(m, onehot)
			}
		}

		out = append(out, m)
	}

	return out
}
