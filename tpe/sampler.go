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

var _ goptuna.Sampler = &TPESampler{}

type TPESampler struct {
	ConsiderPrior     bool
	PriorWeights      float64
	ConsiderMagicClip bool
	ConsiderEndpoints bool
	NStartupTrials    int
	NEICandidates     int
	Gamma             FuncGamma
	Weights           FuncWeights
	random_sampler    *goptuna.RandomSearchSampler
}

func NewTPESampler() *TPESampler {
	sampler := &TPESampler{
		ConsiderPrior:     true,
		PriorWeights:      1.0,
		ConsiderMagicClip: true,
		ConsiderEndpoints: false,
		NStartupTrials:    10,
		NEICandidates:     24,
		Gamma:             DefaultGamma,
		Weights:           DefaultWeights,
		random_sampler:    goptuna.NewRandomSearchSampler(),
	}
	return sampler
}

func (s *TPESampler) Sample(*goptuna.Study, goptuna.FrozenTrial, string, interface{}) (float64, error) {
	panic("implement me")
}

func getObservationPairs(study goptuna.Study, paramName string) ([][]float64, error) {
	var sign float64 = 1
	if study.Direction() == goptuna.StudyDirectionMaximize {
		sign = -1
	}

	pairs := make([][]float64, 0, 8)
	trials, err := study.GetTrials()
	if err != nil {
		return nil, err
	}

	for _, trial := range trials {
		ir, ok := trial.ParamsInIR[paramName]
		if !ok {
			continue
		}

		var first, second, third float64
		first = ir
		if trial.State == goptuna.TrialStateComplete {
			second = math.Inf(-1)
			third = sign * trial.Value
		} else if trial.State == goptuna.TrialStatePruned {
			panic("still be unreachable")
		} else {
			continue
		}
		pairs = append(pairs, []float64{first, second, third})
	}
	return pairs, nil
}
