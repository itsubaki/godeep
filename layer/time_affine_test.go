package layer_test

import (
	"fmt"

	"github.com/itsubaki/neu/layer"
	"github.com/itsubaki/neu/math/matrix"
)

func ExampleTimeAffine() {
	affine := &layer.TimeAffine{
		W: matrix.New([]float64{0.1, 0.3, 0.5}, []float64{0.2, 0.4, 0.6}),
		B: matrix.New([]float64{0.1, 0.2, 0.3}),
	}

	xs := []matrix.Matrix{matrix.New([]float64{1.0, 0.5})}
	fmt.Println(affine)
	fmt.Println(affine.Forward(xs, nil))
	fmt.Println(affine.Backward(xs))

	// Output:
	// *layer.TimeAffine: W(2, 3)*T, B(1, 3)*T: 9*T
	// [[[0.30000000000000004 0.7 1.1]]]
	// [[[0.25 0.4]]]
}

func ExampleTimeAffine_Params() {
	affine := &layer.TimeAffine{}
	affine.SetParams(make([]matrix.Matrix, 2)...)

	fmt.Println(affine.Params())
	fmt.Println(affine.Grads())

	// Output:
	// [[] []]
	// [[] []]
}

func ExampleTimeAffine_state() {
	affine := &layer.TimeAffine{}
	affine.SetState(matrix.New())
	affine.ResetState()

	// Output:
}