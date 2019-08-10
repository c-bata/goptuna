package stats

import (
	"math"
	"sort"
)

func Percentile(a []float64, q float64) float64 {
	// Caution: this modifies the order of given first argument.
	length := len(a)
	if length == 0 {
		return math.NaN()
	}

	if length == 1 {
		return a[0]
	}

	if q <= 0 || q > 100 {
		return math.NaN()
	}

	sort.Float64s(a)
	index := float64(length-1) * (q / 100)
	if index == float64(int64(index)) {
		return a[int(index)]
	}

	i := int(math.Floor(index))
	x := 100 / float64(length-1)
	y := (a[i+1] - a[i]) * (x*float64(i+1) - q) / (x*float64(i+1) - x*float64(i))
	return a[i+1] - y
}
