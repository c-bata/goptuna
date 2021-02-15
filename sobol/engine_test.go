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

func BenchmarkFindRightmostZeroBit(b *testing.B) {
	for n := uint32(0); n < uint32(b.N); n++ {
		findRightmostZeroBit(n)
	}
}

func BenchmarkOldFindRightmostZeroBit(b *testing.B) {
	old := func(n uint32) uint32 {
		c := uint32(0)
		for n&(1<<c) != 0 {
			c++
		}
		return c + 1 // starts from 1
	}

	for n := uint32(0); n < uint32(b.N); n++ {
		old(n)
	}
}

func Test_getNumberOfSkippedPoints(t *testing.T) {
	tests := []struct {
		n        uint32
		expected uint32
	}{
		{
			n:        2,
			expected: 2,
		},
		{
			n:        7,
			expected: 4,
		},
		{
			n:        8,
			expected: 8,
		},
		{
			n:        10,
			expected: 8,
		},
		{
			n:        20,
			expected: 16,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("n=%d", tt.n), func(t *testing.T) {
			if got := getNumberOfSkippedPoints(tt.n); got != tt.expected {
				t.Errorf("skippedPoints(%d) = %v, expected %v", tt.n, got, tt.expected)
			}
		})
	}
}
