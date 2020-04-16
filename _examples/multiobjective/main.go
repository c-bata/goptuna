package main

import (
	"log"
	"math"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/multiobjective"
)

func objective(trial goptuna.Trial) (float64, error) {
	x1, _ := trial.SuggestFloat("x1", -10, 10)
	x2, _ := trial.SuggestFloat("x2", -10, 10)
	f := math.Pow(x1-2, 2) + math.Pow(x2+5, 2)

	err := multiobjective.ReportSubMetrics(trial, map[string]float64{
		"f":  f,
		"x1": x1,
		"x2": x2,
	})
	return f, err
}

func main() {
	study, err := goptuna.CreateStudy(
		"goptuna-example",
	)
	if err != nil {
		log.Fatal("failed to create study:", err)
	}
	if err = study.Optimize(objective, 50); err != nil {
		log.Fatal("failed to optimize:", err)
	}

	paretoFronts, err := multiobjective.GetParetoOptimalTrials(study, map[string]goptuna.StudyDirection{
		"f":  goptuna.StudyDirectionMinimize,
		"x1": goptuna.StudyDirectionMinimize,
		"x2": goptuna.StudyDirectionMaximize,
	})
	if err != nil {
		log.Fatal("failed to get Pareto-optimal solutions:", err)
	}
	for _, pf := range paretoFronts {
		submetrics, err := multiobjective.GetSubMetrics(pf.SystemAttrs)
		if err != nil {
			log.Fatal("failed to get sub metrics:", err)
		}
		log.Printf("Pareto front: %v", submetrics)
	}
}
