package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/itsubaki/neu/dataset/ptb"
	"github.com/itsubaki/neu/math/vector"
	"github.com/itsubaki/neu/model"
	"github.com/itsubaki/neu/optimizer"
	"github.com/itsubaki/neu/optimizer/hook"
	"github.com/itsubaki/neu/trainer"
	"github.com/itsubaki/neu/weight"
)

func main() {
	// flags
	var dir string
	var epochs, corpusSize int
	flag.StringVar(&dir, "dir", "./testdata", "")
	flag.IntVar(&epochs, "epochs", 100, "")
	flag.IntVar(&corpusSize, "corpus-size", 1000, "")
	flag.Parse()

	// data
	train := ptb.Must(ptb.Load(dir, ptb.TrainTxt))
	corpus := train.Corpus[:corpusSize]

	// model
	m := model.NewLSTMLM(&model.LSTMLMConfig{
		RNNLMConfig: model.RNNLMConfig{
			VocabSize:   vector.Max(corpus) + 1,
			WordVecSize: 100,
			HiddenSize:  100,
			WeightInit:  weight.Xavier,
		},
		DropoutRatio: 0.5,
	})

	// summary
	fmt.Println(m.Summary()[0])
	for i, s := range m.Summary()[1:] {
		fmt.Printf("%2d: %v\n", i, s)
	}
	fmt.Println()

	// training
	tr := trainer.NewRNNLM(m, &optimizer.SGD{
		LearningRate: 20,
		Hooks: []optimizer.Hook{
			hook.GradsClipping(0.25),
		},
	})

	now := time.Now()
	tr.Fit(&trainer.RNNLMInput{
		Train:      corpus[:len(corpus)-1],
		TrainLabel: corpus[1:],
		Epochs:     epochs,
		BatchSize:  20,
		TimeSize:   35,
		Verbose: func(epoch, j int, perplexity float64, m trainer.RNNLM) {
			fmt.Printf("%2d, %2d: ppl=%.04f\n", epoch, j, perplexity)
		},
	})

	fmt.Printf("elapsed=%v\n", time.Since(now))
}
