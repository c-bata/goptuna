package rdb_test

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/rdb.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	setupCounter   int
	setupCounterMu sync.Mutex
)

func SetupSQLite3Test() (*rdb.Storage, func(), error) {
	setupCounterMu.Lock()
	defer setupCounterMu.Unlock()
	setupCounter += 1
	sqlitePath := fmt.Sprintf("goptuna-test-%d.db", setupCounter)

	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	err = rdb.RunAutoMigrate(db)
	if err != nil {
		return nil, nil, err
	}
	storage := rdb.NewStorage(db)

	return storage, func() {
		os.Remove(sqlitePath)
	}, nil
}

func TestStorage_CreateNewStudy(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()
	got, err := s.CreateNewStudy("study1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if got != 1 {
		t.Errorf("Storage.CreateNewStudy() = %v, want %v", got, 1)
	}

	// different study name
	got, err = s.CreateNewStudy("study2")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if got != 2 {
		t.Errorf("Storage.CreateNewStudy() = %v, want %v", got, 1)
	}

	// duplicate study name
	got, err = s.CreateNewStudy("study1")
	if err == nil {
		t.Errorf("Storage.CreateNewStudy() error = %v, want duplicate error", err)
		return
	}
}

func TestStorage_DeleteStudy(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	got, err := s.CreateNewStudy("study1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if got != 1 {
		t.Errorf("Storage.CreateNewStudy() = %v, want %v", got, 1)
	}
	got, err = s.CreateNewStudy("study2")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if got != 2 {
		t.Errorf("Storage.CreateNewStudy() = %v, want %v", got, 1)
	}

	summaries, err := s.GetAllStudySummaries()
	if err != nil {
		t.Errorf("Storage.GetAllStudySummaries() error = %v, want nil", err)
		return
	}
	if len(summaries) != 2 {
		t.Errorf("Must be two studies.")
	}

	// duplicate study name
	err = s.DeleteStudy(1)
	if err != nil {
		t.Errorf("Storage.DeleteNewStudy() error = %v, want nil", err)
		return
	}

	summaries, err = s.GetAllStudySummaries()
	if err != nil {
		t.Errorf("Storage.GetAllStudySummaries() error = %v, want nil", err)
		return
	}
	if len(summaries) != 1 {
		t.Errorf("Must be one study.")
	}
}

