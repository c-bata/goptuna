import { jsx, css } from "@emotion/core"
import { FC, useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import { useRecoilState } from "recoil"
import { studyDetailsState } from "../state"
import { updateStudyDetail } from "../action"
import {TrialsTable} from "./trialsTable";

interface ParamTypes {
  studyId: string
}

const style = css``

export const StudyDetail: FC<{}> = () => {
  const { studyId } = useParams<ParamTypes>()
  const studyIdNumber = parseInt(studyId, 10)
  const [ready, setReady] = useState(false)
  const [studyDetails, setStudyDetails] = useRecoilState<StudyDetails>(
    studyDetailsState
  )

  useEffect(() => {
    // fetch immediately
    updateStudyDetail(studyIdNumber, studyDetails, setStudyDetails)
    const intervalId = setInterval(function () {
      updateStudyDetail(studyIdNumber, studyDetails, setStudyDetails)
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
    <div>
      <TrialsTable trials={studyDetail.trials} />
    </div>
) : (
    <p>Now loading...</p>
  )
  return (
    <div css={style}>
      <h1>Study {studyId}</h1>
      {content}
    </div>
  )
}
