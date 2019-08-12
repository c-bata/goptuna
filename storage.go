package goptuna

import (
	"errors"
	"sync"
	"time"
)

// DefaultStudyNamePrefix is a prefix of the default study name.
var DefaultStudyNamePrefix = "no-name-"

// Storage interface abstract a backend database and provide library
// internal interfaces to read/write history of studies and trials.
// This interface is not supposed to be directly accessed by library users.
type Storage interface {
	// Basic study manipulation
	CreateNewStudyID(name string) (int, error)
	SetStudyDirection(studyID int, direction StudyDirection) error
	SetStudyUserAttr(studyID int, key string, value string) error
	SetStudySystemAttr(studyID int, key string, value string) error
	// Basic study access
	GetStudyIDFromName(name string) (int, error)
	GetStudyIDFromTrialID(trialID int) (int, error)
	GetStudyNameFromID(studyID int) (string, error)
	GetStudyDirection(studyID int) (StudyDirection, error)
	GetStudyUserAttrs(studyID int) (map[string]string, error)
	GetStudySystemAttrs(studyID int) (map[string]string, error)
	GetAllStudySummaries() ([]StudySummary, error)
	// Basic trial manipulation
	CreateNewTrialID(studyID int) (int, error)
	SetTrialValue(trialID int, value float64) error
	SetTrialIntermediateValue(trialID int, step int, value float64) error
	SetTrialParam(trialID int, paramName string, paramValueInternal float64,
		distribution interface{}) error
	SetTrialState(trialID int, state TrialState) error
	SetTrialUserAttr(trialID int, key string, value string) error
	SetTrialSystemAttr(trialID int, key string, value string) error
	// Basic trial access
	GetTrialNumberFromID(trialID int) (int, error)
	GetTrialParam(trialID int, paramName string) (float64, error)
	GetTrial(trialID int) (FrozenTrial, error)
	GetAllTrials(studyID int) ([]FrozenTrial, error)
	GetBestTrial(studyID int) (FrozenTrial, error)
	GetTrialParams(trialID int) (map[string]interface{}, error)
	GetTrialUserAttrs(trialID int) (map[string]string, error)
	GetTrialSystemAttrs(trialID int) (map[string]string, error)
}

// StudySummary holds basic attributes and aggregated results of Study.
type StudySummary struct {
	ID            int               `json:"study_id"`
	Name          string            `json:"study_name"`
	Direction     StudyDirection    `json:"direction"`
	BestTrial     FrozenTrial       `json:"best_trial"`
	UserAttrs     map[string]string `json:"user_attrs"`
	SystemAttrs   map[string]string `json:"system_attrs"`
	DatetimeStart time.Time         `json:"datetime_start"`
}

// FrozenTrial holds the status and results of a Trial.
type FrozenTrial struct {
	ID                 int                    `json:"trial_id"`
	StudyID            int                    `json:"study_id"`
	Number             int                    `json:"number"`
	State              TrialState             `json:"state"`
	Value              float64                `json:"value"`
	IntermediateValues map[int]float64        `json:"intermediate_values"`
	DatetimeStart      time.Time              `json:"datetime_start"`
	DatetimeComplete   time.Time              `json:"datetime_complete"`
	Params             map[string]interface{} `json:"params"`
	Distributions      map[string]interface{} `json:"distributions"`
	UserAttrs          map[string]string      `json:"user_attrs"`
	SystemAttrs        map[string]string      `json:"system_attrs"`
	// Note: ParamsInIR is private in Optuna.
	// But we need to keep public because this is accessed by TPE sampler.
	// It couldn't access internal attributes from the external packages.
	// https://github.com/pfnet/optuna/pull/462
	ParamsInIR map[string]float64 `json:"params_in_internal_repr"`
}

var _ Storage = &InMemoryStorage{}

// InMemoryStorageStudyID is a study id for in memory storage backend.
const InMemoryStorageStudyID = 1

// InMemoryStorageStudyUUID is a UUID for in memory storage backend
const InMemoryStorageStudyUUID = "00000000-0000-0000-0000-000000000000"

var (
	// ErrNotFound represents not found.
	ErrNotFound = errors.New("not found")
	// ErrInvalidStudyID represents invalid study id.
	ErrInvalidStudyID = errors.New("invalid study id")
	// ErrInvalidTrialID represents invalid trial id.
	ErrInvalidTrialID = errors.New("invalid trial id")
	// ErrTrialCannotBeUpdated represents trial cannot be updated.
	ErrTrialCannotBeUpdated = errors.New("trial cannot be updated")
	// ErrNoCompletedTrials represents no trials are completed yet.
	ErrNoCompletedTrials = errors.New("no trials are completed yet")
	// ErrUnknownDistribution returns the distribution is unknown.
	ErrUnknownDistribution = errors.New("unknown distribution")
	// ErrTrialPruned represents the pruned.
	ErrTrialPruned = errors.New("trial is pruned")
)

