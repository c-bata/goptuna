package goptuna_test

import (
	"fmt"
	"math"

	"github.com/c-bata/goptuna"
)

func ExampleTrial_SuggestUniform() {
	sampler := goptuna.NewRandomSearchSampler(
		goptuna.RandomSearchSamplerOptionSeed(0),
	)
	study, err := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSampler(sampler),
	)
	if err != nil {
		panic(err)
	}

	objective := func(trial goptuna.Trial) (float64, error) {
		x1, _ := trial.SuggestUniform("x1", -10, 10)
		x2, _ := trial.SuggestUniform("x2", -10, 10)
		fmt.Printf("sampled: %.3f, %.3f\n", x1, x2)
		return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
	}

	err = study.Optimize(objective, 1)
	if err != nil {
		panic(err)
	}
	// Output:
	// sampled: 8.904, -5.101
}

func ExampleTrial_SuggestLogUniform() {
	sampler := goptuna.NewRandomSearchSampler(
		goptuna.RandomSearchSamplerOptionSeed(1),
	)
	study, err := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSampler(sampler),
	)
	if err != nil {
		panic(err)
	}

	objective := func(trial goptuna.Trial) (float64, error) {
		x1, _ := trial.SuggestLogUniform("x1", 1e-5, 1e10)
		x2, _ := trial.SuggestLogUniform("x2", 1e-5, 1e10)
		fmt.Printf("sampled: %f, %f\n", x1, x2)
		return math.Pow(x1-1e5, 2) + math.Pow(x2-1e3, 2), nil
	}

	err = study.Optimize(objective, 1)
	if err != nil {
		panic(err)
	}
	// Output:
	// sampled: 11746.387394, 1281257205.189628
}

func ExampleTrial_SuggestInt() {
	sampler := goptuna.NewRandomSearchSampler(
		goptuna.RandomSearchSamplerOptionSeed(1),
	)
	study, err := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSampler(sampler),
	)
	if err != nil {
		panic(err)
	}

	objective := func(trial goptuna.Trial) (float64, error) {
		x1, _ := trial.SuggestInt("x1", -10, 10)
		x2, _ := trial.SuggestInt("x2", -10, 10)
		fmt.Printf("sampled: %d, %d\n", x1, x2)
		return math.Pow(float64(x1-2), 2) + math.Pow(float64(x2+5), 2), nil
	}

	err = study.Optimize(objective, 1)
	if err != nil {
		panic(err)
	}
	// Output:
	// sampled: -9, -3
}

func ExampleTrial_SuggestDiscreteUniform() {
	sampler := goptuna.NewRandomSearchSampler(
		goptuna.RandomSearchSamplerOptionSeed(1),
	)
	study, err := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSampler(sampler),
	)
	if err != nil {
		panic(err)
	}

	objective := func(trial goptuna.Trial) (float64, error) {
		x1, _ := trial.SuggestDiscreteUniform("x1", -10, 10, 0.1)
		x2, _ := trial.SuggestDiscreteUniform("x2", -10, 10, 0.1)
		fmt.Printf("sampled: %f, %f\n", x1, x2)
		return math.Pow(x1-5, 2) + math.Pow(x2+3, 2), nil
	}
	err = study.Optimize(objective, 1)
	if err != nil {
		panic(err)
	}
	// Output:
	// sampled: 2.100000, 8.900000
}
