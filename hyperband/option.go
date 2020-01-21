package hyperband

// Option to pass the custom option
type Option func(pruner *Pruner) error

// OptionSetMinResource to set the minimum resource.
func OptionSetMinResource(minResource int) Option {
	return func(p *Pruner) error {
		p.MinResource = minResource
		return nil
	}
}

// OptionSetReductionFactor to set the reduction factor.
func OptionSetReductioinFactor(reductionFactor int) Option {
	return func(p *Pruner) error {
		p.ReductionFactor = reductionFactor
		return nil
	}
}

// OptionSetMinEarlyStoppingRateHigh to set the high of min early stopping rate.
func OptionSetMinEarlyStoppingRateHigh(minEarlyStoppingRateHigh int) Option {
	return func(p *Pruner) error {
		p.MinEarlyStoppingRateHigh = minEarlyStoppingRateHigh
		return nil
	}
}

// OptionSetMinEarlyStoppingRate to set the low min early stopping rate.
func OptionSetMinEarlyStoppingRateLow(minEarlyStoppingRateLow int) Option {
	return func(p *Pruner) error {
		p.MinEarlyStoppingRateLow = minEarlyStoppingRateLow
		return nil
	}
}
