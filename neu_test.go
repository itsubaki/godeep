package neu_test

import (
	"fmt"
	"math/rand"

	"github.com/itsubaki/neu/activation"
	"github.com/itsubaki/neu/dataset/mnist"
	"github.com/itsubaki/neu/dataset/ptb"
	"github.com/itsubaki/neu/loss"
	"github.com/itsubaki/neu/math/matrix"
	"github.com/itsubaki/neu/math/numerical"
	"github.com/itsubaki/neu/math/vector"
	"github.com/itsubaki/neu/model"
	"github.com/itsubaki/neu/optimizer"
	"github.com/itsubaki/neu/trainer"
	"github.com/itsubaki/neu/weight"
)

func Example_rNNLM() {
	train := ptb.Must(ptb.Load("./testdata", ptb.TrainTxt))
	corpusSize := 1000
	corpus := train.Corpus[:corpusSize]

	// model
	s := rand.NewSource(1)
	m := model.NewRNNLM(&model.RNNLMConfig{
		VocabSize:   vector.Max(corpus) + 1,
		WordVecSize: 100,
		HiddenSize:  100,
		WeightInit:  weight.Xavier,
	}, s)

	// layer
	fmt.Printf("%T\n", m)
	for i, l := range m.Layers() {
		fmt.Printf("%2d: %v\n", i, l)
	}
	fmt.Println()

	// training
	tr := trainer.NewRNNLM(m, &optimizer.SGD{
		LearningRate: 0.1,
	})

	tr.Fit(&trainer.RNNLMInput{
		Train:      corpus[:len(corpus)-1],
		TrainLabel: corpus[1:],
		Epochs:     3,
		BatchSize:  10,
		TimeSize:   5,
		Verbose: func(epoch, j int, perplexity float64, m trainer.RNNLM) {
			if j != 18 {
				return
			}

			fmt.Printf("%2d, %2d: ppl=%.04f\n", epoch, j, perplexity)
		},
	})

	// Output:
	// *model.RNNLM
	//  0: *layer.TimeEmbedding: W(418, 100): 41800
	//  1: *layer.TimeRNN: Wx(100, 100), Wh(100, 100), B(1, 100): 20100
	//  2: *layer.TimeAffine: W(100, 418), B(1, 418): 42218
	//  3: *layer.TimeSoftmaxWithLoss
	//
	//  0, 18: ppl=372.1515
	//  1, 18: ppl=245.5703
	//  2, 18: ppl=216.1065
}

func Example_mnist() {
	// data
	train, test := mnist.Must(mnist.Load("./testdata"))

	// train data
	x := matrix.New(mnist.Normalize(train.Image)...)
	t := matrix.New(mnist.OneHot(train.Label)...)

	// test data
	xt := matrix.New(mnist.Normalize(test.Image)...)
	tt := matrix.New(mnist.OneHot(test.Label)...)

	// model
	s := rand.NewSource(1)
	m := model.NewMLP(&model.MLPConfig{
		InputSize:         mnist.Width * mnist.Height, // 24 * 24 = 784
		OutputSize:        mnist.Labels,               // 0 ~ 9
		HiddenSize:        []int{50},
		WeightInit:        weight.Std(0.01),
		BatchNormMomentum: 0.9,
	}, s)

	// layer
	fmt.Printf("%T\n", m)
	for i, l := range m.Layers() {
		fmt.Printf("%2d: %v\n", i, l)
	}
	fmt.Println()

	// training
	tr := trainer.New(m, &optimizer.SGD{
		LearningRate: 0.1,
	})

	tr.Fit(&trainer.Input{
		Train:      x[:100],
		TrainLabel: t[:100],
		Epochs:     10,
		BatchSize:  10,
		Verbose: func(epoch, j int, loss float64, m trainer.Model) {
			if j%(train.N/10/10) != 0 {
				return
			}

			acc := trainer.Accuracy(m.Predict(x[:100]), t[:100])
			acct := trainer.Accuracy(m.Predict(xt[:100]), tt[:100])

			fmt.Printf("loss=%.04f, train_acc=%.04f, test_acc=%.04f\n", loss, acc, acct)
		},
	}, s)

	// Output:
	// *model.MLP
	//  0: *layer.Affine: W(784, 50), B(1, 50): 39250
	//  1: *layer.BatchNorm: G(1, 50), B(1, 50): 100
	//  2: *layer.ReLU
	//  3: *layer.Affine: W(50, 10), B(1, 10): 510
	//  4: *layer.SoftmaxWithLoss
	//
	// loss=2.3129, train_acc=0.2700, test_acc=0.1900
	// loss=1.7145, train_acc=0.7200, test_acc=0.5000
	// loss=0.9382, train_acc=0.8600, test_acc=0.6500
	// loss=0.8891, train_acc=0.9300, test_acc=0.6400
	// loss=0.5807, train_acc=0.9300, test_acc=0.7000
	// loss=0.4917, train_acc=0.9000, test_acc=0.7000
	// loss=0.2948, train_acc=0.9900, test_acc=0.7100
	// loss=0.2014, train_acc=1.0000, test_acc=0.7100
	// loss=0.1978, train_acc=1.0000, test_acc=0.7000
	// loss=0.0379, train_acc=1.0000, test_acc=0.6700

}

