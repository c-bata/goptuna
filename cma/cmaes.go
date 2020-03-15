package cma

import (
	"errors"
	"math"
	"math/rand"
	"sort"

	"gonum.org/v1/gonum/mat"
)

type Solution struct {
	// X is a parameter transformed to N(m, σ^2 C) from Z.
	X *mat.VecDense
	// Value represents an evaluation value.
	Value float64
}

// Optimizer is CMA-ES stochastic optimizer class with ask-and-tell interface.
type Optimizer struct {
	mean  *mat.VecDense
	sigma float64
	c     *mat.SymDense

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
	pSigma  *mat.VecDense
	pc      *mat.VecDense
	weights *mat.VecDense

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

	cma := &Optimizer{
		mean:          mat.NewVecDense(dim, mean),
		sigma:         sigma,
		c:             initC(dim),
		dim:           dim,
		popsize:       popsize,
		mu:            mu,
		muEff:         muEff,
		cc:            cc,
		c1:            c1,
		cmu:           cmu,
		cSigma:        cSigma,
		dSigma:        dSigma,
		cm:            cm,
		chiN:          chiN,
		pSigma:        mat.NewVecDense(dim, make([]float64, dim)),
		pc:            mat.NewVecDense(dim, make([]float64, dim)),
		weights:       mat.NewVecDense(popsize, weights),
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

// Generation is incremented when a multivariate normal distribution is updated.
func (c *Optimizer) Generation() int {
	return c.g
}

// PopulationSize returns the population size.
func (c *Optimizer) PopulationSize() int {
	return c.popsize
}

// Ask a next parameter.
func (c *Optimizer) Ask() (*mat.VecDense, error) {
	x, err := c.sampleSolution()
	if err != nil {
		return nil, err
	}
	for i := 0; i < c.maxReSampling; i++ {
		if c.isFeasible(x) {
			return x, nil
		}
		x, err = c.sampleSolution()
		if err != nil {
			return nil, err
		}
	}
	err = c.repairInfeasibleParams(x)
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (c *Optimizer) isFeasible(values *mat.VecDense) bool {
	if c.bounds == nil {
		return true
	}
	if values.Len() != c.dim {
		return false
	}
	for i := 0; i < c.dim; i++ {
		v := values.AtVec(i)
		if !(c.bounds.At(i, 0) < v && c.bounds.At(i, 1) > v) {
			return false
		}
	}
	return true
}

func (c *Optimizer) repairInfeasibleParams(values *mat.VecDense) error {
	if c.bounds == nil {
		return nil
	}
	if values.Len() != c.dim {
		return errors.New("invalid matrix size")
	}

	for i := 0; i < c.dim; i++ {
		v := values.AtVec(i)
		if c.bounds.At(i, 0) > v {
			values.SetVec(i, c.bounds.At(i, 0))
		}
		if c.bounds.At(i, 1) < v {
			values.SetVec(i, c.bounds.At(i, 1))
		}
	}
	return nil
}

func (c *Optimizer) sampleSolution() (*mat.VecDense, error) {
	// TODO(c-bata): Cache B and D
	var eigsym mat.EigenSym
	ok := eigsym.Factorize(c.c, true)
	if !ok {
		return nil, errors.New("symmetric eigendecomposition failed")
	}

	var b mat.Dense
	eigsym.VectorsTo(&b)
	d := make([]float64, c.dim)
	eigsym.Values(d) // d^2
	floatsSqrtTo(d)  // d

	z := make([]float64, c.dim)
	for i := 0; i < c.dim; i++ {
		z[i] = c.rng.NormFloat64()
	}

	var bd mat.Dense
	bd.Mul(&b, mat.NewDiagDense(c.dim, d))

	values := mat.NewVecDense(c.dim, z) // ~ N(0, I)
	values.MulVec(&bd, values)          // ~ N(0, C)
	values.ScaleVec(c.sigma, values)    // ~ N(0, σ^2 C)
	values.AddVec(values, c.mean)       // ~ N(m, σ^2 C)
	return values, nil
}

// Tell evaluation values.
func (c *Optimizer) Tell(solutions []*Solution) error {
	if len(solutions) != c.popsize {
		return errors.New("must tell popsize-length solutions")
	}

	c.g++
	sort.Slice(solutions, func(i, j int) bool {
		return solutions[i].Value < solutions[j].Value
	})

	var eigsym mat.EigenSym
	ok := eigsym.Factorize(c.c, true)
	if !ok {
		return errors.New("symmetric eigendecomposition failed")
	}

	var b mat.Dense
	eigsym.VectorsTo(&b)
	d := make([]float64, c.dim)
	eigsym.Values(d) // d^2
	floatsSqrtTo(d)  // d

	yk := solutionsToX(solutions) // ~ N(m, σ^2 C)
	meank := stackvec(c.popsize, c.dim, c.mean)
	yk.Sub(yk, meank)       // ~ N(0, σ^2 C)
	yk.Scale(1/c.sigma, yk) // ~ N(0, C)

	// Selection and recombination
	ydotw := mat.NewDense(c.mu, c.dim, nil)
	ydotw.Copy(yk.Slice(0, c.mu, 0, c.dim))
	weightsmu := stackvec(c.dim, c.mu, c.weights)
	ydotw.MulElem(ydotw, weightsmu.T())

	yw := sumColumns(ydotw.T())
	meandiff := mat.NewVecDense(c.dim, nil)
	meandiff.CopyVec(yw)
	meandiff.ScaleVec(c.cm*c.sigma, meandiff)
	c.mean.AddVec(c.mean, meandiff)

	// Step-size control
	dinv := mat.NewDiagDense(c.dim, arrinv(d))
	c2 := mat.NewDense(c.dim, c.dim, nil)
	c2.Product(&b, dinv, b.T()) // C^(-1/2) = B D^(-1) B^T

	c2yw := mat.NewDense(c.dim, 1, nil)
	c2yw.Product(c2, yw)
	c2yw.Scale(math.Sqrt(c.cSigma*(2-c.cSigma)*c.muEff), c2yw)
	c.pSigma.ScaleVec(1-c.cSigma, c.pSigma)
	c.pSigma.AddVec(c.pSigma, mat.NewVecDense(c.dim, c2yw.RawMatrix().Data))

	normPSigma := mat.Norm(c.pSigma, 2)
	c.sigma *= math.Exp((c.cSigma / c.dSigma) * (normPSigma/c.chiN - 1))

	hSigmaCondLeft := normPSigma / math.Sqrt(
		1-math.Pow(1-c.cSigma, float64(2*(c.g+1))))
	hSigmaCondRight := (1.4 + 2/float64(c.dim+1)) * c.chiN
	hSigma := 0.0
	if hSigmaCondLeft < hSigmaCondRight {
		hSigma = 1.0
	}

	// eq.45
	c.pc.ScaleVec(1-c.cc, c.pc)
	c.pc.AddScaledVec(c.pc, hSigma*math.Sqrt(c.cc*(2-c.cc)*c.muEff), yw)

	// eq.46
	wio := mat.NewVecDense(c.weights.Len(), nil)
	wio.CopyVec(c.weights)
	c2yk := mat.NewDense(c.dim, c.popsize, nil)
	c2yk.Product(c2, yk.T())
	wio.MulElemVec(wio, vecapply(c.weights, func(i int, a float64) float64 {
		if a > 0 {
			return 1.0
		}
		c2xinorm := mat.Norm(c2yk.ColView(i), 2)
		return float64(c.dim) / math.Pow(c2xinorm, 2)
	}))

	deltaHSigma := (1 - hSigma) * c.cc * (2 - c.cc)
	if deltaHSigma > 1 {
		panic("invalid delta_h_sigma")
	}

	// eq.47
	rankOne := mat.NewSymDense(c.dim, nil)
	rankOne.SymOuterK(1.0, c.pc)

	rankMu := mat.NewSymDense(c.dim, nil)
	for i := 0; i < c.popsize; i++ {
		wi := wio.AtVec(i)
		yi := yk.RowView(i)
		s := mat.NewSymDense(c.dim, nil)
		s.SymOuterK(wi, yi)
		rankMu.AddSym(rankMu, s)
	}

	c.c.ScaleSym(1+c.c1*deltaHSigma-c.c1-c.cmu*mat.Sum(c.weights), c.c)
	rankOne.ScaleSym(c.c1, rankOne)
	rankMu.ScaleSym(c.cmu, rankMu)
	c.c.AddSym(c.c, rankOne)
	c.c.AddSym(c.c, rankMu)

	// Avoid eigendecomposition error by arithmetic overflow
	c.c.AddSym(c.c, initMinC(c.dim))
	return nil
}
