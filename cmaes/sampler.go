package cmaes

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/c-bata/goptuna"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

var _ goptuna.RelativeSampler = &Sampler{}

type popType string

const (
	popTypeSmall = popType("small")
	popTypeLarge = popType("large")
)

// Sampler returns the next search points by using CMA-ES.
type Sampler struct {
	x0               map[string]float64
	sigma0           float64
	rng              *rand.Rand
	nStartUpTrials   int
	optimizerOptions []OptimizerOption
	optimizer        *Optimizer
	optimizerID      string
	// Variables related to IPOP-CMA-ES and BIPOP-CMA-ES
	restartStrategy string
	incPopSize      int
	nRestarts       int // A small restart doesn't count in the nRestarts in BI-POP
	nSmallEval      int
	nLargeEval      int
	popsize0        int
	poptype         popType
}

// SampleRelative samples multiple dimensional parameters in a given search space.
func (s *Sampler) SampleRelative(
	study *goptuna.Study,
	trial goptuna.FrozenTrial,
	searchSpace map[string]interface{},
) (map[string]float64, error) {
	if searchSpace == nil || len(searchSpace) == 0 {
		return nil, nil
	}

	searchSpace = supportedSearchSpace(searchSpace)
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
	if err != nil && err != goptuna.ErrTrialsPartiallyDeleted {
		return nil, err
	}
	completed := make([]goptuna.FrozenTrial, 0, len(trials))
	for i := range trials {
		if trials[i].State == goptuna.TrialStateComplete {
			completed = append(completed, trials[i])
		}
	}
	if len(completed) < s.nStartUpTrials {
		// If catch ErrTrialsPartiallyDeleted, nStartUpTrials should be smaller than len(completed).
		study.GetLogger().Error("Your BlackHoleStorage buffer is too small.",
			fmt.Sprintf("nStartUpTrials:%d", s.nStartUpTrials))
		return nil, err
	}
	if err == goptuna.ErrTrialsPartiallyDeleted && s.optimizer != nil &&
		len(completed) < s.optimizer.PopulationSize() {
		// If catch ErrTrialsPartiallyDeleted, population size should be smaller than len(completed).
		study.GetLogger().Error("Your BlackHoleStorage buffer is too small.",
			fmt.Sprintf("popsize:%d", s.optimizer.PopulationSize()))
		return nil, err
	}
	err = nil

	if s.optimizer == nil {
		s.optimizer, err = s.initOptimizer(searchSpace, orderedKeys)
		if err != nil {
			return nil, err
		}
		s.optimizerID = fmt.Sprintf("%016d", s.rng.Int())
	}

	if s.optimizer.dim != len(orderedKeys) {
		study.GetLogger().Warn("cmaes.Sampler does not support dynamic search space." +
			" All parameters will be sampled by normal sampler.")
		return nil, nil
	}

	solutions := make([]*Solution, 0, s.optimizer.PopulationSize())
	for i := range completed {
		generationID, ok := completed[i].SystemAttrs["goptuna:cmaes:generationId"]
		if !ok || generationID != fmt.Sprintf("%s-%d", s.optimizerID, s.optimizer.Generation()) {
			continue
		}
		x := make([]float64, len(orderedKeys))
		for j := 0; j < len(orderedKeys); j++ {
			p, ok := completed[i].InternalParams[orderedKeys[j]]
			if !ok {
				return nil, errors.New("invalid internal params")
			}
			x[j] = toCMAParam(searchSpace[orderedKeys[j]], p)
		}
		solutions = append(solutions, &Solution{
			Params: x,
			Value:  completed[i].Value,
		})

		if len(solutions) == s.optimizer.PopulationSize() {
			err = s.optimizer.Tell(solutions)
			if err != nil {
				return nil, err
			}

			if s.optimizer.ShouldStop() && s.restartStrategy != "" {
				popsize := s.nextPopsize()
				s.optimizer, err = s.initOptimizer(searchSpace, orderedKeys,
					OptimizerOptionPopulationSize(popsize))
				if err != nil {
					return nil, err
				}
			}
			break
		}
	}

	nextParams, err := s.optimizer.Ask()
	if err != nil {
		return nil, err
	}

	err = study.Storage.SetTrialSystemAttr(
		trial.ID,
		"goptuna:cmaes:generationId",
		fmt.Sprintf("%s-%d", s.optimizerID, s.optimizer.Generation()))
	if err != nil {
		return nil, err
	}

	params := make(map[string]float64, len(orderedKeys))
	for i := range orderedKeys {
		param := nextParams[i]
		params[orderedKeys[i]] = toGoptunaInternalParam(searchSpace[orderedKeys[i]], param)
	}
	return params, nil
}

