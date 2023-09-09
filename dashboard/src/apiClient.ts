import axios from "axios"

const axiosInstance = axios.create({ baseURL: API_ENDPOINT })

interface TrialResponse {
  trial_id: number
  study_id: number
  number: number
  state: TrialState
  value: number
  intermediate_values: TrialIntermediateValue[]
  datetime_start?: string
  datetime_complete?: string
  params: TrialParam[]
  fixed_params: {
    name: string
    param_external_value: string
  }[]
  user_attrs: Attribute[]
}

const convertTrialResponse = (res: TrialResponse): Trial => {
  return {
    trial_id: res.trial_id,
    study_id: res.study_id,
    number: res.number,
    state: res.state,
    value: res.value,
    intermediate_values: res.intermediate_values,
    datetime_start: res.datetime_start
      ? new Date(res.datetime_start)
      : undefined,
    datetime_complete: res.datetime_complete
      ? new Date(res.datetime_complete)
      : undefined,
    params: res.params,
    fixed_params: res.fixed_params,
    user_attrs: res.user_attrs,
  }
}

interface StudyDetailResponse {
  name: string
  datetime_start: string
  direction: StudyDirection
  user_attrs: Attribute[]
  trials: TrialResponse[]
  best_trial: TrialResponse
  intersection_search_space: SearchSpaceItem[]
  union_search_space: SearchSpaceItem[]
  union_user_attrs: AttributeSpec[]
  has_intermediate_values: boolean
}

export const getStudyDetailAPI = (studyId: number): Promise<StudyDetail> => {
  return axiosInstance
    .get<StudyDetailResponse>(`/api/studies/${studyId}`, {})
    .then((res) => {
      const trials = res.data.trials.map((trial): Trial => {
        return convertTrialResponse(trial)
      })
      console.dir(res.data)
      return {
        id: studyId,
        name: res.data.name,
        datetime_start: new Date(res.data.datetime_start),
        direction: res.data.direction,
        user_attrs: res.data.user_attrs,
        trials: trials,
        best_trial: convertTrialResponse(res.data.best_trial),
        union_search_space: res.data.union_search_space,
        intersection_search_space: res.data.intersection_search_space,
        union_user_attrs: res.data.union_user_attrs,
        has_intermediate_values: res.data.has_intermediate_values,
      }
    })
}

interface StudySummariesResponse {
  study_summaries: {
    study_id: number
    study_name: string
    direction: StudyDirection
    user_attrs: Attribute[]
    datetime_start?: string
  }[]
}

export const getStudySummariesAPI = (): Promise<StudySummary[]> => {
  return axiosInstance
    .get<StudySummariesResponse>(`/api/studies`, {})
    .then((res) => {
      return res.data.study_summaries.map((study): StudySummary => {
        return {
          study_id: study.study_id,
          study_name: study.study_name,
          direction: study.direction,
          user_attrs: study.user_attrs,
          datetime_start: study.datetime_start
            ? new Date(study.datetime_start)
            : undefined,
        }
      })
    })
}

interface CreateNewStudyResponse {
  study_summary: {
    study_id: number
    study_name: string
    direction: StudyDirection
    user_attrs: Attribute[]
    is_preferential: boolean
    datetime_start?: string
  }
}

export const createNewStudyAPI = (
  studyName: string,
  direction: StudyDirection
): Promise<StudySummary> => {
  return axiosInstance
    .post<CreateNewStudyResponse>(`/api/studies`, {
      study_name: studyName,
      direction,
    })
    .then((res) => {
      const study_summary = res.data.study_summary
      return {
        study_id: study_summary.study_id,
        study_name: study_summary.study_name,
        direction: study_summary.direction,
        user_attrs: study_summary.user_attrs,
        is_preferential: study_summary.is_preferential,
        datetime_start: study_summary.datetime_start
          ? new Date(study_summary.datetime_start)
          : undefined,
      }
    })
}

export const deleteStudyAPI = (studyId: number): Promise<void> => {
  return axiosInstance.delete(`/api/studies/${studyId}`).then(() => {
    return
  })
}

export const tellTrialAPI = (
  trialId: number,
  state: TrialStateFinished,
  value: number
): Promise<void> => {
  const req: { state: TrialState; value: number } = {
    state: state,
    value: value,
  }

  return axiosInstance
    .post<void>(`/api/trials/${trialId}/tell`, req)
    .then(() => {
      return
    })
}

export const saveTrialUserAttrsAPI = (
  trialId: number,
  user_attrs: { [key: string]: number | string }
): Promise<void> => {
  const req = { user_attrs: user_attrs }

  return axiosInstance
    .post<void>(`/api/trials/${trialId}/user-attrs`, req)
    .then(() => {
      return
    })
}

interface ParamImportancesResponse {
  param_importances: ParamImportance[][]
}

export const getParamImportances = (
  studyId: number
): Promise<ParamImportance[][]> => {
  return axiosInstance
    .get<ParamImportancesResponse>(`/api/studies/${studyId}/param_importances`)
    .then((res) => {
      return res.data.param_importances
    })
}

export const reportPreferenceAPI = (
  studyId: number,
  candidates: number[],
  clicked: number
): Promise<void> => {
  return axiosInstance
    .post<void>(`/api/studies/${studyId}/preference`, {
      candidates: candidates,
      clicked: clicked,
      mode: "ChooseWorst",
    })
    .then(() => {
      return
    })
}

export const skipPreferentialTrialAPI = (
  studyId: number,
  trialId: number
): Promise<void> => {
  return axiosInstance
    .post<void>(`/api/studies/${studyId}/${trialId}/skip`)
    .then(() => {
      return
    })
}
