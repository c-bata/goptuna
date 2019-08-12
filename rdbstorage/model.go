package rdbstorage

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	DirectionNotSet   = 0
	DirectionMINIMIZE = 1
	DirectionMAXIMIZE = 2
)

const (
	TrialStateRunning  = 0
	TrialStateComplete = 1
	TrialStatePruned   = 2
	TrialStateFail     = 3
)

// https://gorm.io/docs/models.html

type StudyModel struct {
	ID        int    `gorm:"column:study_id;PRIMARY_KEY"`
	Name      string `gorm:"column:study_name;type:varchar(512);unique_index;NOT NULL"`
	Direction int    `gorm:"column:direction;NOT NULL"`

	// Associations
	UserAttributes   []StudyUserAttributeModel   `gorm:"foreignkey:UserAttributeReferStudy;association_foreignkey:ID"`
	SystemAttributes []StudySystemAttributeModel `gorm:"foreignkey:SystemAttributeReferStudy;association_foreignkey:ID"`
	Trials           []TrialModel                `gorm:"foreignkey:TrialReferStudy;association_foreignkey:ID"`
}

func (m StudyModel) TableName() string {
	return "studies"
}

type StudyUserAttributeModel struct {
	ID                      int    `gorm:"column:study_user_attribute_id;PRIMARY_KEY"`
	UserAttributeReferStudy int    `gorm:"column:study_id;unique_index:idx_study_user_attr_key"`
	Key                     string `gorm:"column:key;unique_index:idx_study_user_attr_key;type:varchar(512)"`
	ValueJSON               string `gorm:"column:value_json;type:varchar(2048)"`
}

func (m StudyUserAttributeModel) TableName() string {
	return "study_user_attributes"
}

type StudySystemAttributeModel struct {
	ID                        int    `gorm:"column:study_system_attribute_id;PRIMARY_KEY"`
	SystemAttributeReferStudy int    `gorm:"column:study_id;unique_index:idx_study_system_attr_key"`
	Key                       string `gorm:"column:key;unique_index:idx_study_system_attr_key;type:varchar(512)"`
	ValueJSON                 string `gorm:"column:value_json;type:varchar(2048)"`
}

func (m StudySystemAttributeModel) TableName() string {
	return "study_system_attributes"
}

type TrialModel struct {
	ID               int        `gorm:"column:trial_id;PRIMARY_KEY"`
	TrialReferStudy  int        `gorm:"column:study_id"`
	State            int        `gorm:"column:state;NOT NULL"`
	Value            float64    `gorm:"column:value"`
	DatetimeStart    *time.Time `gorm:"column:datetime_start"`
	DatetimeComplete *time.Time `gorm:"column:datetime_complete"`

	// Associations
	UserAttributes   []TrialUserAttributeModel   `gorm:"foreignkey:UserAttributeReferTrial;association_foreignkey:ID"`
	SystemAttributes []TrialSystemAttributeModel `gorm:"foreignkey:SystemAttributeReferTrial;association_foreignkey:ID"`
	TrialParams      []TrialParamModel           `gorm:"foreignkey:TrialParamReferTrial;association_foreignkey:ID"`
	TrialValues      []TrialValueModel           `gorm:"foreignkey:TrialValueReferTrial;association_foreignkey:ID"`
}

func (m TrialModel) TableName() string {
	return "trials"
}

type TrialUserAttributeModel struct {
	ID                      int    `gorm:"column:trial_user_attribute_id;PRIMARY_KEY"`
	UserAttributeReferTrial int    `gorm:"column:trial_id;unique_index:idx_trial_user_attr_key"`
	Key                     string `gorm:"column:key;unique_index:idx_trial_user_attr_key;type:varchar(512)"`
	ValueJSON               string `gorm:"column:value_json;type:varchar(2048)"`
}

func (m TrialUserAttributeModel) TableName() string {
	return "trial_user_attributes"
}

type TrialSystemAttributeModel struct {
	ID                        int    `gorm:"column:trial_system_attribute_id;PRIMARY_KEY"`
	SystemAttributeReferTrial int    `gorm:"column:trial_id;unique_index:idx_trial_system_attr_key"`
	Key                       string `gorm:"column:key;unique_index:idx_trial_system_attr_key;type:varchar(512)"`
	ValueJSON                 string `gorm:"column:value_json;type:varchar(2048)"`
}

func (m TrialSystemAttributeModel) TableName() string {
	return "trial_system_attributes"
}

type TrialParamModel struct {
	ID                   int     `gorm:"column:param_id;PRIMARY_KEY"`
	TrialParamReferTrial int     `gorm:"column:trial_id;unique_index:idx_trial_param_name"`
	Name                 string  `gorm:"column:param_name;unique_index:idx_trial_param_name"`
	Value                float64 `gorm:"column:param_value"`
	DistributionJSON     string  `gorm:"column:distribution_json;type:varchar(2048)"`
}

func (m TrialParamModel) TableName() string {
	return "trial_params"
}

type TrialValueModel struct {
	ID                   int     `gorm:"column:trial_value_id;PRIMARY_KEY"`
	TrialValueReferTrial int     `gorm:"column:trial_id;unique_index:idx_trial_value_step"`
	Step                 int     `gorm:"column:step;unique_index:idx_trial_value_step"`
	Value                float64 `gorm:"column:value"`
}

func (m TrialValueModel) TableName() string {
	return "trial_values"
}

// RunAutoMigrate runs Auto-Migration. This will ONLY create tables,
// missing columns and missing indexes, and WON’T change existing
// column’s type or delete unused columns to protect your data.
func RunAutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&StudyModel{})
	db.AutoMigrate(&StudyUserAttributeModel{})
	db.AutoMigrate(&StudySystemAttributeModel{})
	db.AutoMigrate(&TrialModel{})
	db.AutoMigrate(&TrialUserAttributeModel{})
	db.AutoMigrate(&TrialSystemAttributeModel{})
	db.AutoMigrate(&TrialParamModel{})
	db.AutoMigrate(&TrialValueModel{})
}
