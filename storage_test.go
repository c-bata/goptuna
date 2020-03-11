package goptuna_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/c-bata/goptuna"
)

func ExampleInMemoryStorage_CreateNewStudy() {
	storage := goptuna.NewInMemoryStorage()
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		panic(err)
	}
	fmt.Println(studyID)

	// Output:
	// 1
}

func ExampleInMemoryStorage_SetStudyUserAttr() {
	storage := goptuna.NewInMemoryStorage()
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		panic(err)
	}
	err = storage.SetStudyUserAttr(studyID, "key", "value")
	if err != nil {
		panic(err)
	}
	attrs, err := storage.GetStudyUserAttrs(studyID)
	if err != nil {
		panic(err)
	}
	for k, v := range attrs {
		fmt.Println(k, v)
	}

	// Output:
	// key value
}

func ExampleInMemoryStorage_SetStudySystemAttr() {
	storage := goptuna.NewInMemoryStorage()
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		panic(err)
	}
	err = storage.SetStudySystemAttr(studyID, "key", "value")
	if err != nil {
		panic(err)
	}
	attrs, err := storage.GetStudySystemAttrs(studyID)
	if err != nil {
		panic(err)
	}
	for k, v := range attrs {
		fmt.Println(k, v)
	}
	// Output:
	// key value
}

func ExampleInMemoryStorage_SetTrialUserAttr() {
	storage := goptuna.NewInMemoryStorage()
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		panic(err)
	}
	trialID, err := storage.CreateNewTrial(studyID)
	if err != nil {
		panic(err)
	}

	err = storage.SetTrialUserAttr(trialID, "key", "value")
	if err != nil {
		panic(err)
	}
	attrs, err := storage.GetTrialUserAttrs(trialID)
	if err != nil {
		panic(err)
	}
	for k, v := range attrs {
		fmt.Println(k, v)
	}

	// Output:
	// key value
}

func ExampleInMemoryStorage_SetTrialSystemAttr() {
	storage := goptuna.NewInMemoryStorage()
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		panic(err)
	}
	trialID, err := storage.CreateNewTrial(studyID)
	if err != nil {
		panic(err)
	}

	err = storage.SetTrialSystemAttr(trialID, "key", "value")
	if err != nil {
		panic(err)
	}
	attrs, err := storage.GetTrialSystemAttrs(trialID)
	if err != nil {
		panic(err)
	}
	for k, v := range attrs {
		fmt.Println(k, v)
	}
	// Output:
	// key value
}

func ExampleInMemoryStorage_GetStudyIDFromTrialID() {
	storage := goptuna.NewInMemoryStorage()
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		panic(err)
	}
	trialID, err := storage.CreateNewTrial(studyID)
	if err != nil {
		panic(err)
	}

	actual, err := storage.GetStudyIDFromTrialID(trialID)
	if err != nil {
		panic(err)
	}
	fmt.Println(actual)
	// Output:
	// 1
}

func TestMemoryStorage_GetAllStudySummaries(t *testing.T) {
	storage := goptuna.NewInMemoryStorage()
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	err = storage.SetStudyDirection(studyID, goptuna.StudyDirectionMinimize)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	trialID, err := storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	err = storage.SetTrialValue(trialID, 0.1)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	err = storage.SetTrialState(trialID, goptuna.TrialStateComplete)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	trialID, err = storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	err = storage.SetTrialValue(trialID, 0.5)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	err = storage.SetTrialState(trialID, goptuna.TrialStateComplete)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	summaries, err := storage.GetAllStudySummaries()
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	if len(summaries) != 1 {
		t.Errorf("should get one study summary, but got %d", len(summaries))
		return
	}

	if summaries[0].BestTrial.Value != 0.1 {
		t.Errorf("the value of best trial should be 0.1, but got %f", summaries[0].BestTrial.Value)
		return
	}
}

func TestMemoryStorage_CloneTrial(t *testing.T) {
	storage := goptuna.NewInMemoryStorage()
	now := time.Now()

	baseTrial := goptuna.FrozenTrial{
		ID:     -1, // dummy value (unused)
		Number: -1, // dummy value (unused)
		State:  goptuna.TrialStateComplete,
		Value:  10000,
		IntermediateValues: map[int]float64{
			1: 10,
			2: 100,
			3: 1000,
		},
		DatetimeStart:    now,
		DatetimeComplete: now,
		InternalParams: map[string]float64{
			"x": 0.5,
		},
		Params: map[string]interface{}{
			"x": 0.5,
		},
		Distributions: map[string]interface{}{
			"x": goptuna.UniformDistribution{High: 1, Low: 0},
		},
		UserAttrs: map[string]string{
			"foo": "bar",
		},
		SystemAttrs: map[string]string{
			"baz": "123",
		},
	}

	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	err = storage.SetStudyDirection(studyID, goptuna.StudyDirectionMinimize)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	trialID, err := storage.CloneTrial(studyID, baseTrial)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	trials, err := storage.GetAllTrials(studyID)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	if len(trials) != 1 {
		t.Errorf("should get one trial, but got %d", len(trials))
		return
	}

	if trials[0].ID != trialID {
		t.Errorf("trialID should be %d, but got %d", trialID, trials[0].ID)
	}

	if trials[0].Number != 0 {
		t.Errorf("number should be 0, but got %d", trials[0].Number)
	}

	if trials[0].State != goptuna.TrialStateComplete {
		t.Errorf("state should be complete, but got %s", trials[0].State)
	}

	if !reflect.DeepEqual(trials[0].Distributions, baseTrial.Distributions) {
		t.Errorf("Distributions should be %v, but got %v", trials[0].Distributions, baseTrial.Distributions)
	}

	if !reflect.DeepEqual(trials[0].Params, baseTrial.Params) {
		t.Errorf("Params should be %v, but got %v", trials[0].Params, baseTrial.Params)
	}

	if !reflect.DeepEqual(trials[0].InternalParams, baseTrial.InternalParams) {
		t.Errorf("InternalParams should be %v, but got %v", trials[0].InternalParams, baseTrial.InternalParams)
	}

	if !reflect.DeepEqual(trials[0].IntermediateValues, baseTrial.IntermediateValues) {
		t.Errorf("InternalValues should be %v, but got %v", trials[0].IntermediateValues, baseTrial.IntermediateValues)
	}

	if trials[0].DatetimeStart.Second() != baseTrial.DatetimeStart.Second() {
		t.Errorf("DatetimeStart should be %s, but got %s", trials[0].DatetimeStart, baseTrial.DatetimeStart)
	}

	if trials[0].DatetimeComplete.Second() != baseTrial.DatetimeComplete.Second() {
		t.Errorf("DatetimeComplete should be %s, but got %s", trials[0].DatetimeComplete, baseTrial.DatetimeComplete)
	}
}
