package testutil

func AlmostEqualFloat64(a, b float64, e float64) bool {
	if a+e > b && a-e < b {
		return true
	}
	return false
}

func AlmostEqualFloat641D(a, b []float64, e float64) bool {
	for i := range a {
		if !AlmostEqualFloat64(a[i], b[i], e) {
			return false
		}
	}
	return true
}
