import { SetterOrUpdater } from "recoil"
import { getStudyDetail, getStudySummaries } from "./utils/apiClient"
import { ReactNode } from "react"
import { OptionsObject } from "notistack"

export const actionCreator = (
  enqueueSnackbar: (
    message: ReactNode,
    options?: OptionsObject | undefined
  ) => string | number
) => {
  const updateStudySummaries = (setter: SetterOrUpdater<StudySummary[]>) => {
    getStudySummaries()
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
    getStudyDetail(studyId)
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

  return {
    updateStudyDetail,
    updateStudySummaries,
  }
}

export type Action = ReturnType<typeof actionCreator>
