package cma

func repeat(dst []float64, value float64) []float64 {
	for i := 0; i < len(dst); i++ {
		dst[i] = value
	}
	return dst
}
