package goptuna

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrTrialAlreadyDeleted means that trial is already deleted.
	ErrTrialAlreadyDeleted = errors.New("trial is already deleted")
	// ErrTrialsPartiallyDeleted means that trials are partially deleted.
	ErrTrialsPartiallyDeleted = errors.New("some trials are already deleted")
	// ErrDeleteNonFinishedTrial means that non finished trial is deleted.
	ErrDeleteNonFinishedTrial = errors.New("non finished trial is deleted")
)

var _ Storage = &BlackHoleStorage{}

// NewBlackHoleStorage returns BlackHoleStorage.
func NewBlackHoleStorage(n int) *BlackHoleStorage {
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

// BlackHoleStorage is an in-memory storage, but designed for over 100k trials.
// Please note that this storage just holds the 'nTrials' trials.
//
// Methods to create a trial might return ErrDeleteNonFinishedTrial.
// GetAllTrials method might return ErrTrialsPartiallyDeleted.
// Methods to get or update a trial might return ErrTrialAlreadyDeleted.
//
// Currently, RandomSampler and CMA-ES sampler supports this storage.
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
		return ErrInvalidStudyID
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.direction = StudyDirectionMinimize
	s.trials = make([]FrozenTrial, 0, 128)
	s.userAttrs = make(map[string]string, 8)
	s.systemAttrs = make(map[string]string, 8)
	s.studyName = DefaultStudyNamePrefix + InMemoryStorageStudyUUID
	return nil
}

// SetStudyDirection sets study direction of the objective.
func (s *BlackHoleStorage) SetStudyDirection(studyID int, direction StudyDirection) error {
	if !s.checkStudyID(studyID) {
		return ErrInvalidStudyID
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
		return -1, ErrNotFound
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
	return -1, ErrNotFound
}

// GetStudyNameFromID return the study name from study id.
func (s *BlackHoleStorage) GetStudyNameFromID(studyID int) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.checkStudyID(studyID) {
		return "", ErrNotFound
	}
	return s.studyName, nil
}

