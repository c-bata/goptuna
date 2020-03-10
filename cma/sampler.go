package cma

import (
	"math/rand"

	"github.com/c-bata/goptuna"
)

var _ goptuna.RelativeSampler = &Sampler{}

// Sampler returns the next search points by using CMA-ES.
type Sampler struct {
	rng            *rand.Rand
	nStartUpTrials int
}

func (s *Sampler) SampleRelative(
	study *goptuna.Study,
	trial goptuna.FrozenTrial,
	searchSpace map[string]interface{},
) (map[string]float64, error) {
	if searchSpace == nil || len(searchSpace) == 0 {
		return nil, nil
	}

	if len(searchSpace) == 1 {
		// CMA-ES does not support optimization of 1-D search space.
		return nil, goptuna.ErrUnsupportedSearchSpace
	}

	trials, err := study.GetTrials()
	if err != nil {
		return nil, err
	}
	completed := make([]goptuna.FrozenTrial, 0, len(trials))
	for i := range trials {
		if trials[i].State == goptuna.TrialStateComplete {
			completed = append(completed, trials[i])
		}
	}

	if len(completed) < s.nStartUpTrials {
		return nil, err
	}

	params := make(map[string]float64, len(searchSpace))
	return params, nil
}

func (s *Sampler) initializeMu(searchSpace map[string]interface{}) []float64 {
	distributions := make([]interface{}, 0, len(searchSpace))
	for key := range searchSpace {
		switch _ := searchSpace[key].(type) {
		case goptuna.UniformDistribution:
		case goptuna.LogUniformDistribution:
		case goptuna.IntUniformDistribution:
		case goptuna.DiscreteUniformDistribution:
			distributions = append(distributions, searchSpace[key])
		case goptuna.CategoricalDistribution:
			continue
		}
	}
	dim := len(distributions)
	_ = dim
	panic("wip")
	return nil
}

// NewSampler returns the TPE sampler.
func NewSampler(opts ...SamplerOption) *Sampler {
	sampler := &Sampler{
		rng:            rand.New(rand.NewSource(0)),
		nStartUpTrials: 0,
	}

	for _, opt := range opts {
		opt(sampler)
	}
	return sampler
}
