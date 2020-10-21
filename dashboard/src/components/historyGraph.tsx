import { jsx, css } from "@emotion/core"
import { FC, useEffect } from "react"
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
  const [studyDetails, setStudyDetails] = useRecoilState<StudyDetails>(
    studyDetailsState
  )

  useEffect(() => {
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
  })

  return (
    <div css={style}>
      <h1>{studyId}</h1>
    </div>
  )
}
