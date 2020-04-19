package successivehalving_test

import (
	"testing"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/successivehalving"
)

func TestOptunaPruner_IntermediateValues(t *testing.T) {
	var tests = []struct {
		name              string
		direction         goptuna.StudyDirection
		intermediateValue float64
	}{
		{
			name:              "minimize",
			direction:         goptuna.StudyDirectionMinimize,
			intermediateValue: 2,
		},
		{
			name:              "maximize",
			direction:         goptuna.StudyDirectionMaximize,
			intermediateValue: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pruner := &successivehalving.Pruner{
				MinResource:          1,
				ReductionFactor:      2,
				MinEarlyStoppingRate: 0,
			}
			study, err := goptuna.CreateStudy("optuna-pruner",
				goptuna.StudyOptionDirection(tt.direction),
				goptuna.StudyOptionPruner(pruner))
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
			}

			// A pruner is not activated at a first trial.
			trialID, err := study.Storage.CreateNewTrial(study.ID)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
			}
			trial := goptuna.Trial{
				Study: study,
				ID:    trialID,
			}
			err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 1, 1)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
			}
			ft, err := study.Storage.GetTrial(trialID)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
			}
			prune, err := pruner.Prune(study, ft)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
			}
			if prune {
				t.Errorf("should not be activated at a first trial, but got prune() = %v", prune)
			}

			// A pruner is not activated if a trial has no intermediate values.
			trialID, err = study.Storage.CreateNewTrial(study.ID)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
			}
			trial = goptuna.Trial{
				Study: study,
				ID:    trialID,
			}
			ft, err = study.Storage.GetTrial(trialID)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
			}
			prune, err = pruner.Prune(study, ft)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
			}
			if prune {
				t.Errorf("should not be activated if a trial has no intermediate values, but got prune() = %v", prune)
			}

			// A pruner is activated if a trial has an intermediate value.
			err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 1, tt.intermediateValue)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
			}
			ft, err = study.Storage.GetTrial(trialID)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
			}
			prune, err = pruner.Prune(study, ft)
			if err != nil {
				t.Errorf("should be err=nil, but got %s", err)
			}
			if !prune {
				t.Errorf("should be activated if a trial has no intermediate values, but got prune() = %v", prune)
			}
		})
	}
}

func TestOptunaPruner_RungCheck(t *testing.T) {
	pruner := &successivehalving.Pruner{
		MinResource:          1,
		ReductionFactor:      2,
		MinEarlyStoppingRate: 0,
	}
	study, err := goptuna.CreateStudy("optuna-pruner",
		goptuna.StudyOptionPruner(pruner))
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}

	// Report 7 trials in advance.
	for i := 0; i < 7; i++ {
		trialID, err := study.Storage.CreateNewTrial(study.ID)
		if err != nil {
			t.Errorf("should be err=nil, but got %s", err)
		}
		trial := goptuna.Trial{
			Study: study,
			ID:    trialID,
		}
		err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 7, 0.1 * float64(i+1))
		if err != nil {
			t.Errorf("should be err=nil, but got %s", err)
		}
		ft, err := study.Storage.GetTrial(trialID)
		if err != nil {
			t.Errorf("should be err=nil, but got %s", err)
		}
		_, err = pruner.Prune(study, ft)
		if err != nil {
			t.Errorf("should be err=nil, but got %s", err)
		}
	}

	var isexit = func(x map[string]string, key string) bool {
		for k := range x {
			if k == key {
				return true
			}
		}
		return false
	}

	// Report a trial that has the 7-th value from bottom.
	trialID, err := study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	trial := goptuna.Trial{
		Study: study,
		ID:    trialID,
	}
	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 7, 0.75)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err := study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	_, err = pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if !isexit(ft.SystemAttrs, "completed_rung_0") {
		t.Errorf("completed_rung_0 should be exist")
	}
	if isexit(ft.SystemAttrs, "completed_rung_1") {
		t.Errorf("completed_rung_1 should not be exist")
	}

	// Report a trial that has the third value from bottom.
	trialID, err = study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	trial = goptuna.Trial{
		Study: study,
		ID:    trialID,
	}
	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 7, 0.25)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err = study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	_, err = pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if !isexit(ft.SystemAttrs, "completed_rung_1") {
		t.Errorf("completed_rung_0 should be exist")
	}
	if isexit(ft.SystemAttrs, "completed_rung_2") {
		t.Errorf("completed_rung_1 should not be exist")
	}

	// Report a trial that has the lowest value.
	trialID, err = study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	trial = goptuna.Trial{
		Study: study,
		ID:    trialID,
	}
	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 7, 0.05)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err = study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	_, err = pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if !isexit(ft.SystemAttrs, "completed_rung_2") {
		t.Errorf("completed_rung_0 should be exist")
	}
	if isexit(ft.SystemAttrs, "completed_rung_3") {
		t.Errorf("completed_rung_1 should not be exist")
	}
}

