package cma

import (
	"math/rand"
)

// SamplerOption is a type of the function to customizing CMA-ES sampler.
type SamplerOption func(sampler *Sampler)

// SamplerOptionSeed sets seed number.
func SamplerOptionSeed(seed int64) SamplerOption {
	return func(sampler *Sampler) {
		sampler.rng = rand.New(rand.NewSource(seed))
	}
}
