package optimizer

import "github.com/itsubaki/neu/math/matrix"

type SGD struct {
	LearningRate float64
}

func (o *SGD) Update(params, grads map[string]matrix.Matrix) map[string]matrix.Matrix {
	// TODO
	out := make(map[string]matrix.Matrix)
	return out
}