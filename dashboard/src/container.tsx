import { jsx, css } from "@emotion/core"
import { FC, useEffect } from "react"
import { useRecoilState } from "recoil"
import { BrowserRouter as Router, Link, Switch, Route } from "react-router-dom"
import { fetchStudySummariesAction } from "./api"
import { studySummariesState } from "./atom"
import {HistoryGraph} from "./components/historyGraph";

const style = css`

`

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
    }, 1000)
    return () => clearInterval(intervalId)
  })

  return (
    <Router css={style}>
      <Switch>
        <Route path="/studies/:studyId" children={<HistoryGraph />} />
        <Route path="/about">
          <p>About</p>
        </Route>
        <Route path="/">
          <ul>
            <h1>List of studies.</h1>
            {studySummaries.map((s: StudySummary) => {
              return (
                <li key={s.study_id}>
                  <Link to={`/studies/${s.study_id}`}>{s.study_name}</Link>
                </li>
              )
            })}
          </ul>
        </Route>
      </Switch>
    </Router>
  )
}

export { AppContainer }
