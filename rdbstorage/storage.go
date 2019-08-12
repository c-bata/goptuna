package rdbstorage

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/c-bata/goptuna"
)

var _ goptuna.Storage = &Storage{}

// NewStorage returns new RDB storage.
func NewStorage(db *gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// Storage stores data in your relational databases.
type Storage struct {
	db *gorm.DB
}

// CreateNewStudyID creates study and returns studyID.
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

// SetStudyDirection sets study direction of the objective.
func (s *Storage) SetStudyDirection(studyID int, direction goptuna.StudyDirection) error {
	panic("implement me")
}

// SetStudyUserAttr to store the value for the user.
func (s *Storage) SetStudyUserAttr(studyID int, key string, value interface{}) error {
	panic("implement me")
}

// SetStudySystemAttr to store the value for the system.
func (s *Storage) SetStudySystemAttr(studyID int, key string, value interface{}) error {
	panic("implement me")
}

// GetStudyIDFromName return the study id from study name.
func (s *Storage) GetStudyIDFromName(name string) (int, error) {
	panic("implement me")
}

// GetStudyIDFromTrialID return the study id from trial id.
func (s *Storage) GetStudyIDFromTrialID(trialID int) (int, error) {
	panic("implement me")
}

// GetStudyNameFromID return the study name from study id.
func (s *Storage) GetStudyNameFromID(studyID int) (string, error) {
	panic("implement me")
}

// GetStudyUserAttrs to restore the attributes for the user.
func (s *Storage) GetStudyUserAttrs(studyID int) (map[string]interface{}, error) {
	panic("implement me")
}

// GetStudySystemAttrs to restore the attributes for the system.
func (s *Storage) GetStudySystemAttrs(studyID int) (map[string]interface{}, error) {
	panic("implement me")
}

// GetAllStudySummaries returns all study summaries.
func (s *Storage) GetAllStudySummaries(studyID int) ([]goptuna.StudySummary, error) {
	panic("implement me")
}

// CreateNewTrialID creates trial and returns trialID.
func (s *Storage) CreateNewTrialID(studyID int) (int, error) {
	panic("implement me")
}

// SetTrialValue sets the value of trial.
func (s *Storage) SetTrialValue(trialID int, value float64) error {
	panic("implement me")
}

// SetTrialIntermediateValue sets the intermediate value of trial.
func (s *Storage) SetTrialIntermediateValue(trialID int, step int, value float64) error {
	panic("implement me")
}

// SetTrialParam sets the sampled parameters of trial.
func (s *Storage) SetTrialParam(trialID int, paramName string, paramValueInternal float64,
	distribution goptuna.Distribution) error {
	panic("implement me")
}

// SetTrialState sets the state of trial.
func (s *Storage) SetTrialState(trialID int, state goptuna.TrialState) error {
	panic("implement me")
}

// SetTrialUserAttr to store the value for the user.
func (s *Storage) SetTrialUserAttr(trialID int, key string, value interface{}) error {
	panic("implement me")
}

// SetTrialSystemAttr to store the value for the system.
func (s *Storage) SetTrialSystemAttr(trialID int, key string, value interface{}) error {
	panic("implement me")
}

// GetTrialNumberFromID returns the trial's number.
func (s *Storage) GetTrialNumberFromID(trialID int) (int, error) {
	panic("implement me")
}

// GetTrialParam returns the internal parameter of the trial
func (s *Storage) GetTrialParam(trialID int, paramName string) (float64, error) {
	panic("implement me")
}

// GetTrialParams returns the external parameters in the trial
func (s *Storage) GetTrialParams(trialID int) (map[string]interface{}, error) {
	panic("implement me")
}

// GetTrialUserAttrs to restore the attributes for the user.
func (s *Storage) GetTrialUserAttrs(trialID int) (map[string]interface{}, error) {
	panic("implement me")
}

// GetTrialSystemAttrs to restore the attributes for the system.
func (s *Storage) GetTrialSystemAttrs(trialID int) (map[string]interface{}, error) {
	panic("implement me")
}

// GetBestTrial returns the best trial.
func (s *Storage) GetBestTrial(studyID int) (goptuna.FrozenTrial, error) {
	panic("implement me")
}

// GetAllTrials returns the all trials.
func (s *Storage) GetAllTrials(studyID int) ([]goptuna.FrozenTrial, error) {
	panic("implement me")
}

// GetStudyDirection returns study direction of the objective.
func (s *Storage) GetStudyDirection(studyID int) (goptuna.StudyDirection, error) {
	panic("implement me")
}

// GetTrial returns Trial.
func (s *Storage) GetTrial(trialID int) (goptuna.FrozenTrial, error) {
	panic("implement me")
}