func TestStorage_StudyDirection(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	studyID, err := s.CreateNewStudy("study")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	direction, err := s.GetStudyDirection(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	if direction != goptuna.StudyDirectionMinimize {
		t.Errorf("want %s, but got %s", direction, goptuna.StudyDirectionMinimize)
		return
	}

	err = s.SetStudyDirection(studyID, goptuna.StudyDirectionMaximize)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	direction, err = s.GetStudyDirection(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if direction != goptuna.StudyDirectionMaximize {
		t.Errorf("want %s, but got %s", goptuna.StudyDirectionMaximize, direction)
		return
	}
}

func TestStorage_StudyUserAttrs(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	studyID, err := s.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetStudyUserAttr(studyID, "key", "value")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	got, err := s.GetStudyUserAttrs(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	want := map[string]string{"key": "value"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %#v, but got %#v", want, got)
	}
}

func TestStorage_StudySystemAttrs(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	studyID, err := s.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetStudySystemAttr(studyID, "key", "value")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	got, err := s.GetStudySystemAttrs(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	want := map[string]string{"key": "value"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %#v, but got %#v", want, got)
	}
}

func TestStorage_TrialUserAttrs(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	studyID, err := s.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := s.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetTrialUserAttr(trialID, "key", "value")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	got, err := s.GetTrialUserAttrs(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	want := map[string]string{"key": "value"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %#v, but got %#v", want, got)
	}

	err = s.SetTrialUserAttr(trialID, "key", "value")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
}

func TestStorage_TrialSystemAttrs(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	studyID, err := s.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := s.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetTrialSystemAttr(trialID, "key", "value")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	got, err := s.GetTrialSystemAttrs(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	if v, ok := got["key"]; !ok || v != "value" {
		t.Errorf("want %#v, but got %v %v", "value", ok, got)
	}

	err = s.SetTrialSystemAttr(trialID, "key", "value")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
}

func TestStorage_GetAllTrials(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	studyID, err := s.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := s.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetTrialSystemAttr(trialID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetTrialSystemAttr(trialID, "key2", "value2")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetTrialSystemAttr(trialID, "key3", "value3")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetTrialUserAttr(trialID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	trialID, err = s.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	trials, err := s.GetAllTrials(studyID)
	if len(trials) != 2 {
		t.Errorf("want two trials, but got %d\n  Detail: %#v", len(trials), trials)
		return
	}
}

func TestStorage_SetTrialState(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	studyID, err := s.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	trialID, err := s.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = s.SetTrialValue(trialID, 0.1)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trial, err := s.GetTrial(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	before := trial.DatetimeComplete

	err = s.SetTrialState(trialID, goptuna.TrialStateComplete)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trial, err = s.GetTrial(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	after := trial.DatetimeComplete
	if before.Unix() == after.Unix() {
		t.Errorf("DatetimeComplete should be updated, %s == %s", before.String(), after.String())
		return
	}
}

func TestStorage_GetTrial(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	studyID, err := s.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := s.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetTrialSystemAttr(trialID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetTrialSystemAttr(trialID, "key2", "value2")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetTrialSystemAttr(trialID, "key3", "value3")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = s.SetTrialValue(trialID, 0.1)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = s.SetTrialState(trialID, goptuna.TrialStateComplete)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetTrialUserAttr(trialID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	trial, err := s.GetTrial(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if trial.ID != trialID {
		t.Errorf("want trialID = %d, but got %d", trialID, trial.ID)
	}
}

func TestStorage_GetAllStudySummaries(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	studyID, err := s.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetStudySystemAttr(studyID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = s.SetStudySystemAttr(studyID, "key2", "value2")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = s.SetStudyUserAttr(studyID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	studies, err := s.GetAllStudySummaries()
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if len(studies) != 1 || studies[0].ID != studyID {
		t.Errorf("want studyID = %d, but got %#v", studyID, studies)
	}
}

func TestStorage_GetBestTrial(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	studyID, err := s.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	// trial number 0 (not completed yet)
	trialID, err := s.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = s.SetTrialValue(trialID, 0.1)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	// trial number 1
	trialID, err = s.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = s.SetTrialValue(trialID, 0.3)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = s.SetTrialState(trialID, goptuna.TrialStateComplete)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	// trial number 2
	trialID, err = s.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = s.SetTrialValue(trialID, 0.2)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = s.SetTrialState(trialID, goptuna.TrialStateComplete)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	bestTrial, err := s.GetBestTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if bestTrial.Value == 0.2 && bestTrial.Number != 2 {
		t.Errorf("want Trial(Value=0.2, Number: 2), but got %#v", bestTrial)
	}
}

func TestStorage_SetTrialIntermediateValue(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	studyID, err := s.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	// trial number 1 (not completed yet)
	trialID, err := s.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = s.SetTrialIntermediateValue(trialID, 1, 0.5)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = s.SetTrialIntermediateValue(trialID, 3, 0.7)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	trial, err := s.GetTrial(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if len(trial.IntermediateValues) != 2 {
		t.Errorf("want two intermediate vales, but got %#v", trial.IntermediateValues)
	}
}

func TestStorage_CloneTrial(t *testing.T) {
	s, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

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

	studyID, err := s.CreateNewStudy("")
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	err = s.SetStudyDirection(studyID, goptuna.StudyDirectionMinimize)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	trialID, err := s.CloneTrial(studyID, baseTrial)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	trials, err := s.GetAllTrials(studyID)
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
