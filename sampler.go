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

type RandomSearchSamplerOption func(sampler *RandomSearchSampler)

func RandomSearchSamplerOptionSeed(seed int64) RandomSearchSamplerOption {
	return func(sampler *RandomSearchSampler) {
		sampler.rng = rand.New(rand.NewSource(seed))
	}
}

func NewRandomSearchSampler(opts ...RandomSearchSamplerOption) *RandomSearchSampler {
	s := &RandomSearchSampler{
		rng: rand.New(rand.NewSource(0)),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

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
		return s.rng.Float64()*(d.Max-d.Min) + d.Min, nil
	default:
		return 0.0, errors.New("undefined distribution")
	}
}
