package vector

func Add(v, w []float64) []float64 {
	out := make([]float64, len(v))
	for i := range v {
		out[i] = v[i] + w[i]
	}

	return out
}

func Int(v []float64) []int {
	out := make([]int, len(v))
	for i, e := range v {
		out[i] = int(e)
	}

	return out
}
