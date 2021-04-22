package successivehalving

// Option to pass the custom option
type Option func(pruner *Pruner) error

// OptionMinResource to set the minimum resource.
func OptionMinResource(minResource int) Option {
	return func(p *Pruner) error {
		p.MinResource = minResource
		return nil
	}
}

// OptionReductionFactor to set the reduction factor.
func OptionReductionFactor(reductionFactor int) Option {
	return func(p *Pruner) error {
		p.ReductionFactor = reductionFactor
		return nil
	}
}

// OptionMinEarlyStoppingRate to set the minimum value of the early stopping rate.
func OptionMinEarlyStoppingRate(minEarlyStoppingRate int) Option {
	return func(p *Pruner) error {
		p.MinEarlyStoppingRate = minEarlyStoppingRate
		return nil
	}
}

// OptionSetMinResource to set the minimum resource.
// Deprecated: please use OptionMinResource instead.
var OptionSetMinResource = OptionMinResource

// OptionSetReductionFactor to set the reduction factor.
// Deprecated: please use OptionReductionFactor instead.
var OptionSetReductionFactor = OptionReductionFactor

// OptionSetMinEarlyStoppingRate to set the minimum value of the early stopping rate.
// Deprecated: please use OptionMinEarlyStoppingRate instead.
var OptionSetMinEarlyStoppingRate = OptionMinEarlyStoppingRate
