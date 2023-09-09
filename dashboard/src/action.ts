import { useRecoilState, useSetRecoilState } from "recoil"
import { useSnackbar } from "notistack"
import {
  getStudyDetailAPI,
  getStudySummariesAPI,
  createNewStudyAPI,
  deleteStudyAPI,
  tellTrialAPI,
  saveTrialUserAttrsAPI,
} from "./apiClient"
import {
  studyDetailsState,
  studySummariesState,
  reloadIntervalState,
  trialsUpdatingState,
} from "./state"

const localStorageReloadInterval = "reloadInterval"

type LocalStorageReloadInterval = {
  reloadInterval?: number
}

export const actionCreator = () => {
  const { enqueueSnackbar } = useSnackbar()
  const [studySummaries, setStudySummaries] =
    useRecoilState<StudySummary[]>(studySummariesState)
  const [studyDetails, setStudyDetails] =
    useRecoilState<StudyDetails>(studyDetailsState)
  const setReloadInterval = useSetRecoilState<number>(reloadIntervalState)
  const setTrialsUpdating = useSetRecoilState(trialsUpdatingState)

  const setStudyDetailState = (studyId: number, study: StudyDetail) => {
    setStudyDetails((prevVal) => {
      const newVal = Object.assign({}, prevVal)
      newVal[studyId] = study
      return newVal
    })
  }

  const setTrialUpdating = (trialId: number, updating: boolean) => {
    setTrialsUpdating((prev) => {
      const newVal = Object.assign({}, prev)
      newVal[trialId] = updating
      return newVal
    })
  }

  const setTrialStateValues = (
    studyId: number,
    index: number,
    state: TrialState,
    value?: number
  ) => {
    const newTrial: Trial = Object.assign(
      {},
      studyDetails[studyId].trials[index]
    )
    newTrial.state = state
    newTrial.value = value
    const newTrials: Trial[] = [...studyDetails[studyId].trials]
    newTrials[index] = newTrial
    const newStudy: StudyDetail = Object.assign({}, studyDetails[studyId])
    newStudy.trials = newTrials

    // Update Best Trials
    if (state === "Complete") {
      // Single objective optimization
      const bestValue = newStudy.best_trial?.value
      const currentValue = value
      if (newStudy.best_trial === undefined) {
        newStudy.best_trial = newTrial
      } else if (bestValue !== undefined && currentValue !== undefined) {
        if (newStudy.direction === "minimize" && currentValue < bestValue) {
          newStudy.best_trial = newTrial
        } else if (
          newStudy.direction === "maximize" &&
          currentValue > bestValue
        ) {
          newStudy.best_trial = newTrial
        }
      }
    }
    setStudyDetailState(studyId, newStudy)
  }

  const setTrialUserAttrs = (
    studyId: number,
    index: number,
    user_attrs: { [key: string]: number | string }
  ) => {
    const newTrial: Trial = Object.assign(
      {},
      studyDetails[studyId].trials[index]
    )
    newTrial.user_attrs = Object.keys(user_attrs).map((key) => ({
      key: key,
      value: user_attrs[key].toString(),
    }))
    const newTrials: Trial[] = [...studyDetails[studyId].trials]
    newTrials[index] = newTrial
    const newStudy: StudyDetail = Object.assign({}, studyDetails[studyId])
    newStudy.trials = newTrials
    setStudyDetailState(studyId, newStudy)
  }

  const updateStudySummaries = (successMsg?: string) => {
    getStudySummariesAPI()
      .then((studySummaries: StudySummary[]) => {
        setStudySummaries(studySummaries)

        if (successMsg) {
          enqueueSnackbar(successMsg, { variant: "success" })
        }
      })
      .catch((err) => {
        enqueueSnackbar(`Failed to fetch study list.`, {
          variant: "error",
        })
        console.log(err)
      })
  }

  const updateStudyDetail = (studyId: number) => {
    getStudyDetailAPI(studyId)
      .then((study) => {
        setStudyDetailState(studyId, study)
      })
      .catch((err) => {
        const reason = err.response?.data.reason
        if (reason !== undefined) {
          enqueueSnackbar(`Failed to fetch study (reason=${reason})`, {
            variant: "error",
          })
        }
        console.log(err)
      })
  }

  const createNewStudy = (studyName: string, direction: StudyDirection) => {
    createNewStudyAPI(studyName, direction)
      .then((study_summary) => {
        const newVal = [...studySummaries, study_summary]
        setStudySummaries(newVal)
        enqueueSnackbar(`Success to create a study (study_name=${studyName})`, {
          variant: "success",
        })
      })
      .catch((err) => {
        enqueueSnackbar(`Failed to create a study (study_name=${studyName})`, {
          variant: "error",
        })
        console.log(err)
      })
  }

  const deleteStudy = (studyId: number) => {
    deleteStudyAPI(studyId)
      .then(() => {
        setStudySummaries(studySummaries.filter((s) => s.study_id !== studyId))
        enqueueSnackbar(`Success to delete a study (id=${studyId})`, {
          variant: "success",
        })
      })
      .catch((err) => {
        enqueueSnackbar(`Failed to delete study (id=${studyId})`, {
          variant: "error",
        })
        console.log(err)
      })
  }

  const loadReloadInterval = () => {
    const reloadIntervalJSON = localStorage.getItem(localStorageReloadInterval)
    if (reloadIntervalJSON === null) {
      return
    }
    const gp = JSON.parse(reloadIntervalJSON) as LocalStorageReloadInterval
    if (gp.reloadInterval !== undefined) {
      setReloadInterval(gp.reloadInterval)
    }
  }

  const saveReloadInterval = (interval: number) => {
    setReloadInterval(interval)
    const value: LocalStorageReloadInterval = {
      reloadInterval: interval,
    }
    localStorage.setItem(localStorageReloadInterval, JSON.stringify(value))
  }

  const makeTrialFail = (studyId: number, trialId: number): void => {
    const message = `id=${trialId}, state=Fail`
    setTrialUpdating(trialId, true)
    tellTrialAPI(trialId, "Fail")
      .then(() => {
        const index = studyDetails[studyId].trials.findIndex(
          (t) => t.trial_id === trialId
        )
        if (index === -1) {
          enqueueSnackbar(`Unexpected error happens. Please reload the page.`, {
            variant: "error",
          })
          return
        }
        setTrialStateValues(studyId, index, "Fail")
        enqueueSnackbar(`Successfully updated trial (${message})`, {
          variant: "success",
        })
      })
      .catch((err) => {
        setTrialUpdating(trialId, false)
        const reason = err.response?.data.reason
        enqueueSnackbar(
          `Failed to update trial (${message}). Reason: ${reason}`,
          {
            variant: "error",
          }
        )
        console.log(err)
      })
  }

  const makeTrialComplete = (
    studyId: number,
    trialId: number,
    value: number
  ): void => {
    const message = `id=${trialId}, state=Complete, value=${value}`
    setTrialUpdating(trialId, true)
    tellTrialAPI(trialId, "Complete", value)
      .then(() => {
        const index = studyDetails[studyId].trials.findIndex(
          (t) => t.trial_id === trialId
        )
        if (index === -1) {
          enqueueSnackbar(`Unexpected error happens. Please reload the page.`, {
            variant: "error",
          })
          return
        }
        setTrialStateValues(studyId, index, "Complete", value)
      })
      .catch((err) => {
        setTrialUpdating(trialId, false)
        const reason = err.response?.data.reason
        enqueueSnackbar(
          `Failed to update trial (${message}). Reason: ${reason}`,
          {
            variant: "error",
          }
        )
        console.log(err)
      })
  }

  const saveTrialUserAttrs = (
    studyId: number,
    trialId: number,
    user_attrs: { [key: string]: string | number }
  ): void => {
    const message = `id=${trialId}, user_attrs=${JSON.stringify(user_attrs)}`
    setTrialUpdating(trialId, true)
    saveTrialUserAttrsAPI(trialId, user_attrs)
      .then(() => {
        const index = studyDetails[studyId].trials.findIndex(
          (t) => t.trial_id === trialId
        )
        if (index === -1) {
          enqueueSnackbar(`Unexpected error happens. Please reload the page.`, {
            variant: "error",
          })
          return
        }
        setTrialUserAttrs(studyId, index, user_attrs)
        enqueueSnackbar(`Successfully updated trial (${message})`, {
          variant: "success",
        })
      })
      .catch((err) => {
        setTrialUpdating(trialId, false)
        const reason = err.response?.data.reason
        enqueueSnackbar(
          `Failed to update trial (${message}). Reason: ${reason}`,
          {
            variant: "error",
          }
        )
        console.log(err)
      })
  }

  return {
    updateStudyDetail,
    updateStudySummaries,
    createNewStudy,
    deleteStudy,
    loadReloadInterval,
    saveReloadInterval,
    makeTrialComplete,
    makeTrialFail,
    saveTrialUserAttrs,
  }
}

export type Action = ReturnType<typeof actionCreator>
