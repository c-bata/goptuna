package dashboard

import (
	"fmt"
	"sort"
	"time"

	"github.com/c-bata/goptuna"
)

type TrialParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type IntermediateValue struct {
	Step  int     `json:"step"`
	Value float64 `json:"value"`
}

type FrozenTrial struct {
	ID                 int                 `json:"trial_id"`
	StudyID            int                 `json:"study_id"`
	Number             int                 `json:"number"`
	State              string              `json:"state"`
	Value              float64             `json:"value"`
	IntermediateValues []IntermediateValue `json:"intermediate_values"`
	DatetimeStart      string              `json:"datetime_start"`
	DatetimeComplete   string              `json:"datetime_complete"`
	Params             []TrialParam        `json:"params"`
	UserAttrs          []Attribute         `json:"user_attrs"`
	SystemAttrs        []Attribute         `json:"system_attrs"`
}

func toAttrs(from map[string]string) []Attribute {
	attrs := make([]Attribute, 0, len(from))
	for key := range from {
		attrs = append(attrs, Attribute{
			Key:   key,
			Value: from[key],
		})
	}
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})
	return attrs
}

func toIntermediateValues(from map[int]float64) []IntermediateValue {
	attrs := make([]IntermediateValue, 0, len(from))
	for step := range from {
		attrs = append(attrs, IntermediateValue{
			Step:  step,
			Value: from[step],
		})
	}
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Step < attrs[j].Step
	})
	return attrs
}

func toFrozenTrial(from goptuna.FrozenTrial) FrozenTrial {
	params := make([]TrialParam, 0, len(from.Params))
	for paramName := range from.Params {
		params = append(params, TrialParam{
			Name:  paramName,
			Value: fmt.Sprintf("%v", from.Params[paramName]),
		})
	}
	sort.Slice(params, func(i, j int) bool {
		return params[i].Name < params[j].Name
	})

	return FrozenTrial{
		ID:                 from.ID,
		StudyID:            from.StudyID,
		Number:             from.Number,
		State:              from.State.String(),
		Value:              from.Value,
		IntermediateValues: toIntermediateValues(from.IntermediateValues),
		DatetimeStart:      from.DatetimeStart.UTC().Format(time.RFC3339),
		DatetimeComplete:   from.DatetimeComplete.UTC().Format(time.RFC3339),
		Params:             params,
		UserAttrs:          toAttrs(from.UserAttrs),
		SystemAttrs:        toAttrs(from.SystemAttrs),
	}
}

func toFrozenTrials(from []goptuna.FrozenTrial) []FrozenTrial {
	res := make([]FrozenTrial, len(from))
	for i := 0; i < len(from); i++ {
		res[i] = toFrozenTrial(from[i])
	}
	return res
}

// StudySummary holds basic attributes and aggregated results of Study.
type StudySummary struct {
	ID            int         `json:"study_id"`
	Name          string      `json:"study_name"`
	Direction     string      `json:"direction"`
	BestTrial     FrozenTrial `json:"best_trial"`
	UserAttrs     []Attribute `json:"user_attrs"`
	SystemAttrs   []Attribute `json:"system_attrs"`
	DatetimeStart string      `json:"datetime_start"`
}

func toStudySummary(from goptuna.StudySummary) StudySummary {
	return StudySummary{
		ID:            from.ID,
		Name:          from.Name,
		Direction:     string(from.Direction),
		BestTrial:     toFrozenTrial(from.BestTrial),
		UserAttrs:     toAttrs(from.UserAttrs),
		SystemAttrs:   toAttrs(from.SystemAttrs),
		DatetimeStart: from.DatetimeStart.UTC().Format(time.RFC3339),
	}
}

func toStudySummaries(from []goptuna.StudySummary) []StudySummary {
	res := make([]StudySummary, len(from))
	for i := 0; i < len(from); i++ {
		res[i] = toStudySummary(from[i])
	}
	return res
}
