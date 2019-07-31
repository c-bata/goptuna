package medianstopping

import (
	"errors"
	"math"
	"sort"

	"github.com/c-bata/goptuna"
)

func NewPercentilePruner(q float64) (*PercentilePruner, error) {
	if q > 100 || q < 0 {
		return nil, errors.New("please specify the percentile between 0 and 100")
	}
	return &PercentilePruner{
		Percentile:     q,
		NStartUpTrials: 5,
		NWarmUpSteps:   0,
	}, nil
}

type PercentilePruner struct {
	Percentile     float64
	NStartUpTrials int
	NWarmUpSteps   int
}

func getCompletedTrials(storage goptuna.Storage, studyID int) ([]goptuna.FrozenTrial, error) {
	trials, err := storage.GetAllTrials(studyID)
	if err != nil {
		return nil, err
	}
	completed := make([]goptuna.FrozenTrial, 0, len(trials))
	for i := range trials {
		if trials[i].State == goptuna.TrialStateComplete {
			completed = append(completed, trials[i])
		}
	}
	return completed, nil
}

func getBestIntermediateResultOverSteps(trial goptuna.FrozenTrial, direction goptuna.StudyDirection) float64 {
	keys := make([]int, 0, len(trial.IntermediateValues))
	for k := range trial.IntermediateValues {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	values := make([]float64, len(trial.IntermediateValues))
	for i, k := range keys {
		values[i] = trial.IntermediateValues[k]
	}

	if len(values) == 0 {
		return math.NaN()
	}

	var bestResult float64
	for i := range values {
		if direction == goptuna.StudyDirectionMaximize {
			bestResult = math.Max(bestResult, values[i])
		} else {
			bestResult = math.Min(bestResult, values[i])
		}
	}
	return bestResult
}

func getPercentileIntermediateResultOverSteps(
	trials []goptuna.FrozenTrial,
	step int,
	q float64,
	direction goptuna.StudyDirection,
) float64 {
	if len(trials) == 0 {
		panic("unreachable")
	}

	if direction == goptuna.StudyDirectionMaximize {
		q = 100 - q
	}

	intermediateValues := make([]float64, 0, len(trials))
	for i := range trials {
		value, ok := trials[i].IntermediateValues[step]
		if !ok {
			continue
		}
		intermediateValues = append(intermediateValues, value)
	}

	if len(intermediateValues) == 0 {
		return math.NaN()
	}
	return percentile(intermediateValues, q)
}

func (p *PercentilePruner) Prune(storage goptuna.Storage, studyID, trialID, step int) (bool, error) {
	completedTrials, err := getCompletedTrials(storage, studyID)
	if err != nil {
		return false, err
	}
	ntrials := len(completedTrials)
	if ntrials == 0 {
		return false, nil
	}
	if ntrials < p.NStartUpTrials {
		return false, nil
	}
	if step <= p.NWarmUpSteps {
		return false, nil
	}

	trial, err := storage.GetTrial(trialID)
	if err != nil {
		return false, err
	}
	if len(trial.IntermediateValues) == 0 {
		return false, nil
	}

	direction, err := storage.GetStudyDirection(studyID)
	if err != nil {
		return false, err
	}
	bestIntermediateResult := getBestIntermediateResultOverSteps(trial, direction)
	if math.IsNaN(bestIntermediateResult) {
		return true, nil
	}

	percentileResult := getPercentileIntermediateResultOverSteps(completedTrials, step, p.Percentile, direction)
	if math.IsNaN(percentileResult) {
		return false, nil
	}

	if direction == goptuna.StudyDirectionMaximize {
		return bestIntermediateResult < percentileResult, nil
	}
	return bestIntermediateResult > percentileResult, nil
}
