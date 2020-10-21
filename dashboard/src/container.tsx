import React, { FC } from "react"
import { BrowserRouter as Router, Switch, Route } from "react-router-dom"
import { StudyDetail } from "./components/studyDetail"
import { StudyList } from "./components/studyList"

const AppContainer: FC<{}> = () => {
  return (
    <Router>
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
