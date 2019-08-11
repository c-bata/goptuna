package rdbstorage

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/c-bata/goptuna"
)

var _ goptuna.Storage = &Storage{}

func NewStorage() (*Storage, error) {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

type Storage struct {
	db *gorm.DB
}

func (s *Storage) CreateNewStudyID(name string) (int, error) {
	if name == "" {
		u, err := uuid.NewUUID()
		if err != nil {
			return -1, err
		}
		name = goptuna.DefaultStudyNamePrefix + u.String()
	}
	s.db.Create(&StudyModel{
		Name:      name,
		Direction: DirectionNotSet,
	})
	var study StudyModel
	s.db.First(&study, "study_name = ?", name)
	return study.ID, nil
}

func (s *Storage) SetStudyDirection(studyID int, direction goptuna.StudyDirection) error {
	panic("implement me")
}

func (*Storage) SetStudyUserAttr(studyID int, key string, value interface{}) error {
	panic("implement me")
}

func (*Storage) SetStudySystemAttr(studyID int, key string, value interface{}) error {
	panic("implement me")
}

func (*Storage) GetStudyIDFromName(name string) (int, error) {
	panic("implement me")
}

func (*Storage) GetStudyIDFromTrialID(trialID int) (int, error) {
	panic("implement me")
}

func (*Storage) GetStudyNameFromID(studyID int) (string, error) {
	panic("implement me")
}

func (*Storage) GetStudyDirection(studyID int) (goptuna.StudyDirection, error) {
	panic("implement me")
}

func (*Storage) GetStudyUserAttrs(studyID int) (map[string]interface{}, error) {
	panic("implement me")
}

func (*Storage) GetStudySystemAttrs(studyID int) (map[string]interface{}, error) {
	panic("implement me")
}

func (*Storage) GetAllStudySummaries(studyID int) ([]goptuna.StudySummary, error) {
	panic("implement me")
}

func (*Storage) CreateNewTrialID(studyID int) (int, error) {
	panic("implement me")
}

func (*Storage) SetTrialValue(trialID int, value float64) error {
	panic("implement me")
}

func (*Storage) SetTrialIntermediateValue(trialID int, step int, value float64) error {
	panic("implement me")
}

func (*Storage) SetTrialParam(trialID int, paramName string, paramValueInternal float64,
	distribution goptuna.Distribution) error {
	panic("implement me")
}

func (*Storage) SetTrialState(trialID int, state goptuna.TrialState) error {
	panic("implement me")
}

func (*Storage) SetTrialUserAttr(trialID int, key string, value interface{}) error {
	panic("implement me")
}

func (*Storage) SetTrialSystemAttr(trialID int, key string, value interface{}) error {
	panic("implement me")
}

func (*Storage) GetTrialNumberFromID(trialID int) (int, error) {
	panic("implement me")
}

func (*Storage) GetTrialParam(trialID int, paramName string) (float64, error) {
	panic("implement me")
}

func (*Storage) GetTrial(trialID int) (goptuna.FrozenTrial, error) {
	panic("implement me")
}

func (*Storage) GetAllTrials(studyID int) ([]goptuna.FrozenTrial, error) {
	panic("implement me")
}

func (*Storage) GetBestTrial(studyID int) (goptuna.FrozenTrial, error) {
	panic("implement me")
}

func (*Storage) GetTrialParams(trialID int) (map[string]interface{}, error) {
	panic("implement me")
}

func (*Storage) GetTrialUserAttrs(trialID int) (map[string]interface{}, error) {
	panic("implement me")
}

func (*Storage) GetTrialSystemAttrs(trialID int) (map[string]interface{}, error) {
	panic("implement me")
}
