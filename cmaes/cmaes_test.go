package cmaes_test

import (
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"

	"github.com/c-bata/goptuna/cmaes"
)

func TestNewOptimizer(t *testing.T) {
	mean := []float64{0, 0}
	sigma0 := 1.3
	optimizer, err := cmaes.NewOptimizer(mean, sigma0)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	if dim := cmaes.ExportDim(optimizer); dim != 2 {
		t.Errorf("should be 2, but got %d", dim)
	}
}

func TestOptimizer_Ask(t *testing.T) {
	mean := []float64{1, 2}
	sigma0 := 1.3
	optimizer, err := cmaes.NewOptimizer(
		mean, sigma0,
		cmaes.OptimizerOptionSeed(0),
	)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
	}
	x, err := optimizer.Ask()
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	r, c := x.Dims()
	if r != 2 || c != 1 {
		t.Errorf("should be (2, 1), but got (%d, %d)", r, c)
	}
}

func TestOptimizer_Tell(t *testing.T) {
	mean := []float64{1, 2}
	sigma0 := 1.3
	optimizer, err := cmaes.NewOptimizer(
		mean, sigma0,
		cmaes.OptimizerOptionSeed(0),
	)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
	}
	solutions := make([]*cmaes.Solution, optimizer.PopulationSize())
	for i := 0; i < optimizer.PopulationSize(); i++ {
		x, err := optimizer.Ask()
		if err != nil {
			t.Errorf("should be nil, but got %s", err)
			return
		}
		solutions[i] = &cmaes.Solution{
			X:     x,
			Value: float64(i),
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
			optimizer, err := cmaes.NewOptimizer(
				[]float64{0, 0}, 1.3,
				cmaes.OptimizerOptionBounds(tt.bounds),
			)
			if err != nil {
				t.Errorf("should be nil, but got %s", err)
			}

			feasible := cmaes.ExportOptimizerIsFeasible(optimizer, tt.value)
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
			optimizer, err := cmaes.NewOptimizer(
				[]float64{0, 0}, 1.3,
				cmaes.OptimizerOptionBounds(tt.bounds),
			)
			if err != nil {
				t.Errorf("should be nil, but got %s", err)
			}

			err = cmaes.ExportOptimizerRepairInfeasibleParams(optimizer, tt.value)
			if err != nil {
				t.Errorf("should be nil, but got %s", err)
			}
			if !floats.Same(tt.value.RawVector().Data, tt.repaired.RawVector().Data) {
				t.Errorf("should be %v, but got %v", tt.value.RawVector().Data, tt.repaired.RawVector().Data)
			}
		})
	}
}
