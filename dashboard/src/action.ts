import axios from "axios"

const axiosInstance = axios.create({ baseURL: API_ENDPOINT })

declare interface FrozenTrialResponse {
  trial_id: number
  study_id: number
  number: number
  state: TrialState
  value?: number
  intermediate_value: IntermediateValue[]
  datetime_start: string
  datetime_complete?: string
  params: Param[]
  user_attrs: Attribute[]
  system_attrs: Attribute[]
}

declare interface TrialsResponse {
  trials: FrozenTrialResponse[]
}

export const fetchTrialsAction = (studyId: number) => {
  axiosInstance
    .get<TrialsResponse>(`/api/studies/${studyId}/trials`, {})
    .then((res) => {
      let trials: FrozenTrial[] = res.data.trials.map(
        (trial: FrozenTrialResponse): FrozenTrial => {
          return {
            trial_id: trial.trial_id,
            study_id: trial.study_id,
            number: trial.number,
            state: trial.state,
            value: trial.value,
            intermediate_value: trial.intermediate_value,
            datetime_start: new Date(trial.datetime_start),
            datetime_complete: trial.datetime_complete
              ? new Date(trial.datetime_complete)
              : undefined,
            params: trial.params,
            user_attrs: trial.user_attrs,
            system_attrs: trial.system_attrs,
          }
        }
      )
      // Notify to dispatchers
      console.log("DEBUG")
      console.log(trials)
    })
    .catch((err) => {
      // Notify to error dispatchers
      console.log(err)
    })
}

declare interface StudySummaryResponse {
  study_id: number
  study_name: string
  direction: StudyDirection
  best_trial?: FrozenTrialResponse
  user_attrs: Attribute[]
  system_attrs: Attribute[]
  datetime_start: string
}

declare interface StudySummariesResponse {
  study_summaries: StudySummaryResponse[]
}

export const fetchStudySummariesAction = () => {
  axiosInstance
    .get<StudySummariesResponse>(`/api/studies`, {})
    .then((res) => {
      let studySummaries: StudySummary[] = res.data.study_summaries.map(
        (study: StudySummaryResponse): StudySummary => {
          const best_trial = study.best_trial
            ? {
                trial_id: study.best_trial.trial_id,
                study_id: study.best_trial.study_id,
                number: study.best_trial.number,
                state: study.best_trial.state,
                value: study.best_trial.value,
                intermediate_value: study.best_trial.intermediate_value,
                datetime_start: new Date(study.best_trial.datetime_start),
                datetime_complete: study.best_trial.datetime_complete
                  ? new Date(study.best_trial.datetime_complete)
                  : undefined,
                params: study.best_trial.params,
                user_attrs: study.best_trial.user_attrs,
                system_attrs: study.best_trial.system_attrs,
              }
            : undefined
          return {
            study_id: study.study_id,
            study_name: study.study_name,
            direction: study.direction,
            best_trial: best_trial,
            user_attrs: study.user_attrs,
            system_attrs: study.system_attrs,
            datetime_start: new Date(study.datetime_start),
          }
        }
      )
      // Notify to dispatchers
      console.log("DEBUG")
      console.log(studySummaries)
    })
    .catch((err) => {
      // Notify to error dispatchers
      console.log(err)
    })
}
