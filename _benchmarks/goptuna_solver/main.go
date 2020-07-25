package main

import (
	"os"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/cmaes"
	"github.com/c-bata/goptuna/tpe"
	kurobako "github.com/sile/kurobako-go"
	"github.com/sile/kurobako-go/goptuna/solver"
)

func randomSampler(seed int64) (*goptuna.Study, error) {
	s := goptuna.NewRandomSampler(goptuna.RandomSamplerOptionSeed(seed))
	return goptuna.CreateStudy("example-study",
		goptuna.StudyOptionSampler(s))
}

func tpeSampler(seed int64) (*goptuna.Study, error) {
	s := tpe.NewSampler(tpe.SamplerOptionSeed(seed))
	return goptuna.CreateStudy("example-study",
		goptuna.StudyOptionSampler(s))
}

func cmaSampler(seed int64) (*goptuna.Study, error) {
	s := goptuna.NewRandomSampler(goptuna.RandomSamplerOptionSeed(seed))
	rs := cmaes.NewSampler(cmaes.SamplerOptionSeed(seed))
	return goptuna.CreateStudy("example-study",
		goptuna.StudyOptionSampler(s),
		goptuna.StudyOptionRelativeSampler(rs))
}

func ipopCmaSampler(seed int64) (*goptuna.Study, error) {
	s := goptuna.NewRandomSampler(goptuna.RandomSamplerOptionSeed(seed))
	rs := cmaes.NewSampler(cmaes.SamplerOptionSeed(seed),
		cmaes.SamplerOptionIPop(2))
	return goptuna.CreateStudy("example-study",
		goptuna.StudyOptionSampler(s),
		goptuna.StudyOptionRelativeSampler(rs))
}

func bipopCmaSampler(seed int64) (*goptuna.Study, error) {
	s := goptuna.NewRandomSampler(goptuna.RandomSamplerOptionSeed(seed))
	rs := cmaes.NewSampler(cmaes.SamplerOptionSeed(seed),
		cmaes.SamplerOptionBIPop(2))
	return goptuna.CreateStudy("example-study",
		goptuna.StudyOptionSampler(s),
		goptuna.StudyOptionRelativeSampler(rs))
}

func main() {
	if len(os.Args) != 2 {
		panic("please specify sampler algorithm")
	}

	var factory solver.GoptunaSolverFactory
	if sampler := os.Args[1]; sampler == "random" {
		factory = solver.NewGoptunaSolverFactory("Goptuna (Random search)", randomSampler)
	} else if sampler == "cmaes" {
		factory = solver.NewGoptunaSolverFactory("Goptuna (CMA-ES)", cmaSampler)
	} else if sampler == "ipop-cmaes" {
		factory = solver.NewGoptunaSolverFactory("Goptuna (IPOP-CMA-ES)", ipopCmaSampler)
	} else if sampler == "bipop-cmaes" {
		factory = solver.NewGoptunaSolverFactory("Goptuna (BIPOP-CMA-ES)", bipopCmaSampler)
	} else if sampler == "tpe" {
		factory = solver.NewGoptunaSolverFactory("Goptuna (TPE)", tpeSampler)
	} else {
		panic("invalid sampler")
	}

	runner := kurobako.NewSolverRunner(&factory)
	if err := runner.Run(); err != nil {
		panic(err)
	}
}
