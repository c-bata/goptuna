package successivehalving

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/c-bata/goptuna"
)

var errRungNotFound = errors.New("rung not found")

// NewPruner is a constructor for Pruner.
func NewPruner(opts ...Option) (*Pruner, error) {
	pruner := &Pruner{
		MinResource:          1,
		ReductionFactor:      4,
		MinEarlyStoppingRate: 0,
	}

	for _, opt := range opts {
		if err := opt(pruner); err != nil {
			return nil, err
		}
	}
	return pruner, nil
}

// This is a compile-time assertion to check PercentilePruner implements Pruner interface.
var _ goptuna.Pruner = &Pruner{}

// Pruner using Optuna flavored Asynchronous Successive Halving Algorithm.
//
// Successive Halving (arXiv: https://arxiv.org/abs/1502.07943) is a bandit-based algorithm to identify
// the best one among multiple configurations. This is based on Asynchronous Successive Halving Algorithm
// (arXiv: http://arxiv.org/abs/1810.05934), but currently this only supports Optuna flavored Asynchronous
// Successive Halving Algorithm. See https://github.com/optuna/optuna/pull/404 for more details.
type Pruner struct {
	MinResource          int
	ReductionFactor      int
	MinEarlyStoppingRate int
}

func (p *Pruner) Prune(study *goptuna.Study, trial goptuna.FrozenTrial) (bool, error) {
	step, exist := trial.GetLatestStep()
	if !exist {
		return false, nil
	}
	value := trial.IntermediateValues[step]

	rung := getCurrentRung(trial)

	var allTrials []goptuna.FrozenTrial
	for {
		promotionStep := p.MinResource * (int(math.Pow(
			float64(p.ReductionFactor),
			float64(p.MinEarlyStoppingRate+rung))))

		if step < promotionStep {
			return false, nil
		}

		var err error
		if allTrials == nil {
			allTrials, err = study.GetTrials()
			if err != nil {
				return false, err
			}
		}

		err = study.Storage.SetTrialSystemAttr(
			trial.ID, completedRungKey(rung),
			strconv.FormatFloat(value, 'f', -1, 64))
		if err != nil {
			return false, err
		}

		direction := study.Direction()
		if promotable, err := p.isPromotable(rung, value, allTrials, direction); err != nil {
			return false, err
		} else if !promotable {
			return true, nil
		}

		rung++
	}
}

func (p *Pruner) isPromotable(rung int, value float64, allTrials []goptuna.FrozenTrial, direction goptuna.StudyDirection) (bool, error) {
	competingValues := make([]float64, 0, len(allTrials)+1)
	for i := range allTrials {
		v, err := getValueAtRung(allTrials[i], rung)
		if err == errRungNotFound {
			continue
		} else if err != nil {
			return false, err
		}
		competingValues = append(competingValues, v)
	}
	competingValues = append(competingValues, value)
	sort.Float64s(competingValues)

	promotableIdx := (len(competingValues) / p.ReductionFactor) - 1
	if promotableIdx == -1 {
		// Optuna does not support to suspend/resume ongoing trials.
		//
		// For the first `eta - 1` trials, this implementation promotes a trial if its
		// intermediate value is the smallest one among the trials that have completed the rung.
		promotableIdx = 0
	}

	if direction == goptuna.StudyDirectionMaximize {
		l := len(competingValues)
		reversed := make([]float64, l)
		for i := range competingValues {
			reversed[i] = competingValues[l-i-1]
		}
		competingValues = reversed

		return value >= competingValues[promotableIdx], nil
	}
	return value <= competingValues[promotableIdx], nil
}

func getValueAtRung(trial goptuna.FrozenTrial, rung int) (float64, error) {
	rungkey := completedRungKey(rung)
	for key := range trial.SystemAttrs {
		if key == rungkey {
			valuestr := trial.SystemAttrs[key]
			return strconv.ParseFloat(valuestr, 64)
		}
	}
	return -1, errRungNotFound
}

func getCurrentRung(trial goptuna.FrozenTrial) int {
	var currentRung = 0
	for k := range trial.SystemAttrs {
		if !strings.HasPrefix(k, "completed_rung_") {
			continue
		}

		rungstr := strings.TrimPrefix(k, "completed_rung_")
		rung, err := strconv.Atoi(rungstr)
		if err != nil {
			continue
		}

		if rung > currentRung {
			currentRung = rung
		}
	}
	return currentRung
}

func completedRungKey(rung int) string {
	return fmt.Sprintf("completed_rung_%d", rung)
}
