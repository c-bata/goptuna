package cma

import (
	"errors"
	"math"
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

type Solution struct {
	// X is a parameter transformed to N(m, σ^2 C) from Z.
	X mat.Matrix
	// Value represents an evaluation value.
	Value float64
}

type Optimizer struct {
	mean  []float64
	sigma float64
	c     mat.Symmetric

	dim     int
	mu      int
	muEff   float64
	popsize int
	cc      float64
	c1      float64
	cmu     float64
	cSigma  float64
	dSigma  float64
	cm      float64
	chiN    float64
	pSigma  []float64
	pc      []float64
	weights []float64

	bounds        mat.Matrix
	maxReSampling int

	rng *rand.Rand
	g   int
}

func NewOptimizer(mean []float64, sigma float64, opts ...OptimizerOption) (*Optimizer, error) {
	if sigma <= 0 {
		return nil, errors.New("sigma should be non-zero positive number")
	}
	dim := len(mean)
	popsize := 4 + int(math.Floor(3*math.Log(float64(dim))))
	mu := popsize / 2

	sumWeightsPrimeBeforeMu := 0.
	sumWeightsPrimeSquareBeforeMu := 0.
	sumWeightsPrimeAfterMu := 0.
	sumWeightsPrimeSquareAfterMu := 0.
	weightsPrime := make([]float64, popsize)
	weightsPrimePositiveSum := 0.0
	weightsPrimeNegativeSum := 0.0
	for i := 0; i < popsize; i++ {
		wp := math.Log((float64(popsize)+1)/2 - math.Log(float64(i+1)))
		weightsPrime[i] = wp

		if i < mu {
			sumWeightsPrimeBeforeMu += wp
			sumWeightsPrimeSquareBeforeMu += math.Pow(wp, 2)
		} else {
			sumWeightsPrimeAfterMu += weightsPrime[i]
			sumWeightsPrimeSquareAfterMu += math.Pow(wp, 2)
		}

		if wp > 0 {
			weightsPrimePositiveSum += wp
		} else {
			weightsPrimeNegativeSum -= wp
		}
	}
	muEff := math.Pow(sumWeightsPrimeBeforeMu, 2) / sumWeightsPrimeSquareBeforeMu
	muEffMinus := math.Pow(sumWeightsPrimeAfterMu, 2) / sumWeightsPrimeSquareAfterMu

	alphaCov := 2.0
	// learning rate for the rank-one update
	c1 := alphaCov / (math.Pow(float64(dim)+1.3, 2) + muEff)
	// learning rate for the rank-μ update
	cmu := math.Min(
		1-c1,
		alphaCov*(muEff-2+1/muEff)/(math.Pow(float64(dim+2), 2)+alphaCov*muEff/2),
	)
	if c1+cmu > 1 {
		return nil, errors.New("invalid learning rate for the rank-one and rank-μ update")
	}

	alphaMin := math.Min(
		1+c1/cmu,                   // α_μ-
		1+(2*muEffMinus)/(muEff+2), // α_μ_eff-
	)
	alphaMin = math.Min(alphaMin, (1-c1-cmu)/float64(dim)*cmu) // α_{pos_def}^{minus}

	weights := make([]float64, popsize)
	for i := 0; i < popsize; i++ {
		if weightsPrime[i] > 0 {
			weights[i] = 1 / weightsPrimePositiveSum * weightsPrime[i]
		} else {
			weights[i] = alphaMin / weightsPrimeNegativeSum * weightsPrime[i]
		}
	}
	cm := 1.0

	// learning rate for the cumulation for the step-size control (eq.55)
	cSigma := (muEff + 2) / (float64(dim) + muEff + 5)
	dSigma := 1 + 2*math.Max(0, math.Sqrt((muEff-1)/(float64(dim)+1))-1) + cSigma
	if cSigma >= 1 {
		return nil, errors.New("invalid learning rate for cumulation for the ste-size control")
	}

	// learning rate for cumulation for the rank-one update (eq.56)
	cc := (4 + muEff/float64(dim)) / (float64(dim) + 4 + 2*muEff/float64(dim))
	if cc > 1 {
		return nil, errors.New("invalid learning rate for cumulation for the rank-one update")
	}

	chiN := math.Sqrt(float64(dim)) * (1.0 - (1.0 / (4.0 * float64(dim))) + 1.0/(21.0*(math.Pow(float64(dim), 2))))

	cArray := make([]float64, dim)
	for i := 0; i < dim; i++ {
		cArray[i] = 1.0
	}

	cma := &Optimizer{
		mean:          mean,
		sigma:         sigma,
		c:             mat.NewDiagDense(dim, cArray),
		dim:           dim,
		mu:            mu,
		muEff:         muEff,
		cc:            cc,
		c1:            c1,
		cmu:           cmu,
		cSigma:        cSigma,
		dSigma:        dSigma,
		cm:            cm,
		chiN:          chiN,
		pSigma:        make([]float64, dim),
		pc:            make([]float64, dim),
		weights:       weights,
		bounds:        nil,
		maxReSampling: 100,
		rng:           rand.New(rand.NewSource(0)),
		g:             0,
	}

	for _, opt := range opts {
		opt(cma)
	}
	return cma, nil
}

// Generation is monotonically increased when a multivariate normal distribution is updated.
func (c *Optimizer) Generation() int {
	return c.g
}

// PopulationSize returns the population size
func (c *Optimizer) PopulationSize() int {
	return c.popsize
}

func (c *Optimizer) Ask() (mat.Matrix, error) {
	for i := 0; i < c.maxReSampling; i++ {
		x, err := c.sampleSolution()
		if err != nil {
			return nil, err
		}
		if c.isFeasible(x) {
			return x, nil
		}
	}
	panic("implement me")
	return nil, nil
}

func (c *Optimizer) isFeasible(matrix mat.Matrix) bool {
	if c.bounds == nil {
		return true
	}
	panic("implement me")
	return true
}

func (c *Optimizer) sampleSolution() (x mat.Matrix, err error) {
	var eigsym mat.EigenSym
	ok := eigsym.Factorize(c.c, true)
	if !ok {
		return nil, errors.New("symmetric eigendecomposition failed")
	}

	d2 := eigsym.Values(nil)
	d := make([]float64, len(d2))
	for i := 0; i < len(d2); i++ {
		d[i] = math.Sqrt(d2[i])
	}

	var b mat.Dense
	eigsym.VectorsTo(&b)

	z := make([]float64, c.dim)
	for i := 0; i < c.dim; i++ {
		z[i] = c.rng.NormFloat64()
	}

	var bd mat.Dense
	bd.Mul(&b, mat.NewDiagDense(c.dim, d))
	var y mat.Dense
	y.Mul(&bd, mat.NewVecDense(c.dim, z))
	var x mat.VecDense

	return c.popsize
}

func (c *Optimizer) Tell(solutions []*Solution) error {
	panic("implement me")
	return nil
}
