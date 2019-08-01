package goptuna

// Pruner is a interface for early stopping algorithms.
type Pruner interface {
	// Prune judge whether the trial should be pruned at the given step.
	Prune(storage Storage, studyID, trialID int, step int) (bool, error)
}
