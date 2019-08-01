package medianstopping

func NewMedianPruner() *MedianPruner {
	percentile := &PercentilePruner{
		Percentile:     50,
		NStartUpTrials: 5,
		NWarmUpSteps:   0,
	}
	return &MedianPruner{percentile}
}

// MedianPruner implements a median stopping rule of Google Vizier.
// Prune if the trial's best intermediate result is worse than median of
// intermediate results of previous trials at the same step.
// See https://ai.google/research/pubs/pub46180 for more details.
type MedianPruner struct {
	*PercentilePruner
}
