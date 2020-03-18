package main

import (
	"fmt"
	"math"

	"github.com/c-bata/goptuna/cmaes"
)

func objective(x1, x2 float64) float64 {
	return math.Pow(x1-3, 2) + math.Pow(10*(x2+2), 2)
}

func main() {
	mean := []float64{1, 2}
	sigma0 := 1.3
	optimizer, err := cmaes.NewOptimizer(
		mean, sigma0,
		cmaes.OptimizerOptionSeed(0),
	)
	if err != nil {
		panic(err)
	}

	solutions := make([]*cmaes.Solution, optimizer.PopulationSize())
	for generation := 0; generation < 50; generation++ {
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
			fmt.Printf("generation %d: %f (x1=%f, x2=%f)\n",
				generation, v, x1, x2)
		}

		err = optimizer.Tell(solutions)
		if err != nil {
			panic(err)
		}
	}
}
