package tpe_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/c-bata/goptuna/internal/testutil"

	"github.com/c-bata/goptuna"

	"github.com/c-bata/goptuna/tpe"
)

func TestDefaultGamma(t *testing.T) {
	type args struct {
		x int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test case 5",
			args: args{x: 5},
			want: 1,
		},
		{
			name: "test case 100",
			args: args{x: 100},
			want: 10,
		},
		{
			name: "test case 255",
			args: args{x: 255},
			want: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tpe.DefaultGamma(tt.args.x); got != tt.want {
				t.Errorf("DefaultGamma() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHyperoptDefaultGamma(t *testing.T) {
	type args struct {
		x int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test case 1",
			args: args{x: 5},
			want: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tpe.HyperoptDefaultGamma(tt.args.x); got != tt.want {
				t.Errorf("HyperoptDefaultGamma() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultWeights(t *testing.T) {
	type args struct {
		x int
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "test case 1",
			args: args{
				x: 30,
			},
			want: []float64{
				(1.0-1.0/30)*0/4 + 1.0/30,
				(1.0-1.0/30)*1/4 + 1.0/30,
				(1.0-1.0/30)*2/4 + 1.0/30,
				(1.0-1.0/30)*3/4 + 1.0/30,
				(1.0-1.0/30)*4/4 + 1.0/30,
				1, 1, 1, 1, 1,
				1, 1, 1, 1, 1,
				1, 1, 1, 1, 1,
				1, 1, 1, 1, 1,
				1, 1, 1, 1, 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tpe.DefaultWeights(tt.args.x); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultWeights() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSampler_SampleDiscreteUniform(t *testing.T) {
	sampler := tpe.NewSampler(tpe.SamplerOptionNumberOfStartupTrials(0))
	study, err := goptuna.CreateStudy("", goptuna.StudyOptionSampler(sampler))
	if err != nil {
		t.Errorf("should not be err, but got %s", err)
		return
	}

	distribution := goptuna.DiscreteUniformDistribution{
		Low:  -10,
		High: 10,
		Q:    0.1,
	}

	points := make([]float64, 100)
	for i := 0; i < 100; i++ {
		trialID, err := study.Storage.CreateNewTrialID(study.ID)
		if err != nil {
			t.Errorf("should not be err, but got %s", err)
			return
		}
		trial, err := study.Storage.GetTrial(trialID)
		if err != nil {
			t.Errorf("should not be err, but got %s", err)
			return
		}
		sampled, err := study.Sampler.Sample(study, trial, "x", distribution)
		if err != nil {
			t.Errorf("should not be err, but got %s", err)
			return
		}
		if sampled < distribution.Low || sampled > distribution.High {
			t.Errorf("should not be less than %f, and larger than %f, but got %f",
				distribution.High, distribution.Low, sampled)
			return
		}
		points[i] = sampled
	}

	for i := range points {
		points[i] -= distribution.Low
		points[i] /= distribution.Q
		roundPoint := math.Round(points[i])
		if !testutil.AlmostEqualFloat64(roundPoint, points[i], 1e-6) {
			t.Errorf("should be almost the same, but got %f and %f",
				roundPoint, points[i])
			return
		}
	}
}
