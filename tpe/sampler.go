package tpe

import (
	"math"
	"math/rand"
	"sync"

	"github.com/c-bata/goptuna"
	"gonum.org/v1/gonum/floats"
)

const EPS = 1e-12

type FuncGamma func(int) int

type FuncWeights func(int) []float64

func DefaultGamma(x int) int {
	a := int(math.Ceil(0.1 * float64(x)))
	if a > 25 {
		return 25
	}
	return a
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
		return ones1d(x)
	} else {
		ramp := linspace(1.0/float64(x), 1.0, x-25, true)
		flat := ones1d(25)
		return append(ramp, flat...)
	}
}

var _ goptuna.Sampler = &Sampler{}

type Sampler struct {
	numberOfStartupTrials int
	numberOfEICandidates  int
	gamma                 FuncGamma
	params                ParzenEstimatorParams
	rng                   *rand.Rand
	randomSampler         *goptuna.RandomSearchSampler
	mu                    sync.Mutex
}

func NewSampler(opts ...SamplerOption) *Sampler {
	sampler := &Sampler{
		numberOfStartupTrials: 10,
		numberOfEICandidates:  24,
		gamma:                 DefaultGamma,
		params: ParzenEstimatorParams{
			ConsiderPrior:     true,
			PriorWeight:       1.0,
			ConsiderMagicClip: true,
			ConsiderEndpoints: false,
			Weights:           DefaultWeights,
		},
		rng:           rand.New(rand.NewSource(0)),
		randomSampler: goptuna.NewRandomSearchSampler(),
	}

	for _, opt := range opts {
		opt(sampler)
	}
	return sampler
}

