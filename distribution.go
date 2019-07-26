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
	} else {
		return d.Low <= ir && ir < d.High
	}
}

var _ Distribution = &IntUniformDistribution{}

// IntUniformDistribution is a uniform distribution on integers.
type IntUniformDistribution struct {
	// High is higher endpoint of the range of the distribution (included in the range).
	High int `json:"high"`
	// Low is lower endpoint of the range of the distribution (included in the range).
	Low int `json:"low"`
}

// ToInternalRepr to convert external representation of a parameter value into internal representation.
func (d *IntUniformDistribution) ToInternalRepr(xr interface{}) float64 {
	return xr.(float64)
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
	} else {
		return d.Low <= value && value < d.High
	}
}

// DistributionToJSON serialize a distribution to JSON format.
func DistributionToJSON(distribution interface{}) ([]byte, error) {
	var ir struct {
		Name  string      `json:"name"`
		Attrs interface{} `json:"attributes"`
	}
	switch distribution.(type) {
	case UniformDistribution:
		ir.Name = "UniformDistribution"
	case IntUniformDistribution:
		ir.Name = "IntUniformDistribution"
	default:
		return nil, ErrUnexpectedDistribution
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
	case "UniformDistribution":
		var y UniformDistribution
		var dbytes []byte
		dbytes, err = json.Marshal(x.Attrs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dbytes, &y)
		return y, err
	case "IntUniformDistribution":
		var y IntUniformDistribution
		var dbytes []byte
		dbytes, err = json.Marshal(x.Attrs)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dbytes, &y)
		return y, err
	}
	return nil, ErrUnexpectedDistribution
}
