package hyperband

import (
	"math"

	"github.com/c-bata/goptuna"
	"github.com/c-bata/goptuna/successivehalving"
)

func NewPruner(opts ...Option) (*Pruner, error) {
	pruner := &Pruner{
		MinResource:              1,
		ReductionFactor:          4,
		MinEarlyStoppingRateLow:  0,
		MinEarlyStoppingRateHigh: 4,
	}

	for _, opt := range opts {
		if err := opt(pruner); err != nil {
			return nil, err
		}
	}

	nPruners := pruner.MinEarlyStoppingRateHigh - pruner.MinEarlyStoppingRateLow + 1
	pruner.successiveHalvingPruners = make([]successivehalving.Pruner, nPruners)
	pruner.bracketResourceBudgets = make([]int, nPruners)
	pruner.resourceBudget = 0
	for i := 0; i < nPruners; i++ {
		bracketResourceBudget := calcBracketResourceBudget(pruner.ReductionFactor, i, nPruners)
		pruner.resourceBudget += bracketResourceBudget
		pruner.bracketResourceBudgets[i] = bracketResourceBudget
		pruner.successiveHalvingPruners[i] = successivehalving.Pruner{
			MinResource:          pruner.MinResource,
			ReductionFactor:      pruner.ReductionFactor,
			MinEarlyStoppingRate: pruner.MinEarlyStoppingRateLow + i,
		}
	}
	return pruner, nil
}

func calcBracketResourceBudget(reductionFactor, prunerIndex, nPruners int) int {
	n := int(math.Pow(float64(reductionFactor), float64(nPruners-1)))
	budget := n
	for i := nPruners - 1; i < prunerIndex; i++ {
		budget += n / 2
	}
	return budget
}

// This is a compile-time assertion to check PercentilePruner implements Pruner interface.
var _ goptuna.Pruner = &Pruner{}

// Pruner is Optuna-flavored Hyperband Algorithm.
type Pruner struct {
	MinResource              int
	ReductionFactor          int
	MinEarlyStoppingRateLow  int
	MinEarlyStoppingRateHigh int

	successiveHalvingPruners []successivehalving.Pruner
	bracketResourceBudgets   []int
	resourceBudget           int
}

func (p *Pruner) Prune(study *goptuna.Study, trial goptuna.FrozenTrial) (bool, error) {
	bracketID := p.getBracketID(trial.Number)
	return p.successiveHalvingPruners[bracketID].Prune(study, trial)
}

func (p *Pruner) getBracketID(trialNumber int) int {
	nPruners := p.MinEarlyStoppingRateHigh - p.MinEarlyStoppingRateLow + 1
	n := trialNumber % p.resourceBudget
	for i := 0; i < nPruners; i++ {
		n -= p.bracketResourceBudgets[i]
		if n < 0 {
			return i
		}
	}
	panic("failed to get bracket id")
}
