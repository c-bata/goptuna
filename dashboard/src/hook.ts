import { useEffect } from "react"
import { useRecoilState } from "recoil"
import { studyDetailsState, studySummariesState } from "./state"
import { Action } from "./action"

export const useStudySummaries = (action: Action): StudySummary[] => {
  const [studySummaries, setStudySummaries] = useRecoilState<StudySummary[]>(
    studySummariesState
  )

  useEffect(() => {
    action.updateStudySummaries(setStudySummaries)
    const intervalId = setInterval(function () {
      action.updateStudySummaries(setStudySummaries)
    }, 10 * 1000)
    return () => clearInterval(intervalId)
  }, [])

  return studySummaries
}

export const useStudyDetail = (
  action: Action,
  studyId: number
): StudyDetail | null => {
  const [studyDetails, setStudyDetails] = useRecoilState<StudyDetails>(
    studyDetailsState
  )

  useEffect(() => {
    action.updateStudyDetail(studyId, studyDetails, setStudyDetails)
    const intervalId = setInterval(function () {
      action.updateStudyDetail(studyId, studyDetails, setStudyDetails)
    }, 5 * 1000)
    return () => clearInterval(intervalId)
  }, [])

  return studyDetails[studyId] ? studyDetails[studyId] : null
}