// NewInMemoryStorage returns new memory storage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		direction:   StudyDirectionMinimize,
		trials:      make([]FrozenTrial, 0, 128),
		userAttrs:   make(map[string]string, 8),
		systemAttrs: make(map[string]string, 8),
		studyName:   DefaultStudyNamePrefix + InMemoryStorageStudyUUID,
	}
}

// InMemoryStorage stores data in memory of the Go process.
type InMemoryStorage struct {
	direction   StudyDirection
	trials      []FrozenTrial
	userAttrs   map[string]string
	systemAttrs map[string]string
	studyName   string

	mu sync.RWMutex
}

// CreateNewStudyID creates study and returns studyID.
func (s *InMemoryStorage) CreateNewStudyID(name string) (int, error) {
	if name != "" {
		s.studyName = name
	}
	return InMemoryStorageStudyID, nil
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

// SetStudyUserAttr to store the value for the user.
func (s *InMemoryStorage) SetStudyUserAttr(studyID int, key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userAttrs[key] = value
	return nil
}

// SetStudySystemAttr to store the value for the system.
func (s *InMemoryStorage) SetStudySystemAttr(studyID int, key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.systemAttrs[key] = value
	return nil
}

// GetStudyIDFromName return the study id from study name.
func (s *InMemoryStorage) GetStudyIDFromName(name string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if name != s.studyName {
		return -1, ErrNotFound
	}
	return InMemoryStorageStudyID, nil
}

// GetStudyIDFromTrialID return the study id from trial id.
func (s *InMemoryStorage) GetStudyIDFromTrialID(trialID int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.trials {
		if s.trials[i].ID == trialID {
			return InMemoryStorageStudyID, nil
		}
	}
	return -1, ErrNotFound
}

// GetStudyNameFromID return the study name from study id.
func (s *InMemoryStorage) GetStudyNameFromID(studyID int) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.checkStudyID(studyID) {
		return "", ErrNotFound
	}
	return s.studyName, nil
}

// GetStudyUserAttrs to restore the attributes for the user.
func (s *InMemoryStorage) GetStudyUserAttrs(studyID int) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := make(map[string]string, len(s.userAttrs))
	for k := range s.userAttrs {
		n[k] = s.userAttrs[k]
	}
	return n, nil
}

// GetStudySystemAttrs to restore the attributes for the system.
func (s *InMemoryStorage) GetStudySystemAttrs(studyID int) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := make(map[string]string, len(s.systemAttrs))
	for k := range s.systemAttrs {
		n[k] = s.systemAttrs[k]
	}
	return n, nil
}

// GetAllStudySummaries returns all study summaries.
func (s *InMemoryStorage) GetAllStudySummaries() ([]StudySummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var datetimeStart time.Time

	var foundBestTrial bool
	var bestTrial FrozenTrial
	for i, t := range s.trials {
		if i == 0 {
			datetimeStart = t.DatetimeStart
		}

		if datetimeStart.Unix() > t.DatetimeStart.Unix() {
			datetimeStart = t.DatetimeStart
		}

		if t.State != TrialStateComplete {
			continue
		}

		if !foundBestTrial {
			bestTrial = t
			foundBestTrial = true
			continue
		}

		if s.direction == StudyDirectionMaximize {
			if t.Value > bestTrial.Value {
				bestTrial = t
			}
		} else {
			if t.Value < bestTrial.Value {
				bestTrial = t
			}
		}
	}

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
			BestTrial:     bestTrial,
			UserAttrs:     ua,
			SystemAttrs:   sa,
			DatetimeStart: datetimeStart,
		},
	}, nil
}

func (s *InMemoryStorage) checkStudyID(studyID int) bool {
	return studyID == InMemoryStorageStudyID
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
	s.trials = append(s.trials, FrozenTrial{
		ID:                 number,
		Number:             number,
		State:              TrialStateRunning,
		Value:              0,
		IntermediateValues: make(map[int]float64, 8),
		DatetimeStart:      time.Now(),
		DatetimeComplete:   time.Time{},
		Params:             make(map[string]interface{}, 8),
		Distributions:      make(map[string]interface{}, 8),
		UserAttrs:          make(map[string]string, 8),
		SystemAttrs:        make(map[string]string, 8),
		ParamsInIR:         make(map[string]float64, 8),
	})
	return trialID, nil
}

// SetTrialValue sets the value of trial.
func (s *InMemoryStorage) SetTrialValue(trialID int, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if trialID >= len(s.trials) {
		return ErrInvalidTrialID
	}
	trial := s.trials[trialID]
	if trial.State.IsFinished() {
		return ErrTrialCannotBeUpdated
	}
	trial.Value = value
	s.trials[trialID] = trial
	return nil
}

