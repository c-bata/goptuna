package main

import kurobako "github.com/sile/kurobako-go"

type quadraticProblemFactory struct {
}

func (r *quadraticProblemFactory) Specification() (*kurobako.ProblemSpec, error) {
	spec := kurobako.NewProblemSpec("Quadratic Function")

	x := kurobako.NewVar("x")
	x.Range = kurobako.ContinuousRange{-10.0, 10.0}.ToRange()

	y := kurobako.NewVar("y")
	y.Range = kurobako.DiscreteRange{-3, 3}.ToRange()

	spec.Params = []kurobako.Var{x, y}

	spec.Values = []kurobako.Var{kurobako.NewVar("x**2 + y")}

	return &spec, nil
}

func (r *quadraticProblemFactory) CreateProblem(seed int64) (kurobako.Problem, error) {
	return &quadraticProblem{}, nil
}

type quadraticProblem struct {
}

func (r *quadraticProblem) CreateEvaluator(params []float64) (kurobako.Evaluator, error) {
	x := params[0]
	y := params[1]
	return &quadraticEvaluator{x, y}, nil
}

type quadraticEvaluator struct {
	x float64
	y float64
}

func (r *quadraticEvaluator) Evaluate(nextStep uint64) (uint64, []float64, error) {
	values := []float64{r.x*r.x + r.y}
	return 1, values, nil
}

func main() {
	runner := kurobako.NewProblemRunner(&quadraticProblemFactory{})
	if err := runner.Run(); err != nil {
		panic(err)
	}
}
