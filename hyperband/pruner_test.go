package hyperband

import (
	"testing"

	"github.com/c-bata/goptuna"
)

func TestPruner_Prune(t *testing.T) {
	pruner, err := NewPruner(
		OptionSetMinResource(1),
		OptionSetReductioinFactor(2),
		OptionSetMinEarlyStoppingRateLow(0),
		OptionSetMinEarlyStoppingRateHigh(3),
	)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
	nPruners := 3 - 0 + 1
	study, err := goptuna.CreateStudy(
		"test",
		goptuna.StudyOptionPruner(pruner))
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}

	objective := func(trial goptuna.Trial) (float64, error) {
		for i := 0; i < 10; i++ {
			_ = trial.Report(float64(i), i)
		}
		return 1.0, nil
	}

	err = study.Optimize(objective, nPruners*10)
	if err != nil {
		t.Errorf("should be nil, but got %s", err)
		return
	}
}
