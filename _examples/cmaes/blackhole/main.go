package main

import (
	"log"
	"math"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/cmaes"
)

func objective(trial goptuna.Trial) (float64, error) {
	x1, err := trial.SuggestFloat("x1", -10, 10)
	if err != nil {
		return -1, err
	}
	x2, err := trial.SuggestFloat("x2", -10, 10)
	if err != nil {
		return -1, err
	}
	return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
}

func main() {
	relativeSampler := cmaes.NewSampler(
		cmaes.SamplerOptionNStartupTrials(0))
	study, err := goptuna.CreateStudy(
		"goptuna-example",
		goptuna.StudyOptionStorage(goptuna.NewBlackHoleStorage(20)),
		goptuna.StudyOptionRelativeSampler(relativeSampler),
		goptuna.StudyOptionDefineSearchSpace(map[string]interface{}{
			"x1": goptuna.UniformDistribution{
				High: 10,
				Low:  -10,
			},
			"x2": goptuna.UniformDistribution{
				High: 10,
				Low:  -10,
			},
		}),
	)
	if err != nil {
		log.Fatal("failed to create study:", err)
	}

	if err = study.Optimize(objective, 10000); err != nil {
		log.Fatal("failed to optimize:", err)
	}

	v, err := study.GetBestValue()
	if err != nil {
		log.Fatal("failed to get best value:", err)
	}
	params, err := study.GetBestParams()
	if err != nil {
		log.Fatal("failed to get best params:", err)
	}
	log.Printf("Best evaluation=%f (x1=%f, x2=%f)",
		v, params["x1"].(float64), params["x2"].(float64))
}
