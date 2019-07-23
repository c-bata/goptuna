package tpe

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
