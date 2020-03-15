package cma

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"

	"github.com/c-bata/goptuna"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

var _ goptuna.RelativeSampler = &Sampler{}

// Sampler returns the next search points by using CMA-ES.
type Sampler struct {
	x0             map[string]float64
	sigma0         float64
	rng            *rand.Rand
	nStartUpTrials int
	optimizer      *Optimizer
}

func (s *Sampler) SampleRelative(
	study *goptuna.Study,
	trial goptuna.FrozenTrial,
	searchSpace map[string]interface{},
) (map[string]float64, error) {
	if searchSpace == nil || len(searchSpace) == 0 {
		return nil, nil
	}

	searchSpace = normalizeSearchSpace(searchSpace)
	if len(searchSpace) == 1 {
		// CMA-ES does not support two or more dimensional continuous search space.
		return nil, goptuna.ErrUnsupportedSearchSpace
	}
	orderedKeys := make([]string, 0, len(searchSpace))
	for name := range searchSpace {
		orderedKeys = append(orderedKeys, name)
	}
	sort.Strings(orderedKeys)

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
		return nil, nil
	}

	if s.optimizer == nil {
		s.optimizer, err = s.initOptimizer(searchSpace, orderedKeys)
		if err != nil {
			return nil, err
		}
	}

	if s.optimizer.dim != len(orderedKeys) {
		// TODO(c-bata): Use logger for warning.
		fmt.Println("This optimizer does not support dynamic search space.")
		return nil, nil
	}

	solutions := make([]*Solution, 0, s.optimizer.PopulationSize())
	for i := range completed {
		gstr, ok := completed[i].SystemAttrs["goptuna:cma:generation"]
		if !ok {
			continue
		}
		if g, err := strconv.Atoi(gstr); err != nil || g != s.optimizer.Generation() {
			continue
		}
		x := mat.NewVecDense(len(orderedKeys), nil)
		for i := 0; i < len(orderedKeys); i++ {
			p, ok := completed[i].InternalParams[orderedKeys[i]]
			if !ok {
				return nil, errors.New("invalid internal params")
			}
			x.SetVec(i, p)
		}
		solutions = append(solutions, &Solution{
			X:     x,
			Value: completed[i].Value,
		})

		if len(solutions) == s.optimizer.PopulationSize() {
			err = s.optimizer.Tell(solutions)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	nextParams, err := s.optimizer.Ask()
	if err != nil {
		return nil, err
	}

	err = study.Storage.SetTrialSystemAttr(trial.ID,
		"goptuna:cma:generation", strconv.Itoa(s.optimizer.Generation()))
	if err != nil {
		return nil, err
	}

	params := make(map[string]float64, len(orderedKeys))
	for i := range orderedKeys {
		params[orderedKeys[i]] = nextParams.AtVec(i)
	}
	return params, nil
}

func (s *Sampler) initOptimizer(
	searchSpace map[string]interface{},
	orderedKeys []string,
) (*Optimizer, error) {
	x0, sigma0, err := initialParam(searchSpace)
	if err != nil {
		return nil, err
	}
	if s.x0 != nil {
		x0 = s.x0
	}
	if s.sigma0 > 0 {
		sigma0 = s.sigma0
	}

	mean := make([]float64, len(orderedKeys))
	for i := range orderedKeys {
		mean0, ok := x0[orderedKeys[i]]
		if !ok {
			return nil, errors.New("keys and search_space do not match")
		}
		mean[i] = mean0
	}
	bounds := getSearchSpaceBounds(searchSpace, orderedKeys)
	return NewOptimizer(
		mean, sigma0,
		OptimizerOptionBounds(bounds),
		OptimizerOptionSeed(s.rng.Int63()),
	)
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

func normalizeSearchSpace(searchSpace map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{}, len(searchSpace))
	for name := range searchSpace {
		switch searchSpace[name].(type) {
		case goptuna.UniformDistribution:
			normalized[name] = searchSpace[name]
		case goptuna.DiscreteUniformDistribution:
			normalized[name] = searchSpace[name]
		case goptuna.LogUniformDistribution:
			normalized[name] = searchSpace[name]
		case goptuna.IntUniformDistribution:
			normalized[name] = searchSpace[name]
		}
	}
	return normalized
}

func initialParam(searchSpace map[string]interface{}) (map[string]float64, float64, error) {
	x0 := make(map[string]float64, len(searchSpace))
	sigma0 := make([]float64, 0, len(searchSpace))
	for name := range searchSpace {
		switch d := searchSpace[name].(type) {
		case goptuna.UniformDistribution:
			x0[name] = (d.High + d.Low) / 2
			sigma0 = append(sigma0, math.Abs(d.High-d.Low)/6)
		case goptuna.DiscreteUniformDistribution:
			x0[name] = (d.High + d.Low) / 2
			sigma0 = append(sigma0, math.Abs(d.High-d.Low)/6)
		case goptuna.LogUniformDistribution:
			x0[name] = (d.High + d.Low) / 2
			sigma0 = append(sigma0, math.Abs(d.High-d.Low)/6)
		case goptuna.IntUniformDistribution:
			x0[name] = float64(d.High+d.Low) / 2
			sigma0 = append(sigma0, math.Abs(float64(d.High-d.Low))/6)
		default:
			return nil, 0, goptuna.ErrUnknownDistribution
		}
	}
	return x0, floats.Min(sigma0), nil
}

func getSearchSpaceBounds(
	searchSpace map[string]interface{},
	orderedKeys []string,
) *mat.Dense {
	bounds := mat.NewDense(len(orderedKeys), 2, nil)
	for i, name := range orderedKeys {
		switch d := searchSpace[name].(type) {
		case goptuna.UniformDistribution:
			bounds.Set(i, 0, d.Low)
			bounds.Set(i, 1, d.High)
		case goptuna.DiscreteUniformDistribution:
			bounds.Set(i, 0, d.Low)
			bounds.Set(i, 1, d.High)
		case goptuna.LogUniformDistribution:
			bounds.Set(i, 0, d.Low)
			bounds.Set(i, 1, d.High)
		case goptuna.IntUniformDistribution:
			bounds.Set(i, 0, float64(d.Low))
			bounds.Set(i, 1, float64(d.High))
		default:
			panic("keys and search_space do not match")
		}
	}
	return bounds
}
