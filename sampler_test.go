package goptuna_test

import (
	"testing"

	"github.com/c-bata/goptuna"
)

func TestRandomSearchSamplerOptionSeed(t *testing.T) {
	tests := []struct {
		name         string
		distribution interface{}
	}{
		{
			name: "uniform distribution",
			distribution: goptuna.UniformDistribution{
				High: 10,
				Low:  0,
			},
		},
		{
			name: "int uniform distribution",
			distribution: goptuna.IntUniformDistribution{
				High: 10,
				Low:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler1 := goptuna.NewRandomSearchSampler()
			sampler2 := goptuna.NewRandomSearchSampler(goptuna.RandomSearchSamplerOptionSeed(2))
			sampler3 := goptuna.NewRandomSearchSampler(goptuna.RandomSearchSamplerOptionSeed(2))

			s1, err := sampler1.Sample(nil, goptuna.FrozenTrial{}, "foo", tt.distribution)
			if err != nil {
				t.Errorf("should not be err, but got %s", err)
			}
			s2, err := sampler2.Sample(nil, goptuna.FrozenTrial{}, "foo", tt.distribution)
			if err != nil {
				t.Errorf("should not be err, but got %s", err)
			}
			s3, err := sampler3.Sample(nil, goptuna.FrozenTrial{}, "foo", tt.distribution)
			if err != nil {
				t.Errorf("should not be err, but got %s", err)
			}
			if s1 == s2 {
				t.Errorf("should not be the same but both are %f", s1)
			}
			if s2 != s3 {
				t.Errorf("should be equal, but got %f and %f", s2, s3)
			}
		})
	}
}
