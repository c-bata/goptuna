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
	NStartupTrials        int
	NEICandidates         int
	Gamma                 FuncGamma
	ParzenEstimatorParams ParzenEstimatorParams

	random_sampler *goptuna.RandomSearchSampler
}

func NewTPESampler() *TPESampler {
	sampler := &TPESampler{
		NStartupTrials: 10,
		NEICandidates:  24,
		Gamma:          DefaultGamma,
		ParzenEstimatorParams: ParzenEstimatorParams{
			ConsiderPrior:     true,
			PriorWeight:       1.0,
			ConsiderMagicClip: true,
			ConsiderEndpoints: false,
			Weights:           DefaultWeights,
		},
		random_sampler: goptuna.NewRandomSearchSampler(),
	}
	return sampler
}

func genKeepIdxs(
	lossIdxs []int, lossAscending []int, n int, below bool) []int {
	var l []int
	if below {
		l = lossAscending[:n]
	} else {
		l = lossAscending[n:]
	}

	lossIdxSet := make([]int, 0)
	for _, index := range l {
		item := lossIdxs[index]

		found := false
		for _, j := range lossIdxSet {
			if lossIdxSet[j] == item {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		lossIdxSet = append(lossIdxSet, item)
	}
	return lossIdxSet
}

func genBelowOrAbove(
	keepIdxs []int,
	configIdxs []int,
	configVals []float64,
) []float64 {
	size := len(configIdxs)
	if size > len(configVals) {
		size = len(configVals)
	}
	results := make([]float64, 0, size)

	isExist := func(index int, configIdxs []int) bool {
		for _, idx := range configIdxs {
			if index == idx {
				return true
			}
		}
		return false
	}

	for i := 0; i < size; i++ {
		index := configIdxs[i]
		value := configVals[i]

		if isExist(index, keepIdxs) {
			results = append(results, value)
		}
	}
	return results
}

func (s *TPESampler) splitObservationPairs(
	configIdxs []int,
	configVals []float64,
	lossIdxs []int,
	lossVals [][2]float64,
) ([]float64, []float64) {
	nbelow := s.Gamma(len(configVals))
	lossAscending := ArgSort2DFloat64(lossVals)

	keepIdxs := genKeepIdxs(lossIdxs, lossAscending, nbelow, true)
	below := genBelowOrAbove(keepIdxs, configIdxs, configVals)

	keepIdxs = genKeepIdxs(lossIdxs, lossAscending, nbelow, false)
	above := genBelowOrAbove(keepIdxs, configIdxs, configVals)

	return below, above
}

func (s *TPESampler) sampleFromGMM(parzenEstimator *ParzenEstimator, low, high float64, size int) []float64 {
	panic("not implemented yet")
}

func (s *TPESampler) gmmLogPDF(samples []float64, parzenEstimator *ParzenEstimator, low, high float64) []float64 {
	panic("not implemented yet")
}

func (s *TPESampler) compare(samples []float64, logL []float64, logG []float64) []float64 {
	panic("not implemented yet")
}

func (s *TPESampler) sampleNumerical(low, high float64, below, above []float64) float64 {
	size := s.NEICandidates
	parzenEstimatorBelow := NewParzenEstimator(below, low, high, s.ParzenEstimatorParams)
	sampleBelow := s.sampleFromGMM(parzenEstimatorBelow, low, high, size)
	logLikelihoodsBelow := s.gmmLogPDF(sampleBelow, parzenEstimatorBelow, low, high)

	parzenEstimatorAbove := NewParzenEstimator(above, low, high, s.ParzenEstimatorParams)
	sampleAbove := s.sampleFromGMM(parzenEstimatorAbove, low, high, size)
	logLikelihoodsAbove := s.gmmLogPDF(sampleAbove, parzenEstimatorAbove, low, high)

	return s.compare(sampleBelow, logLikelihoodsBelow, logLikelihoodsAbove)[0]
}

func (s *TPESampler) sampleUniform(distribution goptuna.UniformDistribution, below, above []float64) float64 {
	low := distribution.Min
	high := distribution.Max
	return s.sampleNumerical(low, high, below, above)
}

func (s *TPESampler) Sample(
	study *goptuna.Study,
	trial goptuna.FrozenTrial,
	paramName string,
	paramDistribution interface{},
) (float64, error) {
	observationPairs, err := getObservationPairs(study, paramName)
	if err != nil {
		return 0, err
	}
	n := len(observationPairs)

	if n < s.NStartupTrials {
		return s.random_sampler.Sample(study, trial, paramName, paramDistribution)
	}

	configIdxs := make([]int, n)
	for i := 0; i < n; i++ {
		configIdxs[i] = i
	}
	configVals := make([]float64, n)
	for i := 0; i < n; i++ {
		configVals[i] = observationPairs[i][0]
	}
	lossIdxs := make([]int, n)
	for i := 0; i < n; i++ {
		lossIdxs[i] = i
	}
	lossVals := make([][2]float64, n)
	for i := 0; i < n; i++ {
		lossVals[i] = [2]float64{observationPairs[i][1], observationPairs[i][2]}
	}
	belowParamValues, adoveParamValues := s.splitObservationPairs(
		configIdxs, configVals, lossIdxs, lossVals)

	switch d := paramDistribution.(type) {
	case goptuna.UniformDistribution:
		return s.sampleUniform(d, belowParamValues, adoveParamValues), nil
	}
	return 0, goptuna.ErrUnexpectedDistribution
}

func getObservationPairs(study *goptuna.Study, paramName string) ([][3]float64, error) {
	var sign float64 = 1
	if study.Direction() == goptuna.StudyDirectionMaximize {
		sign = -1
	}

	pairs := make([][3]float64, 0, 8)
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
		pairs = append(pairs, [3]float64{first, second, third})
	}
	return pairs, nil
}
