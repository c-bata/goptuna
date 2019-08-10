package random_test

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/c-bata/goptuna/internal/random"
)

var randomWeightedSelect = func(weights []int, totalWeight int) (int, error) {
	// https://medium.com/@peterkellyonline/weighted-random-selection-3ff222917eb6
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(totalWeight)
	for i, g := range weights {
		r -= g
		if r <= 0 {
			return i, nil
		}
	}
	return 0, errors.New("no item selected")
}

var oldArgMaxApproxMultinomial = func(pvals []float64, precision float64) (int, error) {
	tw := 0
	weights := make([]int, len(pvals))
	for i := range weights {
		w := int(pvals[i] / precision)
		tw += w
		weights[i] = w
	}
	return randomWeightedSelect(weights, tw)
}

func BenchmarkArgMaxMultinomial(b *testing.B) {
	pvals := make([]float64, 100)
	for i := range pvals {
		pvals[i] = 0.01
	}
	for n := 0; n < b.N; n++ {
		_, _ = random.ArgMaxMultinomial(pvals)
	}
}

func BenchmarkArgMaxMultinomialOld(b *testing.B) {
	pvals := make([]float64, 100)
	for i := range pvals {
		pvals[i] = 0.01
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = oldArgMaxApproxMultinomial(pvals, 1e-5)
	}
}
