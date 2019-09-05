package goptuna_test

import (
	"math"
	"testing"

	"github.com/c-bata/goptuna"
)

func TestTrial_Suggest(t *testing.T) {
	tests := []struct {
		name      string
		objective goptuna.FuncObjective
		wantErr   bool
	}{
		{
			name: "SuggestUniform",
			objective: func(trial goptuna.Trial) (float64, error) {
				// low is larger than high
				x1, err := trial.SuggestUniform("x", -10, 10)
				if err != nil {
					return -1, err
				}
				return math.Pow(x1-2, 2), nil
			},
			wantErr: false,
		},
		{
			name: "SuggestUniform: low is larger than high",
			objective: func(trial goptuna.Trial) (float64, error) {
				// low is larger than high
				x1, err := trial.SuggestUniform("x", 10, -10)
				if err != nil {
					return -1, err
				}
				return math.Pow(x1-2, 2), nil
			},
			wantErr: true,
		},
		{
			name: "SuggestLogUniform",
			objective: func(trial goptuna.Trial) (float64, error) {
				x1, err := trial.SuggestLogUniform("x", 1e5, 1e10)
				if err != nil {
					return -1, err
				}
				return math.Pow(x1-2, 2), nil
			},
			wantErr: false,
		},
		{
			name: "SuggestLogUniform: low is larger than high",
			objective: func(trial goptuna.Trial) (float64, error) {
				x1, err := trial.SuggestLogUniform("x", 1e10, 1e5)
				if err != nil {
					return -1, err
				}
				return math.Pow(x1-2, 2), nil
			},
			wantErr: true,
		},
		{
			name: "SuggestDiscreteUniform",
			objective: func(trial goptuna.Trial) (float64, error) {
				x1, err := trial.SuggestDiscreteUniform("x", -10, 10, 0.5)
				if err != nil {
					return -1, err
				}
				return math.Pow(x1-2, 2), nil
			},
			wantErr: false,
		},
		{
			name: "SuggestDiscreteUniform: low is larger than high",
			objective: func(trial goptuna.Trial) (float64, error) {
				x1, err := trial.SuggestDiscreteUniform("x", 10, -10, 0.5)
				if err != nil {
					return -1, err
				}
				return math.Pow(x1-2, 2), nil
			},
			wantErr: true,
		},
		{
			name: "SuggestInt",
			objective: func(trial goptuna.Trial) (float64, error) {
				x1, err := trial.SuggestInt("x", -10, 10)
				if err != nil {
					return -1, err
				}
				return math.Pow(float64(x1-2), 2), nil
			},
			wantErr: false,
		},
		{
			name: "SuggestInt: low is larger than high",
			objective: func(trial goptuna.Trial) (float64, error) {
				x1, err := trial.SuggestInt("x", 10, -10)
				if err != nil {
					return -1, err
				}
				return math.Pow(float64(x1-2), 2), nil
			},
			wantErr: true,
		},
		{
			name: "SuggestCategorical",
			objective: func(trial goptuna.Trial) (float64, error) {
				x1, err := trial.SuggestCategorical("x", []string{"foo", "bar", "baz"})
				if err != nil {
					return -1, err
				}
				if x1 == "foo" {
					return 0, nil
				}
				return 1, nil
			},
			wantErr: false,
		},
		{
			name: "SuggestCategorical: 'choices' must contains one or more elements",
			objective: func(trial goptuna.Trial) (float64, error) {
				x1, err := trial.SuggestCategorical("x", []string{})
				if err != nil {
					return -1, err
				}
				if x1 == "foo" {
					return 0, nil
				}
				return 1, nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler := goptuna.NewRandomSearchSampler(
				goptuna.RandomSearchSamplerOptionSeed(0),
			)
			study, err := goptuna.CreateStudy(tt.name,
				goptuna.StudyOptionIgnoreError(false),
				goptuna.StudyOptionSampler(sampler))

			err = study.Optimize(tt.objective, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Trial.SuggestUniform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTrial_UserAttrs(t *testing.T) {
	study, _ := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionStorage(goptuna.NewInMemoryStorage()),
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSampler(goptuna.NewRandomSearchSampler()),
	)
	trialID, err := study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}
	trial := goptuna.Trial{
		Study: study,
		ID:    trialID,
	}

	err = trial.SetUserAttr("hello", "world")
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}

	attrs, err := trial.GetUserAttrs()
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}

	hello, ok := attrs["hello"]
	if !ok {
		t.Errorf("'hello' doesn't exist in %#v", attrs)
		return
	}
	if hello != "world" {
		t.Errorf("should be 'world', but got '%s'", hello)
	}
}

func TestTrial_SystemAttrs(t *testing.T) {
	study, _ := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionStorage(goptuna.NewInMemoryStorage()),
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSampler(goptuna.NewRandomSearchSampler()),
	)
	trialID, err := study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}
	trial := goptuna.Trial{
		Study: study,
		ID:    trialID,
	}

	err = trial.SetSystemAttr("hello", "world")
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}

	attrs, err := trial.GetSystemAttrs()
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}

	hello, ok := attrs["hello"]
	if !ok {
		t.Errorf("'hello' doesn't exist in %#v", attrs)
		return
	}
	if hello != "world" {
		t.Errorf("should be 'world', but got '%s'", hello)
	}
}
