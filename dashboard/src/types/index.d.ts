declare module "*.css"
declare module "*.png"
declare module "*.jpg"
declare module "*.svg"

declare const API_ENDPOINT: string

declare enum TrialState {
  Running = "Running",
  Complete = "Complete",
  Pruned = "Pruned",
  Fail = "Fail",
  Waiting = "Waiting",
}

declare enum StudyDirection {
  Maximize = "maximize",
  Minimize = "minimize",
}

declare interface TrialIntermediateValue {
  step: number
  value: number
}

declare interface TrialParam {
  name: string
  value: string
}

declare interface Attribute {
  key: string
  value: string
}

declare interface Trial {
  trial_id: number
  study_id: number
  number: number
  state: TrialState
  value?: number
  intermediate_value: TrialIntermediateValue[]
  datetime_start: Date
  datetime_complete?: Date
  params: TrialParam[]
  user_attrs: Attribute[]
  system_attrs: Attribute[]
}

declare interface StudySummary {
  study_id: number
  study_name: string
  direction: StudyDirection
  best_trial?: Trial
  user_attrs: Attribute[]
  system_attrs: Attribute[]
  datetime_start: Date
}

declare interface StudyDetail {
  trials: Trial[]
}

declare interface StudyDetails {
  [study_id: string]: StudyDetail
}
