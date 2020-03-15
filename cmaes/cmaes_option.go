package cmaes

import (
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

// OptimizerOption is a type of the function to customizing CMA-ES.
type OptimizerOption func(*Optimizer)

// OptimizerOptionSeed sets seed number.
func OptimizerOptionSeed(seed int64) OptimizerOption {
	return func(cma *Optimizer) {
		cma.rng = rand.New(rand.NewSource(seed))
	}
}

// OptimizerOptionMaxReSampling sets a number of max re-sampling.
func OptimizerOptionMaxReSampling(n int) OptimizerOption {
	return func(cma *Optimizer) {
		cma.maxReSampling = n
	}
}

// OptimizerOptionBounds sets the range of parameters.
func OptimizerOptionBounds(bounds *mat.Dense) OptimizerOption {
	_, column := bounds.Dims()
	if column != 2 {
		panic("invalid matrix size")
	}

	return func(cma *Optimizer) {
		cma.bounds = bounds
	}
}
