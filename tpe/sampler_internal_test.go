package tpe

import (
	"testing"

	"github.com/c-bata/goptuna"
)

// Following test cases are generated from Optuna's behavior.

func TestSampler_SampleCategorical(t *testing.T) {
	d := goptuna.CategoricalDistribution{
		Choices: []string{"a", "b", "c", "d"},
	}
	below := []float64{1.0}
	above := []float64{1.0, 3.0, 3.0, 2.0, 3.0, 0.0, 2.0, 3.0, 3.0}
	expected := 1.0

	sampler := NewSampler()
	actual := sampler.sampleCategorical(d, below, above)
	if expected != actual {
		t.Errorf("should be %f, but got %f", expected, actual)
	}
}
