package goptuna_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/c-bata/goptuna"
)

func ExampleStudy_Optimize() {
	sampler := goptuna.NewRandomSearchSampler(
		goptuna.RandomSearchSamplerOptionSeed(0),
	)
	study, _ := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionSampler(sampler),
		goptuna.StudyOptionLogger(nil),
	)

	objective := func(trial goptuna.Trial) (float64, error) {
		x1, _ := trial.SuggestFloat("x1", -10, 10)
		x2, _ := trial.SuggestFloat("x2", -10, 10)
		return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
	}

	if err := study.Optimize(objective, 10); err != nil {
		panic(err)
	}
	value, _ := study.GetBestValue()
	params, _ := study.GetBestParams()

	fmt.Printf("Best trial: %.5f\n", value)
	fmt.Printf("x1: %.3f\n", params["x1"].(float64))
	fmt.Printf("x2: %.3f\n", params["x2"].(float64))
	// Output:
	// Best trial: 0.03833
	// x1: 2.182
	// x2: -4.927
}

func TestStudy_SystemAttrs(t *testing.T) {
	study, _ := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSampler(goptuna.NewRandomSearchSampler()),
	)

	err := study.SetSystemAttr("hello", "world")
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}

	attrs, err := study.GetSystemAttrs()
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}

	hello, ok := attrs["hello"]
	if !ok {
		t.Errorf("'hello' doesn't exist")
		return
	}
	if hello != "world" {
		t.Errorf("should be 'world', but got '%s'", hello)
	}
}

func ExampleStudy_EnqueueTrial() {
	study, _ := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionLogger(nil),
	)

	objective := func(trial goptuna.Trial) (float64, error) {
		x1, _ := trial.SuggestFloat("x1", -10, 10)
		x2, _ := trial.SuggestFloat("x2", -10, 10)
		return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
	}

	study.EnqueueTrial(map[string]float64{"x1": 2, "x2": -5})

	err := study.Optimize(objective, 1)
	if err != nil {
		panic(err)
	}

	trials, err := study.GetTrials()
	if err != nil {
		panic(err)
	}

	fmt.Printf("x1: %.3f\n", trials[0].Params["x1"].(float64))
	fmt.Printf("x2: %.3f\n", trials[0].Params["x2"].(float64))

	// Output:
	// x1: 2.000
	// x2: -5.000
}

func TestStudy_EnqueueTrial_WithUnfixedParameter(t *testing.T) {
	study, _ := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionLogger(nil),
	)

	objective := func(trial goptuna.Trial) (float64, error) {
		x1, _ := trial.SuggestFloat("x1", -10, 10)
		x2, _ := trial.SuggestFloat("x2", -10, 10)
		return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
	}

	study.EnqueueTrial(map[string]float64{"x1": 2})

	err := study.Optimize(objective, 1)
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}
	trials, err := study.GetTrials()
	if err != nil {
		t.Errorf("err: %v != nil", err)
	}

	if trials[0].InternalParams["x1"] != 2 {
		t.Errorf("x1 should be 2, but got %f", trials[0].InternalParams["x1"])
	}
}

func TestStudy_UserAttrs(t *testing.T) {
	study, _ := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSampler(goptuna.NewRandomSearchSampler()),
	)

	err := study.SetUserAttr("hello", "world")
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}

	attrs, err := study.GetUserAttrs()
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}

	hello, ok := attrs["hello"]
	if !ok {
		t.Errorf("'hello' doesn't exist")
		return
	}
	if hello != "world" {
		t.Errorf("should be 'world', but got '%s'", hello)
	}
}

func TestStudy_AppendTrial(t *testing.T) {
	study, err := goptuna.CreateStudy(
		"example",
		goptuna.StudyOptionSetDirection(goptuna.StudyDirectionMinimize),
		goptuna.StudyOptionSampler(goptuna.NewRandomSearchSampler()),
	)
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}
	err = goptuna.ExportStudyAppendTrial(
		study,
		0.8,
		nil,
		nil,
		nil,
		nil,
		nil,
		goptuna.TrialStateComplete,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}
	bestValue, err := study.GetBestValue()
	if err != nil {
		t.Errorf("err: %v != nil", err)
		return
	}
	if bestValue != 0.8 {
		t.Errorf("bestValue: %f != 0.8", bestValue)
		return
	}
}
