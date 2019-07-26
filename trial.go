package goptuna

//go:generate stringer -trimprefix TrialState -output stringer_trial_state.go -type=TrialState

// TrialState is a state of Trial
type TrialState int

const (
	// TrialStateRunning means Trial is running.
	TrialStateRunning TrialState = iota
	// TrialStateComplete means Trial has been finished without any error.
	TrialStateComplete
	// TrialStatePruned means Trial has been pruned.
	TrialStatePruned
	// TrialStateFail means Trial has failed due to ann uncaught error.
	TrialStateFail
)

// IsFinished returns true if trial is not running.
func (i TrialState) IsFinished() bool {
	return i != TrialStateRunning
}

// Trial is a process of evaluating an objective function.
//
// This object is passed to an objective function and provides interfaces to get parameter
// suggestion, manage the trial's state of the trial.
// Note that this object is seamlessly instantiated and passed to the objective function behind;
// hence, in typical use cases, library users do not care about instantiation of this object.
type Trial struct {
	study      *Study
	id         string
	state      TrialState
	paramsInIR map[string]float64 // suggested value
	value      float64
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

// SuggestUniform suggests a value for the continuous parameter.
func (t *Trial) SuggestUniform(name string, low, high float64) (float64, error) {
	return t.suggest(name, UniformDistribution{
		High: high, Low: low,
	})
}

// SuggestInt suggests an integer parameter.
func (t *Trial) SuggestInt(name string, low, high int) (int, error) {
	v, err := t.suggest(name, IntUniformDistribution{
		High: high, Low: low,
	})
	return int(v), err
}
