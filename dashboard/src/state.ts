import { atom, useRecoilValue } from "recoil"

export const studySummariesState = atom<StudySummary[]>({
  key: "studySummaries",
  default: [],
})

export const studyDetailsState = atom<StudyDetails>({
  key: "studyDetails",
  default: {},
})

export const trialsUpdatingState = atom<{
  [trialId: string]: boolean
}>({
  key: "trialsUpdating",
  default: {},
})

export const reloadIntervalState = atom<number>({
  key: "reloadInterval",
  default: 10,
})

export const drawerOpenState = atom<boolean>({
  key: "drawerOpen",
  default: false,
})

export const isFileUploading = atom<boolean>({
  key: "isFileUploading",
  default: false,
})

export const artifactIsAvailable = atom<boolean>({
  key: "artifactIsAvailable",
  default: false,
})

export const useStudyDetailValue = (studyId: number): StudyDetail | null => {
  const studyDetails = useRecoilValue<StudyDetails>(studyDetailsState)
  return studyDetails[studyId] || null
}

export const useStudySummaryValue = (studyId: number): StudySummary | null => {
  const studySummaries = useRecoilValue<StudySummary[]>(studySummariesState)
  return studySummaries.find((s) => s.study_id == studyId) || null
}

export const useTrialUpdatingValue = (trialId: number): boolean => {
  const updating = useRecoilValue(trialsUpdatingState)
  return updating[trialId] || false
}

export const useStudyName = (studyId: number): string | null => {
  const studyDetail = useStudyDetailValue(studyId)
  const studySummary = useStudySummaryValue(studyId)
  return studyDetail?.name || studySummary?.study_name || null
}
