package main

import (
	"math"

	"github.com/sile/kurobako-go"
)

type himmelblauProblemFactory struct {
}

func himmelblau(x1 float64, x2 float64) float64 {
	return math.Pow(math.Pow(x1, 2)+x2-11, 2) + math.Pow(x1+math.Pow(x2, 2)-7, 2)
}

func (r *himmelblauProblemFactory) Specification() (*kurobako.ProblemSpec, error) {
	spec := kurobako.NewProblemSpec("Himmelblau Function")

	x1 := kurobako.NewVar("x1")
	x1.Range = kurobako.ContinuousRange{-4, 4}.ToRange()

	x2 := kurobako.NewVar("x2")
	x2.Range = kurobako.ContinuousRange{-4, 4}.ToRange()

	spec.Params = []kurobako.Var{x1, x2}

	spec.Values = []kurobako.Var{kurobako.NewVar("Himmelblau")}

	return &spec, nil
}

func (r *himmelblauProblemFactory) CreateProblem(seed int64) (kurobako.Problem, error) {
	return &himmelblauProblem{}, nil
}

type himmelblauProblem struct {
}

func (r *himmelblauProblem) CreateEvaluator(params []float64) (kurobako.Evaluator, error) {
	x1 := params[0]
	x2 := params[1]
	return &himmelblauEvaluator{x1, x2}, nil
}

type himmelblauEvaluator struct {
	x1 float64
	x2 float64
}

func (r *himmelblauEvaluator) Evaluate(nextStep uint64) (uint64, []float64, error) {
	values := []float64{himmelblau(r.x1, r.x2)}
	return 1, values, nil
}

func main() {
	runner := kurobako.NewProblemRunner(&himmelblauProblemFactory{})
	if err := runner.Run(); err != nil {
		panic(err)
	}
}
