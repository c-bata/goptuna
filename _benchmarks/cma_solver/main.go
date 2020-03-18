package main

import (
	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/cmaes"
	kurobako "github.com/sile/kurobako-go"
	"github.com/sile/kurobako-go/goptuna/solver"
)

func createStudy(seed int64) (*goptuna.Study, error) {
	sampler := goptuna.NewRandomSearchSampler(goptuna.RandomSearchSamplerOptionSeed(seed))
	relativeSampler := cmaes.NewSampler(cmaes.SamplerOptionSeed(seed))
	return goptuna.CreateStudy("example-study",
		goptuna.StudyOptionSampler(sampler),
		goptuna.StudyOptionRelativeSampler(relativeSampler))
}

func main() {
	factory := solver.NewGoptunaSolverFactory("Goptuna (CMA-ES)", createStudy)
	runner := kurobako.NewSolverRunner(&factory)
	if err := runner.Run(); err != nil {
		panic(err)
	}
}
