package tpe

import (
	"math"
	"math/rand"
	"testing"

	"github.com/c-bata/goptuna"
)

func almostEqualFloat64(a, b float64, e float64) bool {
	if a+e > b && a-e < b {
		return true
	}
	return false
}

func almostEqualFloat641D(a, b []float64, e float64) bool {
	for i := range a {
		if !almostEqualFloat64(a[i], b[i], e) {
			return false
		}
	}
	return true
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
		configIdxs []int
		configVals []float64
		lossIdxs   []int
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
				configIdxs: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				configVals: []float64{7.515720606531342, 5.350185623031333, 5.124041307972975, 1.4089387361626944, -2.895952062621281, -8.814621912214118, 7.603846274084024, 5.915757103674883, 8.364607575197955, 1.4694727910185534},
				lossIdxs:   []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
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
			gotBelow, gotAbove := s.splitObservationPairs(tt.args.configIdxs, tt.args.configVals, tt.args.lossIdxs, tt.args.lossVals)
			if !almostEqualFloat641D(gotBelow, tt.wantBelow, 1e-4) {
				t.Errorf("Sampler.splitObservationPairs() gotBelow = %v, want %v", gotBelow, tt.wantBelow)
			}
			if !almostEqualFloat641D(gotAbove, tt.wantAbove, 1e-4) {
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
