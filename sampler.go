package goptuna

import (
	"errors"
	"math/rand"
	"sync"
)

// Sampler returns the next search points
type Sampler interface {
	// Sample a parameter for a given distribution.
	Sample(*Study, FrozenTrial, string, interface{}) (float64, error)
}

var _ Sampler = &RandomSearchSampler{}

// RandomSearchSampler for random search
type RandomSearchSampler struct {
	rng *rand.Rand
	mu  sync.Mutex
}

// RandomSearchSamplerOption is a type of function to set change the option.
type RandomSearchSamplerOption func(sampler *RandomSearchSampler)

// RandomSearchSamplerOptionSeed sets seed number.
func RandomSearchSamplerOptionSeed(seed int64) RandomSearchSamplerOption {
	return func(sampler *RandomSearchSampler) {
		sampler.rng = rand.New(rand.NewSource(seed))
	}
}

// NewRandomSearchSampler implements random search algorithm.
func NewRandomSearchSampler(opts ...RandomSearchSamplerOption) *RandomSearchSampler {
	s := &RandomSearchSampler{
		rng: rand.New(rand.NewSource(0)),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Sample a parameter for a given distribution.
func (s *RandomSearchSampler) Sample(
	study *Study,
	trial FrozenTrial,
	paramName string,
	paramDistribution interface{},
) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch d := paramDistribution.(type) {
	case UniformDistribution:
		if d.Single() {
			return d.Low, nil
		}
		return s.rng.Float64()*(d.High-d.Low) + d.Low, nil
	case IntUniformDistribution:
		if d.Single() {
			return float64(d.Low), nil
		}
		return float64(s.rng.Intn(d.High-d.Low) + d.Low), nil
	case CategoricalDistribution:
		if d.Single() {
			return float64(0), nil
		}
		return float64(rand.Intn(len(d.Choices))), nil
	default:
		return 0.0, errors.New("undefined distribution")
	}
}
