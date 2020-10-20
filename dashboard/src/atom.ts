import { atom } from "recoil"

export const studySummariesState = atom<StudySummary[]>({
  key: "studySummaries",
  default: [],
})

export const trialsState = atom<FrozenTrials>({
  key: "trials",
  default: {},
})
