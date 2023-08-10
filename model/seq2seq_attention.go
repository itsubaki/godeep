package model

import (
	"math/rand"
	"time"

	"github.com/itsubaki/neu/layer"
	"github.com/itsubaki/neu/math/matrix"
)

type AttentionSeq2Seq struct {
	Encoder *AttentionEncoder
	Decoder *AttentionDecoder
	Softmax *layer.TimeSoftmaxWithLoss
	Source  rand.Source
}

func NewAttentionSeq2Seq(c *RNNLMConfig, s ...rand.Source) *AttentionSeq2Seq {
	if len(s) == 0 {
		s = append(s, rand.NewSource(time.Now().UnixNano()))
	}

	return &AttentionSeq2Seq{
		Encoder: NewAttentionEncoder(c, s[0]),
		Decoder: NewAttentionDecoder(c, s[0]),
		Softmax: &layer.TimeSoftmaxWithLoss{},
		Source:  s[0],
	}
}

func (m *AttentionSeq2Seq) Forward(xs, ts []matrix.Matrix) float64 {
	// xs:  ['5', '7', '+', '5', ' ', ' ', ' ']
	// dxs: ['_', '6', '2', ' ']
	// dts: ['6', '2', ' ', ' ']
	dxs, dts := ts[:len(ts)-1], ts[1:]    // dxs(4, 128, 1), dts(4, 128, 1)
	h := m.Encoder.Forward(xs)            // h(128, 128)
	score := m.Decoder.Forward(dxs, h)    // score(4, 128, 13)
	loss := m.Softmax.Forward(score, dts) // (1, 1, 1)
	return loss[0][0][0]
}

func (m *AttentionSeq2Seq) Backward() {
	dout := []matrix.Matrix{{{1}}}     // (1, 1, 1)
	dscore := m.Softmax.Backward(dout) // (4, 128, 13)
	dh := m.Decoder.Backward(dscore)   // (128, 128)
	m.Encoder.Backward(dh)             // (0, 0, 0)
}

func (m *AttentionSeq2Seq) Generate(xs []matrix.Matrix, startID, length int) []int {
	h := m.Encoder.Forward(xs)                        // xs(7, 1, 1), h(1, 128)
	sampeld := m.Decoder.Generate(h, startID, length) //
	return sampeld
}

func (m *AttentionSeq2Seq) Layers() []TimeLayer {
	layers := make([]TimeLayer, 0)
	layers = append(layers, m.Encoder.Layers()...)
	layers = append(layers, m.Decoder.Layers()...)
	layers = append(layers, m.Softmax)
	return layers
}

func (m *AttentionSeq2Seq) Params() [][]matrix.Matrix {
	return [][]matrix.Matrix{
		m.Encoder.Params(),
		m.Decoder.Params(),
	}
}

func (m *AttentionSeq2Seq) Grads() [][]matrix.Matrix {
	return [][]matrix.Matrix{
		m.Encoder.Grads(),
		m.Decoder.Grads(),
	}
}

func (m *AttentionSeq2Seq) SetParams(p [][]matrix.Matrix) {
	m.Encoder.SetParams(p[0]...)
	m.Decoder.SetParams(p[1]...)
}
