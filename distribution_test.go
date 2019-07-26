package goptuna_test

import (
	"reflect"
	"testing"

	"github.com/c-bata/goptuna"
)

func TestConversionBetweenDistributionAndJSON(t *testing.T) {
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
			name: "int uniform distribution",
			distribution: goptuna.IntUniformDistribution{
				High: 20,
				Low:  5,
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

func TestToInternalRepresentation(t *testing.T) {
	tests := []struct {
		name         string
		distribution goptuna.Distribution
		args         interface{}
		want         float64
	}{
		{
			name:         "uniform distribution",
			distribution: &goptuna.UniformDistribution{Low: 0.5, High: 5.5},
			args:         3.5,
			want:         3.5,
		},
		{
			name:         "int uniform distribution",
			distribution: &goptuna.IntUniformDistribution{Low: 0, High: 10},
			args:         3,
			want:         3.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.distribution.ToInternalRepr(tt.args); got != tt.want {
				t.Errorf("UniformDistribution.ToInternalRepr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToExternalRepresentation(t *testing.T) {
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
			name:         "int uniform distribution",
			distribution: &goptuna.IntUniformDistribution{Low: 0, High: 10},
			args:         3.0,
			want:         3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.distribution.ToExternalRepr(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniformDistribution.ToInternalRepr() = %v, want %v", got, tt.want)
			}
		})
	}
}
