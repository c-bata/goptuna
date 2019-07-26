package goptuna

import (
	"errors"
	"fmt"
	"sync"
)

// Storage interface abstract a backend database and provide library
// internal interfaces to read/write history of studies and trials.
// This interface is not supposed to be directly accessed by library users.
type Storage interface {
	CreateNewStudyID(name string) (string, error)
	CreateNewTrialID(studyID string) (string, error)
	GetTrial(trialID string) (FrozenTrial, error)
	GetAllTrials(studyID string) ([]FrozenTrial, error)
	GetBestTrial(studyID string) (FrozenTrial, error)
	SetTrialValue(trialID string, value float64) error
	SetTrialParam(trialID string, paramName string, paramValueInternal float64) error
	SetTrialState(trialID string, state TrialState) error
	GetStudyDirection(studyID string) (StudyDirection, error)
	SetStudyDirection(studyID string, direction StudyDirection) error
}

// FrozenTrial holds the status and results of a Trial.
type FrozenTrial struct {
	ID         string             `json:"trial_id"`
	StudyID    string             `json:"study_id"`
	State      TrialState         `json:"state"`
	Value      float64            `json:"value"`
	ParamsInIR map[string]float64 `json:"params_in_internal_repr"`
}

var _ Storage = &InMemoryStorage{}

const inMemoryStudyID = "in_memory_storage_study_id"

var (
	// ErrInvalidStudyID represents invalid study id.
	ErrInvalidStudyID = errors.New("invalid study id")
	// ErrInvalidTrialID represents invalid trial id.
	ErrInvalidTrialID = errors.New("invalid trial id")
	// ErrTrialIsNotUpdated represents trial cannot be updated.
	ErrTrialIsNotUpdated = errors.New("trial cannot be updated")
	// ErrNoCompletedTrials represents no trials are completed yet.
	ErrNoCompletedTrials = errors.New("no trials are completed yet")
	// ErrUnknownDistribution returns the distribution is unknown.
	ErrUnknownDistribution = errors.New("unknown distribution")
)

// NewInMemoryStorage returns new memory storage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		direction: StudyDirectionMinimize,
		trials:    make(map[string]FrozenTrial, 128),
	}
}

// InMemoryStorage stores data in memory of the Go process.
type InMemoryStorage struct {
	mu sync.RWMutex

	direction StudyDirection
	trials    map[string]FrozenTrial
}

// GetAllTrials returns the all trials.
func (s *InMemoryStorage) GetAllTrials(studyID string) ([]FrozenTrial, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	trials := make([]FrozenTrial, 0, len(s.trials))

	for k := range s.trials {
		trials = append(trials, s.trials[k])
	}
	return trials, nil
}

// SetTrialParam sets the sampled parameters of trial.
func (s *InMemoryStorage) SetTrialParam(trialID string, paramName string, paramValueInternal float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trial, ok := s.trials[trialID]
	if !ok {
		return ErrInvalidTrialID
	}
	if trial.State.IsFinished() {
		return ErrTrialIsNotUpdated
	}
	if trial.ParamsInIR == nil {
		trial.ParamsInIR = make(map[string]float64, 8)
	}
	trial.ParamsInIR[paramName] = paramValueInternal
	s.trials[trialID] = trial
	return nil
}

// SetTrialState sets the state of trial.
func (s *InMemoryStorage) SetTrialState(trialID string, state TrialState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trial, ok := s.trials[trialID]
	if !ok {
		return ErrInvalidTrialID
	}
	if trial.State.IsFinished() {
		return ErrTrialIsNotUpdated
	}
	trial.State = state
	s.trials[trialID] = trial
	return nil
}

// SetTrialValue sets the value of trial.
func (s *InMemoryStorage) SetTrialValue(trialID string, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trial, ok := s.trials[trialID]
	if !ok {
		return ErrInvalidTrialID
	}
	if trial.State.IsFinished() {
		return ErrTrialIsNotUpdated
	}
	trial.Value = value
	s.trials[trialID] = trial
	return nil
}

// CreateNewStudyID creates study and returns studyID.
func (s *InMemoryStorage) CreateNewStudyID(name string) (string, error) {
	return inMemoryStudyID, nil
}

func (s *InMemoryStorage) checkStudyID(studyID string) bool {
	return studyID == inMemoryStudyID
}

// CreateNewTrialID creates trial and returns trialID.
func (s *InMemoryStorage) CreateNewTrialID(studyID string) (string, error) {
	if !s.checkStudyID(studyID) {
		return "", ErrInvalidStudyID
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	trialID := fmt.Sprintf("trial_%d", len(s.trials))
	s.trials[trialID] = FrozenTrial{
		ID:         trialID,
		StudyID:    "",
		State:      TrialStateRunning,
		Value:      0,
		ParamsInIR: nil,
	}
	return trialID, nil
}

// GetBestTrial returns the best trial.
func (s *InMemoryStorage) GetBestTrial(studyID string) (FrozenTrial, error) {
	if !s.checkStudyID(studyID) {
		return FrozenTrial{}, ErrInvalidStudyID
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	var bestTrial FrozenTrial
	var bestTrialIsSet bool
	for i := range s.trials {
		if s.trials[i].State != TrialStateComplete {
			continue
		}

		if s.direction == StudyDirectionMaximize {
			if !bestTrialIsSet {
				bestTrial = s.trials[i]
				bestTrialIsSet = true
			} else if s.trials[i].Value > bestTrial.Value {
				bestTrial = s.trials[i]
			}
		} else if s.direction == StudyDirectionMinimize {
			if !bestTrialIsSet {
				bestTrial = s.trials[i]
				bestTrialIsSet = true
			} else if s.trials[i].Value < bestTrial.Value {
				bestTrial = s.trials[i]
			}
		}
	}
	if !bestTrialIsSet {
		return FrozenTrial{}, ErrNoCompletedTrials
	}
	return bestTrial, nil
}

// SetStudyDirection sets study direction of the objective.
func (s *InMemoryStorage) SetStudyDirection(studyID string, direction StudyDirection) error {
	if !s.checkStudyID(studyID) {
		return ErrInvalidStudyID
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.direction = direction
	return nil
}

// GetStudyDirection returns study direction of the objective.
func (s *InMemoryStorage) GetStudyDirection(studyID string) (StudyDirection, error) {
	if !s.checkStudyID(studyID) {
		return StudyDirectionMinimize, ErrInvalidStudyID
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.direction, nil
}

// GetTrial returns Trial.
func (s *InMemoryStorage) GetTrial(trialID string) (FrozenTrial, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.trials[trialID], nil
}
