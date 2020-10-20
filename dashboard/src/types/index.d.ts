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

declare interface IntermediateValue {
  step: number
  value: number
}

declare interface Param {
  name: string
  value: string
}

declare interface Attribute {
  key: string
  value: string
}

declare interface FrozenTrial {
  trial_id: number
  study_id: number
  number: number
  state: TrialState
  value?: number
  intermediate_value: IntermediateValue[]
  datetime_start: Date
  datetime_complete?: Date
  params: Param[]
  user_attrs: Attribute[]
  system_attrs: Attribute[]
}

declare interface StudySummary {
  study_id: number
  study_name: string
  direction: StudyDirection
  best_trial?: FrozenTrial
  user_attrs: Attribute[]
  system_attrs: Attribute[]
  datetime_start: Date
}
