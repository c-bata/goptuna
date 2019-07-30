package tpe

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/c-bata/goptuna"
)

func TestGetObservationPairs_MINIMIZE(t *testing.T) {
	study, err := goptuna.CreateStudy(
		"", goptuna.StudyOptionIgnoreObjectiveErr(true),
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize))
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	err = study.Optimize(func(trial goptuna.Trial) (float64, error) {
		x, _ := trial.SuggestInt("x", 5, 5)
		number, _ := trial.Number()
		if number == 0 {
			return float64(x), nil
		} else if number == 1 {
			_ = trial.Report(1, 4)
			_ = trial.Report(2, 7)
			return 0.0, goptuna.ErrTrialPruned
		} else if number == 2 {
			return 0.0, goptuna.ErrTrialPruned
		} else {
			return 0.0, errors.New("runtime error")
		}
	}, 4)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	values, scores, err := getObservationPairs(study, "x")
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
	}

	expectedValues := []float64{5.0, 5.0, 5.0}
	if !reflect.DeepEqual(values, expectedValues) {
		t.Errorf("should be %v, but got %v", expectedValues, values)
	}
	expectedScores := [][2]float64{
		{math.Inf(-1), 5},
		{-7, 2},
		{math.Inf(1), 0},
	}
	if !reflect.DeepEqual(scores, expectedScores) {
		t.Errorf("should be %v, but got %v", expectedScores, scores)
	}
}

func TestGetObservationPairs_MAXIMIZE(t *testing.T) {
	study, err := goptuna.CreateStudy(
		"", goptuna.StudyOptionIgnoreObjectiveErr(true),
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMaximize))
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	err = study.Optimize(func(trial goptuna.Trial) (float64, error) {
		x, _ := trial.SuggestInt("x", 5, 5)
		number, _ := trial.Number()
		if number == 0 {
			return float64(x), nil
		} else if number == 1 {
			_ = trial.Report(1, 4)
			_ = trial.Report(2, 7)
			return 0.0, goptuna.ErrTrialPruned
		} else if number == 2 {
			return 0.0, goptuna.ErrTrialPruned
		} else {
			return 0.0, errors.New("runtime error")
		}
	}, 4)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	values, scores, err := getObservationPairs(study, "x")
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
	}

	expectedValues := []float64{5.0, 5.0, 5.0}
	if !reflect.DeepEqual(values, expectedValues) {
		t.Errorf("should be %v, but got %v", expectedValues, values)
	}
	expectedScores := [][2]float64{
		{math.Inf(-1), -5},
		{-7, -2},
		{math.Inf(1), 0},
	}
	if !reflect.DeepEqual(scores, expectedScores) {
		t.Errorf("should be %v, but got %v", expectedScores, scores)
	}
}

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
