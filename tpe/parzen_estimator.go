package tpe

import (
	"gonum.org/v1/gonum/mat"
)

type ParzenEstimatorParams struct {
	ConsiderPrior     bool
	ConsiderMagicClip bool
	ConsiderEndpoints bool
	Weights           FuncWeights
	PriorWeight       float64 // optional
	PriorWeightIsSet  bool    // for PriorWeight
}

type ParzenEstimator struct {
	Weights mat.Vector
	Mus     mat.Vector
	Sigma   mat.Vector
	Params  ParzenEstimatorParams
}

func (e *ParzenEstimator) calculate(
	mus []float64,
	low float64,
	high float64,
	considerPrior bool,
	priorWeight float64,
	ConsiderMagicClip bool,
	ConsiderEndpoints bool,
	Weights FuncWeights,
) (*mat.VecDense, *mat.VecDense, *mat.VecDense) {
	panic("not implemented yet")
}

func NewParzenEstimator(mus []float64, low, high float64, params ParzenEstimatorParams) *ParzenEstimator {
	estimator := &ParzenEstimator{
		Weights: nil,
		Mus:     nil,
		Sigma:   nil,
		Params:  params,
	}

	sWeights, sMus, sigma := estimator.calculate(mus, low, high, params.ConsiderPrior, params.PriorWeight,
		params.ConsiderMagicClip, params.ConsiderEndpoints, params.Weights)
	estimator.Weights = sWeights
	estimator.Mus = sMus
	estimator.Sigma = sigma
	return estimator
}
