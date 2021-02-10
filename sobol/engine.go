package sobol

import (
	"fmt"
	"math"
)

// findRightmostZeroBit returns index from the right of the first zero bit of n.
func findRightmostZeroBit(n uint32) uint32 {
	c := uint32(0)
	for n&(1<<c) != 0 {
		c++
	}
	return c + 1 // starts from 1
}

func getNumberOfSkippedPoints(n uint32) uint32 {
	// In the 3-page notes of Joe & Kuo, they said:
	//
	// > It has been recommended by some that the Sobol0 sequence tends to perform better
	// > if an initial portion of the sequence is dropped: the number of points skipped is
	// > the largest power of 2 smaller than the number of points to be used.
	// > However, we are less persuaded by such recommendation ourselves.
	cnt := uint32(0)
	for {
		n >>= 1
		if n == 0 {
			break
		}
		cnt++
	}
	return uint32(math.Pow(2, float64(cnt)))
}

func initDirectionNumbers(dim uint32) [][]uint32 {
	v := make([][]uint32, dim)
	for i := uint32(0); i < dim; i++ {
		v[i] = make([]uint32, maxBit)
	}

	// First row of sobol state is all '1'.
	for m := 0; m < maxBit; m++ {
		v[0][m] = 1 << (32 - m) // all m's = 1
	}

	// Remaining rows of sobol state (row 2 through dim, indexed by [1:dim])
	for j := uint32(1); j < dim; j++ {
		v[j] = make([]uint32, maxBit+1)

		// Read in parameters from file
		dn := directionNumbers[j]
		m := make([]uint32, len(dn.M)+1)
		for i := uint32(0); i < dn.S; i++ {
			m[i+1] = dn.M[i]
		}
		for i := uint32(1); i <= dn.S; i++ {
			v[j][i] = m[i] << (32 - i)
		}
		for i := dn.S + 1; i <= maxBit; i++ {
			v[j][i] = v[j][i-dn.S] ^ (v[j][i-dn.S] >> dn.S)
			for k := uint32(1); k <= dn.S-1; k++ {
				v[j][i] ^= ((dn.A >> (dn.S - 1 - k)) & 1) * v[j][i-k]
			}
		}

	}
	return v
}

// Engine is Sobol's quasirandom number generator.
type Engine struct {
	dim uint32     // dimensions
	n   uint32     // the number of generate times
	v   [][]uint32 // direction numbers
	x   [][]uint32
}

// NewEngine returns Sobol's quasirandom number generator.
func NewEngine(dimension uint32) *Engine {
	if dimension > maxDim {
		panic(fmt.Errorf("maximum supported dimensionality is %d", maxDim))
	}

	v := initDirectionNumbers(dimension)
	x := make([][]uint32, dimension+1)
	for i := uint32(0); i <= dimension; i++ {
		// Pre-allocate memory to sample 512 points
		x[i] = make([]uint32, 0, 512)
		x[i] = append(x[i], 0)
	}

	return &Engine{
		dim: dimension,
		n:   0,
		v:   v,
		x:   x,
	}
}

// Draw samples from Sobol sequence.
func (e *Engine) Draw() []float64 {
	e.n++
	points := make([]float64, e.dim)

	for j := uint32(0); j < e.dim; j++ {
		c := findRightmostZeroBit(e.n - 1)
		e.x[j] = append(e.x[j], e.x[j][e.n-1]^e.v[j][c])
		points[j] = float64(e.x[j][e.n]) / math.Pow(2.0, 32)
	}
	return points
}
