package tpe

import (
	"math/rand"

	"github.com/c-bata/goptuna"
)

// SamplerOption is a type of the function to customizing TPE sampler.
type SamplerOption func(sampler *Sampler)

// SamplerOptionSeed sets seed number.
func SamplerOptionSeed(seed int64) SamplerOption {
	randomSampler := goptuna.NewRandomSearchSampler(
		goptuna.RandomSearchSamplerOptionSeed(seed))

	return func(sampler *Sampler) {
		sampler.rng = rand.New(rand.NewSource(seed))
		sampler.randomSampler = randomSampler
	}
}

// SamplerOptionConsiderPrior enhance the stability of Parzen estimator
// by imposing a Gaussian prior when True. The prior is only effective
// if the sampling distribution is either `UniformDistribution`,
// `DiscreteUniformDistribution`, `LogUniformDistribution`, or `IntUniformDistribution`.
func SamplerOptionConsiderPrior(considerPrior bool) SamplerOption {
	return func(sampler *Sampler) {
		sampler.params.ConsiderPrior = considerPrior
	}
}

// SamplerOptionPriorWeight sets the weight of the prior.
func SamplerOptionPriorWeight(priorWeight float64) SamplerOption {
	return func(sampler *Sampler) {
		sampler.params.PriorWeight = priorWeight
	}
}

// SamplerOptionPriorWeight enable a heuristic to limit the smallest variances
// of Gaussians used in the Parzen estimator.
func SamplerOptionConsiderMagicClip(considerMagicClip bool) SamplerOption {
	return func(sampler *Sampler) {
		sampler.params.ConsiderMagicClip = considerMagicClip
	}
}

// SamplerOptionConsiderEndpoints take endpoints of domains into account
// when calculating variances of Gaussians in Parzen estimator.
// See the original paper for details on the heuristics to calculate the variances.
func SamplerOptionConsiderEndpoints(considerEndpoints bool) SamplerOption {
	return func(sampler *Sampler) {
		sampler.params.ConsiderEndpoints = considerEndpoints
	}
}

// SamplerOptionWeights sets the function that takes the number of finished trials
// and returns a weight for them. See `Making a Science of Model Search: Hyperparameter
// Optimization in Hundreds of Dimensions for Vision Architectures
// <http://proceedings.mlr.press/v28/bergstra13.pdf>` for more details.
func SamplerOptionWeights(weights func(x int) []float64) SamplerOption {
	return func(sampler *Sampler) {
		sampler.params.Weights = weights
	}
}

// SamplerOptionGammaFunc sets the function that takes the number of
// finished trials and returns the number of trials to form a density
// function for samples with low grains.
func SamplerOptionGammaFunc(gamma FuncGamma) SamplerOption {
	return func(sampler *Sampler) {
		sampler.gamma = gamma
	}
}

// SamplerOptionNumberOfEICandidates sets the number of EI candidates (default 24).
func SamplerOptionNumberOfEICandidates(n int) SamplerOption {
	return func(sampler *Sampler) {
		sampler.numberOfEICandidates = n
	}
}

// SamplerOptionNumberOfStartupTrials sets the number of start up trials (default 10).
func SamplerOptionNumberOfStartupTrials(n int) SamplerOption {
	return func(sampler *Sampler) {
		sampler.numberOfStartupTrials = n
	}
}

// SamplerOptionParzenEstimatorParams sets the parameter of ParzenEstimator.
func SamplerOptionParzenEstimatorParams(params ParzenEstimatorParams) SamplerOption {
	return func(sampler *Sampler) {
		sampler.params = params
	}
}
