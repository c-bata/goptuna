import { jsx, css } from "@emotion/core"
import { FC } from "react"
import { BrowserRouter as Router, Switch, Route } from "react-router-dom"
import { StudyDetail } from "./components/studyDetail"
import { StudyList } from "./components/studyList"

const style = css``

const AppContainer: FC<{}> = () => {
  return (
    <Router css={style}>
      <Switch>
        <Route path="/studies/:studyId" children={<StudyDetail />} />
        <Route path="/">
          <StudyList />
        </Route>
      </Switch>
    </Router>
  )
}

export { AppContainer }
