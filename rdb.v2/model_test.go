package rdb

import (
	"fmt"
	"math"
	"os"
	"sync"
	"testing"

	"github.com/c-bata/goptuna"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	setupCounter   int
	setupCounterMu sync.Mutex
)

func SetupSQLite3Test() (*gorm.DB, func(), error) {
	setupCounterMu.Lock()
	defer setupCounterMu.Unlock()
	setupCounter++
	sqlitePath := fmt.Sprintf("goptuna-test-%d.db", setupCounter)

	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	err = RunAutoMigrate(db)
	if err != nil {
		return nil, nil, err
	}

	// Enable foreign_keys = ON
	if db.Dialector.Name() == "sqlite" {
		err = db.Exec("PRAGMA foreign_keys = ON").Error
		if err != nil {
			return nil, nil, err
		}
	}

	return db, func() {
		os.Remove(sqlitePath)
	}, nil
}

func TestStudyCascadeOnDelete(t *testing.T) {
	db, teardown, err := SetupSQLite3Test()
	if err != nil {
		t.Errorf("failed to setup tests with %s", err)
		return
	}
	defer teardown()

	// Run optimizatoin
	study, err := goptuna.CreateStudy("cascade-test", goptuna.StudyOptionStorage(NewStorage(db)))
	if err != nil {
		t.Errorf("failed to create study: %s", err)
		return
	}
	err = study.SetSystemAttr("foo", "var")
	if err != nil {
		t.Errorf("failed to set system attr: %s", err)
		return
	}
	err = study.SetUserAttr("foo", "var")
	if err != nil {
		t.Errorf("failed to set user attr: %s", err)
		return
	}
	err = study.Optimize(func(trial goptuna.Trial) (float64, error) {
		x1, err := trial.SuggestFloat("x1", -10, 10)
		if err != nil {
			return 0, err
		}
		x2, err := trial.SuggestFloat("x2", -10, 10)
		if err != nil {
			return 0, err
		}
		err = trial.SetUserAttr("foo", "var")
		if err != nil {
			return 0, err
		}
		err = trial.SetSystemAttr("foo", "var")
		if err != nil {
			return 0, err
		}
		return math.Pow(x1-2, 2) + math.Pow(x2+5, 2), nil
	}, 10)
	if err != nil {
		t.Errorf("failed to run optimize: %s", err)
		return
	}

	// Delete a study
	err = db.Delete(&studyModel{
		ID: study.ID,
	}).Error

	// Count study_user_attributes table
	var count int64
	if db.Model(&studyUserAttributeModel{}).Count(&count).Error != nil {
		t.Errorf("failed to count trial model: %s", err)
		return
	}
	if count != 0 {
		t.Errorf("trials table should be empty, but got count=%d", count)
	}

	// Count study_system_attributes table
	if db.Model(&studySystemAttributeModel{}).Count(&count).Error != nil {
		t.Errorf("failed to count trial model: %s", err)
		return
	}
	if count != 0 {
		t.Errorf("trials table should be empty, but got count=%d", count)
	}

	// Count trials table
	if db.Model(&trialModel{}).Count(&count).Error != nil {
		t.Errorf("failed to count trial model: %s", err)
		return
	}
	if count != 0 {
		t.Errorf("trials table should be empty, but got count=%d", count)
	}

	// Count trial_values table
	if db.Model(&trialValueModel{}).Count(&count).Error != nil {
		t.Errorf("failed to count trial model: %s", err)
		return
	}
	if count != 0 {
		t.Errorf("trials table should be empty, but got count=%d", count)
	}

	// Count trial_params table
	if db.Model(&trialParamModel{}).Count(&count).Error != nil {
		t.Errorf("failed to count trial model: %s", err)
		return
	}
	if count != 0 {
		t.Errorf("trials table should be empty, but got count=%d", count)
	}

	// Count trial_user_attributes table
	if db.Model(&trialUserAttributeModel{}).Count(&count).Error != nil {
		t.Errorf("failed to count trial model: %s", err)
		return
	}
	if count != 0 {
		t.Errorf("trials table should be empty, but got count=%d", count)
	}

	// Count trial_system_attributes table
	if db.Model(&trialSystemAttributeModel{}).Count(&count).Error != nil {
		t.Errorf("failed to count trial model: %s", err)
		return
	}
	if count != 0 {
		t.Errorf("trials table should be empty, but got count=%d", count)
	}
}
