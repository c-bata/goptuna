package medianstopping

import "testing"

func almostEqualFloat64(a, b float64, e float64) bool {
	if a+e > b && a-e < b {
		return true
	}
	return false
}

func Test_percentile(t *testing.T) {
	type args struct {
		a []float64
		q float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "0",
			args: args{
				a: []float64{0, 2, 4, 6, 8, 10},
				q: 20.0,
			},
			want: 2,
		},
		{
			name: "1",
			args: args{
				a: []float64{0, 2, 4, 8},
				q: 50.0,
			},
			want: 3,
		},
		{
			name: "2",
			args: args{
				a: []float64{0, 2, 6, 8},
				q: 25.0,
			},
			want: 1.5,
		},
		{
			name: "3",
			args: args{
				a: []float64{0, 2, 6, 8},
				q: 20.0,
			},
			want: 1.20000,
		},
		{
			name: "4",
			args: args{
				a: []float64{0, 3, 6, 8},
				q: 25.0,
			},
			want: 2.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := percentile(tt.args.a, tt.args.q); !almostEqualFloat64(got, tt.want, 1e-5) {
				t.Errorf("percentile() = %v, want %v", got, tt.want)
			}
		})
	}
}
