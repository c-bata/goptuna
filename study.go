package goptuna

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type FuncObjective func(trial Trial) (float64, error)

type StudyDirection string

const (
	StudyDirectionMaximize StudyDirection = "maximize"
	StudyDirectionMinimize StudyDirection = "minimize"
)

type Study struct {
	id                 string
	storage            Storage
	sampler            Sampler
	direction          StudyDirection
	logger             *zap.Logger
	ignoreObjectiveErr bool
	trialNotifyChan    chan FrozenTrial
	mu                 sync.RWMutex
	ctx                context.Context
}

func (s *Study) GetTrials() ([]FrozenTrial, error) {
	return s.storage.GetAllTrials(s.id)
}

func (s *Study) Direction() StudyDirection {
	return s.direction
}

func (s *Study) Report(trialID string, value float64) error {
	return s.storage.SetTrialValue(trialID, value)
}

func (s *Study) WithContext(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ctx = ctx
}

func (s *Study) runTrial(objective FuncObjective) (string, error) {
	trialID, err := s.storage.CreateNewTrialID(s.id)
	if err != nil {
		return "", err
	}

	trial := Trial{
		study: s,
		id:    trialID,
	}
	evaluation, objerr := objective(trial)
	if objerr != nil {
		saveerr := s.storage.SetTrialState(trialID, TrialStateFail)
		if saveerr != nil {
			return trialID, saveerr
		}

		if !s.ignoreObjectiveErr {
			return trialID, fmt.Errorf("objective: %s", objerr)
		}
	}

	if err = s.storage.SetTrialValue(trialID, evaluation); err != nil {
		return trialID, err
	}
	if err = s.Report(trialID, evaluation); err != nil {
		return trialID, err
	}
	if err = s.storage.SetTrialState(trialID, TrialStateComplete); err != nil {
		return trialID, err
	}
	return trialID, nil
}

func (s *Study) notifyFinishedTrial(trialID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.trialNotifyChan == nil && s.logger == nil {
		return nil
	}

	trial, err := s.storage.GetTrial(trialID)
	if err != nil {
		return err
	}

	if s.trialNotifyChan != nil {
		s.trialNotifyChan <- trial
	}
	if s.logger != nil {
		s.logger.Info("Finished trial",
			zap.String("trialID", trialID),
			zap.String("state", trial.State.String()),
			zap.Float64("value", trial.Value),
			zap.String("paramsInIR", fmt.Sprintf("%v", trial.ParamsInIR)))
	}
	return nil
}

func (s *Study) Optimize(objective FuncObjective, evaluateMax int) error {
	evaluateCnt := 0
	for {
		if evaluateCnt > evaluateMax {
			break
		}
		evaluateCnt++

		if s.ctx != nil {
			select {
			case <-s.ctx.Done():
				return s.ctx.Err()
			default:
				// do nothing
			}
		}

		trialID, err := s.runTrial(objective)
		if err != nil {
			return err
		}

		err = s.notifyFinishedTrial(trialID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Study) GetBestValue() (float64, error) {
	trial, err := s.storage.GetBestTrial(s.id)
	if err != nil {
		return 0.0, err
	}
	return trial.Value, nil
}

func (s *Study) GetBestParams() (map[string]float64, error) {
	// TODO: avoid using internal representation value
	trial, err := s.storage.GetBestTrial(s.id)
	if err != nil {
		return nil, err
	}
	return trial.ParamsInIR, nil
}

func CreateStudy(
	name string,
	opts ...StudyOption,
) (*Study, error) {
	storage := NewInMemoryStorage()
	studyID, err := storage.CreateNewStudyID(name)
	if err != nil {
		return nil, err
	}
	sampler := NewRandomSearchSampler()
	study := &Study{
		id:                 studyID,
		storage:            storage,
		sampler:            sampler,
		direction:          StudyDirectionMinimize,
		logger:             nil,
		ignoreObjectiveErr: false,
	}

	for _, opt := range opts {
		if err := opt(study); err != nil {
			return nil, err
		}
	}
	return study, nil
}

type StudyOption func(study *Study) error

func StudyOptionSetDirection(direction StudyDirection) StudyOption {
	return func(s *Study) error {
		s.direction = direction
		return nil
	}
}

func StudyOptionSetLogger(logger *zap.Logger) StudyOption {
	return func(s *Study) error {
		s.logger = logger
		return nil
	}
}

func StudyOptionStorage(storage Storage) StudyOption {
	return func(s *Study) error {
		s.storage = storage
		return nil
	}
}

func StudyOptionSampler(sampler Sampler) StudyOption {
	return func(s *Study) error {
		s.sampler = sampler
		return nil
	}
}

func StudyOptionIgnoreObjectiveErr(ignore bool) StudyOption {
	return func(s *Study) error {
		s.ignoreObjectiveErr = ignore
		return nil
	}
}

func StudyOptionSetTrialNotifyChannel(notify chan FrozenTrial) StudyOption {
	return func(s *Study) error {
		s.trialNotifyChan = notify
		return nil
	}
}