func TestOptunaPruner_FirstTrialIsNotPruned(t *testing.T) {
	pruner := &successivehalving.Pruner{
		MinResource:          1,
		ReductionFactor:      2,
		MinEarlyStoppingRate: 0,
	}
	study, err := goptuna.CreateStudy("optuna-pruner",
		goptuna.StudyOptionPruner(pruner))
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}

	trialID, err := study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	trial := goptuna.Trial{
		Study: study,
		ID:    trialID,
	}
	for i := 0; i < 10; i++ {
		err = trial.Study.Storage.SetTrialIntermediateValue(trialID, i, 1)
		if err != nil {
			t.Errorf("should be err=nil, but got %s", err)
		}
		ft, err := study.Storage.GetTrial(trialID)
		if err != nil {
			t.Errorf("should be err=nil, but got %s", err)
		}
		prune, err := pruner.Prune(study, ft)
		if err != nil {
			t.Errorf("should be err=nil, but got %s", err)
		}
		if prune {
			t.Errorf("should not be activated, but got prune() = %v", prune)
		}
	}

	var isexit = func(x map[string]string, key string) bool {
		for k := range x {
			if k == key {
				return true
			}
		}
		return false
	}

	ft, err := study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if !isexit(ft.SystemAttrs, "completed_rung_0") {
		t.Errorf("completed_rung_0 should be exist")
	}
	if !isexit(ft.SystemAttrs, "completed_rung_1") {
		t.Errorf("completed_rung_1 should be exist")
	}
	if !isexit(ft.SystemAttrs, "completed_rung_2") {
		t.Errorf("completed_rung_2 should be exist")
	}
	if !isexit(ft.SystemAttrs, "completed_rung_3") {
		t.Errorf("completed_rung_3 should be exist")
	}
	if isexit(ft.SystemAttrs, "completed_rung_4") {
		t.Errorf("completed_rung_4 should not be exist")
	}
}

func TestOptunaPruner_MinResource(t *testing.T) {
	var isexit = func(x map[string]string, key string) bool {
		for k := range x {
			if k == key {
				return true
			}
		}
		return false
	}
	study, err := goptuna.CreateStudy("optuna-pruner")
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}

	// min_resource=1: The rung 0 ends at step 1.
	pruner := &successivehalving.Pruner{
		MinResource:          1,
		ReductionFactor:      2,
		MinEarlyStoppingRate: 0,
	}
	trialID, err := study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	trial := goptuna.Trial{
		Study: study,
		ID:    trialID,
	}
	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 1, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err := study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	prune, err := pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if prune {
		t.Errorf("should not be activated if a trial has no intermediate values, but got prune() = %v", prune)
	}

	if !isexit(ft.SystemAttrs, "completed_rung_0") {
		t.Errorf("completed_rung_0 should be exist")
	}
	if isexit(ft.SystemAttrs, "completed_rung_1") {
		t.Errorf("completed_rung_1 should not be exist")
	}

	// min_resource=2: The rung 0 ends at step 2.
	pruner = &successivehalving.Pruner{
		MinResource:          2,
		ReductionFactor:      2,
		MinEarlyStoppingRate: 0,
	}
	trialID, err = study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	trial = goptuna.Trial{
		Study: study,
		ID:    trialID,
	}
	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 1, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err = study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	prune, err = pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if prune {
		t.Errorf("should not be activated if a trial has no intermediate values, but got prune() = %v", prune)
	}
	if isexit(ft.SystemAttrs, "completed_rung_0") {
		t.Errorf("completed_rung_0 should not be exist")
	}

	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 2, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err = study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	prune, err = pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if prune {
		t.Errorf("should not be activated if a trial has no intermediate values, but got prune() = %v", prune)
	}
	if !isexit(ft.SystemAttrs, "completed_rung_0") {
		t.Errorf("completed_rung_0 should not exist")
	}
	if isexit(ft.SystemAttrs, "completed_rung_1") {
		t.Errorf("completed_rung_1 should not be exist")
	}
}