func (s *Sampler) nextPopsize() (popsize int) {
	if s.restartStrategy == restartStrategyIPOP {
		// I-POP-CMA-ES
		s.nRestarts++
		return s.optimizer.PopulationSize() * s.incPopSize
	}

	// BI-POP-CMA-ES
	if s.popsize0 == 0 {
		s.popsize0 = s.optimizer.PopulationSize()
	}

	nEval := s.optimizer.PopulationSize() * s.optimizer.Generation()
	if s.poptype == popTypeSmall {
		s.nSmallEval += nEval
	} else { // large
		s.nLargeEval += nEval
	}

	if s.nSmallEval < s.nLargeEval {
		s.poptype = popTypeSmall
		popsizeMultiplier := math.Pow(float64(s.incPopSize), float64(s.nRestarts))
		r := math.Pow(s.rng.Float64(), 2)
		return int(math.Floor(float64(s.popsize0) * math.Pow(popsizeMultiplier, r)))
	}

	s.poptype = popTypeLarge
	s.nRestarts++
	return s.popsize0 * int(math.Pow(float64(s.incPopSize), float64(s.nRestarts)))
}

func (s *Sampler) initOptimizer(
	searchSpace map[string]interface{},
	orderedKeys []string,
	additionalOpts ...OptimizerOption,
) (*Optimizer, error) {
	x0, sigma0, err := s.initialParam(searchSpace)
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

	options := make([]OptimizerOption, 0, 2+len(s.optimizerOptions)+len(additionalOpts))
	options = append(options, OptimizerOptionBounds(bounds))
	options = append(options, OptimizerOptionSeed(s.rng.Int63()))
	for _, opt := range s.optimizerOptions {
		options = append(options, opt)
	}
	for _, opt := range additionalOpts {
		options = append(options, opt)
	}
	return NewOptimizer(mean, sigma0, options...)
}

// NewSampler returns the TPE sampler.
func NewSampler(opts ...SamplerOption) *Sampler {
	sampler := &Sampler{
		rng:            rand.New(rand.NewSource(0)),
		nStartUpTrials: 0,

		// Initial run is with "normal" population size; it is
		// the large population before first doubling, but its
		// budget accounting is the same as in case of small
		// population.
		poptype: popTypeSmall,
	}

	for _, opt := range opts {
		opt(sampler)
	}
	return sampler
}

func supportedSearchSpace(searchSpace map[string]interface{}) map[string]interface{} {
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
		case goptuna.StepIntUniformDistribution:
			normalized[name] = searchSpace[name]
		}
	}
	return normalized
}

func toCMAParam(distribution interface{}, goptunaParam float64) float64 {
	switch distribution.(type) {
	case goptuna.LogUniformDistribution:
		return math.Log(goptunaParam)
	}
	return goptunaParam
}

func toGoptunaInternalParam(distribution interface{}, cmaParam float64) float64 {
	switch distribution.(type) {
	case goptuna.LogUniformDistribution:
		return math.Exp(cmaParam)
	}
	return cmaParam
}

func (s *Sampler) initialParam(searchSpace map[string]interface{}) (map[string]float64, float64, error) {
	x0 := make(map[string]float64, len(searchSpace))
	sigma0 := make([]float64, 0, len(searchSpace))

	for name := range searchSpace {
		switch d := searchSpace[name].(type) {
		case goptuna.UniformDistribution:
			if s.nRestarts > 0 {
				x0[name] = d.Low + s.rng.Float64()*(d.High-d.Low)
			} else {
				x0[name] = (d.High + d.Low) / 2
			}
			sigma0 = append(sigma0, (d.High-d.Low)/6)
		case goptuna.DiscreteUniformDistribution:
			if s.nRestarts > 0 {
				x0[name] = d.Low + s.rng.Float64()*(d.High-d.Low)
			} else {
				x0[name] = (d.High + d.Low) / 2
			}
			sigma0 = append(sigma0, (d.High-d.Low)/6)
		case goptuna.LogUniformDistribution:
			high := math.Log(d.High)
			low := math.Log(d.Low)
			if s.nRestarts > 0 {
				x0[name] = low + s.rng.Float64()*(high-low)
			} else {
				x0[name] = (high + low) / 2
			}
			sigma0 = append(sigma0, (high-low)/6)
		case goptuna.IntUniformDistribution:
			if s.nRestarts > 0 {
				x0[name] = float64(d.Low + s.rng.Intn(d.High-d.Low))
			} else {
				x0[name] = float64(d.High+d.Low) / 2
			}
			sigma0 = append(sigma0, float64(d.High-d.Low)/6)
		case goptuna.StepIntUniformDistribution:
			if s.nRestarts > 0 {
				x0[name] = float64(d.Low + s.rng.Intn(d.High-d.Low))
			} else {
				x0[name] = float64(d.High+d.Low) / 2
			}
			sigma0 = append(sigma0, float64(d.High-d.Low)/6)
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
			bounds.Set(i, 0, math.Log(d.Low))
			bounds.Set(i, 1, math.Log(d.High))
		case goptuna.IntUniformDistribution:
			bounds.Set(i, 0, float64(d.Low))
			bounds.Set(i, 1, float64(d.High))
		case goptuna.StepIntUniformDistribution:
			bounds.Set(i, 0, float64(d.Low))
			bounds.Set(i, 1, float64(d.High))
		default:
			panic("keys and search_space do not match")
		}
	}
	return bounds
}
