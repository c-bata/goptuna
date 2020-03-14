package cma_test

import (
	"testing"

	"gonum.org/v1/gonum/mat"

	"github.com/c-bata/goptuna/cma"
)

func TestNewOptimizer(t *testing.T) {
	mean := []float64{0, 0}
	sigma0 := 1.3
	optimizer, err := cma.NewOptimizer(mean, sigma0)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	if dim := cma.ExportDim(optimizer); dim != 2 {
		t.Errorf("should be 2, but got %d", dim)
	}
}

func TestOptimizer_Ask(t *testing.T) {
	mean := []float64{0, 0}
	sigma0 := 1.3
	optimizer, err := cma.NewOptimizer(
		mean, sigma0,
		cma.OptimizerOptionSeed(0),
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
			optimizer, err := cma.NewOptimizer(
				[]float64{0, 0}, 1.3,
				cma.OptimizerOptionBounds(tt.bounds),
			)
			if err != nil {
				t.Errorf("should be nil, but got %s", err)
			}

			feasible := cma.ExportOptimizerIsFeasible(optimizer, tt.value)
			if tt.want != feasible {
				t.Errorf("should be %v, but got %v", tt.want, feasible)
			}
		})
	}
}
