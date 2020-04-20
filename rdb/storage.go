package rdb

import (
	"errors"
	"fmt"
	"time"

	"github.com/c-bata/goptuna"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var _ goptuna.Storage = &Storage{}

// NewStorage returns new RDB storage.
func NewStorage(db *gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// Storage stores data in your relational databases.
type Storage struct {
	db *gorm.DB
}

// CreateNewStudy creates study and returns studyID.
func (s *Storage) CreateNewStudy(name string) (int, error) {
	if name == "" {
		u, err := uuid.NewUUID()
		if err != nil {
			return -1, err
		}
		name = goptuna.DefaultStudyNamePrefix + u.String()
	}
	study := &studyModel{
		Name:      name,
		Direction: directionNotSet,
	}
	err := s.db.Create(study).Error
	return study.ID, err
}

// DeleteNewStudy creates study and returns studyID.
func (s *Storage) DeleteStudy(studyID int) error {
	return s.db.Delete(&studyModel{
		ID: studyID,
	}).Error
}

// SetStudyDirection sets study direction of the objective.
func (s *Storage) SetStudyDirection(studyID int, direction goptuna.StudyDirection) error {
	d := directionMinimize
	if direction == goptuna.StudyDirectionMaximize {
		d = directionMaximize
	}

	err := s.db.Model(&studyModel{}).
		Where("study_id = ?", studyID).
		Update("direction", d).Error
	return err
}

// SetStudyUserAttr to store the value for the user.
func (s *Storage) SetStudyUserAttr(studyID int, key string, value string) error {
	var result studyUserAttributeModel
	return s.db.Where(&studyUserAttributeModel{
		UserAttributeReferStudy: studyID,
		Key:                     key,
	}).Assign(&studyUserAttributeModel{
		UserAttributeReferStudy: studyID,
		Key:                     key,
		ValueJSON:               encodeAttrValue(value),
	}).FirstOrCreate(&result).Error
}

// SetStudySystemAttr to store the value for the system.
func (s *Storage) SetStudySystemAttr(studyID int, key string, value string) error {
	var result studySystemAttributeModel
	return s.db.Where(&studySystemAttributeModel{
		SystemAttributeReferStudy: studyID,
		Key:                       key,
	}).Assign(&studySystemAttributeModel{
		SystemAttributeReferStudy: studyID,
		Key:                       key,
		ValueJSON:                 encodeAttrValue(value),
	}).FirstOrCreate(&result).Error
}

// GetStudyIDFromName return the study id from study name.
func (s *Storage) GetStudyIDFromName(name string) (int, error) {
	var study studyModel
	err := s.db.First(&study, "study_name = ?", name).Error
	return study.ID, err
}

// GetStudyIDFromTrialID return the study id from trial id.
func (s *Storage) GetStudyIDFromTrialID(trialID int) (int, error) {
	var trial trialModel
	err := s.db.First(&trial, "trial_id = ?", trialID).Error
	return trial.TrialReferStudy, err
}

// GetStudyNameFromID return the study name from study id.
func (s *Storage) GetStudyNameFromID(studyID int) (string, error) {
	var study studyModel
	s.db.First(&study, "study_id = ?", studyID)
	return study.Name, nil
}

// GetStudyUserAttrs to restore the attributes for the user.
func (s *Storage) GetStudyUserAttrs(studyID int) (map[string]string, error) {
	var attrs []studyUserAttributeModel
	err := s.db.Find(&attrs, "study_id = ?", studyID).Error
	if err != nil {
		return nil, err
	}

	res := make(map[string]string, len(attrs))
	for i := range attrs {
		res[attrs[i].Key] = decodeAttrValue(attrs[i].ValueJSON)
	}
	return res, nil
}

// GetStudySystemAttrs to restore the attributes for the system.
func (s *Storage) GetStudySystemAttrs(studyID int) (map[string]string, error) {
	var attrs []studySystemAttributeModel
	err := s.db.Find(&attrs, "study_id = ?", studyID).Error
	if err != nil {
		return nil, err
	}

	res := make(map[string]string, len(attrs))
	for i := range attrs {
		res[attrs[i].Key] = decodeAttrValue(attrs[i].ValueJSON)
	}
	return res, nil
}

// GetAllStudySummaries returns all study summaries.
func (s *Storage) GetAllStudySummaries() ([]goptuna.StudySummary, error) {
	var err error
	var studies []studyModel
	err = s.db.
		Preload("UserAttributes").
		Preload("SystemAttributes").
		Preload("Trials").
		Find(&studies).Error
	if err != nil {
		return nil, err
	}

	res := make([]goptuna.StudySummary, len(studies))
	for i, study := range studies {
		var best *trialModel
		var start *time.Time
		for i := range study.Trials {
			if study.Trials[i].State == trialStateComplete {
				if best == nil {
					best = &study.Trials[i]
				}
				if study.Direction == directionMaximize {
					if best.Value < study.Trials[i].Value {
						best = &study.Trials[i]
					}
				} else {
					if best.Value > study.Trials[i].Value {
						best = &study.Trials[i]
					}
				}
			}

			if start == nil {
				start = study.Trials[i].DatetimeStart
			}
			if start.Unix() < study.Trials[i].DatetimeStart.Unix() {
				start = study.Trials[i].DatetimeStart
			}
		}

		var ft goptuna.FrozenTrial
		if best != nil {
			ft, err = s.GetTrial(best.ID)
			if err != nil {
				return nil, err
			}
		}
		var studyStart time.Time
		if start != nil {
			studyStart = *start
		}
		ss, err := toStudySummary(studies[i], ft, studyStart)
		if err != nil {
			return nil, err
		}
		res[i] = ss
	}
	return res, nil
}

// CreateNewTrial creates trial and returns trialID.
func (s *Storage) CreateNewTrial(studyID int) (int, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return -1, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create a new trial
	start := time.Now()
	trial := &trialModel{
		TrialReferStudy: studyID,
		State:           trialStateRunning,
		DatetimeStart:   &start,
	}
	if err := tx.Create(trial).Error; err != nil {
		tx.Rollback()
		return -1, err
	}

	// Calculate the trial number
	var number int
	err := tx.Model(&trialModel{}).
		Where("study_id = ?", studyID).
		Where("trial_id < ?", trial.ID).
		Count(&number).Error
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	err = tx.Model(&trialModel{}).
		Where("trial_id = ?", trial.ID).
		Update("number", number).Error
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	err = tx.Commit().Error
	return trial.ID, err
}

// CloneTrial creates new Trial from the given base Trial.
func (s *Storage) CloneTrial(studyID int, baseTrial goptuna.FrozenTrial) (int, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return -1, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Because only `RUNNING` trials can be updated,
	// we temporarily set the state of the new trial to `RUNNING`.
	// After all fields of the trial have been updated,
	// the state is set to `template_trial.state`.
	tempState := trialStateWaiting

	// Avoid to insert zero-time for `NON_ZERO_DATE` mode on MySQL.
	// See https://github.com/jinzhu/gorm/issues/595
	var datetimeStart, datetimeComplete *time.Time
	if !baseTrial.DatetimeStart.IsZero() {
		datetimeStart = &baseTrial.DatetimeStart
	}
	if !baseTrial.DatetimeComplete.IsZero() {
		datetimeComplete = &baseTrial.DatetimeComplete
	}

	trial := &trialModel{
		TrialReferStudy:  studyID,
		State:            tempState,
		Value:            baseTrial.Value,
		DatetimeStart:    datetimeStart,
		DatetimeComplete: datetimeComplete,
	}
	if err := tx.Create(trial).Error; err != nil {
		tx.Rollback()
		return -1, err
	}

	// params
	for name := range baseTrial.InternalParams {
		d, ok := baseTrial.Distributions[name]
		if !ok {
			tx.Rollback()
			return -1, fmt.Errorf("'%s' distribution is not found", name)
		}
		jsonBytes, err := goptuna.DistributionToJSON(d)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
		err = tx.Create(&trialParamModel{
			TrialParamReferTrial: trial.ID,
			Name:                 name,
			Value:                baseTrial.InternalParams[name],
			DistributionJSON:     string(jsonBytes),
		}).Error
	}

	// user attrs
	for key := range baseTrial.UserAttrs {
		err := tx.Create(&trialUserAttributeModel{
			UserAttributeReferTrial: trial.ID,
			Key:                     key,
			ValueJSON:               encodeAttrValue(baseTrial.UserAttrs[key]),
		}).Error
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	// system attrs
	for key := range baseTrial.SystemAttrs {
		err := tx.Create(&trialSystemAttributeModel{
			SystemAttributeReferTrial: trial.ID,
			Key:                       key,
			ValueJSON:                 encodeAttrValue(baseTrial.SystemAttrs[key]),
		}).Error
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	// intermediate values
	for step := range baseTrial.IntermediateValues {
		err := tx.Create(&trialValueModel{
			TrialValueReferTrial: trial.ID,
			Step:                 step,
			Value:                baseTrial.IntermediateValues[step],
		}).Error
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	// state
	state, err := toStateInternalRepresentation(baseTrial.State)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	err = tx.Model(&trialModel{}).
		Where("trial_id = ?", trial.ID).
		Update("state", state).Error
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// trial number
	var number int
	err = tx.Model(&trialModel{}).
		Where("study_id = ?", studyID).
		Where("trial_id < ?", trial.ID).
		Count(&number).Error
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	err = tx.Model(&trialModel{}).
		Where("trial_id = ?", trial.ID).
		Update("number", number).Error
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	err = tx.Commit().Error
	return trial.ID, err
}

// SetTrialValue sets the value of trial.
func (s *Storage) SetTrialValue(trialID int, value float64) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	var trial trialModel
	err := tx.First(&trial, "trial_id = ?", trialID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	state, err := toStateExternalRepresentation(trial.State)
	if err != nil {
		tx.Rollback()
		return err
	}
	if state.IsFinished() {
		tx.Rollback()
		return goptuna.ErrTrialCannotBeUpdated
	}

	err = tx.Model(&trialModel{}).
		Where("trial_id = ?", trialID).
		Update("value", value).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// SetTrialIntermediateValue sets the intermediate value of trial.
// While sets the intermediate value, trial.value is also updated.
func (s *Storage) SetTrialIntermediateValue(trialID int, step int, value float64) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	// Check whether trial is finished.
	var trial trialModel
	err := tx.First(&trial, "trial_id = ?", trialID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	state, err := toStateExternalRepresentation(trial.State)
	if err != nil {
		tx.Rollback()
		return err
	}
	if state.IsFinished() {
		tx.Rollback()
		return goptuna.ErrTrialCannotBeUpdated
	}

	// If trial value is already exist, then do rollback.
	result := tx.First(&trialValueModel{}, "trial_id = ? AND step = ?", trialID, step)
	if !(result.Error != nil && result.RecordNotFound()) {
		tx.Rollback()
		return err
	}

	// Set trial intermediate value
	err = tx.Create(&trialValueModel{
		TrialValueReferTrial: trialID,
		Step:                 step,
		Value:                value,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// SetTrialParam sets the sampled parameters of trial.
func (s *Storage) SetTrialParam(
	trialID int,
	paramName string,
	paramValueInternal float64,
	distribution interface{},
) error {
	j, err := goptuna.DistributionToJSON(distribution)
	if err != nil {
		return err
	}
	err = s.db.Create(&trialParamModel{
		TrialParamReferTrial: trialID,
		Name:                 paramName,
		Value:                paramValueInternal,
		DistributionJSON:     string(j),
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// SetTrialState sets the state of trial.
func (s *Storage) SetTrialState(trialID int, state goptuna.TrialState) error {
	xr, err := toStateInternalRepresentation(state)
	if err != nil {
		return err
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	var trial trialModel
	var result *gorm.DB
	if s.db.Dialect().GetName() == "sqlite3" {
		// TODO(c-bata): Fix concurrency problem on SQLite3.
		// SQLite3 cannot interpret `FOR UPDATE` clause.
		result = tx.First(&trial, "trial_id = ?", trialID)
	} else {
		result = tx.Set("gorm:query_option", "FOR UPDATE").
			First(&trial, "trial_id = ?", trialID)
	}
	if result.RecordNotFound() {
		tx.Rollback()
		return goptuna.ErrInvalidTrialID
	}
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	previousState, err := toStateExternalRepresentation(trial.State)
	if err != nil {
		tx.Rollback()
		return err
	}
	if previousState.IsFinished() || previousState == state {
		tx.Rollback()
		return goptuna.ErrTrialCannotBeUpdated
	}

	err = tx.Model(&trialModel{}).
		Where("trial_id = ?", trialID).
		Update("state", xr).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if state.IsFinished() {
		completedAt := time.Now()
		err = tx.Model(&trialModel{}).
			Where("trial_id = ?", trialID).
			Update("datetime_complete", completedAt).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// SetTrialValueAndState sets the evaluation value and state of trial at the same time.
func (s *Storage) SetTrialValueAndState(trialID int, value float64, state goptuna.TrialState) error {
	if !state.IsFinished() {
		return errors.New("SetTrialValueAndState only supports finished state")
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	// skip to check previous state because it set `state`.
	xr, err := toStateInternalRepresentation(state)
	if err != nil {
		return err
	}

	completedAt := time.Now()
	err = tx.Model(&trialModel{}).
		Where("trial_id = ?", trialID).
		Update("value", value).
		Update("state", xr).
		Update("datetime_complete", completedAt).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// SetTrialUserAttr to store the value for the user.
func (s *Storage) SetTrialUserAttr(trialID int, key string, value string) error {
	var result trialUserAttributeModel
	return s.db.Where(&trialUserAttributeModel{
		UserAttributeReferTrial: trialID,
		Key:                     key,
	}).Assign(&trialUserAttributeModel{
		UserAttributeReferTrial: trialID,
		Key:                     key,
		ValueJSON:               encodeAttrValue(value),
	}).FirstOrCreate(&result).Error
}

// SetTrialSystemAttr to store the value for the system.
func (s *Storage) SetTrialSystemAttr(trialID int, key string, value string) error {
	var result trialSystemAttributeModel
	return s.db.Where(&trialSystemAttributeModel{
		SystemAttributeReferTrial: trialID,
		Key:                       key,
	}).Assign(&trialSystemAttributeModel{
		SystemAttributeReferTrial: trialID,
		Key:                       key,
		ValueJSON:                 encodeAttrValue(value),
	}).FirstOrCreate(&result).Error
}

// GetTrialNumberFromID returns the trial's number.
func (s *Storage) GetTrialNumberFromID(trialID int) (int, error) {
	trial, err := s.GetTrial(trialID)
	if err != nil {
		return -1, err
	}
	return trial.Number, err
}

// GetTrialParam returns the internal parameter of the trial
func (s *Storage) GetTrialParam(trialID int, paramName string) (float64, error) {
	var param trialParamModel
	err := s.db.First(&param, "trial_id = ? AND param_name = ?", trialID, paramName).Error
	if err != nil {
		return -1, err
	}
	return param.Value, nil
}

// GetTrialParams returns the external parameters in the trial
func (s *Storage) GetTrialParams(trialID int) (map[string]interface{}, error) {
	trial, err := s.GetTrial(trialID)
	if err != nil {
		return nil, err
	}
	return trial.Params, nil
}

// GetTrialUserAttrs to restore the attributes for the user.
func (s *Storage) GetTrialUserAttrs(trialID int) (map[string]string, error) {
	var attrs []trialUserAttributeModel
	result := s.db.Find(&attrs, "trial_id = ?", trialID)
	if result.Error != nil {
		return nil, result.Error
	}

	res := make(map[string]string, len(attrs))
	for i := range attrs {
		res[attrs[i].Key] = decodeAttrValue(attrs[i].ValueJSON)
	}
	return res, nil
}

// GetTrialSystemAttrs to restore the attributes for the system.
func (s *Storage) GetTrialSystemAttrs(trialID int) (map[string]string, error) {
	var attrs []trialSystemAttributeModel
	result := s.db.Find(&attrs, "trial_id = ?", trialID)
	if result.Error != nil {
		return nil, result.Error
	}

	res := make(map[string]string, len(attrs))
	for i := range attrs {
		res[attrs[i].Key] = decodeAttrValue(attrs[i].ValueJSON)
	}
	return res, nil
}

// GetBestTrial returns the best trial.
func (s *Storage) GetBestTrial(studyID int) (goptuna.FrozenTrial, error) {
	var err error
	var study studyModel
	err = s.db.
		Preload("Trials").
		First(&study, "study_id = ?", studyID).Error
	if err != nil {
		return goptuna.FrozenTrial{}, err
	}

	if len(study.Trials) == 0 {
		return goptuna.FrozenTrial{}, nil
	}

	var best *trialModel
	for i := range study.Trials {
		if study.Trials[i].State != trialStateComplete {
			continue
		}

		if best == nil {
			best = &study.Trials[i]
		}
		if study.Direction == directionMaximize {
			if best.Value < study.Trials[i].Value {
				best = &study.Trials[i]
			}
		} else {
			if best.Value > study.Trials[i].Value {
				best = &study.Trials[i]
			}
		}
	}
	var ft goptuna.FrozenTrial
	if best != nil {
		ft, err = s.GetTrial(best.ID)
		if err != nil {
			return goptuna.FrozenTrial{}, err
		}
	}
	return ft, err
}

// GetAllTrials returns the all trials.
func (s *Storage) GetAllTrials(studyID int) ([]goptuna.FrozenTrial, error) {
	var trials []trialModel
	var params []trialParamModel
	var values []trialValueModel
	var userAttrs []trialUserAttributeModel
	var systemAttrs []trialSystemAttributeModel
	if err := s.db.
		Order("trial_id").
		Find(&trials, "study_id = ?", studyID).
		Error; err != nil {
		return nil, err
	}
	if err := s.db.
		Joins("JOIN trials ON trial_params.trial_id = trials.trial_id").
		Find(&params, "trials.study_id = ?", studyID).Error; err != nil {
		return nil, err
	}
	if err := s.db.
		Joins("JOIN trials ON trial_values.trial_id = trials.trial_id").
		Find(&values, "trials.study_id = ?", studyID).Error; err != nil {
		return nil, err
	}
	if err := s.db.
		Joins("JOIN trials ON trial_user_attributes.trial_id = trials.trial_id").
		Find(&userAttrs, "trials.study_id = ?", studyID).Error; err != nil {
		return nil, err
	}
	if err := s.db.
		Joins("JOIN trials ON trial_system_attributes.trial_id = trials.trial_id").
		Find(&systemAttrs, "trials.study_id = ?", studyID).Error; err != nil {
		return nil, err
	}

	// Following SQL might raise 'too many SQL variables' error.
	// See https://github.com/c-bata/goptuna/issues/30 for more details.
	// err := s.db.
	// 	Where("study_id = ?", studyID).
	// 	Preload("UserAttributes").
	// 	Preload("SystemAttributes").
	// 	Preload("TrialParams").
	// 	Preload("TrialValues").
	// 	Find(&trials).Error

	return s.mergeTrialsORM(trials, params, values, userAttrs, systemAttrs)
}

func (s *Storage) mergeTrialsORM(
	trials []trialModel,
	params []trialParamModel,
	values []trialValueModel,
	userAttrs []trialUserAttributeModel,
	systemAttrs []trialSystemAttributeModel,
) ([]goptuna.FrozenTrial, error) {
	idToTrials := make(map[int]trialModel, len(trials))
	for i := range trials {
		idToTrials[trials[i].ID] = trials[i]
	}

	defaultSize := 3
	idToParams := make(map[int][]trialParamModel, len(trials))
	for i := range params {
		trialID := params[i].TrialParamReferTrial
		l, ok := idToParams[trialID]
		if !ok {
			idToParams[trialID] = make([]trialParamModel, 0, defaultSize)
		}
		idToParams[trialID] = append(l, params[i])
	}
	idToValues := make(map[int][]trialValueModel, len(trials))
	for i := range values {
		trialID := values[i].TrialValueReferTrial
		l, ok := idToValues[trialID]
		if !ok {
			idToValues[trialID] = make([]trialValueModel, 0, defaultSize)
		}
		idToValues[trialID] = append(l, values[i])
	}
	idToUserAttrs := make(map[int][]trialUserAttributeModel, len(trials))
	for i := range userAttrs {
		trialID := userAttrs[i].UserAttributeReferTrial
		l, ok := idToUserAttrs[trialID]
		if !ok {
			idToUserAttrs[trialID] = make([]trialUserAttributeModel, 0, defaultSize)
		}
		idToUserAttrs[trialID] = append(l, userAttrs[i])
	}
	idToSystemAttrs := make(map[int][]trialSystemAttributeModel, len(trials))
	for i := range systemAttrs {
		trialID := systemAttrs[i].SystemAttributeReferTrial
		l, ok := idToSystemAttrs[trialID]
		if !ok {
			idToSystemAttrs[trialID] = make([]trialSystemAttributeModel, 0, defaultSize)
		}
		idToSystemAttrs[trialID] = append(l, systemAttrs[i])
	}

	merged := make([]goptuna.FrozenTrial, 0, len(trials))
	for i := range trials {
		if v, ok := idToParams[trials[i].ID]; ok {
			trials[i].TrialParams = v
		}
		if v, ok := idToValues[trials[i].ID]; ok {
			trials[i].TrialValues = v
		}
		if v, ok := idToUserAttrs[trials[i].ID]; ok {
			trials[i].UserAttributes = v
		}
		if v, ok := idToSystemAttrs[trials[i].ID]; ok {
			trials[i].SystemAttributes = v
		}
		frozen, err := toFrozenTrial(trials[i])
		if err != nil {
			return nil, err
		}
		merged = append(merged, frozen)
	}
	return merged, nil
}

// GetStudyDirection returns study direction of the objective.
func (s *Storage) GetStudyDirection(studyID int) (goptuna.StudyDirection, error) {
	var study studyModel
	err := s.db.First(&study, "study_id = ?", studyID).Error
	if err != nil {
		return goptuna.StudyDirectionMinimize, err
	}
	return toGoptunaStudyDirection(study.Direction), nil
}

// GetTrial returns Trial.
func (s *Storage) GetTrial(trialID int) (goptuna.FrozenTrial, error) {
	var trial trialModel
	err := s.db.
		Preload("UserAttributes").
		Preload("SystemAttributes").
		Preload("TrialParams").
		Preload("TrialValues").
		First(&trial, "trial_id = ?", trialID).Error
	if err != nil {
		return goptuna.FrozenTrial{}, err
	}
	ft, err := toFrozenTrial(trial)
	return ft, err
}
