package tpe

import (
	"math"

	"github.com/c-bata/goptuna"
)

const EPS = 1e-12

type FuncGamma func(int) int

type FuncWeights func(int) []float64

func DefaultGamma(x int) int {
	a := int(math.Ceil(0.1 * float64(x)))
	if a > 25 {
		return a
	}
	return 25
}

func HyperoptDefaultGamma(x int) int {
	a := int(math.Ceil(0.25 * float64(x)))
	if a > 25 {
		return a
	}
	return 25
}

func DefaultWeights(x int) []float64 {
	if x == 0 {
		return []float64{}
	} else if x < 25 {
		return ones(25)
	} else {
		ramp := linspace(1.0/float64(x), 1.0, x-25, true)
		flat := ones(25)
		return append(ramp, flat...)
	}
}

type TPESampler struct{}

func (s *TPESampler) Sample(*goptuna.Study, goptuna.FrozenTrial, string, interface{}) (float64, error) {
	panic("implement me")
}
