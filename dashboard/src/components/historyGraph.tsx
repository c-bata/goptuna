import { jsx, css } from "@emotion/core"
import { FC, useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import { useRecoilState } from "recoil"
import { studyDetailsState } from "../state"
import { fetchStudyDetailAction } from "../api"

interface ParamTypes {
  studyId: string
}

const style = css``

export const HistoryGraph: FC<{}> = () => {
  const { studyId } = useParams<ParamTypes>()
  const studyIdNumber = parseInt(studyId, 10)
  const [ready, setReady] = useState(false)
  const [studyDetails, setStudyDetails] = useRecoilState<StudyDetails>(
    studyDetailsState
  )

  useEffect(() => {
    // fetch immediately
    fetchStudyDetailAction(studyIdNumber)
      .then((study) => {
        let newStudies = Object.assign({}, studyDetails)
        newStudies[studyIdNumber] = study
        setStudyDetails(newStudies)
      })
      .catch((err) => {
        console.log(err) // Notify to error dispatchers
      })
    const intervalId = setInterval(function () {
      fetchStudyDetailAction(studyIdNumber)
        .then((study) => {
          let newStudies = Object.assign({}, studyDetails)
          newStudies[studyIdNumber] = study
          setStudyDetails(newStudies)
        })
        .catch((err) => {
          console.log(err) // Notify to error dispatchers
        })
    }, 1000)
    return () => clearInterval(intervalId)
  }, [])

  useEffect(() => {
    if (!ready && studyDetails[studyIdNumber]) {
      setReady(true)
    }
  }, [studyDetails])

  const studyDetail = studyDetails[studyIdNumber]
  const content = ready ? (
    studyDetail.trials.map((t) => {
      return (
        <li key={t.trial_id}>
          <p>{t.number}</p>
        </li>
      )
    })
  ) : (
    <p>Now loading...</p>
  )
  return (
    <div css={style}>
      <h1>Study {studyId}</h1>
      <ul>{content}</ul>
    </div>
  )
}
