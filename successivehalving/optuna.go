package successivehalving

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/c-bata/goptuna"
)

var errRungNotFound = errors.New("rung not found")

// NewOptunaPruner is a constructor for OptunaSuccessiveHalvingPruner.
func NewOptunaPruner() *OptunaSuccessiveHalvingPruner {
	return &OptunaSuccessiveHalvingPruner{
		MinResource:          1,
		ReductionFactor:      4,
		MinEarlyStoppingRate: 0,
	}
}

// This is a compile-time assertion to check PercentilePruner implements Pruner interface.
var _ goptuna.Pruner = &OptunaSuccessiveHalvingPruner{}

// OptunaSuccessiveHalvingPruner is Optuna-flavored Asynchronous Successive Halving Algorithm.
// See https://github.com/pfnet/optuna/pull/404 for details.
type OptunaSuccessiveHalvingPruner struct {
	MinResource          int
	ReductionFactor      int
	MinEarlyStoppingRate int
}

// Prune by Optuna-flavored Asynchronous Successive Halving Algorithm.
func (p *OptunaSuccessiveHalvingPruner) Prune(study *goptuna.Study, trial goptuna.FrozenTrial) (bool, error) {
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

		if math.IsNaN(value) {
			// todo(c-bata): need to check this line.
			return true, nil
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
			fmt.Sprintf("%f", value))
		if err != nil {
			return false, err
		}

		direction := study.Direction()
		if promotable, err := p.isPromotable(rung, value, allTrials, direction); err != nil {
			return false, err
		} else if promotable {
			return true, nil
		}

		rung++
	}
}

func (p *OptunaSuccessiveHalvingPruner) isPromotable(rung int, value float64, allTrials []goptuna.FrozenTrial, direction goptuna.StudyDirection) (bool, error) {
	competingValues := make([]float64, 0, len(allTrials))
	for i := range allTrials {
		value, err := getValueAtRung(allTrials[i], rung)
		if err != nil {
			return false, err
		}
		competingValues = append(competingValues, value)
	}

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
