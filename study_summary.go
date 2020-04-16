package goptuna

import "time"

// StudySummary holds basic attributes and aggregated results of Study.
type StudySummary struct {
	ID            int               `json:"study_id"`
	Name          string            `json:"study_name"`
	Direction     StudyDirection    `json:"direction"`
	BestTrial     FrozenTrial       `json:"best_trial"`
	UserAttrs     map[string]string `json:"user_attrs"`
	SystemAttrs   map[string]string `json:"system_attrs"`
	DatetimeStart time.Time         `json:"datetime_start"`
}
