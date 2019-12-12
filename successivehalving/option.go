package successivehalving

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

// OptionSetMinEarlyStoppingRate to set the minimum value of the early stopping rate.
func OptionSetMinEarlyStoppingRate(minEarlyStoppingRate int) Option {
	return func(p *Pruner) error {
		p.MinEarlyStoppingRate = minEarlyStoppingRate
		return nil
	}
}
