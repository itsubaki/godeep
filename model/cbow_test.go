package model_test

import (
	"fmt"
	"math/rand"

	"github.com/itsubaki/neu/math/matrix"
	"github.com/itsubaki/neu/model"
)

func ExampleCBOW() {
	// you, say, goodbye, and, I, hello, .

	// context data
	c := []matrix.Matrix{
		matrix.New(
			[]float64{1, 0, 0, 0, 0, 0, 0}, // you
			[]float64{0, 0, 1, 0, 0, 0, 0}, // goodbye
		),
	}
	t := []matrix.Matrix{
		matrix.New(
			[]float64{0, 1, 0, 0, 0, 0, 0}, // say
		),
	}

	// model
	s := rand.NewSource(1)
	m := model.NewCBOW(&model.CBOWConfig{
		VocabSize:  7,
		HiddenSize: 5,
	}, s)

	// layer
	fmt.Printf("%T\n", m)
	for i, l := range m.Layers() {
		fmt.Printf("%2d: %v\n", i, l)
	}
	fmt.Println()

	y := m.Predict([]matrix.Matrix{})
	loss := m.Forward(c, t)
	dout := m.Backward()

	fmt.Println(y)
	fmt.Println(loss)
	fmt.Println(dout)
	for _, g := range m.Grads() {
		fmt.Println(g)
	}

	// Output:
	// *model.CBOW
	//  0: *layer.Dot: W(7, 5): 35
	//  1: *layer.Dot: W(7, 5): 35
	//  2: *layer.Dot: W(5, 7): 35
	//  3: *layer.SoftmaxWithLoss
	//
	// []
	// [[1.9461398376656527]]
	// []
	// [[[-0.009653507324393098 -0.008984561380889693 -0.004249060448413236 0.0008831124819104136 0.008369624482009563] [0 0 0 0 0] [0 0 0 0 0] [0 0 0 0 0] [0 0 0 0 0] [0 0 0 0 0] [0 0 0 0 0]]]
	// [[[0 0 0 0 0] [0 0 0 0 0] [-0.009653507324393098 -0.008984561380889693 -0.004249060448413236 0.0008831124819104136 0.008369624482009563] [0 0 0 0 0] [0 0 0 0 0] [0 0 0 0 0] [0 0 0 0 0]]]
	// [[[-0.00192072813742285 0.011525859598972725 -0.0019206163092527154 -0.001920913611935071 -0.001920721577145852 -0.0019217493698969338 -0.0019211305933193024] [0.0006047213303225987 -0.0036287973368028205 0.0006046861223832273 0.0006047797250488412 0.0006047192648874358 0.0006050428547738025 0.0006048480393869149] [-0.000582952848866384 0.0034981695524440096 -0.000582918908326287 -0.0005830091415260189 -0.0005829508577818099 -0.0005832627992276676 -0.0005830749967158421] [0.0019968435276115203 -0.011982611017116031 0.001996727267869505 0.0019970363521826027 0.001996836707361216 0.001997905229898444 0.001997261932192742] [-0.00010665569552509538 0.0006400169540403842 -0.00010664948584292713 -0.00010666599469899599 -0.00010665533124089898 -0.0001067124033213181 -0.00010667804341114859]]]

}
