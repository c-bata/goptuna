package rdb_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/c-bata/goptuna"

	"github.com/c-bata/goptuna/rdb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func SetupSQLite3Test(t *testing.T, sqlitePath string) (*gorm.DB, func(), error) {
	db, err := gorm.Open("sqlite3", sqlitePath)
	if err != nil {
		t.Errorf("failed to setup sqlite3 with %s", err)
		return nil, nil, err
	}
	db.LogMode(false)
	rdb.RunAutoMigrate(db)
	if db.Error != nil {
		t.Errorf("failed to setup sqlite3 with %s", err)
		return nil, nil, err
	}

	return db, func() {
		db.Close()
		os.Remove(sqlitePath)
	}, nil
}

func TestStorage_CreateNewStudy(t *testing.T) {
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	s := rdb.NewStorage(db)
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
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	s := rdb.NewStorage(db)
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
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	s := rdb.NewStorage(db)
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
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetStudyUserAttr(studyID, "key", "value")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	got, err := storage.GetStudyUserAttrs(studyID)
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
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetStudySystemAttr(studyID, "key", "value")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	got, err := storage.GetStudySystemAttrs(studyID)
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
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetTrialUserAttr(trialID, "key", "value")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	got, err := storage.GetTrialUserAttrs(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	want := map[string]string{"key": "value"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %#v, but got %#v", want, got)
	}
}

func TestStorage_TrialSystemAttrs(t *testing.T) {
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetTrialSystemAttr(trialID, "key", "value")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	got, err := storage.GetTrialSystemAttrs(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	if v, ok := got["key"]; !ok || v != "value" {
		t.Errorf("want %#v, but got %v %v", "value", ok, got)
	}
}

func TestStorage_GetAllTrials(t *testing.T) {
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetTrialSystemAttr(trialID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetTrialSystemAttr(trialID, "key2", "value2")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetTrialSystemAttr(trialID, "key3", "value3")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetTrialUserAttr(trialID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	trialID, err = storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	trials, err := storage.GetAllTrials(studyID)
	if len(trials) != 2 {
		t.Errorf("want two trials, but got %d\n  Detail: %#v", len(trials), trials)
		return
	}
}

func TestStorage_SetTrialState(t *testing.T) {
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	trialID, err := storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetTrialValue(trialID, 0.1)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trial, err := storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	before := trial.DatetimeComplete

	err = storage.SetTrialState(trialID, goptuna.TrialStateRunning)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trial, err = storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	after := trial.DatetimeComplete
	if before.Unix() != after.Unix() {
		t.Errorf("DatetimeComplete should not be updated, but changed: %s to %s",
			before.String(), after.String())
		return
	}

	err = storage.SetTrialState(trialID, goptuna.TrialStateComplete)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trial, err = storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	after = trial.DatetimeComplete
	if before.Unix() == after.Unix() {
		t.Errorf("DatetimeComplete should be updated, but got %s", after.String())
		return
	}
}

func TestStorage_GetTrial(t *testing.T) {
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetTrialSystemAttr(trialID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetTrialSystemAttr(trialID, "key2", "value2")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetTrialSystemAttr(trialID, "key3", "value3")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetTrialValue(trialID, 0.1)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetTrialState(trialID, goptuna.TrialStateComplete)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetTrialUserAttr(trialID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	trial, err := storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if trial.ID != trialID {
		t.Errorf("want trialID = %d, but got %d", trialID, trial.ID)
	}
}

func TestStorage_GetAllStudySummaries(t *testing.T) {
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetStudySystemAttr(studyID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetStudySystemAttr(studyID, "key2", "value2")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetStudyUserAttr(studyID, "key1", "value1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	studies, err := storage.GetAllStudySummaries()
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if len(studies) != 1 || studies[0].ID != studyID {
		t.Errorf("want studyID = %d, but got %#v", studyID, studies)
	}
}

func TestStorage_GetBestTrial(t *testing.T) {
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	// trial number 0 (not completed yet)
	trialID, err := storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetTrialValue(trialID, 0.1)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	// trial number 1
	trialID, err = storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetTrialValue(trialID, 0.3)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetTrialState(trialID, goptuna.TrialStateComplete)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	// trial number 2
	trialID, err = storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetTrialValue(trialID, 0.2)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetTrialState(trialID, goptuna.TrialStateComplete)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	bestTrial, err := storage.GetBestTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if bestTrial.Value == 0.2 && bestTrial.Number != 2 {
		t.Errorf("want Trial(Value=0.2, Number: 2), but got %#v", bestTrial)
	}
}

func TestStorage_SetTrialIntermediateValue(t *testing.T) {
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudy("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	// trial number 1 (not completed yet)
	trialID, err := storage.CreateNewTrial(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	err = storage.SetTrialIntermediateValue(trialID, 1, 0.5)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetTrialIntermediateValue(trialID, 3, 0.7)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	trial, err := storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if len(trial.IntermediateValues) != 2 {
		t.Errorf("want two intermediate vales, but got %#v", trial.IntermediateValues)
	}
}

func TestStorage_CloneTrial(t *testing.T) {
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
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
