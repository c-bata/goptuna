package tpe

import (
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
