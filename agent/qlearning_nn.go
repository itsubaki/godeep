package agent

import (
	"math/rand"

	"github.com/itsubaki/neu/math/matrix"
	"github.com/itsubaki/neu/math/vector"
	"github.com/itsubaki/neu/model"
	"github.com/itsubaki/neu/optimizer"
)

type QLearningAgentNN struct {
	Gamma      float64
	Epsilon    float64
	ActionSize int
	Q          *model.QNet
	Optimizer  *optimizer.SGD
	Source     rand.Source
}

func (a *QLearningAgentNN) GetAction(state [][]float64) int {
	rng := rand.New(a.Source)
	if a.Epsilon > rng.Float64() {
		return rng.Intn(a.ActionSize)
	}

	qs := a.Q.Predict(matrix.New(state...))
	return vector.Argmax(qs[0])
}

func (a *QLearningAgentNN) Update(state [][]float64, action int, reward float64, next [][]float64, done bool) matrix.Matrix {
	var nextq float64
	if !done {
		nextqs := a.Q.Predict(matrix.New(next...))
		nextq = vector.Max(nextqs[0])
	}

	target := matrix.New([]float64{reward + a.Gamma*nextq})
	qs := a.Q.Predict(matrix.New(state...))
	q := matrix.New([]float64{qs[0][action]})
	loss := a.Q.Loss(target, q)

	a.Q.Backward()
	a.Optimizer.Update(a.Q)
	return loss
}
