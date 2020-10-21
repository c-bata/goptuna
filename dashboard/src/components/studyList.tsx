import { jsx, css } from "@emotion/core"
import { FC, useEffect } from "react"
import { Link } from "react-router-dom"
import { useRecoilState } from "recoil"
import { studySummariesState } from "../state"
import { fetchStudySummariesAction } from "../api"

const style = css``

export const StudyList: FC<{}> = () => {
  const [studySummaries, setStudySummaries] = useRecoilState<StudySummary[]>(
    studySummariesState
  )

  useEffect(() => {
    const intervalId = setInterval(function () {
      fetchStudySummariesAction()
        .then((studySummaries: StudySummary[]) => {
          setStudySummaries(studySummaries)
        })
        .catch((err) => {
          console.log(err) // Notify to error dispatchers
        })
    }, 1000)
    return () => clearInterval(intervalId)
  })

  return (
    <div css={style}>
      <h1>List of studies.</h1>
      <ul>
        {studySummaries.map((s: StudySummary) => {
          return (
            <li key={s.study_id}>
              <Link to={`/studies/${s.study_id}`}>{s.study_name}</Link>
            </li>
          )
        })}
      </ul>
    </div>
  )
}
