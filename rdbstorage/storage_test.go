package rdbstorage_test

import (
	"os"
	"testing"

	"github.com/c-bata/goptuna"

	"github.com/c-bata/goptuna/rdbstorage"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func SetupSQLite3Test(t *testing.T, sqlitePath string) (*gorm.DB, func(), error) {
	db, err := gorm.Open("sqlite3", sqlitePath)
	if err != nil {
		t.Errorf("failed to setup sqlite3 with %s", err)
		return nil, nil, err
	}
	rdbstorage.RunAutoMigrate(db)
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

	s := rdbstorage.NewStorage(db)
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

	s := rdbstorage.NewStorage(db)
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
