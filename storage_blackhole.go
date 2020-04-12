package goptuna

import (
	"errors"
	"sync"
	"time"

	"github.com/c-bata/goptuna"
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

// BlackholeStorage is an in-memory storage, but designed for over 10k trials.
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

// CreateNewStudy creates study and returns studyID.
func (s *BlackHoleStorage) CreateNewStudy(name string) (int, error) {
	if name != "" {
		s.studyName = name
	}
	return InMemoryStorageStudyID, nil
}

// DeleteStudy deletes a study.
func (s *BlackHoleStorage) DeleteStudy(studyID int) error {
	if !s.checkStudyID(studyID) {
		return goptuna.ErrInvalidStudyID
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.direction = goptuna.StudyDirectionMinimize
	s.trials = make([]goptuna.FrozenTrial, 0, 128)
	s.userAttrs = make(map[string]string, 8)
	s.systemAttrs = make(map[string]string, 8)
	s.studyName = DefaultStudyNamePrefix + InMemoryStorageStudyUUID
	return nil
}

// SetStudyDirection sets study direction of the objective.
func (s *BlackHoleStorage) SetStudyDirection(studyID int, direction StudyDirection) error {
	if !s.checkStudyID(studyID) {
		return goptuna.ErrInvalidStudyID
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.direction = direction
	return nil
}

// SetStudyUserAttr to store the value for the user.
func (s *BlackHoleStorage) SetStudyUserAttr(studyID int, key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userAttrs[key] = value
	return nil
}

// SetStudySystemAttr to store the value for the system.
func (s *BlackHoleStorage) SetStudySystemAttr(studyID int, key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.systemAttrs[key] = value
	return nil
}

// GetStudyIDFromName return the study id from study name.
func (s *BlackHoleStorage) GetStudyIDFromName(name string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if name != s.studyName {
		return -1, goptuna.ErrNotFound
	}
	return InMemoryStorageStudyID, nil
}

// GetStudyIDFromTrialID return the study id from trial id.
func (s *BlackHoleStorage) GetStudyIDFromTrialID(trialID int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// counter represents an id for next trial.
	if s.counter > trialID {
		return InMemoryStorageStudyID, nil
	}
	return -1, goptuna.ErrNotFound
}

// GetStudyNameFromID return the study name from study id.
func (s *BlackHoleStorage) GetStudyNameFromID(studyID int) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.checkStudyID(studyID) {
		return "", goptuna.ErrNotFound
	}
	return s.studyName, nil
}

// GetStudyDirection returns study direction of the objective.
func (s *BlackHoleStorage) GetStudyDirection(studyID int) (StudyDirection, error) {
	if !s.checkStudyID(studyID) {
		return goptuna.StudyDirectionMinimize, goptuna.ErrInvalidStudyID
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.direction, nil
}

// GetStudyUserAttrs to restore the attributes for the user.
func (s *BlackHoleStorage) GetStudyUserAttrs(studyID int) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := make(map[string]string, len(s.userAttrs))
	for k := range s.userAttrs {
		n[k] = s.userAttrs[k]
	}
	return n, nil
}

// GetStudySystemAttrs to restore the attributes for the system.
func (s *BlackHoleStorage) GetStudySystemAttrs(studyID int) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := make(map[string]string, len(s.systemAttrs))
	for k := range s.systemAttrs {
		n[k] = s.systemAttrs[k]
	}
	return n, nil
}

// GetAllStudySummaries returns all study summaries.
func (s *BlackHoleStorage) GetAllStudySummaries() ([]StudySummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var datetimeStart time.Time
	sa := make(map[string]string, len(s.systemAttrs))
	for k := range s.systemAttrs {
		sa[k] = s.systemAttrs[k]
	}
	ua := make(map[string]string, len(s.userAttrs))
	for k := range s.userAttrs {
		ua[k] = s.userAttrs[k]
	}

	return []StudySummary{
		{
			ID:            InMemoryStorageStudyID,
			Name:          s.studyName,
			Direction:     s.direction,
			BestTrial:     s.bestTrial,
			UserAttrs:     ua,
			SystemAttrs:   sa,
			DatetimeStart: datetimeStart,
		},
	}, nil
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

func (s *BlackHoleStorage) checkStudyID(studyID int) bool {
	return studyID == InMemoryStorageStudyID
}