func Example_simpleNet() {
	// https://github.com/oreilly-japan/deep-learning-from-scratch/wiki/errata#%E7%AC%AC7%E5%88%B7%E3%81%BE%E3%81%A7

	// weight
	W := matrix.New(
		[]float64{0.47355232, 0.99773930, 0.84668094},
		[]float64{0.85557411, 0.03563661, 0.69422093},
	)

	// data
	x := matrix.New([]float64{0.6, 0.9})
	t := []float64{0, 0, 1}

	// predict
	p := matrix.Dot(x, W)
	y := activation.Softmax(p[0])
	e := loss.CrossEntropyError(y, t)

	fmt.Println(p)
	fmt.Println(e)

	// gradient
	fW := func(w ...float64) float64 {
		p := matrix.Dot(x, W)
		y := activation.Softmax(p[0])
		e := loss.CrossEntropyError(y, t)
		return e
	}

	grad := func(f func(x ...float64) float64, x matrix.Matrix) matrix.Matrix {
		out := make(matrix.Matrix, 0)
		for _, r := range x {
			out = append(out, numerical.Gradient(f, r))
		}

		return out
	}

	dW := grad(fW, W)
	for _, r := range dW {
		fmt.Println(r)
	}

	// Output:
	// [[1.054148091 0.630716529 1.132807401]]
	// 0.9280682857864075
	// [0.2192475712392561 0.14356242984070455 -0.3628100010055757]
	// [0.3288713569016277 0.21534364482433954 -0.5442150014750569]

}

func Example_neuralNet() {
	// weight
	W1 := matrix.New([]float64{0.1, 0.3, 0.5}, []float64{0.2, 0.4, 0.6})
	B1 := matrix.New([]float64{0.1, 0.2, 0.3})
	W2 := matrix.New([]float64{0.1, 0.4}, []float64{0.2, 0.5}, []float64{0.3, 0.6})
	B2 := matrix.New([]float64{0.1, 0.2})
	W3 := matrix.New([]float64{0.1, 0.3}, []float64{0.2, 0.4})
	B3 := matrix.New([]float64{0.1, 0.2})

	// data
	x := matrix.New([]float64{1.0, 0.5})

	// forward
	A1 := matrix.Dot(x, W1).Add(B1)
	Z1 := matrix.F(A1, activation.Sigmoid)
	A2 := matrix.Dot(Z1, W2).Add(B2)
	Z2 := matrix.F(A2, activation.Sigmoid)
	A3 := matrix.Dot(Z2, W3).Add(B3)
	y := matrix.F(A3, activation.Identity)

	// print
	fmt.Println(A1)
	fmt.Println(Z1)

	fmt.Println(A2)
	fmt.Println(Z2)

	fmt.Println(A3)
	fmt.Println(y)

	// Output:
	// [[0.30000000000000004 0.7 1.1]]
	// [[0.574442516811659 0.6681877721681662 0.7502601055951177]]
	// [[0.5161598377933344 1.2140269561658172]]
	// [[0.6262493703990729 0.7710106968556123]]
	// [[0.3168270764110298 0.6962790898619668]]
	// [[0.3168270764110298 0.6962790898619668]]

}

func Example_perceptron() {
	f := func(x, w []float64, b float64) int {
		var sum float64
		for i := range x {
			sum = sum + x[i]*w[i]
		}

		v := sum + b
		if v <= 0 {
			return 0
		}

		return 1
	}

	AND := func(x []float64) int { return f(x, []float64{0.5, 0.5}, -0.7) }
	NAND := func(x []float64) int { return f(x, []float64{-0.5, -0.5}, 0.7) }
	OR := func(x []float64) int { return f(x, []float64{0.5, 0.5}, -0.2) }
	XOR := func(x []float64) int { return AND([]float64{float64(NAND(x)), float64(OR(x))}) }

	fmt.Println("AND")
	fmt.Println(AND([]float64{0, 0}))
	fmt.Println(AND([]float64{1, 0}))
	fmt.Println(AND([]float64{0, 1}))
	fmt.Println(AND([]float64{1, 1}))
	fmt.Println()

	fmt.Println("NAND")
	fmt.Println(NAND([]float64{0, 0}))
	fmt.Println(NAND([]float64{1, 0}))
	fmt.Println(NAND([]float64{0, 1}))
	fmt.Println(NAND([]float64{1, 1}))
	fmt.Println()

	fmt.Println("OR")
	fmt.Println(OR([]float64{0, 0}))
	fmt.Println(OR([]float64{1, 0}))
	fmt.Println(OR([]float64{0, 1}))
	fmt.Println(OR([]float64{1, 1}))
	fmt.Println()

	fmt.Println("XOR")
	fmt.Println(XOR([]float64{0, 0}))
	fmt.Println(XOR([]float64{1, 0}))
	fmt.Println(XOR([]float64{0, 1}))
	fmt.Println(XOR([]float64{1, 1}))
	fmt.Println()

	// Output:
	// AND
	// 0
	// 0
	// 0
	// 1
	//
	// NAND
	// 1
	// 1
	// 1
	// 0
	//
	// OR
	// 0
	// 1
	// 1
	// 1
	//
	// XOR
	// 0
	// 1
	// 1
	// 0
	//

}
