package goptuna_test

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/internal/testutil"
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

func TestRandomSearchSampler_SampleLogUniform(t *testing.T) {
	sampler := goptuna.NewRandomSearchSampler()
	study, err := goptuna.CreateStudy("", goptuna.StudyOptionSampler(sampler))
	if err != nil {
		t.Errorf("should not be err, but got %s", err)
		return
	}

	distribution := goptuna.LogUniformDistribution{
		Low:  1e-7,
		High: 1,
	}

	points := make([]float64, 100)
	for i := 0; i < 100; i++ {
		trialID, err := study.Storage.CreateNewTrial(study.ID)
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
		if points[i] < distribution.Low {
			t.Errorf("should be higher than %f, but got %f",
				distribution.Low, points[i])
			return
		}
		if points[i] > distribution.High {
			t.Errorf("should be lower than %f, but got %f",
				distribution.High, points[i])
			return
		}
	}
}

func TestRandomSearchSampler_SampleDiscreteUniform(t *testing.T) {
	sampler := goptuna.NewRandomSearchSampler()
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
		trialID, err := study.Storage.CreateNewTrial(study.ID)
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

func TestIntersectionSearchSpace(t *testing.T) {
	tests := []struct {
		name         string
		study        func() *goptuna.Study
		expectedKeys []string
		want         map[string]interface{}
		wantErr      bool
	}{
		{
			name: "No trial",
			study: func() *goptuna.Study {
				study, err := goptuna.CreateStudy("sampler_test")
				if err != nil {
					panic(err)
				}
				return study
			},
			want:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "One trial",
			study: func() *goptuna.Study {
				study, err := goptuna.CreateStudy("sampler_test")
				if err != nil {
					panic(err)
				}

				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					x, _ := trial.SuggestInt("x", 0, 10)
					y, _ := trial.SuggestUniform("y", -3, 3)
					return float64(x) + y, nil
				}, 1); err != nil {
					panic(err)
				}
				return study
			},
			want: map[string]interface{}{
				"x": goptuna.IntUniformDistribution{
					High: 10,
					Low:  0,
				},
				"y": goptuna.UniformDistribution{
					High: 3,
					Low:  -3,
				},
			},
			wantErr: false,
		},
		{
			name: "Second trial (only 'y' parameter is suggested in this trial)",
			study: func() *goptuna.Study {
				study, err := goptuna.CreateStudy("sampler_test")
				if err != nil {
					panic(err)
				}

				// First Trial
				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					x, _ := trial.SuggestInt("x", 0, 10)
					y, _ := trial.SuggestUniform("y", -3, 3)
					return float64(x) + y, nil
				}, 1); err != nil {
					panic(err)
				}

				// Second Trial
				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					y, _ := trial.SuggestUniform("y", -3, 3)
					return y, nil
				}, 1); err != nil {
					panic(err)
				}
				return study
			},
			want: map[string]interface{}{
				"y": goptuna.UniformDistribution{
					High: 3,
					Low:  -3,
				},
			},
			wantErr: false,
		},
		{
			name: "Failed or pruned trials are not considered in the calculation of an intersection search space.",
			study: func() *goptuna.Study {
				study, err := goptuna.CreateStudy("sampler_test")
				if err != nil {
					panic(err)
				}

				// First Trial
				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					x, _ := trial.SuggestInt("x", 0, 10)
					y, _ := trial.SuggestUniform("y", -3, 3)
					return float64(x) + y, nil
				}, 1); err != nil {
					panic(err)
				}

				// Second Trial
				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					y, _ := trial.SuggestUniform("y", -3, 3)
					return y, nil
				}, 1); err != nil {
					panic(err)
				}

				// Failed trial (ignore error)
				_ = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					_, _ = trial.SuggestUniform("y", -3, 3)
					return 0.0, errors.New("something error")
				}, 1)
				// Pruned trial
				if err = study.Optimize(func(trial goptuna.Trial) (v float64, e error) {
					_, _ = trial.SuggestUniform("y", -3, 3)
					return 0.0, goptuna.ErrTrialPruned
				}, 1); err != nil {
					panic(err)
				}
				return study
			},
			want: map[string]interface{}{
				"y": goptuna.UniformDistribution{
					High: 3,
					Low:  -3,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := goptuna.IntersectionSearchSpace(tt.study())
			if (err != nil) != tt.wantErr {
				t.Errorf("IntersectionSearchSpace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("IntersectionSearchSpace() return %d items, but want %d", len(got), len(tt.want))
			}
			for key := range tt.want {
				if distribution, ok := got[key]; !ok {
					t.Errorf("IntersectionSearchSpace() should have %s key", key)
				} else if !reflect.DeepEqual(distribution, tt.want[key]) {
					t.Errorf("IntersectionSearchSpace()[%s] = %v, want %v", key, distribution, tt.want[key])
				}
			}
		})
	}
}
