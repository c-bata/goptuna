package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/c-bata/goptuna/cmaes"
	"gonum.org/v1/gonum/mat"
)

const (
	popTypeSmall = "small"
	popTypeLarge = "large"
)

func objective(x1, x2 float64) float64 {
	// Ackley 2D: https://www.sfu.ca/~ssurjano/ackley.html
	v := -20 * math.Exp(-0.2*math.Sqrt(0.5*(math.Pow(x1, 2)+math.Pow(x2, 2))))
	v -= math.Exp(0.5 * (math.Cos(2*math.Pi*x1) + math.Cos(2*math.Pi*x2)))
	v += math.E + 20
	return v
}

func main() {
	seed := int64(0)
	rng := rand.New(rand.NewSource(seed))

	bounds := mat.NewDense(2, 2, []float64{-32.768, 32.768, -32.768, 32.768})
	sigma := 32.768 * 2 / 5 // 1/5 of the domain width
	mean := []float64{
		-32.768 + (rng.Float64() * 32.768 * 2),
		-32.768 + (rng.Float64() * 32.768 * 2),
	}

	optimizer, err := cmaes.NewOptimizer(
		mean, sigma,
		cmaes.OptimizerOptionSeed(seed),
		cmaes.OptimizerOptionBounds(bounds),
	)
	if err != nil {
		panic(err)
	}

	// BIPOP-related variables
	nRestarts := 0
	smallEvaluations := 0
	largeEvaluations := 0
	popsize0 := optimizer.PopulationSize()
	incPopSize := 2

	// Initial run is with "normal" population size; it is
	// the large population before first doubling, but its
	// budget accounting is the same as in case of small
	// population.
	poptype := popTypeSmall

	solutions := make([]*cmaes.Solution, optimizer.PopulationSize())
	for nRestarts <= 5 {
		for i := 0; i < optimizer.PopulationSize(); i++ {
			x, err := optimizer.Ask()
			if err != nil {
				panic(err)
			}
			x1, x2 := x[0], x[1]
			v := objective(x1, x2)
			solutions[i] = &cmaes.Solution{
				Params: x,
				Value:  v,
			}
			// fmt.Printf("f = %f (x1=%f, x2=%f)\n", v, x1, x2)
		}

		err = optimizer.Tell(solutions)
		if err != nil {
			panic(err)
		}

		if optimizer.ShouldStop() {
			seed++
			nEvaluations := optimizer.PopulationSize() * optimizer.Generation()
			if poptype == popTypeSmall {
				smallEvaluations += nEvaluations
			} else { // popTypeLarge
				largeEvaluations += nEvaluations
			}

			var popsize int
			if smallEvaluations < largeEvaluations {
				poptype = popTypeSmall
				popsizeMultiplier := math.Pow(float64(incPopSize), float64(nRestarts))
				r := math.Pow(rng.Float64(), 2)
				popsize = int(math.Floor(float64(popsize0) * math.Pow(popsizeMultiplier, r)))
			} else {
				poptype = popTypeLarge
				nRestarts++
				popsize = popsize0 * int(math.Pow(float64(incPopSize), float64(nRestarts)))
			}

			mean = []float64{
				-32.768 + (rng.Float64() * 32.768 * 2),
				-32.768 + (rng.Float64() * 32.768 * 2),
			}
			optimizer, err = cmaes.NewOptimizer(
				mean, sigma,
				cmaes.OptimizerOptionSeed(seed),
				cmaes.OptimizerOptionPopulationSize(popsize))
			if err != nil {
				panic(fmt.Errorf("failed to restart CMA-ES with popsize=%d", popsize))
			}
			solutions = make([]*cmaes.Solution, popsize)
			fmt.Printf("Restart CMA-ES with popsize=%d (%s)\n", popsize, poptype)
		}
	}
}
