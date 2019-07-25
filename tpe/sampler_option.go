package tpe

import "math/rand"

type SamplerOption func(sampler *Sampler)

func SamplerOptionSeed(seed int64) SamplerOption {
	return func(sampler *Sampler) {
		sampler.rng = rand.New(rand.NewSource(seed))
	}
}

func SamplerOptionGammaFunc(gamma FuncGamma) SamplerOption {
	return func(sampler *Sampler) {
		sampler.gamma = gamma
	}
}

func SamplerOptionNumberOfEICandidates(n int) SamplerOption {
	return func(sampler *Sampler) {
		sampler.numberOfEICandidates = n
	}
}

func SamplerOptionNumberOfStartupTrials(n int) SamplerOption {
	return func(sampler *Sampler) {
		sampler.numberOfStartupTrials = n
	}
}

func SamplerOptionParzenEstimatorParams(params ParzenEstimatorParams) SamplerOption {
	return func(sampler *Sampler) {
		sampler.params = params
	}
}
