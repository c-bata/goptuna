package cmaes

import (
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
