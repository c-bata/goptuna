package main

import (
	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/cma"
	kurobako "github.com/sile/kurobako-go"
	"github.com/sile/kurobako-go/goptuna/solver"
)

func createStudy(seed int64) (*goptuna.Study, error) {
	relativeSampler := cma.NewSampler(cma.SamplerOptionSeed(seed))
	return goptuna.CreateStudy("example-study",
		goptuna.StudyOptionRelativeSampler(relativeSampler))
}

func main() {
	factory := solver.NewGoptunaSolverFactory("Goptuna (CMA)", createStudy)
	runner := kurobako.NewSolverRunner(&factory)
	if err := runner.Run(); err != nil {
		panic(err)
	}
}
