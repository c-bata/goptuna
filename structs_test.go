package goptuna_test

import (
	"testing"
	"time"

	"github.com/c-bata/goptuna"
)

func TestFrozenTrial_Validate(t *testing.T) {
	tests := []struct {
		name    string
		trial   goptuna.FrozenTrial
		wantErr bool
	}{
		{
			name: "Valid case",
			trial: goptuna.FrozenTrial{
				DatetimeStart:    time.Now(),
				DatetimeComplete: time.Now(),
				State:            goptuna.TrialStateComplete,
				InternalParams: map[string]float64{
					"x1": 1,
				},
				Distributions: map[string]interface{}{
					"x1": goptuna.UniformDistribution{
						High: 10,
						Low:  0,
					},
				},
				Params: map[string]interface{}{
					"x1": 1.0,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid when DatetimeStart is zero.",
			trial: goptuna.FrozenTrial{
				DatetimeStart:    time.Time{},
				DatetimeComplete: time.Now(),
				State:            goptuna.TrialStateComplete,
				InternalParams: map[string]float64{
					"x1": 1,
				},
				Distributions: map[string]interface{}{
					"x1": goptuna.UniformDistribution{
						High: 10,
						Low:  0,
					},
				},
				Params: map[string]interface{}{
					"x1": 1.0,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid when DatetimeComplete is zero even if completed.",
			trial: goptuna.FrozenTrial{
				DatetimeStart:    time.Now(),
				DatetimeComplete: time.Time{},
				State:            goptuna.TrialStateComplete,
				InternalParams: map[string]float64{
					"x1": 1,
				},
				Distributions: map[string]interface{}{
					"x1": goptuna.UniformDistribution{
						High: 10,
						Low:  0,
					},
				},
				Params: map[string]interface{}{
					"x1": 1.0,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid when DatetimeComplete is set even if running.",
			trial: goptuna.FrozenTrial{
				DatetimeStart:    time.Now(),
				DatetimeComplete: time.Now(),
				State:            goptuna.TrialStateRunning,
				InternalParams: map[string]float64{
					"x1": 1,
				},
				Distributions: map[string]interface{}{
					"x1": goptuna.UniformDistribution{
						High: 10,
						Low:  0,
					},
				},
				Params: map[string]interface{}{
					"x1": 1.0,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid when Params and Distributions doesn't match",
			trial: goptuna.FrozenTrial{
				DatetimeStart:    time.Now(),
				DatetimeComplete: time.Now(),
				State:            goptuna.TrialStateComplete,
				InternalParams: map[string]float64{
					"x1": 1,
					"x2": 2,
				},
				Distributions: map[string]interface{}{
					"x1": goptuna.UniformDistribution{
						High: 10,
						Low:  0,
					},
				},
				Params: map[string]interface{}{
					"x1": 1,
					"x2": 2,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid external param value",
			trial: goptuna.FrozenTrial{
				DatetimeStart:    time.Now(),
				DatetimeComplete: time.Now(),
				State:            goptuna.TrialStateComplete,
				InternalParams: map[string]float64{
					"x1": 11,
				},
				Distributions: map[string]interface{}{
					"x1": goptuna.UniformDistribution{
						High: 100,
						Low:  0,
					},
				},
				Params: map[string]interface{}{
					"x1": 10.0,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid when the param is out of the distribution range",
			trial: goptuna.FrozenTrial{
				DatetimeStart:    time.Now(),
				DatetimeComplete: time.Now(),
				State:            goptuna.TrialStateComplete,
				InternalParams: map[string]float64{
					"x1": 100,
				},
				Distributions: map[string]interface{}{
					"x1": goptuna.UniformDistribution{
						High: 10,
						Low:  0,
					},
				},
				Params: map[string]interface{}{
					"x1": 100.0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := goptuna.ExportFrozenTrialValidate(&tt.trial)
			if !tt.wantErr && err != nil {
				t.Errorf("should retuurn nil, but got %s", err)
				return
			}
			if tt.wantErr && err == nil {
				t.Error("should return error, but got nil")
				return
			}
		})
	}
}
