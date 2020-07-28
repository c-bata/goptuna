package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/sile/kurobako-go"
)

type rastriginProblemFactory struct {
	dim int
}

func (r *rastriginProblemFactory) Specification() (*kurobako.ProblemSpec, error) {
	spec := kurobako.NewProblemSpec(fmt.Sprintf("Rastrigin function (dim=%d)", r.dim))

	spec.Params = make([]kurobako.Var, r.dim)
	for i := 0; i < r.dim; i++ {
		x := kurobako.NewVar("x" + strconv.Itoa(i+1))
		x.Range = kurobako.ContinuousRange{Low: -5.12, High: 5.12}.ToRange()
		spec.Params[i] = x
	}

	spec.Values = []kurobako.Var{kurobako.NewVar("Rastrigin")}
	return &spec, nil
}

func (r *rastriginProblemFactory) CreateProblem(seed int64) (kurobako.Problem, error) {
	return &rastriginProblem{}, nil
}

type rastriginProblem struct {
}

func (r *rastriginProblem) CreateEvaluator(params []float64) (kurobako.Evaluator, error) {
	return &rastriginEvaluator{params: params}, nil
}

type rastriginEvaluator struct {
	params []float64
}

func (r *rastriginEvaluator) Evaluate(nextStep uint64) (uint64, []float64, error) {
	v := float64(10 * len(r.params))
	for i := range r.params {
		v += math.Pow(r.params[i], 2) - 10*math.Cos(2*math.Pi*r.params[i])
	}
	return 1, []float64{v}, nil
}

func main() {
	var dim int
	if len(os.Args) == 2 {
		a := os.Args[1]
		b, err := strconv.Atoi(a)
		if err != nil {
			panic(err)
		}
		dim = b
	}
	runner := kurobako.NewProblemRunner(&rastriginProblemFactory{
		dim: dim,
	})
	if err := runner.Run(); err != nil {
		panic(err)
	}
}
