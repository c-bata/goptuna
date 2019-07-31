package medianstopping

func NewMedianPruner() *MedianPruner {
	percentile := &PercentilePruner{
		Percentile:     50,
		NStartUpTrials: 5,
		NWarmUpSteps:   0,
	}
	return &MedianPruner{percentile}
}

type MedianPruner struct {
	*PercentilePruner
}
