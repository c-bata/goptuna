package rdb

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/c-bata/goptuna"
)

func ToFrozenTrial(trial trialModel) (goptuna.FrozenTrial, error) {
	userAttrs := make(map[string]string, len(trial.UserAttributes))
	for i := range trial.UserAttributes {
		userAttrs[trial.UserAttributes[i].Key] = trial.UserAttributes[i].ValueJSON
	}

	systemAttrs := make(map[string]string, len(trial.SystemAttributes))
	for i := range trial.SystemAttributes {
		systemAttrs[trial.SystemAttributes[i].Key] = trial.SystemAttributes[i].ValueJSON
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

	// todo: convert all attributes
	return goptuna.FrozenTrial{
		ID:                 trial.ID,
		StudyID:            trial.TrialReferStudy,
		Number:             number,
		State:              state,
		Value:              trial.Value,
		IntermediateValues: nil,
		DatetimeStart:      datetimeStart,
		DatetimeComplete:   datetimeComplete,
		Params:             nil,
		Distributions:      nil,
		UserAttrs:          userAttrs,
		SystemAttrs:        systemAttrs,
		ParamsInIR:         nil,
	}, nil
}

func toStateExternalRepresentation(state int) (goptuna.TrialState, error) {
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

func toGoptunaStudyDirection(direction int) goptuna.StudyDirection {
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
