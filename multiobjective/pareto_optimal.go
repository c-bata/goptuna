package multiobjective

import (
	"fmt"

	"github.com/c-bata/goptuna"
)

func isDominate(self, other goptuna.FrozenTrial, directions map[string]goptuna.StudyDirection) error {
	metricsSelf, err := GetSubMetrics(self.SystemAttrs)
	if err != nil {
		return err
	}
	metricsOther, err := GetSubMetrics(other.SystemAttrs)
	if err != nil {
		return err
	}

	for name := range directions {
		valueSelf, ok := metricsSelf[name]
		if !ok {
			return fmt.Errorf("trial %d does not contain '%s' metric", self.ID, name)
		}
		valueOther, ok := metricsOther[name]
		if !ok {
			return fmt.Errorf("trial %d does not contain '%s' metric", other.ID, name)
		}
		if directions[name] == goptuna.StudyDirectionMaximize {
			valueSelf = -valueSelf
			valueOther = -valueOther
		}
		if valueSelf >= valueOther {
			return fmt.Errorf("trial %d is not dominant in '%s' metric", self.ID, name)
		}
	}
	return nil
}

// GetParetoOptimalTrials returns Pareto-optimal front solutions.
func GetParetoOptimalTrials(
	study *goptuna.Study,
	directions map[string]goptuna.StudyDirection,
) ([]goptuna.FrozenTrial, error) {
	trials, err := study.Storage.GetAllTrials(study.ID)
	if err != nil {
		return nil, err
	}

	paretoFront := make([]goptuna.FrozenTrial, 0, 8)
	for i := range trials {
		if trials[i].State != goptuna.TrialStateComplete {
			continue
		}

		var dominated bool
		for _, other := range trials {
			if other.State != goptuna.TrialStateComplete {
				continue
			}
			if trials[i].ID == other.ID {
				continue
			}

			if err := isDominate(trials[i], other, directions); err == nil {
				dominated = true
				break
			}
		}

		if dominated {
			paretoFront = append(paretoFront, trials[i])
		}
	}
	return paretoFront, nil
}
