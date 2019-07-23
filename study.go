package goptuna

import (
	"fmt"

	"go.uber.org/zap"
)

type FuncObjective func(trial Trial) (float64, error)

type StudyDirection string

const (
	StudyDirectionMaximize StudyDirection = "maximize"
	StudyDirectionMinimize StudyDirection = "minimize"
)

type Study struct {
	id        string
	storage   Storage
	sampler   Sampler
	direction StudyDirection
	logger    *zap.Logger
}

func (s *Study) Report(trialID string, value float64) error {
	return s.storage.SetTrialValue(trialID, value)
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
	evaluation, err := objective(trial)
	if err != nil {
		saveErr := s.storage.SetTrialState(trialID, TrialStateFail)
		return trialID, saveErr
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

func (s *Study) Optimize(objective FuncObjective, evaluateMax int) error {
	evaluateCnt := 0
	for {
		if evaluateCnt > evaluateMax {
			break
		}
		evaluateCnt++
		trialID, err := s.runTrial(objective)
		if err != nil {
			return err
		}
		if s.logger != nil {
			if trial, err := s.storage.GetTrial(trialID); err == nil {
				s.logger.Debug("Finished trial",
					zap.String("trialID", trialID),
					zap.Int("state", int(trial.State)),
					zap.Float64("value", trial.Value),
					zap.String("paramsInIR", fmt.Sprintf("%#v", trial.ParamsInIR)))
			}
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
	storage Storage,
	sampler Sampler,
	opts ...StudyOption,
) (*Study, error) {
	studyID, err := storage.CreateNewStudyID(name)
	if err != nil {
		return nil, err
	}
	study := &Study{
		id:      studyID,
		storage: storage,
		sampler: sampler,
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
