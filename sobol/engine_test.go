package sobol

import (
	"fmt"
	"testing"
)

func Test_findRightmostZeroBit(t *testing.T) {
	tests := []struct {
		n uint32
		c uint32
	}{
		{
			n: 0,
			c: 1,
		},
		{
			n: 1,
			c: 2,
		},
		{
			n: 2,
			c: 1,
		},
		{
			n: 3,
			c: 3,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("n=%d", tt.n), func(t *testing.T) {
			if got := findRightmostZeroBit(tt.n); got != tt.c {
				t.Errorf("findRightmostZeroBit() = %v, expected %v", got, tt.c)
			}
		})
	}
}
