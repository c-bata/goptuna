package goptuna_test

import (
	"fmt"
	"math"

	"github.com/c-bata/goptuna"
)

func ExampleStudy_Optimize() {
	sampler := goptuna.NewRandomSearchSampler(
		goptuna.RandomSearchSamplerOptionSeed(0),
	)
	study, _ := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSampler(sampler),
	)

	objective := func(trial goptuna.Trial) (float64, error) {
		x1, _ := trial.SuggestUniform("x1", -10, 10)
		x2, _ := trial.SuggestUniform("x2", -10, 10)
		return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
	}

	if err := study.Optimize(objective, 10); err != nil {
		panic(err)
	}
	value, _ := study.GetBestValue()
	params, _ := study.GetBestParams()

	fmt.Printf("Best trial: %.5f\n", value)
	for k := range params {
		fmt.Printf("%s: %.3f\n", k, params[k])
	}
	// Output:
	// Best trial: 0.03833
	// x1: 2.182
	// x2: -4.927
}
