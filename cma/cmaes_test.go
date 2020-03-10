package cma_test

import (
	"testing"

	"github.com/c-bata/goptuna/cma"
)

func TestNewCMAES(t *testing.T) {
	mean := []float64{0, 0}
	sigma0 := 1.3
	optimizer, err := cma.NewOptimizer(mean, sigma0)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
	}
	if dim := cma.ExportDim(optimizer); dim != 2 {
		t.Errorf("should be 2, but got %d", dim)
	}
}