func TestOptunaPruner_ReductionFactor(t *testing.T) {
	var isexit = func(x map[string]string, key string) bool {
		for k := range x {
			if k == key {
				return true
			}
		}
		return false
	}
	study, err := goptuna.CreateStudy("optuna-pruner")
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}

	// reduction_factor=2: The rung 0 ends at step 1.
	pruner := &successivehalving.Pruner{
		MinResource:          1,
		ReductionFactor:      2,
		MinEarlyStoppingRate: 0,
	}
	trialID, err := study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	trial := goptuna.Trial{
		Study: study,
		ID:    trialID,
	}
	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 1, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err := study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	prune, err := pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if prune {
		t.Errorf("should not be activated if a trial has no intermediate values, but got prune() = %v", prune)
	}
	if !isexit(ft.SystemAttrs, "completed_rung_0") {
		t.Errorf("completed_rung_0 should be exist")
	}
	if isexit(ft.SystemAttrs, "completed_rung_1") {
		t.Errorf("completed_rung_1 should not be exist")
	}

	// reduction_factor=3: The rung 1 ends at step 3.
	pruner = &successivehalving.Pruner{
		MinResource:          1,
		ReductionFactor:      3,
		MinEarlyStoppingRate: 0,
	}
	trialID, err = study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	trial = goptuna.Trial{
		Study: study,
		ID:    trialID,
	}
	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 1, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err = study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	prune, err = pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if prune {
		t.Errorf("should not be activated if a trial has no intermediate values, but got prune() = %v", prune)
	}
	if !isexit(ft.SystemAttrs, "completed_rung_0") {
		t.Errorf("completed_rung_0 should be exist")
	}
	if isexit(ft.SystemAttrs, "completed_rung_1") {
		t.Errorf("completed_rung_1 should not be exist")
	}

	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 2, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err = study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	prune, err = pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if prune {
		t.Errorf("should not be activated if a trial has no intermediate values, but got prune() = %v", prune)
	}
	if isexit(ft.SystemAttrs, "completed_rung_1") {
		t.Errorf("completed_rung_1 should not be exist")
	}

	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 3, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err = study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	prune, err = pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if prune {
		t.Errorf("should not be activated if a trial has no intermediate values, but got prune() = %v", prune)
	}
	if !isexit(ft.SystemAttrs, "completed_rung_1") {
		t.Errorf("completed_rung_1 should be exist")
	}
	if isexit(ft.SystemAttrs, "completed_rung_2") {
		t.Errorf("completed_rung_2 should not be exist")
	}
}

func TestOptunaPruner_MinEarlyStoppingRate(t *testing.T) {
	var isexit = func(x map[string]string, key string) bool {
		for k := range x {
			if k == key {
				return true
			}
		}
		return false
	}
	study, err := goptuna.CreateStudy("optuna-pruner")
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}

	// min_early_stopping_rate=0: The rung 0 ends at step 1.
	pruner := &successivehalving.Pruner{
		MinResource:          1,
		ReductionFactor:      2,
		MinEarlyStoppingRate: 0,
	}
	trialID, err := study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	trial := goptuna.Trial{
		Study: study,
		ID:    trialID,
	}
	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 1, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err := study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	prune, err := pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if prune {
		t.Errorf("should not be activated if a trial has no intermediate values, but got prune() = %v", prune)
	}
	if !isexit(ft.SystemAttrs, "completed_rung_0") {
		t.Errorf("completed_rung_0 should be exist")
	}

	// min_early_stopping_rate=1: The rung 0 ends at step 2.
	pruner = &successivehalving.Pruner{
		MinResource:          1,
		ReductionFactor:      2,
		MinEarlyStoppingRate: 1,
	}
	trialID, err = study.Storage.CreateNewTrial(study.ID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	trial = goptuna.Trial{
		Study: study,
		ID:    trialID,
	}
	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 1, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err = study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	prune, err = pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if prune {
		t.Errorf("should not be activated if a trial has no intermediate values, but got prune() = %v", prune)
	}
	if isexit(ft.SystemAttrs, "completed_rung_0") {
		t.Errorf("completed_rung_0 should not be exist")
	}
	if isexit(ft.SystemAttrs, "completed_rung_1") {
		t.Errorf("completed_rung_1 should not be exist")
	}

	err = trial.Study.Storage.SetTrialIntermediateValue(trialID, 2, 1)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	ft, err = study.Storage.GetTrial(trialID)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	prune, err = pruner.Prune(study, ft)
	if err != nil {
		t.Errorf("should be err=nil, but got %s", err)
	}
	if prune {
		t.Errorf("should not be activated if a trial has no intermediate values, but got prune() = %v", prune)
	}
	if !isexit(ft.SystemAttrs, "completed_rung_0") {
		t.Errorf("completed_rung_0 should be exist")
	}
	if isexit(ft.SystemAttrs, "completed_rung_1") {
		t.Errorf("completed_rung_1 should not be exist")
	}
}
