package goptuna

// Distribution represents a parameter that can be optimized.
type Distribution interface {
	// GetName returns the name of the parameter.
	GetName() string
}

var _ Distribution = &UniformDistribution{}

// UniformDistribution is a uniform distribution in the linear domain.
type UniformDistribution struct {
	Name string
	// High is higher endpoint of the range of the distribution (included in the range).
	Max float64
	// Low is lower endpoint of the range of the distribution (included in the range).
	Min float64
}

func (d *UniformDistribution) GetName() string {
	return d.Name
}

var _ Distribution = &IntUniformDistribution{}

// IntUniformDistribution is a uniform distribution on integers.
type IntUniformDistribution struct {
	Name string
	// High is higher endpoint of the range of the distribution (included in the range).
	High int
	// Low is lower endpoint of the range of the distribution (included in the range).
	Low int
}

func (d *IntUniformDistribution) GetName() string {
	return d.Name
}
