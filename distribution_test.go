package goptuna_test

import (
	"reflect"
	"testing"

	"github.com/c-bata/goptuna"
)

func TestDistributionConversionBetweenDistributionAndJSON(t *testing.T) {
	tests := []struct {
		name         string
		distribution interface{}
	}{
		{
			name: "uniform distribution",
			distribution: goptuna.UniformDistribution{
				High: 10.0,
				Low:  -5.0,
			},
		},
		{
			name: "log uniform distribution",
			distribution: goptuna.LogUniformDistribution{
				High: 1e2,
				Low:  1e-1,
			},
		},
		{
			name: "int uniform distribution",
			distribution: goptuna.IntUniformDistribution{
				High: 20,
				Low:  5,
			},
		},
		{
			name: "discrete uniform distribution",
			distribution: goptuna.DiscreteUniformDistribution{
				High: 20,
				Low:  5,
				Q:    1,
			},
		},
		{
			name: "step int uniform distribution",
			distribution: goptuna.StepIntUniformDistribution{
				High: 20,
				Low:  5,
				Step: 2,
			},
		},
		{
			name: "categorical distribution",
			distribution: goptuna.CategoricalDistribution{
				Choices: []string{"foo", "bar"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := goptuna.DistributionToJSON(tt.distribution)
			if err != nil {
				t.Errorf("DistributionToJSON should not be err, but got %s", err)
			}
			d, err := goptuna.JSONToDistribution(got)
			if err != nil {
				t.Errorf("JSONToDistribution should not be err, but got %s", err)
			}
			if !reflect.DeepEqual(tt.distribution, d) {
				t.Errorf("Must be the same, but %#v != %#v", tt.distribution, d)
			}
		})
	}
}

func TestDistributionToExternalRepresentation(t *testing.T) {
	tests := []struct {
		name         string
		distribution goptuna.Distribution
		args         float64
		want         interface{}
	}{
		{
			name:         "uniform distribution",
			distribution: &goptuna.UniformDistribution{Low: 0.5, High: 5.5},
			args:         3.5,
			want:         3.5,
		},
		{
			name:         "log uniform distribution",
			distribution: &goptuna.LogUniformDistribution{Low: 1e-2, High: 1e3},
			args:         1e2,
			want:         1e2,
		},
		{
			name:         "int uniform distribution 1",
			distribution: &goptuna.IntUniformDistribution{Low: 0, High: 10},
			args:         3.0,
			want:         3,
		},
		{
			name:         "int uniform distribution 2",
			distribution: &goptuna.IntUniformDistribution{Low: 0, High: 10},
			args:         1.6,
			want:         2,
		},
		{
			name:         "step int uniform distribution 1",
			distribution: &goptuna.StepIntUniformDistribution{Low: 0, High: 10, Step: 2},
			args:         2.2,
			want:         2,
		},
		{
			name:         "step int uniform distribution 2",
			distribution: &goptuna.StepIntUniformDistribution{Low: 0, High: 10, Step: 3},
			args:         2.2,
			want:         3,
		},
		{
			name:         "step int uniform distribution 2",
			distribution: &goptuna.StepIntUniformDistribution{Low: -1, High: 10, Step: 3},
			args:         0.6,
			want:         2,
		},
		{
			name:         "discrete uniform distribution 1",
			distribution: &goptuna.DiscreteUniformDistribution{Low: 0.5, High: 5.5, Q: 0.5},
			args:         3.5,
			want:         3.5,
		},
		{
			name:         "discrete uniform distribution 2",
			distribution: &goptuna.DiscreteUniformDistribution{Low: 0.5, High: 5.5, Q: 0.5},
			args:         3.3,
			want:         3.5,
		},
		{
			name:         "discrete uniform distribution 3",
			distribution: &goptuna.DiscreteUniformDistribution{Low: 0.5, High: 5.5, Q: 0.05},
			args:         3.52,
			want:         3.5,
		},
		{
			name:         "categorical distribution",
			distribution: &goptuna.CategoricalDistribution{Choices: []string{"a", "b", "c"}},
			args:         2.0,
			want:         "c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.distribution.ToExternalRepr(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Distribution.ToExternalRepr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDistributionSingle(t *testing.T) {
	tests := []struct {
		name         string
		distribution goptuna.Distribution
		want         bool
	}{
		{
			name:         "uniform distribution true",
			distribution: &goptuna.UniformDistribution{Low: 0.5, High: 0.5},
			want:         true,
		},
		{
			name:         "uniform distribution false",
			distribution: &goptuna.UniformDistribution{Low: 0.5, High: 5.5},
			want:         false,
		},
		{
			name:         "log uniform distribution true",
			distribution: &goptuna.LogUniformDistribution{Low: 1e3, High: 1e3},
			want:         true,
		},
		{
			name:         "log uniform distribution false",
			distribution: &goptuna.LogUniformDistribution{Low: 1e-3, High: 1e2},
			want:         false,
		},
		{
			name:         "int uniform distribution true",
			distribution: &goptuna.IntUniformDistribution{Low: 10, High: 10},
			want:         true,
		},
		{
			name:         "int uniform distribution false",
			distribution: &goptuna.IntUniformDistribution{Low: 0, High: 10},
			want:         false,
		},
		{
			name:         "discrete uniform distribution true",
			distribution: &goptuna.DiscreteUniformDistribution{Low: 10, High: 10, Q: 1},
			want:         true,
		},
		{
			name:         "discrete uniform distribution false",
			distribution: &goptuna.DiscreteUniformDistribution{Low: 0, High: 10, Q: 1},
			want:         false,
		},
		{
			name:         "step int uniform distribution true",
			distribution: &goptuna.StepIntUniformDistribution{Low: 10, High: 10, Step: 2},
			want:         true,
		},
		{
			name:         "step int uniform distribution true when step is greater than high-low",
			distribution: &goptuna.StepIntUniformDistribution{Low: 0, High: 1, Step: 2},
			want:         true,
		},
		{
			name:         "step int uniform distribution false",
			distribution: &goptuna.StepIntUniformDistribution{Low: 0, High: 10, Step: 2},
			want:         false,
		},
		{
			name:         "categorical distribution true",
			distribution: &goptuna.CategoricalDistribution{Choices: []string{"a"}},
			want:         true,
		},
		{
			name:         "categorical distribution false",
			distribution: &goptuna.CategoricalDistribution{Choices: []string{"a", "b", "c"}},
			want:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.distribution.Single(); got != tt.want {
				t.Errorf("Distribution.Single() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDistributionContains(t *testing.T) {
	tests := []struct {
		name         string
		distribution goptuna.Distribution
		args         float64
		want         bool
	}{
		{
			name:         "uniform distribution true",
			distribution: &goptuna.UniformDistribution{Low: 0.5, High: 5.5},
			args:         3.5,
			want:         true,
		},
		{
			name:         "uniform distribution upper bound",
			distribution: &goptuna.UniformDistribution{Low: 0.5, High: 5.5},
			args:         5.5,
			want:         true,
		},
		{
			name:         "uniform distribution lower",
			distribution: &goptuna.UniformDistribution{Low: 0.5, High: 5.5},
			args:         -0.5,
			want:         false,
		},
		{
			name:         "uniform distribution higher",
			distribution: &goptuna.UniformDistribution{Low: 0.5, High: 5.5},
			args:         7.5,
			want:         false,
		},
		{
			name:         "log uniform distribution true",
			distribution: &goptuna.LogUniformDistribution{Low: 1e-1, High: 1e3},
			args:         1e2,
			want:         true,
		},
		{
			name:         "log uniform distribution upper bound",
			distribution: &goptuna.LogUniformDistribution{Low: 1e-1, High: 1e3},
			args:         1e2,
			want:         true,
		},
		{
			name:         "log uniform distribution lower",
			distribution: &goptuna.LogUniformDistribution{Low: 1e-1, High: 1e3},
			args:         1e-3,
			want:         false,
		},
		{
			name:         "log uniform distribution higher",
			distribution: &goptuna.LogUniformDistribution{Low: 1e-1, High: 1e3},
			args:         1e5,
			want:         false,
		},
		{
			name:         "int uniform distribution true 1",
			distribution: &goptuna.IntUniformDistribution{Low: 0, High: 10},
			args:         3,
			want:         true,
		},
		{
			name:         "int uniform distribution true 2",
			distribution: &goptuna.IntUniformDistribution{Low: 0, High: 10},
			args:         -0.4999,
			want:         true,
		},
		{
			name:         "int uniform distribution true 3",
			distribution: &goptuna.IntUniformDistribution{Low: 0, High: 10},
			args:         10.4999,
			want:         true,
		},
		{
			name:         "int uniform distribution lower",
			distribution: &goptuna.IntUniformDistribution{Low: 0, High: 10},
			args:         -3,
			want:         false,
		},
		{
			name:         "int uniform distribution higher",
			distribution: &goptuna.IntUniformDistribution{Low: 0, High: 10},
			args:         15,
			want:         false,
		},
		{
			name:         "step int uniform distribution true",
			distribution: &goptuna.StepIntUniformDistribution{Low: 0, High: 10, Step: 2},
			args:         3,
			want:         true,
		},
		{
			name:         "step int uniform distribution lower",
			distribution: &goptuna.StepIntUniformDistribution{Low: 0, High: 10, Step: 2},
			args:         -3,
			want:         false,
		},
		{
			name:         "step int uniform distribution higher",
			distribution: &goptuna.StepIntUniformDistribution{Low: 0, High: 10, Step: 2},
			args:         15,
			want:         false,
		},
		{
			name:         "discrete uniform distribution true",
			distribution: &goptuna.DiscreteUniformDistribution{Low: 0.5, High: 5.5, Q: 1},
			args:         3.5,
			want:         true,
		},
		{
			name:         "discrete uniform distribution true",
			distribution: &goptuna.DiscreteUniformDistribution{Low: 0.0, High: 1.0, Q: 0.3},
			args:         0.3,
			want:         true,
		},
		{
			name:         "discrete uniform distribution true aaa",
			distribution: &goptuna.DiscreteUniformDistribution{Low: 0.0, High: 1.0, Q: 0.3},
			args:         0.7,
			want:         true,
		},
		{
			name:         "discrete uniform distribution true",
			distribution: &goptuna.DiscreteUniformDistribution{Low: 0.5, High: 3.5, Q: 1},
			args:         3.0,
			want:         false,
		},
		{
			name:         "discrete uniform distribution lower",
			distribution: &goptuna.DiscreteUniformDistribution{Low: 0.5, High: 5.5, Q: 0.5},
			args:         -3,
			want:         false,
		},
		{
			name:         "discrete uniform distribution higher",
			distribution: &goptuna.DiscreteUniformDistribution{Low: 0.5, High: 5.5, Q: 0.5},
			args:         15,
			want:         false,
		},
		{
			name:         "categorical distribution true",
			distribution: &goptuna.CategoricalDistribution{Choices: []string{"a", "b", "c"}},
			args:         1,
			want:         true,
		},
		{
			name:         "categorical distribution lower",
			distribution: &goptuna.CategoricalDistribution{Choices: []string{"a", "b", "c"}},
			args:         -1,
			want:         false,
		},
		{
			name:         "categorical distribution higher",
			distribution: &goptuna.CategoricalDistribution{Choices: []string{"a", "b", "c"}},
			args:         3,
			want:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.distribution.Contains(tt.args); got != tt.want {
				t.Errorf("Distribution.ToInternalRepr() = %v, want %v", got, tt.want)
			}
		})
	}
}
