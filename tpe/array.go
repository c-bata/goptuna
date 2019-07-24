package tpe

import "sort"

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