// GetStudyDirection returns study direction of the objective.
func (s *BlackHoleStorage) GetStudyDirection(studyID int) (StudyDirection, error) {
	if !s.checkStudyID(studyID) {
		return StudyDirectionMinimize, ErrInvalidStudyID
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

// CreateNewTrial creates trial and returns trialID.
func (s *BlackHoleStorage) CreateNewTrial(studyID int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.checkStudyID(studyID) {
		return -1, ErrInvalidStudyID
	}
	number := s.counter
	trialID := number
	s.counter++

	var err error
	idx := s.getTrialIndex(trialID)
	if s.isPartiallyDeleted() && !s.trials[idx].State.IsFinished() {
		err = ErrDeleteNonFinishedTrial
	}

	s.trials[idx] = FrozenTrial{
		ID:                 trialID,
		Number:             number,
		State:              TrialStateRunning,
		Value:              0,
		IntermediateValues: make(map[int]float64, 8),
		DatetimeStart:      time.Now(),
		DatetimeComplete:   time.Time{},
		InternalParams:     make(map[string]float64, 8),
		Params:             make(map[string]interface{}, 8),
		Distributions:      make(map[string]interface{}, 8),
		UserAttrs:          make(map[string]string, 8),
		SystemAttrs:        make(map[string]string, 8),
	}
	return trialID, err
}

// CloneTrial creates new Trial from the given base Trial.
func (s *BlackHoleStorage) CloneTrial(studyID int, baseTrial FrozenTrial) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.checkStudyID(studyID) {
		return -1, ErrInvalidStudyID
	}
	number := s.counter
	trialID := number
	s.counter++

	var err error
	idx := s.getTrialIndex(trialID)
	if !s.trials[idx].State.IsFinished() {
		err = ErrDeleteNonFinishedTrial
	}
	s.trials[idx] = FrozenTrial{
		ID:                 trialID,
		StudyID:            studyID,
		Number:             number,
		State:              baseTrial.State,
		Value:              baseTrial.Value,
		IntermediateValues: baseTrial.IntermediateValues,
		DatetimeStart:      baseTrial.DatetimeStart,
		DatetimeComplete:   baseTrial.DatetimeComplete,
		InternalParams:     baseTrial.InternalParams,
		Params:             baseTrial.Params,
		Distributions:      baseTrial.Distributions,
		UserAttrs:          baseTrial.UserAttrs,
		SystemAttrs:        baseTrial.SystemAttrs,
	}
	return trialID, err
}

// SetTrialValue sets the value of trial.
func (s *BlackHoleStorage) SetTrialValue(trialID int, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkTrialID(trialID); err != nil {
		return err
	}
	idx := s.getTrialIndex(trialID)
	trial := s.trials[idx]
	if trial.State.IsFinished() {
		return ErrTrialCannotBeUpdated
	}
	trial.Value = value
	s.trials[idx] = trial
	return nil
}

// SetTrialIntermediateValue sets the intermediate value of trial.
func (s *BlackHoleStorage) SetTrialIntermediateValue(trialID int, step int, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkTrialID(trialID); err != nil {
		return err
	}

	idx := s.getTrialIndex(trialID)
	trial := s.trials[idx]
	if trial.State.IsFinished() {
		return ErrTrialCannotBeUpdated
	}

	for key := range trial.IntermediateValues {
		if key == step {
			return errors.New("step value is already exist")
		}
	}

	trial.IntermediateValues[step] = value
	s.trials[idx] = trial
	return nil
}

// SetTrialParam sets the sampled parameters of trial.
func (s *BlackHoleStorage) SetTrialParam(
	trialID int,
	paramName string,
	paramValueInternal float64,
	distribution interface{},
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkTrialID(trialID); err != nil {
		return err
	}

	idx := s.getTrialIndex(trialID)
	trial := s.trials[idx]

	// Check trial is able to update
	if trial.State.IsFinished() {
		return ErrTrialCannotBeUpdated
	}

	paramValueExternal, err := ToExternalRepresentation(distribution, paramValueInternal)
	if err != nil {
		return err
	}

	trial.Distributions[paramName] = distribution
	trial.InternalParams[paramName] = paramValueInternal
	trial.Params[paramName] = paramValueExternal
	s.trials[idx] = trial
	return nil
}

// SetTrialState sets the state of trial.
func (s *BlackHoleStorage) SetTrialState(trialID int, state TrialState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkTrialID(trialID); err != nil {
		return err
	}
	idx := s.getTrialIndex(trialID)
	trial := s.trials[idx]
	if trial.State.IsFinished() {
		return ErrTrialCannotBeUpdated
	}
	trial.State = state
	if trial.State.IsFinished() {
		trial.DatetimeComplete = time.Now()
		s.updateBestTrial(trial)
	}
	s.trials[idx] = trial
	return nil
}

// SetTrialUserAttr to store the value for the user.
func (s *BlackHoleStorage) SetTrialUserAttr(trialID int, key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkTrialID(trialID); err != nil {
		return err
	}
	idx := s.getTrialIndex(trialID)
	if s.trials[idx].State.IsFinished() {
		return ErrTrialCannotBeUpdated
	}
	s.trials[idx].UserAttrs[key] = value
	return nil
}

// SetTrialSystemAttr to store the value for the system.
func (s *BlackHoleStorage) SetTrialSystemAttr(trialID int, key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.checkTrialID(trialID); err != nil {
		return err
	}
	idx := s.getTrialIndex(trialID)
	if s.trials[idx].State == TrialStateComplete {
		return ErrTrialCannotBeUpdated
	}
	s.trials[idx].SystemAttrs[key] = value
	return nil
}

// GetTrialNumberFromID returns the trial's number.
func (s *BlackHoleStorage) GetTrialNumberFromID(trialID int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	err := s.checkTrialID(trialID)
	if err == ErrTrialAlreadyDeleted {
		return trialID, err
	}
	return -1, err
}

// GetTrialParam returns the internal parameter of the trial
func (s *BlackHoleStorage) GetTrialParam(trialID int, paramName string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkTrialID(trialID); err != nil {
		return -1, err
	}
	idx := s.getTrialIndex(trialID)
	ir, ok := s.trials[idx].InternalParams[paramName]
	if !ok {
		return -1.0, errors.New("param doesn't exist")
	}
	return ir, nil
}

// GetTrial returns Trial.
func (s *BlackHoleStorage) GetTrial(trialID int) (FrozenTrial, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if err := s.checkTrialID(trialID); err != nil {
		return FrozenTrial{}, err
	}
	idx := s.getTrialIndex(trialID)
	return s.trials[idx], nil
}

// GetAllTrials returns the all trials.
func (s *BlackHoleStorage) GetAllTrials(studyID int) ([]FrozenTrial, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var err error
	n := s.nTrials
	if s.isPartiallyDeleted() {
		err = ErrTrialsPartiallyDeleted
	} else {
		n = s.counter
	}

	trials := make([]FrozenTrial, 0, n)
	for i := 0; i < len(s.trials); i++ {
		idx := s.getTrialIndex(s.counter + i)
		trials = append(trials, s.trials[idx])
	}
	return trials, err
}

// GetBestTrial returns the best trial.
func (s *BlackHoleStorage) GetBestTrial(studyID int) (FrozenTrial, error) {
	var err error
	if s.bestTrial.State != TrialStateComplete {
		err = ErrNoCompletedTrials
	}
	return s.bestTrial, err
}

// GetTrialParams returns the external parameters in the trial
func (s *BlackHoleStorage) GetTrialParams(trialID int) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkTrialID(trialID); err != nil {
		return nil, err
	}
	idx := s.getTrialIndex(trialID)
	return s.trials[idx].Params, nil
}

