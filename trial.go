package goptuna

//go:generate stringer -trimprefix TrialState -output stringer_trial_state.go -type=TrialState
type TrialState int

const (
	TrialStateRunning TrialState = iota
	TrialStateComplete
	TrialStatePruned
	TrialStateFail
)

func (t TrialState) IsFinished() bool {
	return t != TrialStateRunning
}

type Trial struct {
	study      *Study
	id         string
	state      TrialState
	paramsInIR map[string]float64 // suggested value
	value      float64
}

func (t *Trial) Frozen() FrozenTrial {
	return FrozenTrial{
		ID:         t.id,
		StudyID:    t.study.id,
		State:      t.state,
		Value:      t.value,
		ParamsInIR: t.paramsInIR,
	}
}

func (t *Trial) suggest(name string, distribution interface{}) (float64, error) {
	trial, err := t.study.storage.GetTrial(t.id)
	if err != nil {
		return 0.0, err
	}

	v, err := t.study.sampler.Sample(t.study, trial, name, distribution)
	if err != nil {
		return 0.0, err
	}

	if t.paramsInIR == nil {
		t.paramsInIR = make(map[string]float64, 8)
	}
	t.paramsInIR[name] = v
	err = t.study.storage.SetTrialParam(trial.ID, name, v)
	return v, err
}

func (t *Trial) SuggestUniform(name string, low, high float64) (float64, error) {
	return t.suggest(name, UniformDistribution{
		Name: name, High: high, Low: low,
	})
}

func (t *Trial) SuggestInt(name string, low, high int) (int, error) {
	v, err := t.suggest(name, IntUniformDistribution{
		Name: name, High: high, Low: low,
	})
	return int(v), err
}
