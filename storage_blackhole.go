package goptuna

import (
	"errors"
	"sync"
)

var (
	// ErrTrialsAlreadyDeleted means that trial is already deleted.
	ErrTrialsAlreadyDeleted = errors.New("trial is already deleted")
	// ErrTrialsPartiallyDeleted means that trials are partially deleted.
	ErrTrialsPartiallyDeleted = errors.New("some trials are already deleted")
	// ErrDeleteNonFinishedTrial means that non finished trial is deleted.
	ErrDeleteNonFinishedTrial = errors.New("non finished trial is deleted")
)

var _ Storage = &BlackHoleStorage{}

func NewBlackholeStorage(n int) *BlackHoleStorage {
	return &BlackHoleStorage{
		direction:   StudyDirectionMinimize,
		counter:     0,
		nTrials:     n,
		trials:      make([]FrozenTrial, n),
		bestTrial:   FrozenTrial{},
		userAttrs:   make(map[string]string, 8),
		systemAttrs: make(map[string]string, 8),
		studyName:   DefaultStudyNamePrefix + InMemoryStorageStudyUUID,
	}
}

// BlackholeStorage is an in-memory storage, but designed for over 10k+ trials.
// This storage just holds the latest N trials you specified.
type BlackHoleStorage struct {
	direction   StudyDirection
	counter     int
	nTrials     int
	trials      []FrozenTrial
	bestTrial   FrozenTrial
	userAttrs   map[string]string
	systemAttrs map[string]string
	studyName   string
	mu          sync.RWMutex
}

func (s *BlackHoleStorage) CreateNewStudy(name string) (int, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) DeleteStudy(studyID int) error {
	panic("implement me")
}

func (s *BlackHoleStorage) SetStudyDirection(studyID int, direction StudyDirection) error {
	panic("implement me")
}

func (s *BlackHoleStorage) SetStudyUserAttr(studyID int, key string, value string) error {
	panic("implement me")
}

func (s *BlackHoleStorage) SetStudySystemAttr(studyID int, key string, value string) error {
	panic("implement me")
}

func (s *BlackHoleStorage) GetStudyIDFromName(name string) (int, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetStudyIDFromTrialID(trialID int) (int, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetStudyNameFromID(studyID int) (string, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetStudyDirection(studyID int) (StudyDirection, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetStudyUserAttrs(studyID int) (map[string]string, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetStudySystemAttrs(studyID int) (map[string]string, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetAllStudySummaries() ([]StudySummary, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) CreateNewTrial(studyID int) (int, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) CloneTrial(studyID int, baseTrial FrozenTrial) (int, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) SetTrialValue(trialID int, value float64) error {
	panic("implement me")
}

func (s *BlackHoleStorage) SetTrialIntermediateValue(trialID int, step int, value float64) error {
	panic("implement me")
}

func (s *BlackHoleStorage) SetTrialParam(trialID int, paramName string, paramValueInternal float64,
	distribution interface{}) error {
	panic("implement me")
}

func (s *BlackHoleStorage) SetTrialState(trialID int, state TrialState) error {
	panic("implement me")
}

func (s *BlackHoleStorage) SetTrialUserAttr(trialID int, key string, value string) error {
	panic("implement me")
}

func (s *BlackHoleStorage) SetTrialSystemAttr(trialID int, key string, value string) error {
	panic("implement me")
}

func (s *BlackHoleStorage) GetTrialNumberFromID(trialID int) (int, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetTrialParam(trialID int, paramName string) (float64, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetTrial(trialID int) (FrozenTrial, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetAllTrials(studyID int) ([]FrozenTrial, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetBestTrial(studyID int) (FrozenTrial, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetTrialParams(trialID int) (map[string]interface{}, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetTrialUserAttrs(trialID int) (map[string]string, error) {
	panic("implement me")
}

func (s *BlackHoleStorage) GetTrialSystemAttrs(trialID int) (map[string]string, error) {
	panic("implement me")
}