// GetTrialUserAttrs to restore the attributes for the user.
func (s *BlackHoleStorage) GetTrialUserAttrs(trialID int) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkTrialID(trialID); err != nil {
		return nil, err
	}

	idx := s.getTrialIndex(trialID)
	n := make(map[string]string, len(s.trials[idx].UserAttrs))
	for k := range s.trials[idx].UserAttrs {
		n[k] = s.trials[idx].UserAttrs[k]
	}
	return n, nil
}

// GetTrialSystemAttrs to restore the attributes for the system.
func (s *BlackHoleStorage) GetTrialSystemAttrs(trialID int) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkTrialID(trialID); err != nil {
		return nil, err
	}

	idx := s.getTrialIndex(trialID)
	n := make(map[string]string, len(s.trials[idx].SystemAttrs))
	for k := range s.trials[idx].SystemAttrs {
		n[k] = s.trials[idx].SystemAttrs[k]
	}
	return n, nil
}

func (s *BlackHoleStorage) checkStudyID(studyID int) bool {
	return studyID == InMemoryStorageStudyID
}

func (s *BlackHoleStorage) getTrialIndex(trialID int) int {
	return trialID % s.nTrials
}

func (s *BlackHoleStorage) checkTrialID(trialID int) error {
	// | nTrials | counter |  trials |
	// |       3 |       0 |      [] |
	// |       3 |       3 | [0,1,2] |
	// |       3 |       4 | [1,2,3] |
	if trialID < 0 || trialID >= s.counter {
		// counter represents an id for next trial.
		return ErrInvalidTrialID
	}
	if s.counter-s.nTrials < trialID {
		return nil
	}
	return ErrTrialAlreadyDeleted
}

func (s *BlackHoleStorage) isPartiallyDeleted() bool {
	// | nTrials | counter |  trials |
	// |       3 |       0 |      [] |
	// |       3 |       3 | [0,1,2] |
	// |       3 |       4 | [1,2,3] |
	return s.counter > s.nTrials
}

func (s *BlackHoleStorage) updateBestTrial(trial FrozenTrial) {
	if trial.State != TrialStateComplete {
		return
	}

	if s.bestTrial.State != TrialStateComplete {
		s.bestTrial = trial
		return
	}

	if s.direction == StudyDirectionMaximize && trial.Value > s.bestTrial.Value {
		s.bestTrial = trial
		return
	}
	if s.direction == StudyDirectionMinimize && trial.Value < s.bestTrial.Value {
		s.bestTrial = trial
		return
	}
}
