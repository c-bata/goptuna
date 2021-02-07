package sobol

import (
	"math"
	"sort"

	"github.com/c-bata/goptuna"
)

var _ goptuna.RelativeSampler = &Sampler{}

// Sampler for quasi-Monte Carlo Sampling based on Sobol sequence.
type Sampler struct {
	engine *Engine
}

// SampleRelative samples multiple dimensional parameters in a given search space.
func (s *Sampler) SampleRelative(study *goptuna.Study, trial goptuna.FrozenTrial, searchSpace map[string]interface{}) (map[string]float64, error) {
	dim := len(searchSpace)
	if s.engine == nil {
		s.engine = NewEngine(uint32(dim))
	} else {
		// Detect dynamic search space.
		if s.engine.dim != uint32(dim) {
			return nil, nil
		}
	}
	points := s.engine.Draw()

	orderedKeys := make([]string, 0, len(searchSpace))
	for name := range searchSpace {
		orderedKeys = append(orderedKeys, name)
	}
	sort.Strings(orderedKeys)
	params := make(map[string]float64, len(orderedKeys))
	for i, name := range orderedKeys {
		switch d := searchSpace[name].(type) {
		case goptuna.UniformDistribution:
			params[name] = points[i]*(d.High-d.Low) + d.Low
		case goptuna.DiscreteUniformDistribution:
			q := d.Q
			r := d.High - d.Low
			// [low, high] is shifted to [0, r] to align sampled values at regular intervals.
			low := 0 - 0.5*q
			high := r + 0.5*q
			x := points[i]*(high-low) + low
			v := math.Round(x/q)*q + d.Low
			params[name] = math.Min(math.Max(v, d.Low), d.High)
		case goptuna.LogUniformDistribution:
			logLow := math.Log(d.Low)
			logHigh := math.Log(d.High)
			params[name] = math.Exp(points[i]*(logHigh-logLow) + logLow)
		case goptuna.IntUniformDistribution:
			params[name] = math.Floor(points[i]*float64(d.High-d.Low)) + float64(d.Low)
		case goptuna.StepIntUniformDistribution:
			r := (d.High - d.Low) / d.Step
			v := (int(math.Floor(points[i]*float64(r))) * d.Step) + d.Low
			params[name] = float64(v)
		case goptuna.CategoricalDistribution:
			params[name] = math.Floor(points[i] * float64(len(d.Choices)))
		default:
			return nil, goptuna.ErrUnknownDistribution
		}
	}
	return params, nil
}

// NewSampler returns the Sobol sampler.
func NewSampler() *Sampler {
	sampler := &Sampler{
		engine: nil,
	}
	return sampler
}
