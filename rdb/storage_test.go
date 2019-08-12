package rdb_test

import (
	"os"
	"reflect"
	"testing"

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

func TestStorage_CreateNewStudyID(t *testing.T) {
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	s := rdb.NewStorage(db)
	got, err := s.CreateNewStudyID("study1")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if got != 1 {
		t.Errorf("Storage.CreateNewStudyID() = %v, want %v", got, 1)
	}

	// different study name
	got, err = s.CreateNewStudyID("study2")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	if got != 2 {
		t.Errorf("Storage.CreateNewStudyID() = %v, want %v", got, 1)
	}

	// duplicate study name
	got, err = s.CreateNewStudyID("study1")
	if err == nil {
		t.Errorf("Storage.CreateNewStudyID() error = %v, want duplicate error", err)
		return
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
	studyID, err := s.CreateNewStudyID("study")
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
	studyID, err := storage.CreateNewStudyID("")
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
	studyID, err := storage.CreateNewStudyID("")
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
	studyID, err := storage.CreateNewStudyID("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := storage.CreateNewTrialID(studyID)
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
	studyID, err := storage.CreateNewStudyID("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := storage.CreateNewTrialID(studyID)
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
	studyID, err := storage.CreateNewStudyID("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := storage.CreateNewTrialID(studyID)
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

	trialID, err = storage.CreateNewTrialID(studyID)
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

func TestStorage_GetTrial(t *testing.T) {
	db, teardown, err := SetupSQLite3Test(t, "goptuna-test.db")
	defer teardown()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}

	storage := rdb.NewStorage(db)
	studyID, err := storage.CreateNewStudyID("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	trialID, err := storage.CreateNewTrialID(studyID)
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
	studyID, err := storage.CreateNewStudyID("")
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
	studyID, err := storage.CreateNewStudyID("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	// trial number 1 (not completed yet)
	trialID, err := storage.CreateNewTrialID(studyID)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	err = storage.SetTrialValue(trialID, 0.1)
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}

	// trial number 2
	trialID, err = storage.CreateNewTrialID(studyID)
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

	// trial number 3
	trialID, err = storage.CreateNewTrialID(studyID)
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
	if bestTrial.Value == 0.2 && bestTrial.Number != 3 {
		t.Errorf("want Trial(Value=0.2, Number: 3), but got %#v", bestTrial)
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
	studyID, err := storage.CreateNewStudyID("")
	if err != nil {
		t.Errorf("error: %v != nil", err)
		return
	}
	// trial number 1 (not completed yet)
	trialID, err := storage.CreateNewTrialID(studyID)
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
