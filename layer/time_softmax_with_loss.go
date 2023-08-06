package layer

import (
	"fmt"

	"github.com/itsubaki/neu/math/matrix"
)

type TimeSoftmaxWithLoss struct {
	layer []*SoftmaxWithLoss
	xs    []matrix.Matrix
}

func (l *TimeSoftmaxWithLoss) Params() []matrix.Matrix      { return make([]matrix.Matrix, 0) }
func (l *TimeSoftmaxWithLoss) Grads() []matrix.Matrix       { return make([]matrix.Matrix, 0) }
func (l *TimeSoftmaxWithLoss) SetParams(p ...matrix.Matrix) {}
func (l *TimeSoftmaxWithLoss) SetState(h ...matrix.Matrix)  {}
func (l *TimeSoftmaxWithLoss) ResetState()                  {}

func (l *TimeSoftmaxWithLoss) Forward(xs, ts []matrix.Matrix, _ ...Opts) []matrix.Matrix {
	T, V := len(xs), len(xs[0][0])
	ots := oneHot(ts, V)
	l.layer = make([]*SoftmaxWithLoss, T)
	l.xs = xs

	var loss float64
	for t := 0; t < T; t++ {
		l.layer[t] = &SoftmaxWithLoss{}
		loss += l.layer[t].Forward(xs[t], ots[t])[0][0]
	}

	return []matrix.Matrix{matrix.New([]float64{loss / float64(T)})}
}

func (l *TimeSoftmaxWithLoss) Backward(dout []matrix.Matrix) []matrix.Matrix {
	T := len(l.xs)
	dx := make([]matrix.Matrix, T)
	do := dout[0].MulC(1.0 / float64(T))

	// naive
	for t := 0; t < T; t++ {
		dx[t], _ = l.layer[t].Backward(do)
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
