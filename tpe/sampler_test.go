package tpe_test

import (
	"reflect"
	"testing"

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
			name: "test case 1",
			args: args{x: 5},
			want: 25,
		},
		{
			name: "test case 2",
			args: args{x: 255},
			want: 26,
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
