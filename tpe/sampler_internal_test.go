package tpe

import (
	"errors"
	"math"
	"math/rand"
	"reflect"
	"testing"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/internal/testutil"
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

func TestSampler_splitObservationPairs(t *testing.T) {
	type fields struct {
		NStartupTrials        int
		NEICandidates         int
		Gamma                 FuncGamma
		ParzenEstimatorParams ParzenEstimatorParams
		rng                   *rand.Rand
		randomSampler         *goptuna.RandomSearchSampler
	}
	type args struct {
		configVals []float64
		lossVals   [][2]float64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantBelow []float64
		wantAbove []float64
	}{
		{
			name: "test case 1",
			fields: fields{
				Gamma: DefaultGamma,
			},
			args: args{
				configVals: []float64{7.515720606531342, 5.350185623031333, 5.124041307972975, 1.4089387361626944, -2.895952062621281, -8.814621912214118, 7.603846274084024, 5.915757103674883, 8.364607575197955, 1.4694727910185534},
				lossVals: [][2]float64{
					{math.Inf(-1), 51.07650573447907},
					{math.Inf(-1), 100.79007507622603},
					{math.Inf(-1), 20.712990047058412},
					{math.Inf(-1), 142.49871053544777},
					{math.Inf(-1), 61.74467260557292},
					{math.Inf(-1), 116.44303200021926},
					{math.Inf(-1), 132.8075795417795},
					{math.Inf(-1), 25.243709057350483},
					{math.Inf(-1), 141.5303287376019},
					{math.Inf(-1), 33.64359889992425},
				},
			},
			wantBelow: []float64{5.12404131},
			wantAbove: []float64{
				7.51572061, 5.35018562, 1.40893874, -2.89595206, -8.81462191,
				7.60384627, 5.9157571, 8.36460758, 1.46947279},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sampler{
				numberOfStartupTrials: tt.fields.NStartupTrials,
				numberOfEICandidates:  tt.fields.NEICandidates,
				gamma:                 tt.fields.Gamma,
				params:                tt.fields.ParzenEstimatorParams,
				rng:                   tt.fields.rng,
				randomSampler:         tt.fields.randomSampler,
			}
			gotBelow, gotAbove := s.splitObservationPairs(tt.args.configVals, tt.args.lossVals)
			if !testutil.AlmostEqualFloat641D(gotBelow, tt.wantBelow, 1e-6) {
				t.Errorf("Sampler.splitObservationPairs() gotBelow = %v, want %v", gotBelow, tt.wantBelow)
			}
			if !testutil.AlmostEqualFloat641D(gotAbove, tt.wantAbove, 1e-6) {
				t.Errorf("Sampler.splitObservationPairs() gotAbove = %v, want %v", gotAbove, tt.wantAbove)
			}
		})
	}
}

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
