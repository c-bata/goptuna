package cmaes

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

func TestNewOptimizer(t *testing.T) {
	mean := []float64{0, 0}
	sigma0 := 1.3
	optimizer, err := NewOptimizer(mean, sigma0)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	if optimizer.dim != 2 {
		t.Errorf("should be 2, but got %d", optimizer.dim)
	}
	if optimizer.popsize != 6 {
		t.Errorf("should be 6, but got %d", optimizer.popsize)
	}
	if optimizer.mu != 3 {
		t.Errorf("should be 3, but got %d", optimizer.mu)
	}
	if math.Abs(optimizer.muEff-2.0286114646100617) > 0.0001 {
		t.Errorf("should be 2.0286114646100617, but got %f", optimizer.muEff)
	}
	if math.Abs(optimizer.c1-0.1548153998964136) > 0.0001 {
		t.Errorf("should be 0.1548153998964136, but got %f", optimizer.c1)
	}
	if math.Abs(optimizer.cmu-0.05785908507191633) > 0.0001 {
		t.Errorf("should be 0.05785908507191633, but got %f", optimizer.cmu)
	}
	expectedWeights := []float64{0.63704257, 0.28457026, 0.07838717, -0.28638378, -0.76495809, -1.15598178}
	if !floats.EqualApprox(optimizer.weights.RawVector().Data, expectedWeights, 0.0001) {
		t.Errorf("should be %#v, but got %#v", expectedWeights, optimizer.weights.RawVector().Data)
	}
	if math.Abs(optimizer.cSigma-0.44620498737831715) > 0.0001 {
		t.Errorf("should be 0.44620498737831715, but got %f", optimizer.cSigma)
	}
	if math.Abs(optimizer.dSigma-1.4462049873783172) > 0.0001 {
		t.Errorf("should be 1.4462049873783172, but got %f", optimizer.dSigma)
	}
	if math.Abs(optimizer.cc-0.6245545390268264) > 0.0001 {
		t.Errorf("should be 0.6245545390268264, but got %f", optimizer.cc)
	}
	if math.Abs(optimizer.chiN-1.254272742818995) > 0.0001 {
		t.Errorf("should be 1.254272742818995, but got %f", optimizer.chiN)
	}
}

func TestOptimizer_Ask(t *testing.T) {
	mean := []float64{1, 2}
	sigma0 := 1.3
	optimizer, err := NewOptimizer(
		mean, sigma0,
		OptimizerOptionSeed(0),
	)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
	}
	x, err := optimizer.Ask()
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	if len(x) != 2 {
		t.Errorf("dim should be 2, but got %d", len(x))
	}
}

func TestOptimizer_Tell(t *testing.T) {
	mean := []float64{1, 2}
	sigma0 := 1.3
	optimizer, err := NewOptimizer(
		mean, sigma0,
		OptimizerOptionSeed(0),
	)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
	}
	solutions := make([]*Solution, optimizer.PopulationSize())
	for i := 0; i < optimizer.PopulationSize(); i++ {
		x, err := optimizer.Ask()
		if err != nil {
			t.Errorf("should be nil, but got %s", err)
			return
		}
		solutions[i] = &Solution{
			Params: x,
			Value:  float64(i),
		}
	}
	err = optimizer.Tell(solutions)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
}

func TestOptimizer_IsFeasible(t *testing.T) {
	tests := []struct {
		name   string
		bounds *mat.Dense
		value  *mat.VecDense
		want   bool
	}{
		{
			name: "feasible",
			bounds: mat.NewDense(2, 2, []float64{
				-1, 1,
				-2, -1,
			}),
			value: mat.NewVecDense(2, []float64{-0.5, -1.5}),
			want:  true,
		},
		{
			name: "out of lower bound",
			bounds: mat.NewDense(2, 2, []float64{
				-1, 1,
				-2, -1,
			}),
			value: mat.NewVecDense(2, []float64{-1.5, -1.5}),
			want:  false,
		},
		{
			name: "out of upper bound",
			bounds: mat.NewDense(2, 2, []float64{
				-1, 1,
				-2, -1,
			}),
			value: mat.NewVecDense(2, []float64{-0.5, 1}),
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer, err := NewOptimizer(
				[]float64{0, 0}, 1.3,
				OptimizerOptionBounds(tt.bounds),
			)
			if err != nil {
				t.Errorf("should be nil, but got %s", err)
			}

			feasible := optimizer.isFeasible(tt.value)
			if tt.want != feasible {
				t.Errorf("should be %v, but got %v", tt.want, feasible)
			}
		})
	}
}

func TestOptimizer_RepairInfeasibleParams(t *testing.T) {
	tests := []struct {
		name     string
		bounds   *mat.Dense
		value    *mat.VecDense
		repaired *mat.VecDense
	}{
		{
			name: "out of lower bound",
			bounds: mat.NewDense(2, 2, []float64{
				-1, 1,
				-2, -1,
			}),
			value:    mat.NewVecDense(2, []float64{-1.5, -1.5}),
			repaired: mat.NewVecDense(2, []float64{-1, -1.5}),
		},
		{
			name: "out of upper bound",
			bounds: mat.NewDense(2, 2, []float64{
				-1, 1,
				-2, -1,
			}),
			value:    mat.NewVecDense(2, []float64{-0.5, 1}),
			repaired: mat.NewVecDense(2, []float64{-0.5, -1}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer, err := NewOptimizer(
				[]float64{0, 0}, 1.3,
				OptimizerOptionBounds(tt.bounds),
			)
			if err != nil {
				t.Errorf("should be nil, but got %s", err)
			}

			err = optimizer.repairInfeasibleParams(tt.value)
			if err != nil {
				t.Errorf("should be nil, but got %s", err)
			}
			if !floats.Same(tt.value.RawVector().Data, tt.repaired.RawVector().Data) {
				t.Errorf("should be %v, but got %v", tt.value.RawVector().Data, tt.repaired.RawVector().Data)
			}
		})
	}
}
