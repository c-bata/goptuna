package medianstopping_test

import (
	"reflect"
	"testing"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/medianstopping"
)

func TestPercentilePruner_PruneWithOneTrial(t *testing.T) {
	study, err := goptuna.CreateStudy("")
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	trialID, err := study.Storage.CreateNewTrialID(study.ID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}

	trial := goptuna.Trial{
		Study: study,
		ID:    trialID,
	}
	err = trial.Report(1, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}

	pruner := &medianstopping.PercentilePruner{
		Percentile:     25.0,
		NStartUpTrials: 0,
		NWarmUpSteps:   0,
	}
	ft, err := study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	prune, err := pruner.Prune(study, ft, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}

	// pruner not activated at a first trial
	if prune {
		t.Errorf("should be false, but got true")
	}
}

func TestNewPercentilePruner(t *testing.T) {
	type args struct {
		q float64
	}
	tests := []struct {
		name    string
		args    args
		want    *medianstopping.PercentilePruner
		wantErr bool
	}{
		{
			name: "under bound",
			args: args{
				q: -1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "on under bound",
			args: args{
				q: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "in range",
			args: args{
				q: 25.0,
			},
			want: &medianstopping.PercentilePruner{
				Percentile:     25.0,
				NStartUpTrials: 5,
				NWarmUpSteps:   0,
			},
			wantErr: false,
		},
		{
			name: "on upper bound",
			args: args{
				q: 100,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "over upper bound",
			args: args{
				q: 101,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := medianstopping.NewPercentilePruner(tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPercentilePruner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPercentilePruner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPercentilePruner_Prune(t *testing.T) {
	type fields struct {
		Percentile     float64
		NStartUpTrials int
		NWarmUpSteps   int
	}
	tests := []struct {
		name               string
		direction          goptuna.StudyDirection
		intermediateValues []float64
		latestValue        float64
		fields             fields
		step               int
		want               bool
	}{
		{
			name:               "minimize",
			direction:          goptuna.StudyDirectionMinimize,
			intermediateValues: []float64{1, 2, 3, 4, 5},
			latestValue:        2.1,
			fields: fields{
				Percentile:     25.0,
				NStartUpTrials: 0,
				NWarmUpSteps:   0,
			},
			want: true,
		},
		{
			name:               "maximize",
			direction:          goptuna.StudyDirectionMaximize,
			intermediateValues: []float64{1, 2, 3, 4, 5},
			latestValue:        3.9,
			fields: fields{
				Percentile:     25.0,
				NStartUpTrials: 0,
				NWarmUpSteps:   0,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			study, err := goptuna.CreateStudy("", goptuna.StudyOptionSetDirection(tt.direction))
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
				return
			}
			for _, v := range tt.intermediateValues {
				trialID, err := study.Storage.CreateNewTrialID(study.ID)
				if err != nil {
					t.Errorf("should be err=nil, but got %s", err)
					return
				}
				trial := goptuna.Trial{
					Study: study,
					ID:    trialID,
				}
				err = trial.Report(v, 1)
				if err != nil {
					t.Errorf("should be err=nil, but got %s", err)
					return
				}
				err = study.Storage.SetTrialState(trialID, goptuna.TrialStateComplete)
				if err != nil {
					t.Errorf("should be err=nil, but got %s", err)
					return
				}
			}

			p := &medianstopping.PercentilePruner{
				Percentile:     tt.fields.Percentile,
				NStartUpTrials: tt.fields.NStartUpTrials,
				NWarmUpSteps:   tt.fields.NWarmUpSteps,
			}

			// A pruner is not activated if a trial has no intermediate values.
			trialID, err := study.Storage.CreateNewTrialID(study.ID)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
				return
			}
			trial := goptuna.Trial{
				Study: study,
				ID:    trialID,
			}
			ft, err := study.Storage.GetTrial(trialID)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
				return
			}
			prune, err := p.Prune(study, ft, 1)
			if err != nil || prune {
				t.Errorf("A pruner is not activated if a trial has no intermediate values., %v %s", prune, err)
				return
			}

			// A pruner is activated if a trial has an intermediate value.
			err = trial.Report(tt.latestValue, 1)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
				return
			}

			ft, err = study.Storage.GetTrial(trialID)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
				return
			}
			got, err := p.Prune(study, ft, 1)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
				return
			}
			if got != tt.want {
				t.Errorf("PercentilePruner.Prune() = %v, want %v", got, tt.want)
			}
		})
	}
}
