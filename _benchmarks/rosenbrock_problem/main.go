package main

import (
	"math"

	"github.com/sile/kurobako-go"
)

type rosenbrockProblemFactory struct {
}

func rosenbrock(x1 float64, x2 float64) float64 {
	return 100*math.Pow(x2-math.Pow(x1, 2), 2) + math.Pow(x1-1, 2)
}

func (r *rosenbrockProblemFactory) Specification() (*kurobako.ProblemSpec, error) {
	spec := kurobako.NewProblemSpec("Rosenbrock Function")

	x1 := kurobako.NewVar("x1")
	x1.Range = kurobako.ContinuousRange{Low: -4, High: 4}.ToRange()

	x2 := kurobako.NewVar("x2")
	x2.Range = kurobako.ContinuousRange{Low: -4, High: 4}.ToRange()

	spec.Params = []kurobako.Var{x1, x2}

	spec.Values = []kurobako.Var{kurobako.NewVar("Rosenbrock")}

	return &spec, nil
}

func (r *rosenbrockProblemFactory) CreateProblem(seed int64) (kurobako.Problem, error) {
	return &rosenbrockProblem{}, nil
}

type rosenbrockProblem struct {
}

func (r *rosenbrockProblem) CreateEvaluator(params []float64) (kurobako.Evaluator, error) {
	x1 := params[0]
	x2 := params[1]
	return &rosenbrockEvaluator{x1, x2}, nil
}

type rosenbrockEvaluator struct {
	x1 float64
	x2 float64
}

func (r *rosenbrockEvaluator) Evaluate(nextStep uint64) (uint64, []float64, error) {
	values := []float64{rosenbrock(r.x1, r.x2)}
	return 1, values, nil
}

func main() {
	runner := kurobako.NewProblemRunner(&rosenbrockProblemFactory{})
	if err := runner.Run(); err != nil {
		panic(err)
	}
}
