package rdb

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/c-bata/goptuna"
)

func toFrozenTrial(trial trialModel) (goptuna.FrozenTrial, error) {
	var err error
	userAttrs := make(map[string]string, len(trial.UserAttributes))
	for i := range trial.UserAttributes {
		userAttrs[trial.UserAttributes[i].Key], err = decodeAttrFromJSON(trial.UserAttributes[i].ValueJSON)
		if err != nil {
			return goptuna.FrozenTrial{}, err
		}
	}

	systemAttrs := make(map[string]string, len(trial.SystemAttributes))
	for i := range trial.SystemAttributes {
		systemAttrs[trial.SystemAttributes[i].Key], err = decodeAttrFromJSON(trial.SystemAttributes[i].ValueJSON)
		if err != nil {
			return goptuna.FrozenTrial{}, err
		}
	}

	paramsInIR := make(map[string]float64, len(trial.TrialParams))
	distributions := make(map[string]interface{}, len(trial.TrialParams))
	paramsInXR := make(map[string]interface{}, len(trial.TrialParams))
	for i := range trial.TrialParams {
		// paramsInIR
		paramsInIR[trial.TrialParams[i].Name] = trial.TrialParams[i].Value
		// distributions
		d, err := goptuna.JSONToDistribution([]byte(trial.TrialParams[i].DistributionJSON))
		if err != nil {
			return goptuna.FrozenTrial{}, err
		}
		distributions[trial.TrialParams[i].Name] = d
		// external representations
		paramsInXR[trial.TrialParams[i].Name], err = goptuna.ToExternalRepresentation(d, trial.TrialParams[i].Value)
		if err != nil {
			return goptuna.FrozenTrial{}, err
		}
	}

	numberStr, ok := systemAttrs["_number"]
	if !ok {
		return goptuna.FrozenTrial{}, errors.New("number is not exist in system attrs")
	}
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return goptuna.FrozenTrial{}, fmt.Errorf("invalid trial number: %s", err)
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

	intermediateValue := make(map[int]float64, len(trial.TrialValues))
	for i := range trial.TrialValues {
		intermediateValue[trial.TrialValues[i].Step] = trial.TrialValues[i].Value
	}

	return goptuna.FrozenTrial{
		ID:                 trial.ID,
		StudyID:            trial.TrialReferStudy,
		Number:             number,
		State:              state,
		Value:              trial.Value,
		IntermediateValues: intermediateValue,
		DatetimeStart:      datetimeStart,
		DatetimeComplete:   datetimeComplete,
		Params:             paramsInXR,
		Distributions:      distributions,
		UserAttrs:          userAttrs,
		SystemAttrs:        systemAttrs,
		ParamsInIR:         paramsInIR,
	}, nil
}

func toStudySummary(study studyModel, bestTrial goptuna.FrozenTrial, start time.Time) (goptuna.StudySummary, error) {
	var err error
	userAttrs := make(map[string]string, len(study.UserAttributes))
	for i := range study.UserAttributes {
		userAttrs[study.UserAttributes[i].Key], err = decodeAttrFromJSON(study.UserAttributes[i].ValueJSON)
		if err != nil {
			return goptuna.StudySummary{}, err
		}
	}

	systemAttrs := make(map[string]string, len(study.SystemAttributes))
	for i := range study.SystemAttributes {
		systemAttrs[study.SystemAttributes[i].Key], err = decodeAttrFromJSON(study.SystemAttributes[i].ValueJSON)
		if err != nil {
			return goptuna.StudySummary{}, err
		}
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