func genKeepIdxs(
	lossIdxs []int,
	lossAscending []int,
	n int, below bool) []int {
	var l []int
	if below {
		l = lossAscending[:n]
	} else {
		l = lossAscending[n:]
	}

	set := make([]int, 0, len(l))
	isExist := func(l []int, item int) bool {
		for i := range l {
			if l[i] == item {
				return true
			}
		}
		return false
	}
	for _, index := range l {
		item := lossIdxs[index]
		if !isExist(set, item) {
			set = append(set, item)
		}
	}
	return set
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

func (s *Sampler) splitObservationPairs(
	configIdxs []int,
	configVals []float64,
	lossIdxs []int,
	lossVals [][2]float64,
) ([]float64, []float64) {
	nbelow := s.gamma(len(configVals))
	lossAscending := argSort2d(lossVals)

	keepIdxs := genKeepIdxs(lossIdxs, lossAscending, nbelow, true)
	below := genBelowOrAbove(keepIdxs, configIdxs, configVals)

	keepIdxs = genKeepIdxs(lossIdxs, lossAscending, nbelow, false)
	above := genBelowOrAbove(keepIdxs, configIdxs, configVals)

	return below, above
}

func (s *Sampler) sampleFromGMM(parzenEstimator *ParzenEstimator, low, high float64, size int, q float64) []float64 {
	weights := parzenEstimator.Weights
	mus := parzenEstimator.Mus
	sigmas := parzenEstimator.Sigmas
	nsamples := size

	if low > high {
		panic("the low should be lower than the high")
	}

	samples := make([]float64, 0, nsamples)
	for {
		if len(samples) == nsamples {
			break
		}
		active, err := argMaxMultinomial(weights)
		if err != nil {
			panic(err)
		}
		x := s.rng.NormFloat64()
		draw := x*sigmas[active] + mus[active]
		if low <= draw && draw < high {
			samples = append(samples, draw)
		}
	}
	if q > 0 {
		for i := range samples {
			samples[i] = math.Round(samples[i]/q) * q
		}
	}
	return samples
}

func (s *Sampler) normalCDF(x float64, mu []float64, sigma []float64) []float64 {
	l := len(mu)
	results := make([]float64, l)
	for i := 0; i < l; i++ {
		denominator := x - mu[i]
		numerator := math.Max(math.Sqrt(2)*sigma[i], EPS)
		z := denominator / numerator
		results[i] = 0.5 * (1 + math.Erf(z))
	}
	return results
}

func (s *Sampler) logsumRows(x [][]float64) []float64 {
	y := make([]float64, len(x))
	for i := range x {
		m := floats.Max(x[i])

		sum := 0.0
		for j := range x[i] {
			sum += math.Log(math.Exp(x[i][j] - m))
		}
		y[i] = sum + m
	}
	return y
}

func (s *Sampler) gmmLogPDF(samples []float64, parzenEstimator *ParzenEstimator, low, high float64, q float64) []float64 {
	weights := parzenEstimator.Weights
	mus := parzenEstimator.Mus
	sigmas := parzenEstimator.Sigmas

	if len(samples) == 0 {
		return []float64{}
	}

	highNormalCdf := s.normalCDF(high, mus, sigmas)
	lowNormalCdf := s.normalCDF(low, mus, sigmas)
	if len(weights) != len(highNormalCdf) {
		panic("the length should be the same with weights")
	}

	paccept := 0.0
	for i := 0; i < len(highNormalCdf); i++ {
		paccept += highNormalCdf[i]*weights[i] - lowNormalCdf[i]
	}

	if q > 0 {
		probabilities := make([]float64, len(samples))
		if len(weights) != len(mus) || len(weights) != len(sigmas) {
			panic("should be the same length of weights, mus and sigmas")
		}
		for i := range weights {
			w := weights[i]
			mu := mus[i]
			sigma := sigmas[i]
			upperBound := make([]float64, len(samples))
			lowerBound := make([]float64, len(samples))
			for i := range upperBound {
				upperBound[i] = math.Min(samples[i]+q/2.0, high)
				lowerBound[i] = math.Max(samples[i]-q/2.0, low)

			}

			incAmt := make([]float64, len(samples))
			for j := range upperBound {
				incAmt[j] = w * s.normalCDF(upperBound[j], []float64{mu}, []float64{sigma})[0]
				incAmt[j] -= w * s.normalCDF(lowerBound[j], []float64{mu}, []float64{sigma})[0]
			}
			for j := range probabilities {
				probabilities[j] += incAmt[j]
			}
		}
		returnValue := make([]float64, len(samples))
		for i := range probabilities {
			returnValue[i] = math.Log(probabilities[i]+EPS) + math.Log(paccept+EPS)
		}
		return returnValue
	} else {
		jacobian := ones1d(len(samples))
		distance := make([][]float64, len(samples))
		for i := range samples {
			distance[i] = make([]float64, len(mus))
			for j := range mus {
				distance[i][j] = samples[i] - mus[j]
			}
		}
		mahalanobis := make([][]float64, len(distance))
		for i := range distance {
			mahalanobis[i] = make([]float64, len(distance[i]))
			for j := range distance[i] {
				mahalanobis[i][j] = distance[i][j] / math.Pow(math.Max(sigmas[j], EPS), 2)
			}
		}
		z := make([][]float64, len(distance))
		for i := range distance {
			z[i] = make([]float64, len(distance[i]))
			for j := range distance[i] {
				z[i][j] = math.Sqrt(2*math.Pi) * sigmas[j] * jacobian[i]
			}
		}
		coefficient := make([][]float64, len(distance))
		for i := range distance {
			coefficient[i] = make([]float64, len(distance[i]))
			for j := range distance[i] {
				coefficient[i][j] = weights[j] / z[i][j] / paccept
			}
		}

		y := make([][]float64, len(distance))
		for i := range distance {
			y[i] = make([]float64, len(distance[i]))
			for j := range distance[i] {
				y[i][j] = -0.5*mahalanobis[i][j] + math.Log(coefficient[i][j])
			}
		}
		return s.logsumRows(y)
	}
}

func (s *Sampler) compare(samples []float64, logL []float64, logG []float64) []float64 {
	if len(samples) > 0 {
		if len(logL) != len(logG) {
			panic("the size of the log_l and log_g should be same")
		}
		score := make([]float64, len(logL))
		for i := range score {
			score[i] = logL[i] - logG[i]
		}
		if len(samples) != len(score) {
			panic("the size of the samples and score should be same")
		}

		argMax := func(s []float64) int {
			max := s[0]
			maxIdx := 0
			for i := range s {
				if i == 0 {
					continue
				}
				if s[i] > max {
					max = s[i]
					maxIdx = i
				}
			}
			return maxIdx
		}
		best := argMax(score)
		results := make([]float64, len(samples))
		for i := range results {
			results[i] = samples[best]
		}
		return results
	} else {
		return []float64{}
	}
}

func (s *Sampler) sampleNumerical(low, high float64, below, above []float64, q float64) float64 {
	size := s.numberOfEICandidates
	parzenEstimatorBelow := NewParzenEstimator(below, low, high, s.params)
	sampleBelow := s.sampleFromGMM(parzenEstimatorBelow, low, high, size, q)
	logLikelihoodsBelow := s.gmmLogPDF(sampleBelow, parzenEstimatorBelow, low, high, q)

	parzenEstimatorAbove := NewParzenEstimator(above, low, high, s.params)
	sampleAbove := s.sampleFromGMM(parzenEstimatorAbove, low, high, size, q)
	logLikelihoodsAbove := s.gmmLogPDF(sampleAbove, parzenEstimatorAbove, low, high, q)

	return s.compare(sampleBelow, logLikelihoodsBelow, logLikelihoodsAbove)[0]
}

func (s *Sampler) sampleUniform(distribution goptuna.UniformDistribution, below, above []float64) float64 {
	low := distribution.Low
	high := distribution.High
	return s.sampleNumerical(low, high, below, above, 0)
}

func (s *Sampler) sampleInt(distribution goptuna.IntUniformDistribution, below, above []float64) float64 {
	q := 1.0
	low := float64(distribution.Low) - 0.5*q
	high := float64(distribution.High) + 0.5*q
	return s.sampleNumerical(low, high, below, above, q)
}

func (s *Sampler) Sample(
	study *goptuna.Study,
	trial goptuna.FrozenTrial,
	paramName string,
	paramDistribution interface{},
) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	observationPairs, err := getObservationPairs(study, paramName)
	if err != nil {
		return 0, err
	}
	n := len(observationPairs)

	if n < s.numberOfStartupTrials {
		return s.randomSampler.Sample(study, trial, paramName, paramDistribution)
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
	belowParamValues, aboveParamValues := s.splitObservationPairs(
		configIdxs, configVals, lossIdxs, lossVals)

	switch d := paramDistribution.(type) {
	case goptuna.UniformDistribution:
		return s.sampleUniform(d, belowParamValues, aboveParamValues), nil
	case goptuna.IntUniformDistribution:
		return s.sampleInt(d, belowParamValues, aboveParamValues), nil
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
