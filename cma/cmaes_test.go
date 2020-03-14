package cma_test

import (
	"testing"

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
