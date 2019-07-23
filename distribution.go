package goptuna

// Distribution represents a parameter that can be optimized.
type Distribution interface {
	// GetName returns the name of the parameter.
	GetName() string
}

var _ Distribution = &UniformDistribution{}

type UniformDistribution struct {
	Name     string
	Max, Min float64
}

func (d *UniformDistribution) GetName() string {
	return d.Name
}
