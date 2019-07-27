package tpe_test

import (
	"reflect"
	"testing"

	"github.com/c-bata/goptuna/tpe"
)

func TestNewParzenEstimatorShapeCheck(t *testing.T) {
	type args struct {
		mus    []float64
		low    float64
		high   float64
		params tpe.ParzenEstimatorParams
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "buildEstimator shape check 1",
			args: args{
				mus:  []float64{},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     true,
					ConsiderMagicClip: true,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
		},
		{
			name: "buildEstimator shape check 1-1",
			args: args{
				mus:  []float64{},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     false,
					ConsiderMagicClip: true,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
		},
		{
			name: "buildEstimator shape check 1-2",
			args: args{
				mus:  []float64{},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     true,
					ConsiderMagicClip: false,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
		},
		{
			name: "buildEstimator shape check 1-3",
			args: args{
				mus:  []float64{},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     true,
					ConsiderMagicClip: true,
					ConsiderEndpoints: false,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
		},
		{
			name: "buildEstimator shape check 2",
			args: args{
				mus:  []float64{0.4},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     true,
					ConsiderMagicClip: true,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
		},
		{
			name: "buildEstimator shape check 2-1",
			args: args{
				mus:  []float64{0.4},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     false,
					ConsiderMagicClip: true,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
		},
		{
			name: "buildEstimator shape check 2-2",
			args: args{
				mus:  []float64{0.4},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     true,
					ConsiderMagicClip: false,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
		},
		{
			name: "buildEstimator shape check 2-3",
			args: args{
				mus:  []float64{0.4},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     true,
					ConsiderMagicClip: true,
					ConsiderEndpoints: false,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
		},
		{
			name: "buildEstimator shape check 3",
			args: args{
				mus:  []float64{-0.4, 0.4},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     true,
					ConsiderMagicClip: true,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimator := tpe.NewParzenEstimator(tt.args.mus, tt.args.low, tt.args.high, tt.args.params)

			actual := len(estimator.Weights)
			expected := len(tt.args.mus)
			if tt.args.params.ConsiderPrior {
				expected += 1
			}
			if actual != expected {
				t.Errorf("length of NewParzenEstimator().Weights = %d, want %v", actual, expected)
			}
		})
	}
}

func TestNewParzenEstimator(t *testing.T) {
	type args struct {
		mus    []float64
		low    float64
		high   float64
		params tpe.ParzenEstimatorParams
	}
	type expected struct {
		weights []float64
		mus     []float64
		sigmas  []float64
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "buildEstimator 1",
			args: args{
				mus:  []float64{},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     false,
					ConsiderMagicClip: false,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
			expected: expected{
				weights: []float64{},
				mus:     []float64{},
				sigmas:  []float64{},
			},
		},
		{
			name: "buildEstimator 2",
			args: args{
				mus:  []float64{},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     true,
					ConsiderMagicClip: false,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
			expected: expected{
				weights: []float64{1.0},
				mus:     []float64{0.0},
				sigmas:  []float64{2.0},
			},
		},
		{
			name: "buildEstimator 3",
			args: args{
				mus:  []float64{0.4},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     true,
					ConsiderMagicClip: false,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
			expected: expected{
				weights: []float64{0.5, 0.5},
				mus:     []float64{0.0, 0.4},
				sigmas:  []float64{2.0, 0.6},
			},
		},
		{
			name: "buildEstimator 4",
			args: args{
				mus:  []float64{-0.4},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     true,
					ConsiderMagicClip: false,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
			expected: expected{
				weights: []float64{0.5, 0.5},
				mus:     []float64{-0.4, 0.0},
				sigmas:  []float64{0.6, 2.0},
			},
		},
		{
			name: "buildEstimator 5",
			args: args{
				mus:  []float64{-0.4, 0.4},
				low:  -1.0,
				high: 1.0,
				params: tpe.ParzenEstimatorParams{
					ConsiderPrior:     true,
					ConsiderMagicClip: false,
					ConsiderEndpoints: true,
					Weights:           tpe.DefaultWeights,
					PriorWeight:       1.0,
				},
			},
			expected: expected{
				weights: []float64{1.0 / 3, 1.0 / 3, 1.0 / 3},
				mus:     []float64{-0.4, 0.0, 0.4},
				sigmas:  []float64{0.6, 2.0, 0.6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			estimator := tpe.NewParzenEstimator(tt.args.mus, tt.args.low, tt.args.high, tt.args.params)

			if !reflect.DeepEqual(estimator.Weights, tt.expected.weights) {
				t.Errorf("NewParzenEstimator() Weights = %v, want %v", estimator.Weights, tt.expected.weights)
			}
			if !reflect.DeepEqual(estimator.Mus, tt.expected.mus) {
				t.Errorf("NewParzenEstimator() Mus = %v, want %v", estimator.Mus, tt.expected.mus)
			}
			// to pass test case 0.
			if len(estimator.Sigmas) != 0 || len(tt.expected.sigmas) != 0 {
				if !reflect.DeepEqual(estimator.Sigmas, tt.expected.sigmas) {
					t.Errorf("NewParzenEstimator() Sigmas = %v, want %v", estimator.Sigmas, tt.expected.sigmas)
				}
			}
		})
	}
}
