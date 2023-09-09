declare module "*.css"
declare module "*.png"
declare module "*.jpg"
declare module "*.svg"

declare const APP_BAR_TITLE: string
declare const API_ENDPOINT: string
declare const URL_PREFIX: string

type TrialState = "Running" | "Complete" | "Pruned" | "Fail" | "Waiting"
type TrialStateFinished = "Complete" | "Fail" | "Pruned"
type StudyDirection = "maximize" | "minimize"

type FloatDistribution = {
  type: "FloatDistribution"
  low: number
  high: number
  step: number
  log: boolean
}

type IntDistribution = {
  type: "IntDistribution"
  low: number
  high: number
  step: number
  log: boolean
}

type CategoricalDistribution = {
  type: "CategoricalDistribution"
  choices: { pytype: string; value: string }[]
}

type Distribution =
  | FloatDistribution
  | IntDistribution
  | CategoricalDistribution

type TrialIntermediateValue = {
  step: number
  value: number
}

type TrialParam = {
  name: string
  param_internal_value: number
  param_external_value: string
  param_external_type: string
  // distribution: Distribution
}

type ParamImportance = {
  name: string
  importance: number
  distribution: Distribution
}

type SearchSpaceItem = {
  name: string
  distribution: Distribution
}

type Attribute = {
  key: string
  value: string
}

type AttributeSpec = {
  key: string
  sortable: boolean
}

type Trial = {
  trial_id: number
  study_id: number
  number: number
  state: TrialState
  value?: number
  intermediate_values: TrialIntermediateValue[]
  datetime_start?: Date
  datetime_complete?: Date
  params: TrialParam[]
  fixed_params: {
    name: string
    param_external_value: string
  }[]
  user_attrs: Attribute[]
}

type StudySummary = {
  study_id: number
  study_name: string
  direction: StudyDirection
  user_attrs: Attribute[]
  datetime_start?: Date
}

type StudyDetail = {
  id: number
  name: string
  direction: StudyDirection
  user_attrs: Attribute[]
  datetime_start: Date
  best_trial: Trial
  trials: Trial[]
  intersection_search_space: SearchSpaceItem[]
  union_search_space: SearchSpaceItem[]
  union_user_attrs: AttributeSpec[]
  has_intermediate_values: boolean
}

type StudyDetails = {
  [study_id: string]: StudyDetail
}
