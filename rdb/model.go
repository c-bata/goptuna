package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	directionNotSet   = "NOT_SET"
	directionMinimize = "MINIMIZE"
	directionMaximize = "MAXIMIZE"
)

const (
	trialStateRunning  = "RUNNING"
	trialStateComplete = "COMPLETE"
	trialStatePruned   = "PRUNED"
	trialStateFail     = "FAIL"
)

// https://gorm.io/docs/models.html

type studyModel struct {
	ID        int    `gorm:"column:study_id;PRIMARY_KEY"`
	Name      string `gorm:"column:study_name;type:varchar(512);unique_index;NOT NULL"`
	Direction string `gorm:"column:direction;NOT NULL"`

	// Associations
	UserAttributes   []studyUserAttributeModel   `gorm:"foreignkey:UserAttributeReferStudy;association_foreignkey:ID"`
	SystemAttributes []studySystemAttributeModel `gorm:"foreignkey:SystemAttributeReferStudy;association_foreignkey:ID"`
	Trials           []trialModel                `gorm:"foreignkey:TrialReferStudy;association_foreignkey:ID"`
}

func (m studyModel) TableName() string {
	return "studies"
}

type studyUserAttributeModel struct {
	ID                      int    `gorm:"column:study_user_attribute_id;PRIMARY_KEY"`
	UserAttributeReferStudy int    `gorm:"column:study_id;unique_index:idx_study_user_attr_key"`
	Key                     string `gorm:"column:key;unique_index:idx_study_user_attr_key;type:varchar(512)"`
	ValueJSON               string `gorm:"column:value_json;type:varchar(2048)"`
}

func (m studyUserAttributeModel) TableName() string {
	return "study_user_attributes"
}

type studySystemAttributeModel struct {
	ID                        int    `gorm:"column:study_system_attribute_id;PRIMARY_KEY"`
	SystemAttributeReferStudy int    `gorm:"column:study_id;unique_index:idx_study_system_attr_key"`
	Key                       string `gorm:"column:key;unique_index:idx_study_system_attr_key;type:varchar(512)"`
	ValueJSON                 string `gorm:"column:value_json;type:varchar(2048)"`
}

func (m studySystemAttributeModel) TableName() string {
	return "study_system_attributes"
}

type trialModel struct {
	ID               int        `gorm:"column:trial_id;PRIMARY_KEY"`
	TrialReferStudy  int        `gorm:"column:study_id"`
	State            string     `gorm:"column:state;NOT NULL"`
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
	ID                      int    `gorm:"column:trial_user_attribute_id;PRIMARY_KEY"`
	UserAttributeReferTrial int    `gorm:"column:trial_id;unique_index:idx_trial_user_attr_key"`
	Key                     string `gorm:"column:key;unique_index:idx_trial_user_attr_key;type:varchar(512)"`
	ValueJSON               string `gorm:"column:value_json;type:varchar(2048)"`
}

func (m trialUserAttributeModel) TableName() string {
	return "trial_user_attributes"
}

type trialSystemAttributeModel struct {
	ID                        int    `gorm:"column:trial_system_attribute_id;PRIMARY_KEY"`
	SystemAttributeReferTrial int    `gorm:"column:trial_id;unique_index:idx_trial_system_attr_key"`
	Key                       string `gorm:"column:key;unique_index:idx_trial_system_attr_key;type:varchar(512)"`
	ValueJSON                 string `gorm:"column:value_json;type:varchar(2048)"`
}

func (m trialSystemAttributeModel) TableName() string {
	return "trial_system_attributes"
}

type trialParamModel struct {
	ID                   int     `gorm:"column:param_id;PRIMARY_KEY"`
	TrialParamReferTrial int     `gorm:"column:trial_id;unique_index:idx_trial_param_name"`
	Name                 string  `gorm:"column:param_name;unique_index:idx_trial_param_name"`
	Value                float64 `gorm:"column:param_value"`
	DistributionJSON     string  `gorm:"column:distribution_json;type:varchar(2048)"`
}

func (m trialParamModel) TableName() string {
	return "trial_params"
}

type trialValueModel struct {
	ID                   int     `gorm:"column:trial_value_id;PRIMARY_KEY"`
	TrialValueReferTrial int     `gorm:"column:trial_id;unique_index:idx_trial_value_step"`
	Step                 int     `gorm:"column:step;unique_index:idx_trial_value_step"`
	Value                float64 `gorm:"column:value"`
}

func (m trialValueModel) TableName() string {
	return "trial_values"
}

// RunAutoMigrate runs Auto-Migration. This will ONLY create tables,
// missing columns and missing indexes, and WON’T change existing
// column’s type or delete unused columns to protect your data.
func RunAutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&studyModel{})
	db.AutoMigrate(&studyUserAttributeModel{})
	db.AutoMigrate(&studySystemAttributeModel{})
	db.AutoMigrate(&trialModel{})
	db.AutoMigrate(&trialUserAttributeModel{})
	db.AutoMigrate(&trialSystemAttributeModel{})
	db.AutoMigrate(&trialParamModel{})
	db.AutoMigrate(&trialValueModel{})
}
