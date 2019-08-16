package goptuna

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

// FuncObjective is a type of objective function
type FuncObjective func(trial Trial) (float64, error)

// StudyDirection represents the direction of the optimization
type StudyDirection string

const (
	// StudyDirectionMaximize maximizes objective function value
	StudyDirectionMaximize StudyDirection = "maximize"
	// StudyDirectionMinimize minimizes objective function value
	StudyDirectionMinimize StudyDirection = "minimize"
)

// Study corresponds to an optimization task, i.e., a set of trials.
type Study struct {
	ID                 int
	Storage            Storage
	Sampler            Sampler
	Pruner             Pruner
	direction          StudyDirection
	logger             *zap.Logger
	ignoreObjectiveErr bool
	trialNotifyChan    chan FrozenTrial
	mu                 sync.RWMutex
	ctx                context.Context
}

// GetTrials returns all trials in this study.
func (s *Study) GetTrials() ([]FrozenTrial, error) {
	return s.Storage.GetAllTrials(s.ID)
}

// Direction returns the direction of objective function value
func (s *Study) Direction() StudyDirection {
	return s.direction
}

// Report reports an objective function value
func (s *Study) Report(trialID int, value float64) error {
	return s.Storage.SetTrialValue(trialID, value)
}

// WithContext sets a context and it might cancel the execution of Optimize.
func (s *Study) WithContext(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ctx = ctx
}

func (s *Study) runTrial(objective FuncObjective) error {
	trialID, err := s.Storage.CreateNewTrialID(s.ID)
	if err != nil {
		return err
	}

	trial := Trial{
		Study: s,
		ID:    trialID,
	}

	var state TrialState
	evaluation, objerr := objective(trial)
	if objerr == ErrTrialPruned {
		state = TrialStatePruned
	} else if objerr != nil {
		state = TrialStateFail
	} else {
		state = TrialStateComplete
	}

	if s.logger != nil {
		s.logger.Info("Trial finished",
			zap.Int("trialID", trialID),
			zap.String("state", state.String()),
			zap.Float64("evaluationValue", evaluation),
			zap.Error(objerr))
	}

	if state == TrialStateComplete {
		err = s.Storage.SetTrialValue(trialID, evaluation)
		if err != nil && s.logger != nil {
			s.logger.Error("Failed to set trial value",
				zap.Int("trialID", trialID),
				zap.Float64("evaluationValue", evaluation),
				zap.Error(err))
		}
	}

	err = s.Storage.SetTrialState(trialID, state)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("Failed to set trial state",
				zap.Int("trialID", trialID),
				zap.String("state", state.String()),
				zap.Error(err))
		}
		return err
	}

	err = s.notifyFinishedTrial(trialID)
	if err != nil {
		return err
	}
	if !s.ignoreObjectiveErr && objerr != nil {
		return objerr
	}
	return nil
}

func (s *Study) notifyFinishedTrial(trialID int) error {
	if s.trialNotifyChan == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	trial, err := s.Storage.GetTrial(trialID)
	if err != nil {
		return err
	}

	if s.trialNotifyChan != nil {
		s.trialNotifyChan <- trial
	}
	return nil
}

// Optimize optimizes an objective function.
func (s *Study) Optimize(objective FuncObjective, evaluateMax int) error {
	evaluateCnt := 0
	for {
		if evaluateCnt >= evaluateMax {
			break
		}
		evaluateCnt++

		if s.ctx != nil {
			select {
			case <-s.ctx.Done():
				err := s.ctx.Err()
				if err != nil && s.logger != nil {
					s.logger.Debug("context is canceled", zap.Error(err))
				}
				return err
			default:
				// do nothing
			}
		}
		if err := s.runTrial(objective); err != nil {
			return err
		}
	}
	return nil
}

// GetBestValue return the best objective value
func (s *Study) GetBestValue() (float64, error) {
	trial, err := s.Storage.GetBestTrial(s.ID)
	if err != nil {
		return 0.0, err
	}
	return trial.Value, nil
}

// GetBestParams return parameters of the best trial
func (s *Study) GetBestParams() (map[string]interface{}, error) {
	trial, err := s.Storage.GetBestTrial(s.ID)
	if err != nil {
		return nil, err
	}
	return trial.Params, nil
}

