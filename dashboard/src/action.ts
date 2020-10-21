import { SetterOrUpdater } from "recoil"
import { getStudyDetail, getStudySummaries } from "./apiClient"

export const updateStudySummaries = (
  setter: SetterOrUpdater<StudySummary[]>
) => {
  getStudySummaries()
    .then((studySummaries: StudySummary[]) => {
      setter(studySummaries)
    })
    .catch((err) => {
      console.log(err) // Notify to error dispatchers
    })
}

export const updateStudyDetail = (
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
      console.log(err) // Notify to error dispatchers
    })
}
