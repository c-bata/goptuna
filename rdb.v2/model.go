package rdb

import (
	"time"

	"gorm.io/gorm"
)

const (
	directionMinimize = "MINIMIZE"
	directionMaximize = "MAXIMIZE"
)

const (
	trialStateRunning  = "RUNNING"
	trialStateComplete = "COMPLETE"
	trialStatePruned   = "PRUNED"
	trialStateFail     = "FAIL"
	trialStateWaiting  = "WAITING"
)

// https://gorm.io/docs/models.html

type studyModel struct {
	ID        int    `gorm:"column:study_id;primaryKey"`
	Name      string `gorm:"column:study_name;type:varchar(512);unique;not null"`
	Direction string `gorm:"column:direction;not null"`

	// Associations
	UserAttributes   []studyUserAttributeModel   `gorm:"foreignkey:UserAttributeReferStudy;association_foreignkey:ID"`
	SystemAttributes []studySystemAttributeModel `gorm:"foreignkey:SystemAttributeReferStudy;association_foreignkey:ID"`
	Trials           []trialModel                `gorm:"foreignkey:TrialReferStudy;association_foreignkey:ID"`
}

func (m studyModel) TableName() string {
	return "studies"
}

type studyUserAttributeModel struct {
	ID                      int    `gorm:"column:study_user_attribute_id;primaryKey"`
	UserAttributeReferStudy int    `gorm:"column:study_id;uniqueIndex:idx_study_user_attr_key"`
	Key                     string `gorm:"column:key;uniqueIndex:idx_study_user_attr_key;type:varchar(512)"`
	Value                   string `gorm:"column:value;type:varchar(2048)"`
}

func (m studyUserAttributeModel) TableName() string {
	return "study_user_attributes"
}

type studySystemAttributeModel struct {
	ID                        int    `gorm:"column:study_system_attribute_id;primaryKey"`
	SystemAttributeReferStudy int    `gorm:"column:study_id;uniqueIndex:idx_study_system_attr_key"`
	Key                       string `gorm:"column:key;uniqueIndex:idx_study_system_attr_key;type:varchar(512)"`
	Value                     string `gorm:"column:value;type:varchar(2048)"`
}

func (m studySystemAttributeModel) TableName() string {
	return "study_system_attributes"
}

type trialModel struct {
	ID               int        `gorm:"column:trial_id;primaryKey"`
	Number           int        `gorm:"column:number"`
	TrialReferStudy  int        `gorm:"column:study_id"`
	State            string     `gorm:"column:state;not null"`
	Value            float64    `gorm:"column:value"`
	DatetimeStart    *time.Time `gorm:"column:datetime_start"`
	DatetimeComplete *time.Time `gorm:"column:datetime_complete"`

	// Associations
	UserAttributes   []trialUserAttributeModel   `gorm:"foreignkey:UserAttributeReferTrial;association_foreignkey:ID"`
	SystemAttributes []trialSystemAttributeModel `gorm:"foreignkey:SystemAttributeReferTrial;association_foreignkey:ID"`
	TrialParams      []trialParamModel           `gorm:"foreignkey:TrialParamReferTrial;association_foreignkey:ID"`
	TrialValues      []trialValueModel           `gorm:"foreignkey:TrialValueReferTrial;association_foreignkey:ID"`
}

func (m trialModel) TableName() string {
	return "trials"
}

type trialUserAttributeModel struct {
	ID                      int    `gorm:"column:trial_user_attribute_id;primaryKey"`
	UserAttributeReferTrial int    `gorm:"column:trial_id;uniqueIndex:idx_trial_user_attr_key"`
	Key                     string `gorm:"column:key;uniqueIndex:idx_trial_user_attr_key;type:varchar(512)"`
	Value                   string `gorm:"column:value;type:varchar(2048)"`
}

func (m trialUserAttributeModel) TableName() string {
	return "trial_user_attributes"
}

type trialSystemAttributeModel struct {
	ID                        int    `gorm:"column:trial_system_attribute_id;primaryKey"`
	SystemAttributeReferTrial int    `gorm:"column:trial_id;uniqueIndex:idx_trial_system_attr_key"`
	Key                       string `gorm:"column:key;uniqueIndex:idx_trial_system_attr_key;type:varchar(512)"`
	Value                     string `gorm:"column:value;type:varchar(2048)"`
}

func (m trialSystemAttributeModel) TableName() string {
	return "trial_system_attributes"
}

type trialParamModel struct {
	ID                   int     `gorm:"column:param_id;primaryKey"`
	TrialParamReferTrial int     `gorm:"column:trial_id;uniqueIndex:idx_trial_param_name"`
	Name                 string  `gorm:"column:param_name;uniqueIndex:idx_trial_param_name"`
	Value                float64 `gorm:"column:param_value"`
	DistributionJSON     string  `gorm:"column:distribution_json;type:varchar(2048)"`
}

func (m trialParamModel) TableName() string {
	return "trial_params"
}

type trialValueModel struct {
	ID                   int     `gorm:"column:trial_value_id;primaryKey"`
	TrialValueReferTrial int     `gorm:"column:trial_id;uniqueIndex:idx_trial_value_step"`
	Step                 int     `gorm:"column:step;uniqueIndex:idx_trial_value_step"`
	Value                float64 `gorm:"column:value"`
}

func (m trialValueModel) TableName() string {
	return "trial_values"
}

// RunAutoMigrate runs Auto-Migration. This will ONLY create tables,
// missing columns and missing indexes, and WON’T change existing
// column’s type or delete unused columns to protect your data.
func RunAutoMigrate(db *gorm.DB) error {
	var err error
	err = db.AutoMigrate(&studyModel{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&studyUserAttributeModel{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&studySystemAttributeModel{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&trialModel{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&trialUserAttributeModel{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&trialSystemAttributeModel{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&trialParamModel{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&trialValueModel{})
	if err != nil {
		return err
	}
	return nil
}