// SetUserAttr to store the value for the user.
func (s *Study) SetUserAttr(key, value string) error {
	return s.Storage.SetStudyUserAttr(s.ID, key, value)
}

// SetSystemAttr to store the value for the system.
func (s *Study) SetSystemAttr(key, value string) error {
	return s.Storage.SetStudySystemAttr(s.ID, key, value)
}

// GetUserAttrs to store the value for the user.
func (s *Study) GetUserAttrs() (map[string]string, error) {
	return s.Storage.GetStudyUserAttrs(s.ID)
}

// GetSystemAttrs to store the value for the system.
func (s *Study) GetSystemAttrs() (map[string]string, error) {
	return s.Storage.GetStudySystemAttrs(s.ID)
}

// CreateStudy creates a new Study object.
func CreateStudy(
	name string,
	opts ...StudyOption,
) (*Study, error) {
	storage := NewInMemoryStorage()
	sampler := NewRandomSearchSampler()
	study := &Study{
		ID:                 0,
		Storage:            storage,
		Sampler:            sampler,
		Pruner:             nil,
		direction:          StudyDirectionMinimize,
		logger:             nil,
		ignoreObjectiveErr: false,
	}

	for _, opt := range opts {
		if err := opt(study); err != nil {
			return nil, err
		}
	}

	studyID, err := study.Storage.CreateNewStudyID(name)
	if err != nil {
		return nil, err
	}
	err = study.Storage.SetStudyDirection(studyID, study.direction)
	if err != nil {
		return nil, err
	}
	study.ID = studyID
	return study, nil
}

// LoadStudy loads an existing study.
func LoadStudy(
	name string,
	opts ...StudyOption,
) (*Study, error) {
	storage := NewInMemoryStorage()
	sampler := NewRandomSearchSampler()
	study := &Study{
		ID:                 0,
		Storage:            storage,
		Sampler:            sampler,
		Pruner:             nil,
		direction:          "",
		logger:             nil,
		ignoreObjectiveErr: false,
	}

	for _, opt := range opts {
		if err := opt(study); err != nil {
			return nil, err
		}
	}

	studyID, err := study.Storage.GetStudyIDFromName(name)
	if err != nil {
		return nil, err
	}
	study.ID = studyID
	direction, err := study.Storage.GetStudyDirection(studyID)
	if err != nil {
		return nil, err
	}
	study.direction = direction
	return study, nil
}

// StudyOption to pass the custom option
type StudyOption func(study *Study) error

// StudyOptionSetDirection change the direction of optimize
func StudyOptionSetDirection(direction StudyDirection) StudyOption {
	return func(s *Study) error {
		s.direction = direction
		return nil
	}
}

// StudyOptionSetLogger sets zap.Logger.
func StudyOptionSetLogger(logger *zap.Logger) StudyOption {
	return func(s *Study) error {
		s.logger = logger
		return nil
	}
}

// StudyOptionStorage sets the storage object.
func StudyOptionStorage(storage Storage) StudyOption {
	return func(s *Study) error {
		s.Storage = storage
		return nil
	}
}

// StudyOptionSampler sets the sampler object.
func StudyOptionSampler(sampler Sampler) StudyOption {
	return func(s *Study) error {
		s.Sampler = sampler
		return nil
	}
}

// StudyOptionPruner sets the pruner object.
func StudyOptionPruner(pruner Pruner) StudyOption {
	return func(s *Study) error {
		s.Pruner = pruner
		return nil
	}
}

// StudyOptionIgnoreObjectiveErr sets the option to ignore error returned from objective function
// If true, Optimize method continues to run new trial.
func StudyOptionIgnoreObjectiveErr(ignore bool) StudyOption {
	return func(s *Study) error {
		s.ignoreObjectiveErr = ignore
		return nil
	}
}

// StudyOptionSetTrialNotifyChannel to subscribe the finished trials.
func StudyOptionSetTrialNotifyChannel(notify chan FrozenTrial) StudyOption {
	return func(s *Study) error {
		s.trialNotifyChan = notify
		return nil
	}
}
