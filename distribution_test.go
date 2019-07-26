package goptuna

import (
	"reflect"
	"testing"
)

func TestConversionBetweenDistributionAndJSON(t *testing.T) {
	tests := []struct {
		name         string
		distribution interface{}
	}{
		{
			name: "uniform distribution",
			distribution: UniformDistribution{
				High: 10.0,
				Low:  -5.0,
			},
		},
		{
			name: "int uniform distribution",
			distribution: IntUniformDistribution{
				High: 20,
				Low:  5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DistributionToJSON(tt.distribution)
			if err != nil {
				t.Errorf("DistributionToJSON should not be err, but got %s", err)
			}
			d, err := JSONToDistribution(got)
			if err != nil {
				t.Errorf("JSONToDistribution should not be err, but got %s", err)
			}
			if !reflect.DeepEqual(tt.distribution, d) {
				t.Errorf("Must be the same, but %#v != %#v", tt.distribution, d)
			}
		})
	}
}
