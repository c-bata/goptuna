package testutil

// AlmostEqualFloat64 returns true if given float64 values are almost equals.
func AlmostEqualFloat64(a, b float64, e float64) bool {
	if a+e > b && a-e < b {
		return true
	}
	return false
}

// AlmostEqualFloat641D returns true if given float64 array are almost equals.
func AlmostEqualFloat641D(a, b []float64, e float64) bool {
	for i := range a {
		if !AlmostEqualFloat64(a[i], b[i], e) {
			return false
		}
	}
	return true
}
