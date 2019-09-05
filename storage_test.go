package goptuna_test

import (
	"fmt"
	"testing"

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
