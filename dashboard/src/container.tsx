import { jsx, css } from "@emotion/core"
import { FC, useEffect } from "react"
import { useRecoilState } from "recoil"
import { fetchStudySummariesAction } from "./api"
import { studySummariesState } from "./atom"

const style = css``

const AppContainer: FC<{}> = () => {
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
    }, 10 * 1000)
    return () => clearInterval(intervalId)
  })

  return (
    <div css={style}>
      <h1>Hello world</h1>
      {studySummaries.map((s: StudySummary) => {
        return (
          <p key={s.study_id}>
            {s.study_name} {s.best_trial?.value}
          </p>
        )
      })}
    </div>
  )
}

export { AppContainer }
