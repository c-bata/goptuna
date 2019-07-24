package tpe

import (
	"math"
	"reflect"
	"testing"
)

func Test_linspace(t *testing.T) {
	type args struct {
		start    float64
		stop     float64
		num      int
		endPoint bool
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "test case 1",
			args: args{
				start:    1.0 / 30,
				stop:     1.0,
				num:      30 - 25,
				endPoint: true,
			},
			want: []float64{
				(1.0-1.0/30)*0/4 + 1.0/30,
				(1.0-1.0/30)*1/4 + 1.0/30,
				(1.0-1.0/30)*2/4 + 1.0/30,
				(1.0-1.0/30)*3/4 + 1.0/30,
				(1.0-1.0/30)*4/4 + 1.0/30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := linspace(tt.args.start, tt.args.stop, tt.args.num, tt.args.endPoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("linspace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgSort2DFloat64(t *testing.T) {
	type args struct {
		lossVals [][2]float64
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "test case 1",
			args: args{
				lossVals: [][2]float64{
					{math.Inf(-1), 93.80856756}, {math.Inf(-1), 85.64538195},
					{math.Inf(-1), 44.58783514}, {math.Inf(-1), 4.23458368},
					{math.Inf(-1), 42.17125041}, {math.Inf(-1), 62.14283937},
					{math.Inf(-1), 94.45778947}, {math.Inf(-1), 64.66469149},
					{math.Inf(-1), 36.1033201}, {math.Inf(-1), 105.69868952},
				},
			},
			want: []int{3, 8, 4, 2, 5, 7, 1, 0, 6, 9},
		},
		{
			name: "test case 2",
			args: args{
				lossVals: [][2]float64{
					{3.0, 93.80856756}, {5.0, 85.64538195},
					{math.Inf(-1), 44.58783514}, {math.Inf(-1), 4.23458368},
					{math.Inf(-1), 42.17125041}, {math.Inf(-1), 62.14283937},
					{math.Inf(-1), 94.45778947}, {math.Inf(-1), 64.66469149},
					{math.Inf(-1), 36.1033201}, {math.Inf(-1), 105.69868952},
				},
			},
			want: []int{3, 8, 4, 2, 5, 7, 6, 9, 0, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArgSort2DFloat64(tt.args.lossVals); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArgSort2DFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestApproximateMutlinomial(t *testing.T) {
//	for i := 0; i < 10; i++ {
//		fmt.Println(argMaxApproxMultinomial([]float64{0.1, 0.4, 0.5}, 0.001))
//	}
//}
