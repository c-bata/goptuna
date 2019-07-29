package goptuna

import (
	"errors"
	"sync"
	"time"
)

// Storage interface abstract a backend database and provide library
// internal interfaces to read/write history of studies and trials.
// This interface is not supposed to be directly accessed by library users.
type Storage interface {
	CreateNewStudyID(name string) (int, error)
	CreateNewTrialID(studyID int) (int, error)
	GetTrial(trialID int) (FrozenTrial, error)
	GetAllTrials(studyID int) ([]FrozenTrial, error)
	GetBestTrial(studyID int) (FrozenTrial, error)
	SetTrialValue(trialID int, value float64) error
	SetTrialParam(trialID int, paramName string, paramValueInternal float64,
		distribution Distribution) error
	SetTrialState(trialID int, state TrialState) error
	GetStudyDirection(studyID int) (StudyDirection, error)
	SetStudyDirection(studyID int, direction StudyDirection) error
}

// StudySummary holds basic attributes and aggregated results of Study.
type StudySummary struct {
	ID            int                    `json:"study_id"`
	Direction     StudyDirection         `json:"direction"`
	BestTrial     FrozenTrial            `json:"best_trial"`
	UserAttrs     map[string]interface{} `json:"user_attrs"`
	SystemAttrs   map[string]interface{} `json:"system_attrs"`
	DatetimeStart *time.Time             `json:"datetime_start"`
}

// FrozenTrial holds the status and results of a Trial.
type FrozenTrial struct {
	ID                 int                     `json:"trial_id"`
	Number             int                     `json:"number"`
	State              TrialState              `json:"state"`
	Value              float64                 `json:"value"`
	DatetimeStart      time.Time               `json:"datetime_start"`
	DatetimeComplete   time.Time               `json:"datetime_complete"`
	Params             map[string]interface{}  `json:"params"`
	Distributions      map[string]Distribution `json:"distributions"`
	UserAttrs          map[string]interface{}  `json:"user_attrs"`
	SystemAttrs        map[string]interface{}  `json:"system_attrs"`
	IntermediateValues map[int]float64         `json:"intermediate_values"`
	// Note: ParamsInIR is private in Optuna.
	// But we need to keep public because this is accessed by TPE sampler.
	// It couldn't access internal attributes from the external packages.
	// https://github.com/pfnet/optuna/pull/462
	ParamsInIR map[string]float64 `json:"params_in_internal_repr"`
}

var _ Storage = &InMemoryStorage{}

const inMemoryStudyID = 1

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
		trials:    make(map[int]FrozenTrial, 128),
	}
}

// InMemoryStorage stores data in memory of the Go process.
type InMemoryStorage struct {
	mu sync.RWMutex

	direction StudyDirection
	trials    map[int]FrozenTrial
}

// GetAllTrials returns the all trials.
func (s *InMemoryStorage) GetAllTrials(studyID int) ([]FrozenTrial, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	trials := make([]FrozenTrial, 0, len(s.trials))

	for k := range s.trials {
		trials = append(trials, s.trials[k])
	}
	return trials, nil
}

// SetTrialParam sets the sampled parameters of trial.
func (s *InMemoryStorage) SetTrialParam(
	trialID int,
	paramName string,
	paramValueInternal float64,
	distribution Distribution) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check param has not been set; otherwise, return error
	trial, ok := s.trials[trialID]
	if !ok {
		return ErrInvalidTrialID
	}

	// Check trial is able to update
	if trial.State.IsFinished() {
		return ErrTrialIsNotUpdated
	}

	// Set param distribution
	if trial.Distributions == nil {
		trial.Distributions = make(map[string]Distribution, 8)
	}
	trial.Distributions[paramName] = distribution
	// Set parameter in internal representations
	if trial.ParamsInIR == nil {
		trial.ParamsInIR = make(map[string]float64, 8)
	}
	trial.ParamsInIR[paramName] = paramValueInternal
	// Set parameter in external representations
	if trial.Params == nil {
		trial.Params = make(map[string]interface{}, 8)
	}
	trial.Params[paramName] = distribution.ToExternalRepr(paramValueInternal)

	s.trials[trialID] = trial
	return nil
}

// SetTrialState sets the state of trial.
func (s *InMemoryStorage) SetTrialState(trialID int, state TrialState) error {
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
func (s *InMemoryStorage) SetTrialValue(trialID int, value float64) error {
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
func (s *InMemoryStorage) CreateNewStudyID(name string) (int, error) {
	return inMemoryStudyID, nil
}

func (s *InMemoryStorage) checkStudyID(studyID int) bool {
	return studyID == inMemoryStudyID
}

// CreateNewTrialID creates trial and returns trialID.
func (s *InMemoryStorage) CreateNewTrialID(studyID int) (int, error) {
	if !s.checkStudyID(studyID) {
		return -1, ErrInvalidStudyID
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	number := len(s.trials)
	// trialID equals the number because InMemoryStorage has only 1 study.
	trialID := number
	s.trials[trialID] = FrozenTrial{
		ID:                 number,
		Number:             number,
		State:              TrialStateRunning,
		Value:              0,
		DatetimeStart:      time.Now(),
		DatetimeComplete:   time.Time{},
		Params:             nil,
		Distributions:      nil,
		UserAttrs:          nil,
		SystemAttrs:        nil,
		IntermediateValues: nil,
		ParamsInIR:         nil,
	}
	return trialID, nil
}

// GetBestTrial returns the best trial.
func (s *InMemoryStorage) GetBestTrial(studyID int) (FrozenTrial, error) {
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
func (s *InMemoryStorage) SetStudyDirection(studyID int, direction StudyDirection) error {
	if !s.checkStudyID(studyID) {
		return ErrInvalidStudyID
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.direction = direction
	return nil
}

// GetStudyDirection returns study direction of the objective.
func (s *InMemoryStorage) GetStudyDirection(studyID int) (StudyDirection, error) {
	if !s.checkStudyID(studyID) {
		return StudyDirectionMinimize, ErrInvalidStudyID
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.direction, nil
}

// GetTrial returns Trial.
func (s *InMemoryStorage) GetTrial(trialID int) (FrozenTrial, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.trials[trialID], nil
}
