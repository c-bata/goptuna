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
	Study *Study
	ID    int
	state TrialState
	value float64
}

func (t *Trial) suggest(name string, distribution interface{}) (float64, error) {
	trial, err := t.Study.Storage.GetTrial(t.ID)
	if err != nil {
		return 0.0, err
	}

	v, err := t.Study.Sampler.Sample(t.Study, trial, name, distribution)
	if err != nil {
		return 0.0, err
	}

	var d Distribution
	switch x := distribution.(type) {
	case UniformDistribution:
		d = &x
	case IntUniformDistribution:
		d = &x
	case CategoricalDistribution:
		d = &x
	default:
		return -1, ErrUnknownDistribution
	}

	if trial.Params == nil {
		trial.Params = make(map[string]interface{}, 8)
	}
	trial.Params[name] = d.ToExternalRepr(v)
	if trial.ParamsInIR == nil {
		trial.ParamsInIR = make(map[string]float64, 8)
	}
	trial.ParamsInIR[name] = v

	err = t.Study.Storage.SetTrialParam(trial.ID, name, v, d)
	return v, err
}

// Report an intermediate value of an objective function
func (t *Trial) Report(value float64, step int) error {
	err := t.Study.Storage.SetTrialIntermediateValue(t.ID, step, value)
	if err != nil {
		return err
	}
	return t.Study.Storage.SetTrialValue(t.ID, value)
}

// Number return trial's number which is consecutive and unique in a study.
func (t *Trial) Number() (int, error) {
	return t.Study.Storage.GetTrialNumberFromID(t.ID)
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

// SuggestCategorical suggests an categorical parameter.
func (t *Trial) SuggestCategorical(name string, choices []string) (string, error) {
	v, err := t.suggest(name, CategoricalDistribution{
		Choices: choices,
	})
	return choices[int(v)], err
}
