package goptuna

import (
	"encoding/json"
)

// Distribution represents a parameter that can be optimized.
type Distribution interface {
	// ToInternalRepr to convert external representation of a parameter value into internal representation.
	ToInternalRepr(interface{}) float64
	// ToExternalRepr to convert internal representation of a parameter value into external representation.
	ToExternalRepr(float64) interface{}
	// Single to test whether the range of this distribution contains just a single value.
	Single() bool
	// Contains to check a parameter value is contained in the range of this distribution.
	Contains(float64) bool
}

var _ Distribution = &UniformDistribution{}

// UniformDistribution is a uniform distribution in the linear domain.
type UniformDistribution struct {
	// High is higher endpoint of the range of the distribution (included in the range).
	High float64 `json:"high"`
	// Low is lower endpoint of the range of the distribution (included in the range).
	Low float64 `json:"low"`
}

// UniformDistributionName is the identifier name of UniformDistribution
const UniformDistributionName = "UniformDistribution"

// ToInternalRepr to convert external representation of a parameter value into internal representation.
func (d *UniformDistribution) ToInternalRepr(xr interface{}) float64 {
	return xr.(float64)
}

// ToExternalRepr to convert internal representation of a parameter value into external representation.
func (d *UniformDistribution) ToExternalRepr(ir float64) interface{} {
	return ir
}

// Single to test whether the range of this distribution contains just a single value.
func (d *UniformDistribution) Single() bool {
	return d.High == d.Low
}

// Contains to check a parameter value is contained in the range of this distribution.
func (d *UniformDistribution) Contains(ir float64) bool {
	if d.Single() {
		return ir == d.Low
	}
	return d.Low <= ir && ir < d.High
}

var _ Distribution = &IntUniformDistribution{}

// IntUniformDistribution is a uniform distribution on integers.
type IntUniformDistribution struct {
	// High is higher endpoint of the range of the distribution (included in the range).
	High int `json:"high"`
	// Low is lower endpoint of the range of the distribution (included in the range).
	Low int `json:"low"`
}

// IntUniformDistributionName is the identifier name of IntUniformDistribution
const IntUniformDistributionName = "IntUniformDistribution"

// ToInternalRepr to convert external representation of a parameter value into internal representation.
func (d *IntUniformDistribution) ToInternalRepr(xr interface{}) float64 {
	x := xr.(int)
	return float64(x)
}

// ToExternalRepr to convert internal representation of a parameter value into external representation.
func (d *IntUniformDistribution) ToExternalRepr(ir float64) interface{} {
	return int(ir)
}

// Single to test whether the range of this distribution contains just a single value.
func (d *IntUniformDistribution) Single() bool {
	return d.High == d.Low
}

// Contains to check a parameter value is contained in the range of this distribution.
func (d *IntUniformDistribution) Contains(ir float64) bool {
	value := int(ir)
	if d.Single() {
		return value == d.Low
	}
	return d.Low <= value && value < d.High
}

var _ Distribution = &CategoricalDistribution{}

// CategoricalDistribution is a distribution for categorical parameters
type CategoricalDistribution struct {
	// Choices is a candidates of parameter values
	Choices []string `json:"choices"`
}

// CategoricalDistributionName is the identifier name of CategoricalDistribution
const CategoricalDistributionName = "CategoricalDistribution"

// ToInternalRepr to convert external representation of a parameter value into internal representation.
func (d *CategoricalDistribution) ToInternalRepr(er interface{}) float64 {
	value := er.(string)
	for i := range d.Choices {
		if d.Choices[i] == value {
			return float64(i)
		}
	}
	panic("must not reach here")
}

// ToExternalRepr to convert internal representation of a parameter value into external representation.
func (d *CategoricalDistribution) ToExternalRepr(ir float64) interface{} {
	return d.Choices[int(ir)]
}

// Single to test whether the range of this distribution contains just a single value.
func (d *CategoricalDistribution) Single() bool {
	return len(d.Choices) == 1
}

// Contains to check a parameter value is contained in the range of this distribution.
func (d *CategoricalDistribution) Contains(ir float64) bool {
	index := int(ir)
	return 0 <= index && index < len(d.Choices)
}

// DistributionToJSON serialize a distribution to JSON format.
func DistributionToJSON(distribution interface{}) ([]byte, error) {
	var ir struct {
		Name  string      `json:"name"`
		Attrs interface{} `json:"attributes"`
	}
	switch distribution.(type) {
	case UniformDistribution:
		ir.Name = UniformDistributionName
	case IntUniformDistribution:
		ir.Name = IntUniformDistributionName
	default:
		return nil, ErrUnknownDistribution
	}
	ir.Attrs = distribution
	return json.Marshal(&ir)
}

// JSONToDistribution deserialize a distribution in JSON format.
func JSONToDistribution(jsonBytes []byte) (interface{}, error) {
	var x struct {
		Name  string      `json:"name"`
		Attrs interface{} `json:"attributes"`
	}
	err := json.Unmarshal(jsonBytes, &x)
	if err != nil {
		return nil, err
	}
	switch x.Name {
	case UniformDistributionName:
		var y UniformDistribution
		var dbytes []byte
		dbytes, err = json.Marshal(x.Attrs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dbytes, &y)
		return y, err
	case IntUniformDistributionName:
		var y IntUniformDistribution
		var dbytes []byte
		dbytes, err = json.Marshal(x.Attrs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dbytes, &y)
		return y, err
	case CategoricalDistributionName:
		var y CategoricalDistribution
		var dbytes []byte
		dbytes, err = json.Marshal(x.Attrs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dbytes, &y)
		return y, err
	}
	return nil, ErrUnknownDistribution
}
