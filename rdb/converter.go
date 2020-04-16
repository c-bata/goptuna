package rdb

import (
	"errors"
	"time"

	"github.com/c-bata/goptuna"
)

func toFrozenTrial(trial trialModel) (goptuna.FrozenTrial, error) {
	var err error
	userAttrs := make(map[string]string, len(trial.UserAttributes))
	for i := range trial.UserAttributes {
		userAttrs[trial.UserAttributes[i].Key] = decodeAttrValue(trial.UserAttributes[i].ValueJSON)
	}

	systemAttrs := make(map[string]string, len(trial.SystemAttributes))
	for i := range trial.SystemAttributes {
		systemAttrs[trial.SystemAttributes[i].Key] = decodeAttrValue(trial.SystemAttributes[i].ValueJSON)
	}

	distributions := make(map[string]interface{}, len(trial.TrialParams))
	paramsInIR := make(map[string]float64, len(trial.TrialParams))
	paramsInXR := make(map[string]interface{}, len(trial.TrialParams))
	for i := range trial.TrialParams {
		// distributions
		d, err := goptuna.JSONToDistribution([]byte(trial.TrialParams[i].DistributionJSON))
		if err != nil {
			return goptuna.FrozenTrial{}, err
		}
		distributions[trial.TrialParams[i].Name] = d

		// internal representation
		paramsInIR[trial.TrialParams[i].Name] = trial.TrialParams[i].Value

		// external representation
		paramsInXR[trial.TrialParams[i].Name], err = goptuna.ToExternalRepresentation(d, trial.TrialParams[i].Value)
		if err != nil {
			return goptuna.FrozenTrial{}, err
		}
	}

	state, err := toStateExternalRepresentation(trial.State)
	if err != nil {
		return goptuna.FrozenTrial{}, err
	}

	var datetimeStart, datetimeComplete time.Time
	if trial.DatetimeStart != nil {
		datetimeStart = *trial.DatetimeStart
	}
	if trial.DatetimeComplete != nil {
		datetimeComplete = *trial.DatetimeComplete
	}

	var value float64
	intermediateValue := make(map[int]float64, len(trial.TrialValues))
	for i := range trial.TrialValues {
		if trial.TrialValues[i].Step == finalValueStep {
			value = trial.TrialValues[i].Value
		} else {
			intermediateValue[trial.TrialValues[i].Step] = trial.TrialValues[i].Value
		}
	}

	return goptuna.FrozenTrial{
		ID:                 trial.ID,
		StudyID:            trial.TrialReferStudy,
		Number:             trial.Number,
		State:              state,
		Value:              value,
		IntermediateValues: intermediateValue,
		DatetimeStart:      datetimeStart,
		DatetimeComplete:   datetimeComplete,
		InternalParams:     paramsInIR,
		Params:             paramsInXR,
		Distributions:      distributions,
		UserAttrs:          userAttrs,
		SystemAttrs:        systemAttrs,
	}, nil
}

func toStudySummary(study studyModel, bestTrial goptuna.FrozenTrial, start time.Time) (goptuna.StudySummary, error) {
	userAttrs := make(map[string]string, len(study.UserAttributes))
	for i := range study.UserAttributes {
		userAttrs[study.UserAttributes[i].Key] = decodeAttrValue(study.UserAttributes[i].ValueJSON)
	}

	systemAttrs := make(map[string]string, len(study.SystemAttributes))
	for i := range study.SystemAttributes {
		systemAttrs[study.SystemAttributes[i].Key] = decodeAttrValue(study.SystemAttributes[i].ValueJSON)
	}
	return goptuna.StudySummary{
		ID:            study.ID,
		Name:          study.Name,
		Direction:     toGoptunaStudyDirection(study.Direction),
		BestTrial:     bestTrial,
		UserAttrs:     userAttrs,
		SystemAttrs:   systemAttrs,
		DatetimeStart: start,
	}, nil
}

func toStateExternalRepresentation(state string) (goptuna.TrialState, error) {
	switch state {
	case trialStateRunning:
		return goptuna.TrialStateRunning, nil
	case trialStateComplete:
		return goptuna.TrialStateComplete, nil
	case trialStatePruned:
		return goptuna.TrialStatePruned, nil
	case trialStateFail:
		return goptuna.TrialStateFail, nil
	case trialStateWaiting:
		return goptuna.TrialStateWaiting, nil
	default:
		return goptuna.TrialStateRunning, errors.New("invalid trial state")
	}
}

func toStateInternalRepresentation(state goptuna.TrialState) (string, error) {
	switch state {
	case goptuna.TrialStateRunning:
		return trialStateRunning, nil
	case goptuna.TrialStateComplete:
		return trialStateComplete, nil
	case goptuna.TrialStatePruned:
		return trialStatePruned, nil
	case goptuna.TrialStateFail:
		return trialStateFail, nil
	case goptuna.TrialStateWaiting:
		return trialStateWaiting, nil
	default:
		return "", errors.New("invalid trial state")
	}
}

func toGoptunaStudyDirection(direction string) goptuna.StudyDirection {
	switch direction {
	case directionMaximize:
		return goptuna.StudyDirectionMaximize
	case directionNotSet:
		fallthrough
	case directionMinimize:
		return goptuna.StudyDirectionMinimize
	default:
		return goptuna.StudyDirectionMinimize
	}
}
