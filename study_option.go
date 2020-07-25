package goptuna

// StudyOption to pass the custom option
type StudyOption func(study *Study) error

// StudyOptionDirection change the direction of optimize
func StudyOptionDirection(direction StudyDirection) StudyOption {
	return func(s *Study) error {
		s.direction = direction
		return nil
	}
}

// StudyOptionLogger sets Logger.
func StudyOptionLogger(logger Logger) StudyOption {
	return func(s *Study) error {
		if logger == nil {
			s.logger = &StdLogger{Logger: nil}
		} else {
			s.logger = logger
		}
		return nil
	}
}

// StudyOptionStorage sets the storage object.
func StudyOptionStorage(storage Storage) StudyOption {
	return func(s *Study) error {
		s.Storage = storage
		return nil
	}
}

// StudyOptionSampler sets the sampler object.
func StudyOptionSampler(sampler Sampler) StudyOption {
	return func(s *Study) error {
		s.Sampler = sampler
		return nil
	}
}

// StudyOptionRelativeSampler sets the relative sampler object.
func StudyOptionRelativeSampler(sampler RelativeSampler) StudyOption {
	return func(s *Study) error {
		s.RelativeSampler = sampler
		return nil
	}
}

// StudyOptionPruner sets the pruner object.
func StudyOptionPruner(pruner Pruner) StudyOption {
	return func(s *Study) error {
		s.Pruner = pruner
		return nil
	}
}

// StudyOptionIgnoreError is an option to continue even if
// it receive error while running Optimize method.
func StudyOptionIgnoreError(ignore bool) StudyOption {
	return func(s *Study) error {
		s.ignoreErr = ignore
		return nil
	}
}

// StudyOptionSetTrialNotifyChannel to subscribe the finished trials.
func StudyOptionSetTrialNotifyChannel(notify chan FrozenTrial) StudyOption {
	return func(s *Study) error {
		s.trialNotification = notify
		return nil
	}
}

// StudyOptionLoadIfExists to load the study if exists.
func StudyOptionLoadIfExists(loadIfExists bool) StudyOption {
	return func(s *Study) error {
		s.loadIfExists = loadIfExists
		return nil
	}
}

// StudyOptionInitialSearchSpace to use RelativeSampler from the first trial.
// This option is useful for Define-and-Run interface.
func StudyOptionDefineSearchSpace(space map[string]interface{}) StudyOption {
	return func(s *Study) error {
		s.definedSearchSpace = space
		return nil
	}
}

// StudyOptionSetLogger sets Logger.
// Deprecated: please use StudyOptionLogger instead.
var StudyOptionSetLogger = StudyOptionLogger

// StudyOptionSetDirection change the direction of optimize
// Deprecated: please use StudyOptionDirection instead.
var StudyOptionSetDirection = StudyOptionDirection
