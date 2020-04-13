package goptuna

import "errors"

// Pruner is a interface for early stopping algorithms.
type Pruner interface {
	// Prune judge whether the trial should be pruned at the given step.
	Prune(study *Study, trial FrozenTrial) (bool, error)
}

var (
	// ErrTrialPruned represents the pruned.
	ErrTrialPruned = errors.New("trial is pruned")
)
