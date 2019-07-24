package tpe

import (
	"errors"
	"math/rand"
	"sort"
	"time"
)

func ones(size int) []float64 {
	ones := make([]float64, size)
	for i := 0; i < size; i++ {
		ones[i] = 1
	}
	return ones
}

func linspace(start, stop float64, num int, endPoint bool) []float64 {
	step := 0.
	if endPoint {
		if num == 1 {
			return []float64{start}
		}
		step = (stop - start) / float64(num-1)
	} else {
		if num == 0 {
			return []float64{}
		}
		step = (stop - start) / float64(num)
	}
	r := make([]float64, num, num)
	for i := 0; i < num; i++ {
		r[i] = start + float64(i)*step
	}
	return r
}

func choice(array []float64, idxs []int) []float64 {
	results := make([]float64, len(idxs))
	for i, idx := range idxs {
		results[i] = array[idx]
	}
	return results
}

func location(array []float64, key float64) int {
	i := 0
	size := len(array)
	for {
		mid := (i + size) / 2
		if i == size {
			break
		}
		if array[mid] < key {
			i = mid + 1
		} else {
			size = mid
		}
	}
	return i
}

func Searchsorted(array, values []float64) []int {
	var indexes []int
	for _, val := range values {
		indexes = append(indexes, location(array, val))
	}
	return indexes
}

func randomWeightedSelect(weights []int, totalWeight int) (int, error) {
	// https://medium.com/@peterkellyonline/weighted-random-selection-3ff222917eb6
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(totalWeight)
	for i, g := range weights {
		r -= g
		if r <= 0 {
			return i, nil
		}
	}
	return 0, errors.New("no item selected")
}

func argMaxApproxMultinomial(pvals []float64, precision float64) (int, error) {
	tw := 0
	weights := make([]int, len(pvals))
	for i := range weights {
		w := int(pvals[i] / precision)
		tw += w
		weights[i] = w
	}
	return randomWeightedSelect(weights, tw)
}

func clip(array []float64, min, max float64) {
	for i := range array {
		if array[i] < min {
			array[i] = min
		} else if array[i] > max {
			array[i] = max
		}
	}
}

func ArgSort2DFloat64(lossVals [][2]float64) []int {
	type sortable struct {
		index   int
		lossVal [2]float64
	}
	x := make([]sortable, len(lossVals))
	for i := 0; i < len(lossVals); i++ {
		x[i] = sortable{
			index:   i,
			lossVal: lossVals[i],
		}
	}

	sort.SliceStable(x, func(i, j int) bool {
		if x[i].lossVal[0] == x[j].lossVal[0] {
			return x[i].lossVal[1] < x[j].lossVal[1]
		}
		return x[i].lossVal[0] < x[j].lossVal[0]
	})

	results := make([]int, len(x))
	for i := 0; i < len(x); i++ {
		results[i] = x[i].index
	}
	return results
}
