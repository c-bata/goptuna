package goptuna

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

var errCreateNewTrial = errors.New("failed to create a new trial")

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
	ID                int
	Storage           Storage
	Sampler           Sampler
	Pruner            Pruner
	direction         StudyDirection
	logger            Logger
	ignoreErr         bool
	trialNotification chan FrozenTrial
	mu                sync.RWMutex
	ctx               context.Context
}

// GetTrials returns all trials in this study.
func (s *Study) GetTrials() ([]FrozenTrial, error) {
	return s.Storage.GetAllTrials(s.ID)
}

// Direction returns the direction of objective function value
func (s *Study) Direction() StudyDirection {
	return s.direction
}

// WithContext sets a context and it might cancel the execution of Optimize.
func (s *Study) WithContext(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ctx = ctx
}

func (s *Study) runTrial(objective FuncObjective) (int, error) {
	trialID, err := s.Storage.CreateNewTrialID(s.ID)
	if err != nil {
		s.logger.Error("failed to create a new trial",
			fmt.Sprintf("err=%s", err))
		return -1, errCreateNewTrial
	}

	trial := Trial{
		Study: s,
		ID:    trialID,
	}
	evaluation, objerr := objective(trial)
	var state TrialState
	if objerr == ErrTrialPruned {
		state = TrialStatePruned
		objerr = nil
	} else if objerr != nil {
		state = TrialStateFail
	} else {
		state = TrialStateComplete
	}

	if objerr != nil {
		s.logger.Error("Objective function returns error",
			fmt.Sprintf("trialID=%d", trialID),
			fmt.Sprintf("state=%s", state.String()),
			fmt.Sprintf("err=%s", objerr))
	} else {
		s.logger.Info("Trial finished",
			fmt.Sprintf("trialID=%d", trialID),
			fmt.Sprintf("state=%s", state.String()),
			fmt.Sprintf("evaluation=%f", evaluation))
	}

	if state == TrialStateComplete {
		// The trial.value of pruned trials are already set at trial.Report().
		err = s.Storage.SetTrialValue(trialID, evaluation)
		if err != nil {
			s.logger.Error("Failed to set trial value",
				fmt.Sprintf("trialID=%d", trialID),
				fmt.Sprintf("state=%s", state.String()),
				fmt.Sprintf("evaluation=%f", evaluation),
				fmt.Sprintf("err=%s", err))
			return trialID, err
		}
	}

	err = s.Storage.SetTrialState(trialID, state)
	if err != nil {
		s.logger.Error("Failed to set trial state",
			fmt.Sprintf("trialID=%d", trialID),
			fmt.Sprintf("state=%s", state.String()),
			fmt.Sprintf("evaluation=%f", evaluation),
			fmt.Sprintf("err=%s", err))
		return trialID, err
	}
	return trialID, objerr
}

// Optimize optimizes an objective function.
func (s *Study) Optimize(objective FuncObjective, evaluateMax int) error {
	evaluateCnt := 0
	for {
		if evaluateCnt >= evaluateMax {
			break
		}

		if s.ctx != nil {
			select {
			case <-s.ctx.Done():
				err := s.ctx.Err()
				s.logger.Debug("context is canceled", err)
				return err
			default:
				// do nothing
			}
		}
		// Evaluate an objective function
		trialID, err := s.runTrial(objective)
		if err == errCreateNewTrial {
			continue
		}
		evaluateCnt++

		// Send trial notification
		if s.trialNotification != nil {
			frozen, gerr := s.Storage.GetTrial(trialID)
			if gerr != nil {
				s.logger.Error("Failed to send trial notification",
					fmt.Sprintf("trialID=%d", trialID),
					fmt.Sprintf("err=%s", gerr))
				if !s.ignoreErr {
					return gerr
				}
			}
			s.trialNotification <- frozen
		}

		if !s.ignoreErr && err != nil {
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
		ID:        0,
		Storage:   storage,
		Sampler:   sampler,
		Pruner:    nil,
		direction: StudyDirectionMinimize,
		logger: &StdLogger{
			Logger: log.New(os.Stdout, "", log.LstdFlags),
			Level:  LoggerLevelDebug,
			Color:  true,
		},
		ignoreErr: false,
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
		ID:        0,
		Storage:   storage,
		Sampler:   sampler,
		Pruner:    nil,
		direction: "",
		logger: &StdLogger{
			Logger: log.New(os.Stdout, "", log.LstdFlags),
			Level:  LoggerLevelDebug,
			Color:  true,
		},
		ignoreErr: false,
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

// StudyOptionLogger sets Logger.
func StudyOptionLogger(logger Logger) StudyOption {
	return func(s *Study) error {
		if logger == nil {
			s.logger = &StdLogger{Logger: nil}
		} else {
			s.logger = logger
		}
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

// StudyOptionIgnoreError is an option to continue even if
// it receive error while running Optimize method.
func StudyOptionIgnoreError(ignore bool) StudyOption {
	return func(s *Study) error {
		s.ignoreErr = ignore
		return nil
	}
}

// StudyOptionSetTrialNotifyChannel to subscribe the finished trials.
func StudyOptionSetTrialNotifyChannel(notify chan FrozenTrial) StudyOption {
	return func(s *Study) error {
		s.trialNotification = notify
		return nil
	}
}

// StudyOptionSetLogger sets Logger.
// Deprecated: please use StudyOptionLogger instead.
var StudyOptionSetLogger = StudyOptionLogger