// SetTrialIntermediateValue sets the intermediate value of trial.
func (s *InMemoryStorage) SetTrialIntermediateValue(trialID int, step int, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if trialID >= len(s.trials) {
		return ErrInvalidTrialID
	}
	trial := s.trials[trialID]
	if trial.State.IsFinished() {
		return ErrTrialCannotBeUpdated
	}

	for key := range trial.IntermediateValues {
		if key == step {
			return errors.New("step value is already exist")
		}
	}
	trial.IntermediateValues[step] = value
	s.trials[trialID] = trial
	return nil
}

// SetTrialParam sets the sampled parameters of trial.
func (s *InMemoryStorage) SetTrialParam(
	trialID int,
	paramName string,
	paramValueInternal float64,
	distribution interface{}) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check param has not been set; otherwise, return error
	if trialID >= len(s.trials) {
		return ErrInvalidTrialID
	}
	trial := s.trials[trialID]

	// Check trial is able to update
	if trial.State.IsFinished() {
		return ErrTrialCannotBeUpdated
	}

	// Set param distribution
	trial.Distributions[paramName] = distribution
	trial.ParamsInIR[paramName] = paramValueInternal
	var err error
	trial.Params[paramName], err = ToExternalRepresentation(distribution, paramValueInternal)
	if err != nil {
		return err
	}

	s.trials[trialID] = trial
	return nil
}

// SetTrialState sets the state of trial.
func (s *InMemoryStorage) SetTrialState(trialID int, state TrialState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if trialID >= len(s.trials) {
		return ErrInvalidTrialID
	}
	trial := s.trials[trialID]
	if trial.State.IsFinished() {
		return ErrTrialCannotBeUpdated
	}
	trial.State = state
	if trial.State.IsFinished() {
		trial.DatetimeComplete = time.Now()
	}
	s.trials[trialID] = trial
	return nil
}

// SetTrialUserAttr to store the value for the user.
func (s *InMemoryStorage) SetTrialUserAttr(trialID int, key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.trials {
		if s.trials[i].ID == trialID && s.trials[i].State != TrialStateComplete {
			s.trials[i].UserAttrs[key] = value
			return nil
		}
	}
	return ErrInvalidTrialID
}

// SetTrialSystemAttr to store the value for the system.
func (s *InMemoryStorage) SetTrialSystemAttr(trialID int, key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.trials {
		if s.trials[i].ID == trialID && s.trials[i].State != TrialStateComplete {
			s.trials[i].SystemAttrs[key] = value
			return nil
		}
	}
	return ErrInvalidTrialID
}

// GetTrialNumberFromID returns the trial's number.
func (s *InMemoryStorage) GetTrialNumberFromID(trialID int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.trials {
		if s.trials[i].ID == trialID {
			return trialID, nil
		}
	}
	return -1, ErrInvalidTrialID
}

// GetTrialParam returns the internal parameter of the trial
func (s *InMemoryStorage) GetTrialParam(trialID int, paramName string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.trials {
		if s.trials[i].ID == trialID {
			paramIR, ok := s.trials[i].ParamsInIR[paramName]
			if !ok {
				return -1.0, errors.New("param doesn't exist")
			}
			return paramIR, nil
		}
	}
	return -1, ErrInvalidTrialID
}

// GetTrialParams returns the external parameters in the trial
func (s *InMemoryStorage) GetTrialParams(trialID int) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.trials {
		if s.trials[i].ID == trialID {
			return s.trials[i].Params, nil
		}
	}
	return nil, ErrInvalidTrialID
}

// GetTrialUserAttrs to restore the attributes for the user.
func (s *InMemoryStorage) GetTrialUserAttrs(trialID int) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, t := range s.trials {
		if t.ID == trialID {
			n := make(map[string]string, len(t.UserAttrs))
			for k := range t.UserAttrs {
				n[k] = t.UserAttrs[k]
			}
			return n, nil
		}
	}
	return nil, ErrNotFound
}

// GetTrialSystemAttrs to restore the attributes for the system.
func (s *InMemoryStorage) GetTrialSystemAttrs(trialID int) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, t := range s.trials {
		if t.ID == trialID {
			n := make(map[string]string, len(t.SystemAttrs))
			for k := range t.SystemAttrs {
				n[k] = t.SystemAttrs[k]
			}
			return n, nil
		}
	}
	return nil, ErrNotFound
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

// GetAllTrials returns the all trials.
func (s *InMemoryStorage) GetAllTrials(studyID int) ([]FrozenTrial, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	trials := make([]FrozenTrial, 0, len(s.trials))

	for i := range s.trials {
		trials = append(trials, s.trials[i])
	}
	return trials, nil
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
