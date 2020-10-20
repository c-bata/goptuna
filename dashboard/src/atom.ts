import { atom } from "recoil"

export const studySummariesState = atom<StudySummary[]>({
  key: "studySummaries",
  default: [],
})

export const trialsState = atom<Studies>({
  key: "trials",
  default: {},
})
