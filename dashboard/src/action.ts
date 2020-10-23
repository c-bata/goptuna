import { SetterOrUpdater } from "recoil"
import {
  getStudyDetailAPI,
  getStudySummariesAPI,
  createNewStudyAPI,
} from "./apiClient"
import { ReactNode } from "react"
import { OptionsObject } from "notistack"

export const actionCreator = (
  enqueueSnackbar: (
    message: ReactNode,
    options?: OptionsObject | undefined
  ) => string | number
) => {
  const updateStudySummaries = (setter: SetterOrUpdater<StudySummary[]>) => {
    getStudySummariesAPI()
      .then((studySummaries: StudySummary[]) => {
        setter(studySummaries)
      })
      .catch((err) => {
        enqueueSnackbar(`Failed to fetch study list.`, {
          variant: "error",
        })
        console.log(err)
      })
  }

  const updateStudyDetail = (
    studyId: number,
    oldVal: StudyDetails,
    setter: SetterOrUpdater<StudyDetails>
  ) => {
    getStudyDetailAPI(studyId)
      .then((study) => {
        let newVal = Object.assign({}, oldVal)
        newVal[studyId] = study
        setter(newVal)
      })
      .catch((err) => {
        enqueueSnackbar(`Failed to fetch study (id=${studyId})`, {
          variant: "error",
        })
        console.log(err)
      })
  }

  const createNewStudy = (
    study_name: string,
    direction: StudyDirection,
    oldVal: StudySummary[],
    setter: SetterOrUpdater<StudySummary[]>
  ) => {
    createNewStudyAPI(study_name, direction)
      .then((study_summary) => {
        const newVal = [...oldVal, study_summary]
        setter(newVal)
        enqueueSnackbar(
          `Success to create a study (study_name=${study_name})`,
          {
            variant: "success",
          }
        )
      })
      .catch((err) => {
        enqueueSnackbar(`Failed to create a study (study_name=${study_name})`, {
          variant: "error",
        })
        console.log(err)
      })
  }

  return {
    updateStudyDetail,
    updateStudySummaries,
    createNewStudy,
  }
}

export type Action = ReturnType<typeof actionCreator>
