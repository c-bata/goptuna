package main

import (
	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/tpe"
	kurobako "github.com/sile/kurobako-go"
	"github.com/sile/kurobako-go/goptuna/solver"
)

func createStudy(seed int64) (*goptuna.Study, error) {
	sampler := tpe.NewSampler(tpe.SamplerOptionSeed(seed))
	return goptuna.CreateStudy("example-study", goptuna.StudyOptionSampler(sampler))
}

func main() {
	factory := solver.NewGoptunaSolverFactory("Goptuna (TPE)", createStudy)
	runner := kurobako.NewSolverRunner(&factory)
	if err := runner.Run(); err != nil {
		panic(err)
	}
}
